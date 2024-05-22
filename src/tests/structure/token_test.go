package structure

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"tp/src/tokenizer/tokens"
)

func Test_CreateUnidentifiedToken(t *testing.T) {
	token := tokens.CreateUnidentifiedToken("the text", 1, 2)
	assert.Equal(t, 1, token.LineNumber)
	assert.Equal(t, 2, token.TabNumber)
	assert.Equal(t, "Unidentified", token.SymbolicName)
	assert.Equal(t, "unidentified", token.RuleName)
	assert.Equal(t, "the text", token.Text)
}

func Test_Equal_True(t *testing.T) {
	token1 := tokens.CreateUnidentifiedToken("the text", 1, 2)
	token1.SetValues("KEYWORD", "IF")

	token2 := tokens.CreateUnidentifiedToken("the text", 1, 2)
	token2.SetValues("KEYWORD", "IF")

	assert.True(t, token1.Equal(token2))
}

func Test_Equal_False(t *testing.T) {
	token1 := tokens.CreateUnidentifiedToken("the text", 1, 2)
	token1.SetValues("KEYWORD", "IF")

	token2 := tokens.CreateUnidentifiedToken("the text", 1, 2)
	token2.SetValues("KEYWORD", "ELSE")

	assert.False(t, token1.Equal(token2))
}

func Test_Equal_False_2(t *testing.T) {
	token1 := tokens.CreateUnidentifiedToken("the text", 1, 2)
	token1.SetValues("KEYWORD", "IF")

	token2 := tokens.CreateUnidentifiedToken("the text", 1, 2)
	token2.SetValues("SYMBOL", "IF")

	assert.False(t, token1.Equal(token2))
}

func Test_Equal_False_3(t *testing.T) {
	token1 := tokens.CreateUnidentifiedToken("the text", 1, 2)
	token1.SetValues("KEYWORD", "IF")

	token2 := tokens.CreateUnidentifiedToken("the text", 1, 2)
	token2.SetValues("SYMBOL", "ELSE")

	assert.False(t, token1.Equal(token2))
}

func Test_ValidScopeToken_True(t *testing.T) {
	token := tokens.InitScopeToken()
	assert.True(t, token.ValidScopeToken())
}

func Test_ValidScopeToken_False(t *testing.T) {
	token := tokens.CreateUnidentifiedToken("the text", 1, 2)
	token.SetValues("KEYWORD", "IF")
	assert.False(t, token.ValidScopeToken())
}

func Test_GetScope(t *testing.T) {
	token := tokens.InitScopeToken().GetScopeToken()
	assert.Equal(t, 0, token.Size())
}

func Test_GetScope_DoesNotExist(t *testing.T) {
	// Along with passing this test, an error message should be produced by doing this. Check
	// results of Test_GetScope_DoesNotExist for error message
	token := tokens.CreateUnidentifiedToken("the text", 1, 2)
	token.SetValues("KEYWORD", "IF")

	nullScope := token.GetScopeToken()
	var nullExampleScope *tokens.ScopeObj
	assert.Equal(t, nullExampleScope, nullScope)
}
