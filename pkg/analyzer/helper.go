package analyzer

import (
	"go/ast"
	"go/token"

	"golang.org/x/tools/go/analysis"
)

type helper struct {
	blockPos   []token.Pos
	globalVars []string
}

func initHelper() *helper {
	return &helper{
		blockPos:   []token.Pos{-1, -1},
		globalVars: []string{},
	}
}

func (h *helper) checkLocalFuncVariables(node ast.Node, pass *analysis.Pass) bool {
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
	h.setBlockPos(funcBody)

	return true
}

func (h *helper) checkDeclarations(node ast.Node) bool {
	declaration, ok := node.(*ast.ValueSpec)
	if !ok {
		return true
	}

	for _, n := range declaration.Names {
		h.addToGlobalVarsIfExported(n.Name)
	}

	return true
}

func (h *helper) checkAssignments(node ast.Node, pass *analysis.Pass) bool {
	assignment, ok := node.(*ast.AssignStmt)
	if !ok {
		return true
	}

	if h.blockPos[0] < 0 || h.blockPos[1] < 0 {
		return h.checkGlobalAssignments(assignment)
	}

	assignmentPos := assignment.Pos()
	if assignmentPos < h.blockPos[0] || assignmentPos > h.blockPos[1] {
		return h.checkGlobalAssignments(assignment)
	}

	for _, l := range assignment.Lhs {
		switch l.(type) {
		case *ast.Ident:
			variableName := l.(*ast.Ident).Name
			h.reportIfExported(variableName, node, pass)
		}
	}

	return true
}

func (h *helper) checkGlobalAssignments(assignment *ast.AssignStmt) bool {
	for _, l := range assignment.Lhs {
		switch l.(type) {
		case *ast.Ident:
			variableName := l.(*ast.Ident).Name
			h.addToGlobalVarsIfExported(variableName)
		}
	}

	return true
}

func (h *helper) setBlockPos(body *ast.BlockStmt) {
	h.blockPos = []token.Pos{body.Lbrace, body.Rbrace}
}

func (h *helper) addToGlobalVarsIfExported(varName string) {
	if token.IsExported(varName) {
		h.globalVars = append(h.globalVars, varName)
	}
}

func (h *helper) isDeclaredGlobally(varName string) bool {
	for _, v := range h.globalVars {
		if v != varName {
			continue
		}

		return true
	}

	return false
}

func (h *helper) reportIfExported(varName string, node ast.Node, pass *analysis.Pass) {
	if token.IsExported(varName) {
		if h.isDeclaredGlobally(varName) {
			return
		}

		pass.Reportf(node.Pos(), "local variable %s should not be exported\n",
			varName)
	}
}
