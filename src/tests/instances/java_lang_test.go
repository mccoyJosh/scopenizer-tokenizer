package instances_test

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	javaTokenizer "tp/src/instances/langs/java"
	"tp/src/tests"
	tz "tp/src/tokenizer"
	"tp/src/util"
)

func Test_javaTokenizer(t *testing.T) {
	tokenizer := javaTokenizer.GetJavaTokenizer()
	filepath := "../exampleFiles/hello.java"
	text, err := util.GetTextOfFile(filepath)
	if err != nil {
		util.Error(fmt.Sprintf("Failed to find file: %s", filepath), err)
		assert.Fail(t, "No file found")
	}

	tokensScope := tokenizer.Tokenize(text)

	// FOR DEBUGGING
	//jsonString := tokensScope.ToJsonString("testTagJava")
	//util.MakeDir("../../../output")
	//_ = util.CreateFileWithInfo("../../../output/java_output.json", jsonString)

	assert.Equal(t, 7, tokensScope.Size())
	for i := 0; i < tokensScope.Size(); i++ {
		st1, _ := tokensScope.At(i)
		switch i {
		case 0:
			tests.ValidateToken(t, st1, 1, 0, tz.RULENAME_OTHER, tz.SYMBOLIC_NAME_COMMENT, "/*\nThis class is called hello!\nit is used to print out \"Hello World\"\n */")
		case 1:
			tests.ValidateToken(t, st1, 6, 0, tz.RULENAME_KEYWORD, "PUBLIC", "public")
		case 2:
			tests.ValidateToken(t, st1, 6, 0, tz.RULENAME_KEYWORD, "CLASS", "class")
		case 3:
			tests.ValidateToken(t, st1, 6, 0, tz.RULENAME_KEYWORD, tz.SYMBOLIC_NAME_NON_KEYWORD, "hello")
		case 4:
			tests.ValidateToken(t, st1, 6, 0, tz.RULENAME_SYMBOL, "LCURLY", "{")
		case 5:
			{
				assert.True(t, st1.ValidScopeToken())
				sc1 := st1.GetScopeToken()
				assert.Equal(t, 13, sc1.Size())
				for j := 0; j < sc1.Size(); j++ {
					st2, _ := sc1.At(j)
					switch j {
					case 0:
						tests.ValidateToken(t, st2, 7, 1, tz.RULENAME_KEYWORD, "PUBLIC", "public")
					case 1:
						tests.ValidateToken(t, st2, 7, 1, tz.RULENAME_KEYWORD, "STATIC", "static")
					case 2:
						tests.ValidateToken(t, st2, 7, 1, tz.RULENAME_KEYWORD, "VOID", "void")
					case 3:
						tests.ValidateToken(t, st2, 7, 1, tz.RULENAME_KEYWORD, tz.SYMBOLIC_NAME_NON_KEYWORD, "main")
					case 4:
						tests.ValidateToken(t, st2, 7, 1, tz.RULENAME_SYMBOL, "LPAREN", "(")
					case 5:
						tests.ValidateToken(t, st2, 7, 1, tz.RULENAME_KEYWORD, tz.SYMBOLIC_NAME_NON_KEYWORD, "String")
					case 6:
						tests.ValidateToken(t, st2, 7, 1, tz.RULENAME_SYMBOL, "LBRACKET", "[")
					case 7:
						tests.ValidateToken(t, st2, 7, 1, tz.RULENAME_SYMBOL, "RBRACKET", "]")
					case 8:
						tests.ValidateToken(t, st2, 7, 1, tz.RULENAME_KEYWORD, tz.SYMBOLIC_NAME_NON_KEYWORD, "args")
					case 9:
						tests.ValidateToken(t, st2, 7, 1, tz.RULENAME_SYMBOL, "RPAREN", ")")
					case 10:
						tests.ValidateToken(t, st2, 7, 1, tz.RULENAME_SYMBOL, "LCURLY", "{")
					case 11:
						{
							assert.True(t, st2.ValidScopeToken())
							sc2 := st2.GetScopeToken()
							assert.Equal(t, 10, sc2.Size())
							for k := 0; k < sc2.Size(); k++ {
								st3, _ := sc2.At(k)
								switch k {
								case 0:
									tests.ValidateToken(t, st3, 8, 2, tz.RULENAME_OTHER, tz.SYMBOLIC_NAME_COMMENT, "// This prints out stuff\n")
								case 1:
									tests.ValidateToken(t, st3, 9, 2, tz.RULENAME_KEYWORD, tz.SYMBOLIC_NAME_NON_KEYWORD, "System")
								case 2:
									tests.ValidateToken(t, st3, 9, 2, tz.RULENAME_SYMBOL, "PERIOD", ".")
								case 3:
									tests.ValidateToken(t, st3, 9, 2, tz.RULENAME_KEYWORD, tz.SYMBOLIC_NAME_NON_KEYWORD, "out")
								case 4:
									tests.ValidateToken(t, st3, 9, 2, tz.RULENAME_SYMBOL, "PERIOD", ".")
								case 5:
									tests.ValidateToken(t, st3, 9, 2, tz.RULENAME_KEYWORD, tz.SYMBOLIC_NAME_NON_KEYWORD, "println")
								case 6:
									tests.ValidateToken(t, st3, 9, 2, tz.RULENAME_SYMBOL, "LPAREN", "(")
								case 7:
									tests.ValidateToken(t, st3, 9, 2, tz.RULENAME_OTHER, tz.SYMBOLIC_NAME_STRING, "\"Hello World\"")
								case 8:
									tests.ValidateToken(t, st3, 9, 2, tz.RULENAME_SYMBOL, "RPAREN", ")")
								case 9:
									tests.ValidateToken(t, st3, 9, 2, tz.RULENAME_SYMBOL, "SEMICOLON", ";")
								}
							}
						}
					case 12:
						tests.ValidateToken(t, st2, 10, 1, tz.RULENAME_SYMBOL, "RCURLY", "}")
					}
				}
			}
		case 6:
			tests.ValidateToken(t, st1, 11, 0, tz.RULENAME_SYMBOL, "RCURLY", "}")
		}
	}
}
