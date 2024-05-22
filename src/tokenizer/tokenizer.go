package tokenizer

import (
	"errors"
	"fmt"
	"strings"
	tk "tp/src/tokenizer/tokens"
	"tp/src/util"
)

type Tokenizer struct {
	LanguageType string
	Symbols      [][]string
	Keywords     []string

	// Temp Info
	tempIgnoreChangesFromIncrement bool
	Text                           *string
	currentTabLevel                int
	currentLineNumber              int
	PotentialKeyword               string
	StartInfo                      string
	EndInfo                        string
	currentIndex                   int
	CurrentScope                   *tk.ScopeObj

	// Scope Info
	ScopeStartFunction func(tkzr *Tokenizer) bool
	ScopeEndFunction   func(tkzr *Tokenizer) bool
	// Note: These scope functions are intended to find MOST scopes... not all scopes

	// String Info
	StringStartFunction func(tkzr *Tokenizer) bool
	StringEndFunction   func(tkzr *Tokenizer) bool
	IncludeStrings      bool

	// Comment Info
	CommentStartFunction func(tkzr *Tokenizer) bool
	CommentEndFunction   func(tkzr *Tokenizer) bool
	IncludeComments      bool

	// Keyword Info
	IsKeywordCharacter func(c rune) bool

	// Whitespace Info
	spaceSizeString       string
	NumOfSpacesEquallyTab int
	IgnoreWhitespace      bool
	IgnoreNewLines        bool
}

func GenerateBasicTokenizerObject() Tokenizer {
	return Tokenizer{
		LanguageType: "",
		Symbols:      nil,
		Keywords:     nil,

		tempIgnoreChangesFromIncrement: false,
		Text:                           nil,
		currentTabLevel:                0,
		currentLineNumber:              0,
		PotentialKeyword:               "",
		StartInfo:                      "",
		EndInfo:                        "",
		currentIndex:                   0,

		CurrentScope:       nil,
		ScopeStartFunction: nil,
		ScopeEndFunction:   nil,

		StringStartFunction: nil,
		StringEndFunction:   nil,
		IncludeStrings:      false,

		CommentStartFunction: nil,
		CommentEndFunction:   nil,
		IncludeComments:      false,

		IsKeywordCharacter: nil,

		spaceSizeString:       "",
		NumOfSpacesEquallyTab: 0,
		IgnoreWhitespace:      true,
		IgnoreNewLines:        true,
	}
}

func (tkzr *Tokenizer) initTempVariables(text *string) {
	tkzr.tempIgnoreChangesFromIncrement = false
	tkzr.initSpaceSizeString()
	tkzr.PotentialKeyword = ""
	tkzr.currentTabLevel = 0
	tkzr.currentIndex = 0
	tkzr.currentLineNumber = 0
	tkzr.Text = text
	tkzr.StartInfo = ""
	tkzr.EndInfo = ""
}

func (tkzr *Tokenizer) Tokenize(text string) tk.ScopeObj {
	tkzr.initTempVariables(&text)

	finalScope := tk.InitScope()
	finalScope.SetType("File")
	tkzr.CurrentScope = &finalScope

	for tkzr.IndexInBound() {
		foundString := tkzr.StringStartFunction(tkzr)
		foundComment := tkzr.CommentStartFunction(tkzr)
		foundStartScope := tkzr.ScopeStartFunction(tkzr)
		foundEndScope := tkzr.ScopeEndFunction(tkzr)
		foundWhiteSpace := whiteSpaceStartFunction(tkzr)

		if foundString || foundComment || foundStartScope {
			tkzr.tempIgnoreChangesFromIncrement = true
			tkzr.IncrementIndex() // TODO: determine if this is necessary or should be left up yo anonymous functions

			if tkzr.PotentialKeyword != "" {
				tkzr.CurrentScope.Push(tkzr.createKeywordToken(tkzr.PotentialKeyword))
				tkzr.PotentialKeyword = ""
			}

			if foundWhiteSpace {
				resultingToken := tkzr.applyFunctionUntilFailureTokenCreation(whiteSpaceEndFunction, "WHITESPACE")
				tkzr.CurrentScope.Push(resultingToken)
			} else if foundString {
				resultingToken := tkzr.applyFunctionUntilFailureTokenCreation(tkzr.StringEndFunction, "STRING")
				tkzr.CurrentScope.Push(resultingToken)
			} else if foundComment {
				resultingToken := tkzr.applyFunctionUntilFailureTokenCreation(tkzr.CommentEndFunction, "COMMENT")
				tkzr.CurrentScope.Push(resultingToken)
			} else if foundStartScope {
				newScopeTkn := tk.InitScopeToken()
				tkzr.CurrentScope.Push(newScopeTkn)
				tkzr.CurrentScope = newScopeTkn.GetScopeToken()
			} else if foundEndScope {
				parentScope := tkzr.CurrentScope.GetScopeParent()
				if parentScope == nil {
					err := errors.New(fmt.Sprintf("Either malformed data attempted to be Tokenized or anonymous functions provided to tokenizers incorrectly defined when scopes being/end"))
					util.Error(err.Error(), err)
				} else {
					tkzr.CurrentScope = parentScope
				}
			}
			tkzr.tempIgnoreChangesFromIncrement = false
		} else { // Not a scope identifier, not a comment, not a string
			char := text[tkzr.Index()]
			if tkzr.IsKeywordCharacter(rune(char)) {
				tkzr.PotentialKeyword += string(char)
			} else { // Found a symbol, which needs to be added and
				if tkzr.PotentialKeyword != "" {
					tkzr.CurrentScope.Push(tkzr.createKeywordToken(tkzr.PotentialKeyword))
					tkzr.PotentialKeyword = ""
				}
				tkzr.CurrentScope.Push(tkzr.createSymbolToken(string(char)))
			}
		}
		tkzr.IncrementIndex()
	}

	if tkzr.PotentialKeyword != "" { // TODO: SEE IF ONE LAST CHECK IS NECESSARY
		tkzr.CurrentScope.Push(tkzr.createKeywordToken(tkzr.PotentialKeyword))
		tkzr.PotentialKeyword = ""
	}

	return finalScope
}

