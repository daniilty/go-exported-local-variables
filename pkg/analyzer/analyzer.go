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
	inspect := func(node ast.Node) bool {
		switch node.(type) {
		case *ast.FuncDecl:
			return checkLocalFuncVariables(node, pass)
		}

		return checkLocalAssignments(node, pass)
	}

	for _, f := range pass.Files {
		ast.Inspect(f, inspect)
	}

	return nil, nil
}

func checkLocalFuncVariables(node ast.Node, pass *analysis.Pass) bool {
	funcDecl, ok := node.(*ast.FuncDecl)
	if !ok {
		return true
	}

	for _, l := range funcDecl.Type.Params.List {
		for _, n := range l.Names {
			if token.IsExported(n.Name) {
				pass.Reportf(node.Pos(), "param %s in function %s should not be exported\n",
					n.Name, funcDecl.Name.Name)
			}
		}
	}

	funcBody := funcDecl.Body
	blockPos = []token.Pos{funcBody.Lbrace, funcBody.Rbrace}

	return true
}

func checkLocalAssignments(node ast.Node, pass *analysis.Pass) bool {
	if blockPos[0] < 0 || blockPos[1] < 0 {
		return true
	}

	assignment, ok := node.(*ast.AssignStmt)
	if !ok {
		return true
	}

	assignmentPos := assignment.Pos()
	if assignmentPos < blockPos[0] || assignmentPos > blockPos[1] {
		return true
	}

	variableName := assignment.Lhs[0].(*ast.Ident).Name
	if token.IsExported(variableName) {
		pass.Reportf(node.Pos(), "local variable %s should not be exported\n",
			assignment.Lhs[0])
	}

	return true
}

func changeBlockPosByIfBrackets(block *ast.BlockStmt) {
	blockPos = []token.Pos{block.Lbrace, block.Rbrace}
}
