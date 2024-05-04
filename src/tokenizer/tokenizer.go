package tokenizer

import (
	"fmt"
	"strings"
	tk "tp/src/tokenizer/tokens"
)

type Tokenizer struct {
	Type                   string
	Symbols                [][]string
	Keywords               []string
	SpaceSizeString        string
	BracketIdentifierStart string
	BracketIdentifierEnd   string
	BracketCountRunner     int
	AllowTabs              bool
	GetBracketCount        bool
	StringAndCommentFinder func(text *string, pos int) (bool, string)
	StringAndCommentEnder  func(text *string, selectionText string, lookForThisString string, pos int) bool
}

func (tkzr *Tokenizer) Tokenize(text string) []tk.Token {
	LineTokens := tkzr.splitIntoNewLines(tkzr.removeCommentsAndStrings(text))
	FinalTokens := make([]tk.Token, 0)
	for i := 0; i < len(LineTokens); i++ {
		tks := tkzr.splitLineTokensToIndividualTokens(LineTokens[i])
		for i := 0; i < len(tks); i++ {
			FinalTokens = append(FinalTokens, tks[i])
		}
	}
	return FinalTokens
}

func (tkzr *Tokenizer) splitLineTokensToIndividualTokens(t tk.Token) []tk.Token {
	retTokens := make([]tk.Token, 0)
	textOfLine := t.Text
	/* so, identifiers/keywords and such can only consist of:
	 a-z    A-Z    0-9   and underscores (_)
	97-122  65-90  48-57                  95
		     If this is wrong, fix it HERE
	*/
	currentWord := ""

	for i := 0; i < len(textOfLine); i++ {
		isValidChar := false
		if i != len(textOfLine) {
			numChar := (int)(textOfLine[i])
			isValidChar = (numChar == 95) || (numChar >= 97 && numChar <= 122) || (numChar >= 65 && numChar <= 90) || (numChar >= 48 && numChar <= 57)
		}
		if isValidChar {
			currentWord += textOfLine[i : i+1]
		} else {
			if currentWord != "" {
				token := tkzr.identifyToken(currentWord, t.LineNumber, t.TabNumber, true)
				retTokens = append(retTokens, token)
				currentWord = ""
			}
			if i < len(textOfLine) {
				token := tkzr.identifyToken(textOfLine[i:i+1], t.LineNumber, t.TabNumber, false)
				if token.SymbolicName != "SPACE" {
					retTokens = append(retTokens, token)
				}
			}
		}
	}
	if currentWord != "" {
		token := tkzr.identifyToken(currentWord, t.LineNumber, t.TabNumber, true)
		retTokens = append(retTokens, token)
		currentWord = ""
	}

	return retTokens
}

func (tkzr *Tokenizer) identifyToken(str string, ln int, tab int, isKeyword bool) tk.Token {
	symbolicName := ""
	ruleName := ""

	if isKeyword {
		ruleName = "KEYWORD"
		for i := 0; i < len(tkzr.Keywords); i++ {
			if str == tkzr.Keywords[i] {
				symbolicName = strings.ToUpper(tkzr.Keywords[i])
			}
		}

		if symbolicName == "" {
			symbolicName = "IDENTIFIER"
		}
	} else {
		ruleName = "SYMBOL"
		if tkzr.GetBracketCount {
			if str == tkzr.BracketIdentifierStart {
				tkzr.BracketCountRunner++
			} else if str == tkzr.BracketIdentifierEnd {
				tkzr.BracketCountRunner--
			}
		}

		for i := 0; i < len(tkzr.Symbols); i++ {
			if str == tkzr.Symbols[i][0] {
				symbolicName = strings.ToUpper(tkzr.Symbols[i][1])
			}
		}

		if symbolicName == "" {
			symbolicName = "OTHER_SYMBOL"
		}
	}

	return tk.Token{
		LineNumber:    ln,
		TabNumber:     tab,
		BracketNumber: tkzr.BracketCountRunner,
		SymbolicName:  symbolicName,
		RuleName:      ruleName,
		Text:          str,
	}
}

