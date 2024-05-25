package langs_test

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	pyTokenizer "tp/src/instances/langs/python/tokenizer"
	"tp/src/util"
)

func Test_pythonTokenizer(t *testing.T) {
	tokenizer := pyTokenizer.GetPythonTokenizer()
	filepath := "../exampleFiles/hello.py"
	text, err := util.GetTextOfFile(filepath)
	if err != nil {
		util.Error(fmt.Sprintf("Failed to find file: %s", filepath), err)
		assert.Fail(t, "No file found")
	}

	tokensScope := tokenizer.Tokenize(text)

	jsonString := tokensScope.ToJsonString("testTagPython")
	util.MakeDir("../../../output")
	_ = util.CreateFileWithInfo("../../../output/python_output.json", jsonString)
	//for i := 0; i < tokensScope.Size(); i++ {
	//	st1, _ := tokensScope.At(i)
	//	switch i {
	//	case 0:
	//		assert.Equal(t, "PUBLIC", st1.SymbolicName)
	//	}
	//
	//}
}
