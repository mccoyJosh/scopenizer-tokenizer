package java

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	javaTokenizer "tp/src/instances/langs/java/tokenizer"
	tz "tp/src/tokenizer"
	"tp/src/util"
)

func Test_javaTokenizer(t *testing.T) {
	tokenizer := javaTokenizer.GetJavaTokenizer()
	filepath := "../../exampleFiles/hello.java"
	text, err := util.GetTextOfFile(filepath)
	if err != nil {
		util.Error(fmt.Sprintf("Failed to find file: %s", filepath), err)
		assert.Fail(t, "No file found")
	}

	tokensScope := tokenizer.Tokenize(text)

	jsonString := tokensScope.ToCustomTokenArrayTypeString("testTag")
	util.MakeDir("../../../../output")
	_ = util.CreateFileWithInfo("../../../../output/java_output.json", jsonString)
	//for i := 0; i < tokensScope.Size(); i++ {
	//	st1, _ := tokensScope.At(i)
	//	switch i {
	//	case 0:
	//		assert.Equal(t, "PUBLIC", st1.SymbolicName)
	//	}
	//
	//}
}

func Test_DullTokenizer(t *testing.T) {
	tokenizer := tz.CreateDullTokenizer()
	filepath := "../../exampleFiles/hello.java"
	text, err := util.GetTextOfFile(filepath)
	if err != nil {
		util.Error(fmt.Sprintf("Failed to find file: %s", filepath), err)
		assert.Fail(t, "No file found")
	}

	tokensScope := tokenizer.Tokenize(text)
	for i := 0; i < tokensScope.Size(); i++ {
		//st1, _ := tokensScope.At(i)
		//assert.Equal(t, "IDENTIFIER", st1.SymbolicName)
		// TODO: This fails because symbols are not identifiers
	}
}
