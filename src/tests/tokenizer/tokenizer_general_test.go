package tokenizer_test

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"tp/src/tests"
	tz "tp/src/tokenizer"
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

	// FOR DEBUGGING
	//jsonString := tokensScope.ToJsonString("testTagDull")
	//util.MakeDir("../../../output")
	//_ = util.CreateFileWithInfo("../../../output/dull_output.json", jsonString)

	assert.Equal(t, 55, tokensScope.Size())
	for i := 0; i < tokensScope.Size(); i++ {
		st1, _ := tokensScope.At(i)
		switch i {
		// Line 1
		case 0:
			tests.VerifyUnknownSymbol(t, st1, 1, 0, "/")
		case 1:
			tests.VerifyUnknownSymbol(t, st1, 1, 0, "*")
		// Line 2
		case 2:
			tests.VerifyUnknownKeyword(t, st1, 2, 0, "This")
		case 3:
			tests.VerifyUnknownKeyword(t, st1, 2, 0, "class")
		case 4:
			tests.VerifyUnknownKeyword(t, st1, 2, 0, "is")
		case 5:
			tests.VerifyUnknownKeyword(t, st1, 2, 0, "called")
		case 6:
			tests.VerifyUnknownKeyword(t, st1, 2, 0, "hello")
		case 7:
			tests.VerifyUnknownSymbol(t, st1, 2, 0, "!")
		// Line 3
		case 8:
			tests.VerifyUnknownKeyword(t, st1, 3, 0, "it")
		case 9:
			tests.VerifyUnknownKeyword(t, st1, 3, 0, "is")
		case 10:
			tests.VerifyUnknownKeyword(t, st1, 3, 0, "used")
		case 11:
			tests.VerifyUnknownKeyword(t, st1, 3, 0, "to")
		case 12:
			tests.VerifyUnknownKeyword(t, st1, 3, 0, "print")
		case 13:
			tests.VerifyUnknownKeyword(t, st1, 3, 0, "out")
		case 14:
			tests.VerifyUnknownSymbol(t, st1, 3, 0, "\"")
		case 15:
			tests.VerifyUnknownKeyword(t, st1, 3, 0, "Hello")
		case 16:
			tests.VerifyUnknownKeyword(t, st1, 3, 0, "World")
		case 17:
			tests.VerifyUnknownSymbol(t, st1, 3, 0, "\"")
		// Line 4
		case 18:
			tests.VerifyUnknownSymbol(t, st1, 4, 0, "*")
		case 19:
			tests.VerifyUnknownSymbol(t, st1, 4, 0, "/")
		// Line 5
		// Line 6
		case 20:
			tests.VerifyUnknownKeyword(t, st1, 6, 0, "public")
		case 21:
			tests.VerifyUnknownKeyword(t, st1, 6, 0, "class")
		case 22:
			tests.VerifyUnknownKeyword(t, st1, 6, 0, "hello")
		case 23:
			tests.VerifyUnknownSymbol(t, st1, 6, 0, "{")
		// Line 7
		case 24:
			tests.VerifyUnknownKeyword(t, st1, 7, 1, "public")
		case 25:
			tests.VerifyUnknownKeyword(t, st1, 7, 1, "static")
		case 26:
			tests.VerifyUnknownKeyword(t, st1, 7, 1, "void")
		case 27:
			tests.VerifyUnknownKeyword(t, st1, 7, 1, "main")
		case 28:
			tests.VerifyUnknownSymbol(t, st1, 7, 1, "(")
		case 29:
			tests.VerifyUnknownKeyword(t, st1, 7, 1, "String")
		case 30:
			tests.VerifyUnknownSymbol(t, st1, 7, 1, "[")
		case 31:
			tests.VerifyUnknownSymbol(t, st1, 7, 1, "]")
		case 32:
			tests.VerifyUnknownKeyword(t, st1, 7, 1, "args")
		case 33:
			tests.VerifyUnknownSymbol(t, st1, 7, 1, ")")
		case 34:
			tests.VerifyUnknownSymbol(t, st1, 7, 1, "{")
		// Line 8
		case 35:
			tests.VerifyUnknownSymbol(t, st1, 8, 2, "/")
		case 36:
			tests.VerifyUnknownSymbol(t, st1, 8, 2, "/")
		case 37:
			tests.VerifyUnknownKeyword(t, st1, 8, 2, "This")
		case 38:
			tests.VerifyUnknownKeyword(t, st1, 8, 2, "prints")
		case 39:
			tests.VerifyUnknownKeyword(t, st1, 8, 2, "out")
		case 40:
			tests.VerifyUnknownKeyword(t, st1, 8, 2, "stuff")
		// Line 9
		case 41:
			tests.VerifyUnknownKeyword(t, st1, 9, 2, "System")
		case 42:
			tests.VerifyUnknownSymbol(t, st1, 9, 2, ".")
		case 43:
			tests.VerifyUnknownKeyword(t, st1, 9, 2, "out")
		case 44:
			tests.VerifyUnknownSymbol(t, st1, 9, 2, ".")
		case 45:
			tests.VerifyUnknownKeyword(t, st1, 9, 2, "println")
		case 46:
			tests.VerifyUnknownSymbol(t, st1, 9, 2, "(")
		case 47:
			tests.VerifyUnknownSymbol(t, st1, 9, 2, "\"")
		case 48:
			tests.VerifyUnknownKeyword(t, st1, 9, 2, "Hello")
		case 49:
			tests.VerifyUnknownKeyword(t, st1, 9, 2, "World")
		case 50:
			tests.VerifyUnknownSymbol(t, st1, 9, 2, "\"")
		case 51:
			tests.VerifyUnknownSymbol(t, st1, 9, 2, ")")
		case 52:
			tests.VerifyUnknownSymbol(t, st1, 9, 2, ";")
		// Line 10
		case 53:
			tests.VerifyUnknownSymbol(t, st1, 10, 1, "}")
		// Line 11
		case 54:
			tests.VerifyUnknownSymbol(t, st1, 11, 0, "}")

		}
	}
}

func Test_dullTokenizer_dotText(t *testing.T) {
	tokenizer := tz.CreateDullTokenizer()
	filepath := "../exampleFiles/words.txt"
	text, err := util.GetTextOfFile(filepath)
	if err != nil {
		util.Error(fmt.Sprintf("Failed to find file: %s", filepath), err)
		assert.Fail(t, "No file found")
	}

	tokensScope := tokenizer.Tokenize(text)

	// FOR DEBUGGING
	//jsonString := tokensScope.ToJsonString("testTagDull")
	//util.MakeDir("../../../output")
	//_ = util.CreateFileWithInfo("../../../output/dull_output_txt.json", jsonString)
	assert.Equal(t, 7, tokensScope.Size())
}