func (tkzr *Tokenizer) splitIntoNewLines(s string) []tk.Token {
	lines := strings.Split(s, "\n")
	Tokens := make([]tk.Token, 0)
	for i := 0; i < len(lines); i++ {
		if lines[i] != "" {
			tabNumber, noTabText := tkzr.getNumberOfAndRemoveTabs(lines[i])
			Tokens = append(Tokens, tk.CreateUnidentifiedToken(noTabText, i+1, tabNumber, -1))
		}
	}
	return Tokens
}

func (tkzr *Tokenizer) getNumberOfAndRemoveTabs(text string) (int, string) {

	// So, typically, 4 spaces == 1 tab
	// space in ascii is 32
	// tab in ascii is 9
	// newline is 10
	tabs := 0
	for i := 0; i < len(text); i++ {
		if i+4 < len(text) {
			if text[0:len(tkzr.SpaceSizeString)] == tkzr.SpaceSizeString { // NUMBER OF SPACES FOR TAB COUNT
				text = text[4:]
				tabs++
			} else if (int)(text[i]) == 9 && tkzr.AllowTabs { // if it is the tab character
				tabs++
				text = text[1:]
			} else {
				break
			}
		} else if (int)(text[i]) == 9 && tkzr.AllowTabs { // if it is the tab character
			tabs++
			text = text[1:]
		} else {
			break
		}
	}

	return tabs, text
}

func (tkzr *Tokenizer) removeCommentsAndStrings(text string) string {
	lengthOfText := len(text)
	removeParts := make([][]int, 0)
	skip := false

	for i := 0; i < lengthOfText; { // Removes comments and strings
		foundItem, lookForThisString := tkzr.StringAndCommentFinder(&text, i)

		if foundItem {
			initialIndex := i
			endex := -1
			i++

			for endex == -1 && i+len(lookForThisString) < lengthOfText {
				selectionOfText := text[i : i+len(lookForThisString)]
				if tkzr.StringAndCommentEnder(&text, selectionOfText, lookForThisString, i) {
					endex = i + len(lookForThisString)
				}
				i++
			}

			skip = true
			if endex == -1 {
				endex = lengthOfText
				skip = false
			}

			item := []int{initialIndex, endex}
			removeParts = append(removeParts, item)
		}

		if skip {
			skip = false
		} else {
			i++
		}
	}

	for i := len(removeParts) - 1; i >= 0; i-- {
		removedPortion := text[removeParts[i][0]:removeParts[i][1]]
		numberOfNewlines := strings.Count(removedPortion, "\n")
		addNewLines := ""
		for i := 0; i < numberOfNewlines; i++ {
			addNewLines += "\n"
		}

		text = text[:removeParts[i][0]] + addNewLines + text[removeParts[i][1]:]

	}
	return text
}

func PrintOutTextOfTokens(tokens []tk.Token) {
	currentLine := 0
	for i := 0; i < len(tokens); i++ {
		t := tokens[i]
		if t.LineNumber != currentLine {
			fmt.Println()
			spacesStr := fmt.Sprintf("%d", t.LineNumber)
			spacesNum := 6 - len(spacesStr)
			for i := 0; i < spacesNum; i++ {
				spacesStr += " "
			}
			fmt.Print(spacesStr + "| ")
			printNumberOfTabs(t.TabNumber)
			currentLine = t.LineNumber
		}
		fmt.Print(t.Text + " ")
	}
}

func printNumberOfTabs(n int) {
	for i := 0; i < n; i++ {
		fmt.Print("    ")
	}
}

func PrintAllTokens(FinalTokens []tk.Token) {
	for i := 0; i < len(FinalTokens); i++ {
		FinalTokens[i].PrintToken()
	}
}
