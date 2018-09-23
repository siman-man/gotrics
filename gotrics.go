package gotrics

import (
	"fmt"
	"go/ast"
	"go/token"
	"math"
	"strconv"
)

type (
	GoMetrics struct {
		Name    string
		Length  int
		Level   int
		Count   int
		ABCSize float64
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
					// not count blank identifier
					if i.Name != "_" {
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
		case *ast.CaseClause:
			// nil means default case, not count condition
			if r.List != nil {
				condition++
			}
		}
		return true
	})

	k := math.Sqrt(assignment*assignment + branch*branch + condition*condition)
	i := fmt.Sprintf("%.2f", k)
	b, _ := strconv.ParseFloat(i, 64)
	return b
}

func MethodLength(fset *token.FileSet, n *ast.FuncDecl) int {
	sp := fset.Position(n.Body.Lbrace)
	ep := fset.Position(n.Body.Rbrace)
	length := ep.Line - sp.Line + 1

	return length
}

// `switch`, `type switch`, `select` are not nesting in formatted code (gofmt)
func MethodNesting(f *ast.FuncDecl) int {
	var level = 0

	ast.Inspect(f, func(n ast.Node) bool {
		switch n.(type) {
		case *ast.BlockStmt:
			level++
		case *ast.SwitchStmt:
			level--
		case *ast.SelectStmt:
			level--
		case *ast.TypeSwitchStmt:
			level--
		}
		return true
	})

	return level
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
			gm := GoMetrics{}
			gm.Name = r.Name.Name
			gm.Length = MethodLength(fset, r)
			gm.Level = MethodNesting(r)
			gm.Count = ParameterList(r)
			gm.ABCSize = ABCSize(r)
			report = append(report, gm)
		}
		return true
	})

	return report
}
