package gotrics

import (
	"bytes"
	"go/ast"
	"go/format"
	"go/token"
	"math"
	"strings"
)

type (
	GoMetrics struct {
		Name           string
		PosLine        int
		PosColumn      int
		FuncLength     int
		NestingLevel   int
		ParameterCount int
		ABCSize        float64
	}
)

func Analyze(fset *token.FileSet, f *ast.File) []GoMetrics {
	report := make([]GoMetrics, 0)

	ast.Inspect(f, func(n ast.Node) bool {
		if r, ok := n.(*ast.FuncDecl); ok {
			namePos := fset.Position(r.Name.NamePos)
			gm := GoMetrics{}
			gm.Name = r.Name.Name
			gm.PosLine = namePos.Line
			gm.PosColumn = namePos.Column
			gm.FuncLength = FuncLength(fset, r)
			gm.NestingLevel = FuncNesting(fset, r)
			gm.ParameterCount = ParameterCount(r)
			gm.ABCSize = ABCSize(r)
			report = append(report, gm)
		}
		return true
	})

	return report
}

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

func FuncLength(fset *token.FileSet, n *ast.FuncDecl) int {
	sp := fset.Position(n.Body.Lbrace)
	ep := fset.Position(n.Body.Rbrace)
	length := ep.Line - sp.Line + 1

	return length
}

// `switch`, `type switch`, `select` are not nesting in formatted code (gofmt)
func FuncNesting(fset *token.FileSet, f *ast.FuncDecl) int {
	var level = 0.0
	var out = new(bytes.Buffer)
	format.Node(out, fset, f)

	for _, l := range strings.Split(out.String(), "\n") {
		level = math.Max(float64(countLeadingTab(l)), level)
	}

	return int(level)
}

func ParameterCount(n *ast.FuncDecl) int {
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

func countLeadingTab(line string) int {
	var count = 0

	for _, runeValue := range line {
		if runeValue == '\t' {
			count++
		} else {
			break
		}
	}

	return count
}
