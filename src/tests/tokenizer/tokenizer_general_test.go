package tokenizer_test

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	tz "tp/src/tokenizer"
	tk "tp/src/tokenizer/tokens"
	"tp/src/util"
)

func Test_dullTokenizer(t *testing.T) {
	tokenizer := tz.CreateDullTokenizer()
	filepath := "../exampleFiles/hello.java"
	text, err := util.GetTextOfFile(filepath)
	if err != nil {
		util.Error(fmt.Sprintf("Failed to find file: %s", filepath), err)
		assert.Fail(t, "No file found")
	}

	tokensScope := tokenizer.Tokenize(text)

	//jsonString := tokensScope.ToJsonString("testTagDull")
	//util.MakeDir("../../../output")
	//_ = util.CreateFileWithInfo("../../../output/dull_output.json", jsonString)

	for i := 0; i < tokensScope.Size(); i++ {
		st1, _ := tokensScope.At(i)
		switch i {
		// Line 1
		case 0:
			VerifyUnknownSymbol(t, st1, 1, 0, "/")
		case 1:
			VerifyUnknownSymbol(t, st1, 1, 0, "*")
		// Line 2
		case 2:
			VerifyUnknownKeyword(t, st1, 2, 0, "This")
		case 3:
			VerifyUnknownKeyword(t, st1, 2, 0, "class")
		case 4:
			VerifyUnknownKeyword(t, st1, 2, 0, "is")
		case 5:
			VerifyUnknownKeyword(t, st1, 2, 0, "called")
		case 6:
			VerifyUnknownKeyword(t, st1, 2, 0, "hello")
		case 7:
			VerifyUnknownSymbol(t, st1, 2, 0, "!")
		// Line 3
		case 8:
			VerifyUnknownKeyword(t, st1, 3, 0, "it")
		case 9:
			VerifyUnknownKeyword(t, st1, 3, 0, "is")
		case 10:
			VerifyUnknownKeyword(t, st1, 3, 0, "used")
		case 11:
			VerifyUnknownKeyword(t, st1, 3, 0, "to")
		case 12:
			VerifyUnknownKeyword(t, st1, 3, 0, "print")
		case 13:
			VerifyUnknownKeyword(t, st1, 3, 0, "out")
		case 14:
			VerifyUnknownSymbol(t, st1, 3, 0, "\"")
		case 15:
			VerifyUnknownKeyword(t, st1, 3, 0, "Hello")
		case 16:
			VerifyUnknownKeyword(t, st1, 3, 0, "World")
		case 17:
			VerifyUnknownSymbol(t, st1, 3, 0, "\"")
		// Line 4
		case 18:
			VerifyUnknownSymbol(t, st1, 4, 0, "*")
		case 19:
			VerifyUnknownSymbol(t, st1, 4, 0, "/")
		// Line 5
		// Line 6
		case 20:
			VerifyUnknownKeyword(t, st1, 6, 0, "public")
		case 21:
			VerifyUnknownKeyword(t, st1, 6, 0, "class")
		case 22:
			VerifyUnknownKeyword(t, st1, 6, 0, "hello")
		case 23:
			VerifyUnknownSymbol(t, st1, 6, 0, "{")
		// Line 7
		case 24:
			VerifyUnknownKeyword(t, st1, 7, 1, "public")
		case 25:
			VerifyUnknownKeyword(t, st1, 7, 1, "static")
		case 26:
			VerifyUnknownKeyword(t, st1, 7, 1, "void")
		case 27:
			VerifyUnknownKeyword(t, st1, 7, 1, "main")
		case 28:
			VerifyUnknownSymbol(t, st1, 7, 1, "(")
		case 29:
			VerifyUnknownKeyword(t, st1, 7, 1, "String")
		case 30:
			VerifyUnknownSymbol(t, st1, 7, 1, "[")
		case 31:
			VerifyUnknownSymbol(t, st1, 7, 1, "]")
		case 32:
			VerifyUnknownKeyword(t, st1, 7, 1, "args")
		case 33:
			VerifyUnknownSymbol(t, st1, 7, 1, ")")
		case 34:
			VerifyUnknownSymbol(t, st1, 7, 1, "{")
		// Line 8
		case 35:
			VerifyUnknownSymbol(t, st1, 8, 2, "/")
		case 36:
			VerifyUnknownSymbol(t, st1, 8, 2, "/")
		case 37:
			VerifyUnknownKeyword(t, st1, 8, 2, "This")
		case 38:
			VerifyUnknownKeyword(t, st1, 8, 2, "prints")
		case 39:
			VerifyUnknownKeyword(t, st1, 8, 2, "out")
		case 40:
			VerifyUnknownKeyword(t, st1, 8, 2, "stuff")
		// Line 9
		case 41:
			VerifyUnknownKeyword(t, st1, 9, 2, "System")
		case 42:
			VerifyUnknownSymbol(t, st1, 9, 2, ".")
		case 43:
			VerifyUnknownKeyword(t, st1, 9, 2, "out")
		case 44:
			VerifyUnknownSymbol(t, st1, 9, 2, ".")
		case 45:
			VerifyUnknownKeyword(t, st1, 9, 2, "println")
		case 46:
			VerifyUnknownSymbol(t, st1, 9, 2, "(")
		case 47:
			VerifyUnknownSymbol(t, st1, 9, 2, "\"")
		case 48:
			VerifyUnknownKeyword(t, st1, 9, 2, "Hello")
		case 49:
			VerifyUnknownKeyword(t, st1, 9, 2, "World")
		case 50:
			VerifyUnknownSymbol(t, st1, 9, 2, "\"")
		case 51:
			VerifyUnknownSymbol(t, st1, 9, 2, ")")
		case 52:
			VerifyUnknownSymbol(t, st1, 9, 2, ";")
		// Line 10
		case 53:
			VerifyUnknownSymbol(t, st1, 10, 1, "}")
		// Line 11
		case 54:
			VerifyUnknownSymbol(t, st1, 11, 0, "}")

		}
	}
}

func VerifyUnknownSymbol(t *testing.T, tkn *tk.Token, expectedLineNum int, expectedTabNum int, expectedText string) {
	assert.Equal(t, tz.RULENAME_SYMBOL, tkn.RuleName)
	assert.Equal(t, tz.SYMBOLIC_NAME_UNKNOWN_SYMBOL, tkn.SymbolicName)
	assert.Equal(t, expectedText, tkn.Text)
	assert.Equal(t, expectedLineNum, tkn.LineNumber)
	assert.Equal(t, expectedTabNum, tkn.TabNumber)
}

func VerifyUnknownKeyword(t *testing.T, tkn *tk.Token, expectedLineNum int, expectedTabNum int, expectedText string) {
	assert.Equal(t, tz.RULENAME_KEYWORD, tkn.RuleName)
	assert.Equal(t, tz.SYMBOLIC_NAME_NON_KEYWORD, tkn.SymbolicName)
	assert.Equal(t, expectedText, tkn.Text)
	assert.Equal(t, expectedLineNum, tkn.LineNumber)
	assert.Equal(t, expectedTabNum, tkn.TabNumber)
}
