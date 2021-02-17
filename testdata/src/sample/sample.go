package sample

func AllFine() string {
	a := "That's nicely initialized local variable"

	return a
}

func Incorrect() string {
	A := "Oh it seems not so good" // want `local variable A should not be exported`

	return A
}

func ExportedInIfStatement() string {
	if true {
		A := "Hehe try to find me here" // want `local variable A should not be exported`
		return A
	}

	return ""
}

func UnExportedInFuncDeclaration(a, b string) {}

func ExportedInFuncDeclaration(A, B string) {} // want `param A in function ExportedInFuncDeclaration should not be exported` `param B in function ExportedInFuncDeclaration should not be exported`

var DeclaredGlobally string

var AssignedGlobally = "It will pass"

func ChangingDeclaredGloballyVar() {
	DeclaredGlobally = "It's so tight here"
}
