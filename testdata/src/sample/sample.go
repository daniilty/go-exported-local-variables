package sample

func AllFine() string {
	a := "That's nicely initialized local variable"

	return a
}

func Incorrect() string {
	A := "Oh it seems not so good"

	return A
}

func ExportedInIfStatement() string {
	if true {
		A := "Hehe try to find me here"
		return A
	}

	return ""
}

func UnExportedInFuncDeclaration(a, b string) {}

func ExportedInFuncDeclaration(A, B string) {}
