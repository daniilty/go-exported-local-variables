package analyzer

import (
	"go/ast"
	"go/token"

	"golang.org/x/tools/go/analysis"
)

type helper struct {
	blockPos    []token.Pos
	globalVars  []string
	ignoredVars []string
}

func initHelper() *helper {
	return &helper{
		blockPos:    []token.Pos{-1, -1},
		globalVars:  []string{},
		ignoredVars: []string{},
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

func (h *helper) checkDeclarations(node ast.Node, pass *analysis.Pass) bool {
	declaration, ok := node.(*ast.ValueSpec)
	if !ok {
		return true
	}

	names := declaration.Names

	declarationPos := declaration.Pos()
	if declarationPos >= h.blockPos[0] && declarationPos <= h.blockPos[1] {
		return h.reportAndIgnoreNames(names, node, pass)
	}

	for _, n := range names {
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
		switch l := l.(type) {
		case *ast.Ident:
			variableName := l.Name
			h.reportIfExported(variableName, node, pass)
		}
	}

	return true
}

func (h *helper) checkGlobalAssignments(assignment *ast.AssignStmt) bool {
	for _, l := range assignment.Lhs {
		switch l := l.(type) {
		case *ast.Ident:
			variableName := l.Name
			h.addToGlobalVarsIfExported(variableName)
		}
	}

	return true
}

func (h *helper) reportAndIgnoreNames(names []*ast.Ident, node ast.Node, pass *analysis.Pass) bool {
	for _, n := range names {
		h.reportIfExported(n.Name, node, pass)
		h.ignore(n.Name)
	}

	return true
}

func (h *helper) ignore(varName string) {
	h.ignoredVars = append(h.ignoredVars, varName)
}

func (h *helper) checkIgnoredAndClear(node ast.Node) {
	if node == nil {
		return
	}

	if needsToBeIgnored(node) {
		return
	}

	_, ok := node.(*ast.BlockStmt)
	if !ok {
		return
	}

	h.clearIgnoredList()
}

func (h *helper) clearIgnoredList() {
	h.ignoredVars = []string{}
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

func (h *helper) isIgnored(varName string) bool {
	for _, v := range h.ignoredVars {
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

		if h.isIgnored(varName) {
			return
		}

		pass.Reportf(node.Pos(), "local variable %s should not be exported\n",
			varName)
	}
}

func needsToBeIgnored(node ast.Node) bool {
	_, ok := node.(*ast.ForStmt)
	if ok {
		return true
	}
	_, ok = node.(*ast.SelectStmt)
	if ok {
		return true
	}
	_, ok = node.(*ast.SwitchStmt)
	if ok {
		return true
	}
	_, ok = node.(*ast.IfStmt)
	if ok {
		return true
	}
	_, ok = node.(*ast.TypeSwitchStmt)

	return ok
}
