package tokens

import (
	"errors"
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
