package gotrics

import (
	"go/ast"
	"go/token"
	"math"
)

type (
	GoMetrics struct {
		Name           string
		PosLine        int
		PosColumn      int
		MethodLength   int
		NestingLevel   int
		ParameterCount int
		ABCSize        float64
	}
)

func ABCSize(f *ast.FuncDecl) float64 {
	var assignment = 0.0
	var branch = 0.0
	var condition = 0.0

	ast.Inspect(f, func(n ast.Node) bool {
		switch r := n.(type) {
		case *ast.ValueSpec:
			if r.Values != nil {
				for _, i := range r.Names {
					// not count blank identifier or const declarations
					if i.Name != "_" && i.Obj.Kind != ast.Con {
						assignment++
					}
				}
			}
		case *ast.AssignStmt:
			for _, e := range r.Lhs {
				if i, ok := e.(*ast.Ident); ok {
					// not count blank identifier
					if i.Name != "_" {
						assignment++
					}
				}
			}
		case *ast.IncDecStmt:
			assignment++
		case *ast.CallExpr:
			branch++
		case *ast.BranchStmt:
			if r.Tok == token.GOTO {
				branch++
			}
		case *ast.IfStmt:
			condition++

			if _, ok := r.Else.(*ast.BlockStmt); ok {
				condition++
			}
		case *ast.ForStmt:
			if r.Cond != nil {
				condition++
			}
		case *ast.CaseClause, *ast.CommClause:
			condition++
		}
		return true
	})

	k := math.Sqrt(assignment*assignment + branch*branch + condition*condition)
	return math.Round(k*100) / 100
}

func MethodLength(fset *token.FileSet, n *ast.FuncDecl) int {
	sp := fset.Position(n.Body.Lbrace)
	ep := fset.Position(n.Body.Rbrace)
	length := ep.Line - sp.Line + 1

	return length
}

// `switch`, `type switch`, `select` are not nesting in formatted code (gofmt)
func MethodNesting(f *ast.FuncDecl) int {
	return nestWalk(f, 0)
}

func ParameterList(n *ast.FuncDecl) int {
	var count = 0

	for _, l := range n.Type.Params.List {
		for _, m := range l.Names {
			// not count blank identifier
			if m.Name != "_" {
				count++
			}
		}
	}

	return count
}

func Analyze(fset *token.FileSet, f *ast.File) []GoMetrics {
	report := make([]GoMetrics, 0)

	ast.Inspect(f, func(n ast.Node) bool {
		if r, ok := n.(*ast.FuncDecl); ok {
			namePos := fset.Position(r.Name.NamePos)
			gm := GoMetrics{}
			gm.Name = r.Name.Name
			gm.PosLine = namePos.Line
			gm.PosColumn = namePos.Column
			gm.MethodLength = MethodLength(fset, r)
			gm.NestingLevel = MethodNesting(r)
			gm.ParameterCount = ParameterList(r)
			gm.ABCSize = ABCSize(r)
			report = append(report, gm)
		}
		return true
	})

	return report
}

func nestWalk(node ast.Node, level int) int {
	ast.Inspect(node, func(n ast.Node) bool {
		var currentLevel = level

		switch r := n.(type) {
		case *ast.BlockStmt:
			for _, x := range r.List {
				level = int(math.Max(float64(nestWalk(x, currentLevel+1)), float64(level)))
			}
			return false
		case *ast.SwitchStmt:
			if r.Init != nil {
				level = int(math.Max(float64(nestWalk(r.Init, currentLevel-1)), float64(level)))
			}
			if r.Tag != nil {
				level = int(math.Max(float64(nestWalk(r.Tag, currentLevel-1)), float64(level)))
			}
			level = int(math.Max(float64(nestWalk(r.Body, currentLevel-1)), float64(level)))
			return false
		case *ast.CaseClause:
			for _, x := range r.List {
				level = int(math.Max(float64(nestWalk(x, currentLevel+1)), float64(level)))
			}
			for _, x := range r.Body {
				level = int(math.Max(float64(nestWalk(x, currentLevel+1)), float64(level)))
			}
			return false
		case *ast.SelectStmt:
			level = int(math.Max(float64(nestWalk(r.Body, currentLevel-1)), float64(level)))
			return false
		case *ast.CommClause:
			if r.Comm != nil {
				level = int(math.Max(float64(nestWalk(r.Comm, currentLevel-1)), float64(level)))
			}
			for _, x := range r.Body {
				level = int(math.Max(float64(nestWalk(x, currentLevel+1)), float64(level)))
			}
			return false
		case *ast.TypeSwitchStmt:
			if r.Init != nil {
				level = int(math.Max(float64(nestWalk(r.Init, currentLevel-1)), float64(level)))
			}
			level = int(math.Max(float64(nestWalk(r.Assign, currentLevel-1)), float64(level)))
			level = int(math.Max(float64(nestWalk(r.Body, currentLevel-1)), float64(level)))
			return false
		}
		return true
	})

	return level
}
