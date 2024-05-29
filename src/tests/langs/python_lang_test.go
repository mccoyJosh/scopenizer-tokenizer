package langs_test

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	pyTokenizer "tp/src/instances/langs/python/tokenizer"
	"tp/src/tests"
	tz "tp/src/tokenizer"
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

	for i := 0; i < tokensScope.Size(); i++ {
		st1, _ := tokensScope.At(i)
		switch i {
		case 0:
			tests.ValidateToken(t, st1, 1, 0, tz.RULENAME_OTHER, tz.SYMBOLIC_NAME_COMMENT, "# Here is an example python file\n")
		case 1:
			tests.ValidateToken(t, st1, 3, 0, tz.RULENAME_KEYWORD, "DEF", "def")
		case 2:
			tests.ValidateToken(t, st1, 3, 0, tz.RULENAME_KEYWORD, tz.SYMBOLIC_NAME_NON_KEYWORD, "print_test")
		case 3:
			tests.ValidateToken(t, st1, 3, 0, tz.RULENAME_SYMBOL, "LPAREN", "(")
		case 4:
			tests.ValidateToken(t, st1, 3, 0, tz.RULENAME_KEYWORD, tz.SYMBOLIC_NAME_NON_KEYWORD, "str")
		case 5:
			tests.ValidateToken(t, st1, 3, 0, tz.RULENAME_SYMBOL, "RPAREN", ")")
		case 6:
			tests.ValidateToken(t, st1, 3, 0, tz.RULENAME_SYMBOL, "COLON", ":")
		case 7:
			{
				//SCOPE
				assert.True(t, st1.ValidScopeToken())
				scope_a := st1.GetScopeToken()
				for j := 0; j < scope_a.Size(); j++ {
					st_a, _ := scope_a.At(j)
					switch j {
					case 0:
						tests.ValidateToken(t, st_a, 4, 1, tz.RULENAME_KEYWORD, "IF", "if")
					case 1:
						tests.ValidateToken(t, st_a, 4, 1, tz.RULENAME_KEYWORD, tz.SYMBOLIC_NAME_NON_KEYWORD, "str")
					case 2:
						tests.ValidateToken(t, st_a, 4, 1, tz.RULENAME_SYMBOL, "EQUAL", "=")
					case 3:
						tests.ValidateToken(t, st_a, 4, 1, tz.RULENAME_SYMBOL, "EQUAL", "=")
					case 4:
						tests.ValidateToken(t, st_a, 4, 1, tz.RULENAME_OTHER, tz.SYMBOLIC_NAME_STRING, "'not hello world'")
					case 5:
						tests.ValidateToken(t, st_a, 4, 1, tz.RULENAME_SYMBOL, "COLON", ":")
					case 6:
						{
							// SCOPE
							assert.True(t, st_a.ValidScopeToken())
							scope_b := st_a.GetScopeToken()
							for k := 0; k < scope_b.Size(); k++ {
								st_b, _ := scope_b.At(k)
								switch j {
								case 0:
									tests.ValidateToken(t, st_b, 5, 2, tz.RULENAME_KEYWORD, tz.SYMBOLIC_NAME_NON_KEYWORD, "print_test")
								case 1:
									tests.ValidateToken(t, st_b, 5, 2, tz.RULENAME_SYMBOL, "LPAREN", "(")
								case 2:
									tests.ValidateToken(t, st_a, 5, 2, tz.RULENAME_SYMBOL, tz.SYMBOLIC_NAME_NON_KEYWORD, "str")
								case 3:
									tests.ValidateToken(t, st_b, 5, 2, tz.RULENAME_SYMBOL, "RPAREN", ")")
								}
							}
						}
					}
				}
			}
		case 8:
			tests.ValidateToken(t, st1, 7, 0, tz.RULENAME_KEYWORD, "DEF", "def")
		case 9:
			tests.ValidateToken(t, st1, 7, 0, tz.RULENAME_KEYWORD, tz.SYMBOLIC_NAME_NON_KEYWORD, "main")
		case 10:
			tests.ValidateToken(t, st1, 7, 0, tz.RULENAME_SYMBOL, "LPAREN", "(")
		case 11:
			tests.ValidateToken(t, st1, 7, 0, tz.RULENAME_SYMBOL, "RPAREN", ")")
		case 12:
			tests.ValidateToken(t, st1, 7, 0, tz.RULENAME_SYMBOL, "COLON", ":")
		case 13:
			{
				//SCOPE
				assert.True(t, st1.ValidScopeToken())
				scope_a := st1.GetScopeToken()
				for j := 0; j < scope_a.Size(); j++ {
					st_a, _ := scope_a.At(j)
					switch j {
					case 0:
						tests.ValidateToken(t, st_a, 8, 1, tz.RULENAME_KEYWORD, tz.SYMBOLIC_NAME_NON_KEYWORD, "print_test")
					case 1:
						tests.ValidateToken(t, st_a, 8, 1, tz.RULENAME_SYMBOL, "LPAREN", "(")
					case 2:
						tests.ValidateToken(t, st_a, 8, 1, tz.RULENAME_OTHER, tz.SYMBOLIC_NAME_STRING, "\"hello world\"")
					case 3:
						tests.ValidateToken(t, st_a, 8, 1, tz.RULENAME_SYMBOL, "RPAREN", ")")

					}

				}
			}
		case 14:
			tests.ValidateToken(t, st1, 10, 0, tz.RULENAME_KEYWORD, tz.SYMBOLIC_NAME_NON_KEYWORD, "main")
		case 15:
			tests.ValidateToken(t, st1, 10, 0, tz.RULENAME_SYMBOL, "LPAREN", "(")
		case 16:
			tests.ValidateToken(t, st1, 10, 0, tz.RULENAME_SYMBOL, "RPAREN", ")")
		}

	}
}
