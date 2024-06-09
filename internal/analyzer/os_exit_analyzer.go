package analyzer

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
)

var OsExitAnalyzer = &analysis.Analyzer{
	Name: "errcheck",
	Doc:  "check for using os exit in main",
	Run:  run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		ast.Inspect(file, func(node ast.Node) bool {
			switch x := node.(type) {
			case *ast.File:
				if x.Name.Name != "main" {
					return false
				}
			case *ast.FuncDecl:
				if x.Name.Name != "main" {
					return false
				}
			case *ast.SelectorExpr: // выражение
				if x.Sel.Name == "Exit" {
					if ident, ok := x.X.(*ast.Ident); ok {
						if ident.Name == "os" {
							pass.Reportf(x.Pos(), "expression os.exit deprecated in main function")
						}
					}
				}
			}
			return true
		})
	}
	return nil, nil
}
