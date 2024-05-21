package tokens

import (
	"errors"
	"fmt"
	"tp/src/util"
)

type Token struct {
	LineNumber    int
	TabNumber     int
	BracketNumber int
	SymbolicName  string
	RuleName      string
	Text          string
	scopeToken    *ScopeObj
}

func CreateUnidentifiedToken(text string, lineNumber int, tabNum int, bracketNumber int) Token {
	return Token{
		LineNumber:    lineNumber,
		TabNumber:     tabNum,
		BracketNumber: bracketNumber,
		SymbolicName:  "Unidentified",
		RuleName:      "unidentified",
		Text:          text,
		scopeToken:    nil,
	}
}

func (t *Token) PrintToken() {
	// TODO Fix this
	fmt.Printf("Token:\n\ttxt:%s\tln:%d\ttb:%d\tbrkc:%d\tsn:%s\trn:%s\n", t.Text, t.LineNumber, t.TabNumber, t.BracketNumber, t.SymbolicName, t.RuleName)
}

func (t *Token) Equal(t2 Token) bool {
	return t.SymbolicName == t2.SymbolicName && t.RuleName == t2.RuleName
}

func (t *Token) SetValues(ruleName string, symbolicName string) {
	t.RuleName = ruleName
	t.SymbolicName = symbolicName
}

// ValidScopeToken
// Returns true if the provided token is a scope token; false otherwise.
func (t *Token) ValidScopeToken() bool {
	return t.scopeToken != nil
}

func (t *Token) GetScopeToken() *ScopeObj {
	if !t.ValidScopeToken() {
		err := errors.New("attempting to get scope token from non-scope-token Token")
		util.Error(err.Error(), err)
	}
	return t.scopeToken
}
