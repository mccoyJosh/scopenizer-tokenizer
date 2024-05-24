package javaTokenizer

import (
	tz "tp/src/tokenizer"
	"unicode"
)

func GetJavaTokenizer() *tz.Tokenizer {
	tkzr := tz.GenerateDefaultTokenizerObject()
	tkzr.ConfigureGeneral("java", symbols, keywords,
		// Function To Determine if a letter can be part of a keyword in this language
		func(c rune) bool {
			return unicode.IsLetter(c) || unicode.IsDigit(c) || c == '_'
		},
	)
	tkzr.ConfigureComment(
		// Comment Start
		func(tkzr *tz.Tokenizer) bool {
			if tkzr.CurrentChar() == '/' {
				if tkzr.TextSize() > tkzr.Index()+1 {
					nextChar := tkzr.GetChar(tkzr.Index() + 1)
					if nextChar == '/' {
						tkzr.IncrementIndex()

						tkzr.StartInfo = "//"
						tkzr.EndInfo = "\n"
						return true
					}
					if nextChar == '*' {
						tkzr.IncrementIndex()

						tkzr.StartInfo = "/*"
						tkzr.EndInfo = "*/"
						return true
					}
				}
			}
			return false
		},
		// Comment End
		func(tkzr *tz.Tokenizer) bool {
			if len(tkzr.EndInfo) > 1 {
				currentIndex := tkzr.Index()
				if tkzr.TextSize() > currentIndex+1 {
					nextTwoCharacter, err := tkzr.TextRange(currentIndex, currentIndex+2)
					if err != nil {
						return false
					}
					if nextTwoCharacter == tkzr.EndInfo {
						return true
					}
				}
			}
			return tkzr.CurrentChar() == rune(tkzr.EndInfo[0])
		},
	)
	tkzr.ConfigureScope(
		// Scope Start
		func(tkzr *tz.Tokenizer) bool {
			if tkzr.CurrentChar() == '{' {
				tkzr.StartInfo = "{"
				return true
			}

			return false
		},
		// Scope End
		func(tkzr *tz.Tokenizer) bool {
			if tkzr.CurrentChar() == '}' {
				tkzr.EndInfo = "}"
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
				return true
			}
			if tkzr.CurrentChar() == '\'' {
				tkzr.StartInfo = "'"
				tkzr.EndInfo = "'"
				return true
			}
			return false
		},
		// String End
		func(tkzr *tz.Tokenizer) bool {
			return tkzr.GetChar(tkzr.Index()-1) != '\\' && tkzr.CurrentChar() == rune(tkzr.EndInfo[0])
		},
	)

	return &tkzr
}