func (tkzr *Tokenizer) Index() int {
	return tkzr.currentIndex
}

func (tkzr *Tokenizer) IndexInBound() bool {
	return tkzr.DetermineIfIndexInBound(tkzr.Index())
}

func (tkzr *Tokenizer) DetermineIfIndexInBound(index int) bool {
	return index < len(*tkzr.Text)
}

func (tkzr *Tokenizer) applyFunctionUntilFailureTokenCreation(BooleanEndFunction func(tkzr *Tokenizer) bool, symbolicName string) *tk.Token {
	lineNumber := tkzr.currentLineNumber
	tabLevel := tkzr.currentTabLevel
	tokenText := ""
	for BooleanEndFunction(tkzr) {
		tokenText += string(tkzr.CurrentChar())
		tkzr.IncrementIndex()
	}
	tokenText = tkzr.StartInfo + tokenText + tkzr.EndInfo
	finalToken := tk.CreateUnidentifiedToken(tokenText, lineNumber, tabLevel)
	finalToken.SetValues("OTHER", symbolicName)
	return &finalToken
}

func whiteSpaceStartFunction(tkzr *Tokenizer) bool {
	return util.IsWhitespaceCharacter(tkzr.CurrentChar())
}

func whiteSpaceEndFunction(tkzr *Tokenizer) bool {
	return util.IsWhitespaceCharacter(tkzr.CurrentChar())
}

func (tkzr *Tokenizer) IncrementIndex() {
	tkzr.currentIndex++
	if tkzr.CurrentChar() == '\n' { // Encounters a newline
		if !tkzr.IgnoreNewLines && !tkzr.tempIgnoreChangesFromIncrement { // (potentially) adding the new line character token
			newToken := tk.CreateUnidentifiedToken("\n", tkzr.currentLineNumber, tkzr.currentTabLevel)
			newToken.SetValues("OTHER", "NEWLINE")
			tkzr.CurrentScope.Push(&newToken)
		}
		tkzr.currentLineNumber++

		// Gets the tab level for this line
		gatheredWhitespace := tkzr.GatherWhitespace(!tkzr.tempIgnoreChangesFromIncrement)
		numOfTabs := util.DetermineNumberOfTabs(gatheredWhitespace, tkzr.NumOfSpacesEquallyTab, true)
		tkzr.currentTabLevel = numOfTabs
		if !tkzr.IgnoreWhitespace && !tkzr.tempIgnoreChangesFromIncrement {
			newToken := tk.CreateUnidentifiedToken(gatheredWhitespace, tkzr.currentLineNumber, tkzr.currentTabLevel)
			newToken.SetValues("OTHER", "WHITESPACE")
			tkzr.CurrentScope.Push(&newToken)
		}
	}

}

func (tkzr *Tokenizer) GatherWhitespace(updateCurrentIndex bool) string {
	gatheredWhitespace := ""
	index := tkzr.currentIndex
	var char rune
	for tkzr.DetermineIfIndexInBound(index) {
		char = tkzr.GetChar(index)
		if util.IsWhitespaceCharacter(char) {
			gatheredWhitespace += string(char)
		} else {
			break
		}
		index++
	}
	if updateCurrentIndex {
		tkzr.currentIndex = index
	}
	return gatheredWhitespace
}

func (tkzr *Tokenizer) CurrentChar() rune {
	return tkzr.GetChar(tkzr.currentIndex)
}

func (tkzr *Tokenizer) GetChar(index int) rune {
	return rune((*tkzr.Text)[index])
}

func (tkzr *Tokenizer) initSpaceSizeString() {
	tkzr.spaceSizeString = ""
	for i := 0; i < tkzr.NumOfSpacesEquallyTab; i++ {
		tkzr.spaceSizeString += " "
	}
}

func (tkzr *Tokenizer) createKeywordToken(keywordString string) *tk.Token {
	newToken := tk.CreateUnidentifiedToken(keywordString, tkzr.currentLineNumber, tkzr.currentTabLevel)
	newToken.SetValues("KEYWORD", tkzr.identifyKeyword(keywordString))
	return &newToken
}

func (tkzr *Tokenizer) createSymbolToken(symbolString string) *tk.Token {
	newToken := tk.CreateUnidentifiedToken(symbolString, tkzr.currentLineNumber, tkzr.currentTabLevel)
	newToken.SetValues("SYMBOL", tkzr.identifySymbol(symbolString))
	return &newToken
}

func (tkzr *Tokenizer) identifyKeyword(keywordString string) string {
	symbolicName := ""
	for i := 0; i < len(tkzr.Keywords); i++ {
		if keywordString == tkzr.Keywords[i] {
			symbolicName = strings.ToUpper(tkzr.Keywords[i])
		}
	}
	if symbolicName == "" {
		symbolicName = "IDENTIFIER"
	}
	return symbolicName
}

func (tkzr *Tokenizer) identifySymbol(symbol string) string {
	symbolicName := ""
	for i := 0; i < len(tkzr.Symbols); i++ {
		if symbol == tkzr.Symbols[i][0] {
			symbolicName = strings.ToUpper(tkzr.Symbols[i][1])
		}
	}
	if symbolicName == "" {
		symbolicName = "UNKNOWN"
	}
	return symbolicName
}
