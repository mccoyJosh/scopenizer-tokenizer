package tokens

import "fmt"

func CreateUnidentifiedToken(text string, lineNumber int, tabNum int, bracketNumber int) Token {
	return Token{
		LineNumber:    lineNumber,
		TabNumber:     tabNum,
		BracketNumber: bracketNumber,
		SymbolicName:  "Unidentified",
		RuleName:      "unidentified",
		Text:          text,
	}
}

func (t *Token) PrintToken() {
	fmt.Printf("Token:\n\ttxt:%s\tln:%d\ttb:%d\tbrkc:%d\tsn:%s\trn:%s\n", t.Text, t.LineNumber, t.TabNumber, t.BracketNumber, t.SymbolicName, t.RuleName)
}

func (t *Token) Equal(t2 *Token) bool {
	return t.SymbolicName == t2.SymbolicName && t.RuleName == t2.RuleName
}
