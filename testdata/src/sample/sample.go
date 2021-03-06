package sample

// AllFine - will pass
func AllFine() string {
	a := "That's nicely initialized local variable"

	return a
}

// Incorrect - reports error
func Incorrect() string {
	A := "Oh it seems not so good" // want `local variable A should not be exported`

	return A
}

// ExportedInIfStatement - reports error
func ExportedInIfStatement() string {
	if true {
		A := "Hehe try to find me here" // want `local variable A should not be exported`
		return A
	}

	return ""
}

// UnExportedInFuncDeclaration - will pass
func UnExportedInFuncDeclaration(a, b string) {}

// ExportedInFuncDeclaration - will pass
func ExportedInFuncDeclaration(A, B string) {} // want `param A in function ExportedInFuncDeclaration should not be exported` `param B in function ExportedInFuncDeclaration should not be exported`

// DeclaredGlobally - invisible to linter
var DeclaredGlobally string

// AssignedGlobally - invisible to linter
var AssignedGlobally = "It will pass"

// ChangingDeclaredGloballyVar - when declared globally it's invisible to linter
func ChangingDeclaredGloballyVar() {
	DeclaredGlobally = "It's so tight here"
}

// Name says all for it
func DeclaredInsideFuncAndInsideFuncDecl(A, b string) []string { // want `param A in function DeclaredInsideFuncAndInsideFuncDecl should not be exported`
	var F, N string // want `local variable F should not be exported` `local variable N should not be exported`
	F = "Skip"
	N = "Skip"
	B := "Err" // want `local variable B should not be exported`
	a := "Fine"

	slice := []string{F, N, B, a} // Fine

	return slice
}
