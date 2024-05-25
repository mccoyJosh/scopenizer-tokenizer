package tokenizer_test

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
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

	jsonString := tokensScope.ToJsonString("testTagDull")
	util.MakeDir("../../../output")
	_ = util.CreateFileWithInfo("../../../output/dull_output.json", jsonString)

	for i := 0; i < tokensScope.Size(); i++ {
		//st1, _ := tokensScope.At(i)
		//assert.Equal(t, "IDENTIFIER", st1.SymbolicName)
		// TODO: This fails because symbols are not identifiers
	}
}
