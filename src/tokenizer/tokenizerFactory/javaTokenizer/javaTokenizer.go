package javaTokenizer

import "TokenAndParse/src/tokenizer"

func Tokenizer() *tokenizer.Tokenizer {
	tokenizerObj := tokenizer.Tokenizer{
		Type:                   "java",
		Symbols:                Symbols,
		Keywords:               Keywords,
		SpaceSizeString:        "  ",
		GetBracketCount:        true,
		BracketCountRunner:     0,
		BracketIdentifierStart: "{",
		BracketIdentifierEnd:   "}",
		AllowTabs:              true,
		StringAndCommentFinder: func(text *string, pos int) (bool, string) {
			lengthOfText := len(*text)
			charStr := (*text)[pos : pos+1]
			lookForThisString := ""

			switch charStr {
			case "/":
				{
					if pos+2 < lengthOfText {
						if (*text)[pos:pos+2] == "//" {
							lookForThisString = "\n"
						} else if (*text)[pos:pos+2] == "/*" {
							lookForThisString = "*/"
						}
					}
				}
			case "\"":
				{
					lookForThisString = "\""
				}
			case "'":
				{
					lookForThisString = "'"
				}
			}
			return lookForThisString != "", lookForThisString
		},
		StringAndCommentEnder: func(text *string, selectionOfText string, lookForThisString string, pos int) bool {
			return selectionOfText == lookForThisString && !((lookForThisString == "\"" || lookForThisString == "'") && pos > 0 && (*text)[pos-1:pos] == "\\")
		},
	}
	return &tokenizerObj
}
