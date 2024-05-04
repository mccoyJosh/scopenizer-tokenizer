package tests

import (
	"TokenAndParse/src/tokenizer/tokenizerFactory"
	"TokenAndParse/src/util"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_javaTokenizer(t *testing.T) {
	tokenizer := tokenizerFactory.GetTokenizer("java")
	filename := "exampleFiles/file.java"
	text, err := util.GetTextOfFile(filename)
	if err != nil {
		util.Error(fmt.Sprintf("Failed to find file: %s", filename), err)
		assert.Fail(t, "No file found")
	}

	tokens := tokenizer.Tokenize(text)
	for _, token := range tokens {
		token.PrintToken()
	}
}