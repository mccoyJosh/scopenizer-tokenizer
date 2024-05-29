package tests

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	tk "tp/src/tokenizer/tokens"
)

func ValidateToken(t *testing.T, tkn *tk.Token, expectedLineNum int, expectedTabNum int, expectedRuleName string, expectedSymbolicName string, expectedText string) {
	expectedToken := tk.CreateUnidentifiedToken(expectedText, expectedLineNum, expectedTabNum)
	expectedToken.SetValues(expectedRuleName, expectedSymbolicName)
	invalidTokenStr := fmt.Sprintf("This token was invalid: %s\nThis was the expected token: %s", tkn.ToString(), expectedToken.ToString())

	assert.Equal(t, expectedRuleName, tkn.RuleName, invalidTokenStr)
	assert.Equal(t, expectedSymbolicName, tkn.SymbolicName, invalidTokenStr)
	assert.Equal(t, expectedText, tkn.Text, invalidTokenStr)
	assert.Equal(t, expectedLineNum, tkn.LineNumber, invalidTokenStr)
	assert.Equal(t, expectedTabNum, tkn.TabNumber, invalidTokenStr)
}
