package javaTokenizer

var Symbols = [][]string{
	{"\"", "DoubleQ"},
	{"'", "SingleQ"},
	{":", "Colon"},
	{",", "Comma"},
	{".", "Period"},
	{"*", "Star"},
	{"+", "Addition"},
	{"-", "Subtraction"},
	{"/", "ForwardSlash"},
	{"<", "LessThan"},
	{">", "GreaterThan"},
	{"\\", "Backslash"},
	{"?", "Question"},
	{";", "SemiColon"},
	{"^", "Exponent"},
	{"=", "Equal"},
	{"(", "LParen"},
	{")", "RParen"},
	{"[", "LBracket"},
	{"]", "RBracket"},
	{" ", "Space"}, // Maybe remove this and ignore spaces, essentially (split by spaces)
}
