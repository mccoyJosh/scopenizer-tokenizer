package tokens

type Token struct {
	LineNumber    int
	TabNumber     int
	BracketNumber int
	SymbolicName  string
	RuleName      string
	Text          string
}
