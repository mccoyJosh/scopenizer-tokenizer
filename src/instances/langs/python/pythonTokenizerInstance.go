package pythonTokenizer

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	tz "tp/src/tokenizer"
	"tp/src/util"
	"unicode"
)

func GetPythonTokenizer() *tz.Tokenizer {
	tkzr := tz.GenerateDefaultTokenizerObject()

	setCommentLineNumber := func(tkzr *tz.Tokenizer) {
		stopIndex := strings.Index(tkzr.FunctionSharedInfo, "|")
		if stopIndex == -1 { // Just add |
			tkzr.FunctionSharedInfo = "|" + tkzr.FunctionSharedInfo
		} else { // Remove previous line number
			tkzr.FunctionSharedInfo = tkzr.FunctionSharedInfo[stopIndex:]
		}
		tkzr.FunctionSharedInfo = strconv.Itoa(tkzr.GetCurrentLineNumber()) + tkzr.FunctionSharedInfo
	}

	getCommentLineNumber := func(tkzr *tz.Tokenizer) (int, error) {
		stopIndex := strings.Index(tkzr.FunctionSharedInfo, "|")
		if stopIndex == -1 { // No line number was ever assigned
			return -1, errors.New("attempted to get line number when none was assigned")
		}
		lineNumber, err := strconv.Atoi(tkzr.FunctionSharedInfo[0:stopIndex])
		if err != nil {
			return -1, errors.New("failed to convert line number string to int. attempted string: " + tkzr.FunctionSharedInfo[0:stopIndex])
		}
		return lineNumber, nil
	}

	pushScopeInfo := func(tkzr *tz.Tokenizer) {
		tkzr.FunctionSharedInfo += fmt.Sprintf("&%d,%d", tkzr.GetCurrentLineNumber(), tkzr.GetCurrentTabLevel())
	}

	const NoScopeInfoFoundErrorString = "no scope info pair found"

	frontScopeInfo := func(tkzr *tz.Tokenizer) (int, int, error) {
		stopIndex := strings.LastIndex(tkzr.FunctionSharedInfo, "&")
		if stopIndex == -1 {
			return -1, -1, errors.New(NoScopeInfoFoundErrorString)
		}
		pairString := tkzr.FunctionSharedInfo[stopIndex+1:]
		pairs := strings.Split(pairString, ",")
		if len(pairs) != 2 {
			return -1, -1, errors.New("scope info pair malformed (no comma): " + pairString)
		}

		lineNum, err := strconv.Atoi(pairs[0])
		if err != nil {
			return -1, -1, errors.New("scope info pair malformed (could not convert line number): " + pairString)
		}

		tabLevel, err := strconv.Atoi(pairs[1])
		if err != nil {
			return -1, -1, errors.New("scope info pair malformed (could not convert tab level): " + pairString)
		}

		return lineNum, tabLevel, nil
	}

	popScopeInfo := func(tkzr *tz.Tokenizer) error {
		stopIndex := strings.LastIndex(tkzr.FunctionSharedInfo, "&")
		if stopIndex == -1 {
			return errors.New(NoScopeInfoFoundErrorString)
		}
		tkzr.FunctionSharedInfo = tkzr.FunctionSharedInfo[:stopIndex]
		return nil
	}

	tkzr.ConfigureGeneral("python", symbols, keywords,
		// Function To Determine if a letter can be part of a keyword in this language
		func(c rune) bool {
			return unicode.IsLetter(c) || unicode.IsDigit(c) || c == '_'
		},
	)
	tkzr.ConfigureComment(
		// Comment Start
		func(tkzr *tz.Tokenizer) bool {
			if tkzr.CurrentChar() == '#' {
				tkzr.StartInfo = "#"
				tkzr.EndInfo = "\n"
				setCommentLineNumber(tkzr)
				return true
			}
			return false
		},
		// Comment End
		func(tkzr *tz.Tokenizer) bool {
			num, err := getCommentLineNumber(tkzr)
			if err != nil {
				util.Error("FunctionSharedInfo in CommentEndFunction error (python)", err)
				panic(err)
			}
			return tkzr.GetCurrentLineNumber() != num
		},
	)
	tkzr.ConfigureScope(
		// Scope Start
		func(tkzr *tz.Tokenizer) bool {
			if tkzr.CurrentChar() == ':' {
				tkzr.StartInfo = ":"
				pushScopeInfo(tkzr)
				return true
			}

			return false
		},
		// Scope End
		func(tkzr *tz.Tokenizer) bool {
			lineNum, tabNum, err := frontScopeInfo(tkzr)
			if err != nil {
				if err.Error() != NoScopeInfoFoundErrorString {
					util.Error("Python frontScopeInfo Error (see python tokenizer instance)", err)
					panic(err)
				}
				return false
			}

			if tkzr.GetCurrentLineNumber() != lineNum && tkzr.GetCurrentTabLevel() <= tabNum {
				tkzr.EndInfo = ""
				err = popScopeInfo(tkzr)
				if err != nil {
					util.Error("Failed to pop value that should exist"+fmt.Sprintf(" line number: %d     tab number: %d   ", lineNum, tabNum), err)
					panic(err)
				}
				tkzr.SkipIncrement() // THIS SHOULD MAKE SURE THAT IF MULTIPLE SCOPES NEED TO BE STOPPED AT THE SAME TOKEN, THAT IT STOPS MOVING TO DO ALL THE STOPS
				return true
			}

			return false
		},
	)
	tkzr.ConfigureString(
		// String Start
		func(tkzr *tz.Tokenizer) bool {
			if tkzr.CurrentChar() == '"' {
				tkzr.StartInfo = "\""
				tkzr.EndInfo = "\""

				substring, err := tkzr.TextRange(tkzr.Index(), tkzr.Index()+3)
				if err == nil && substring == "\"\"\"" {
					tkzr.StartInfo = "\"\"\""
					tkzr.EndInfo = "\"\"\""
				}
				return true
			}

			if tkzr.CurrentChar() == '\'' {
				tkzr.StartInfo = "'"
				tkzr.EndInfo = "'"

				substring, err := tkzr.TextRange(tkzr.Index(), tkzr.Index()+3)
				if err == nil && substring == "'''" {
					tkzr.StartInfo = "'''"
					tkzr.EndInfo = "'''"
				}
				return true
			}
			return false
		},
		// String End
		func(tkzr *tz.Tokenizer) bool {
			if len(tkzr.EndInfo) == 3 {
				substring, err := tkzr.TextRange(tkzr.Index(), tkzr.Index()+3)
				if err == nil && tkzr.GetChar(tkzr.Index()-1) != '\\' && substring == tkzr.EndInfo {
					return true
				}
				if err != nil { // TODO: this would only happen if a string never ended
					return true
				}
			}
			return tkzr.GetChar(tkzr.Index()-1) != '\\' && tkzr.CurrentChar() == rune(tkzr.EndInfo[0])
		},
	)

	return &tkzr
}
