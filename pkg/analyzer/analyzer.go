package analyzer

import (
	"go/ast"
	"go/token"

	"golang.org/x/tools/go/analysis"
)

var blockPos = []token.Pos{-1, -1}

// Analyzer - represents our linter
var Analyzer = &analysis.Analyzer{
	Name: "goexportedlocalvariables",
	Doc:  "Checks that there no exported local variables.",
	Run:  run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	helper := initHelper()

	inspect := func(node ast.Node) bool {
		switch node.(type) {
		case *ast.ValueSpec:
			return helper.checkDeclarations(node)
		case *ast.FuncDecl:
			return helper.checkLocalFuncVariables(node, pass)
		case *ast.AssignStmt:
			return helper.checkAssignments(node, pass)
		}

		return true
	}

	for _, f := range pass.Files {
		ast.Inspect(f, inspect)
	}

	return nil, nil
}
