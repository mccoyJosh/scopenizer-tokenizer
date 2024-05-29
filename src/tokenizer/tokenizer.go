package tokenizer

import (
	"errors"
	"fmt"
	"strings"
	tk "tp/src/tokenizer/tokens"
	"tp/src/util"
	"unicode"
)

func CreateDullTokenizer() *Tokenizer {
	tkzr := GenerateDefaultTokenizerObject()
	tkzr.ConfigureGeneral("dull", make([][]string, 0), make([]string, 0),
		// Function To Determine if a letter can be part of a keyword in this language
		func(c rune) bool {
			return unicode.IsLetter(c) || unicode.IsDigit(c) || c == '_'
		},
	)
	tkzr.ConfigureComment(
		// Comment Start
		func(tkzr *Tokenizer) bool {
			tkzr.StartInfo = "("
			tkzr.EndInfo = ")"
			return false
		},
		// Comment End
		func(tkzr *Tokenizer) bool {
			return true
		},
	)
	tkzr.ConfigureScope(
		// Scope Start
		func(tkzr *Tokenizer) bool {
			tkzr.StartInfo = "("
			tkzr.EndInfo = ")"
			return false
		},
		// Scope End
		func(tkzr *Tokenizer) bool {
			return false
		},
	)
	tkzr.ConfigureString(
		// String Start
		func(tkzr *Tokenizer) bool {
			tkzr.StartInfo = "("
			tkzr.EndInfo = ")"
			return false
		},
		// String End
		func(tkzr *Tokenizer) bool {
			return true
		},
	)
	return &tkzr
}

type Tokenizer struct {
	LanguageType string
	Symbols      [][]string
	Keywords     []string

	// Temp Info
	tempIgnoreChangesFromIncrement bool
	Text                           *string
	currentTabLevel                int
	currentLineNumber              int
	potentialKeyword               string
	StartInfo                      string
	EndInfo                        string
	FunctionSharedInfo             string
	currentIndex                   int
	currentScope                   *tk.ScopeObj
	skipIncrement                  bool

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

func GenerateDefaultTokenizerObject() Tokenizer {
	return Tokenizer{
		LanguageType: "",
		Symbols:      nil,
		Keywords:     nil,

		tempIgnoreChangesFromIncrement: false,
		Text:                           nil,
		currentTabLevel:                0,
		currentLineNumber:              0,
		potentialKeyword:               "",
		StartInfo:                      "",
		EndInfo:                        "",
		FunctionSharedInfo:             "",
		currentIndex:                   0,
		skipIncrement:                  false,

		currentScope:       nil,
		ScopeStartFunction: nil,
		ScopeEndFunction:   nil,

		StringStartFunction: nil,
		StringEndFunction:   nil,
		IncludeStrings:      true,

		CommentStartFunction: nil,
		CommentEndFunction:   nil,
		IncludeComments:      true,

		IsKeywordCharacter: nil,

		spaceSizeString:       "",
		NumOfSpacesEquallyTab: 4,
		IgnoreWhitespace:      true,
		IgnoreNewLines:        true,
	}
}

func (tkzr *Tokenizer) initTempVariables(text *string) {
	tkzr.tempIgnoreChangesFromIncrement = false
	tkzr.initSpaceSizeString()
	tkzr.potentialKeyword = ""
	tkzr.currentTabLevel = 0
	tkzr.currentIndex = 0
	tkzr.currentLineNumber = 1
	tkzr.Text = text
	tkzr.StartInfo = ""
	tkzr.EndInfo = ""
	tkzr.FunctionSharedInfo = ""
}

func (tkzr *Tokenizer) Tokenize(text string) tk.ScopeObj {
	err := tkzr.IsConfigured()
	if err != nil {
		util.Error("TOKENIZER NOT CONFIGURED CORRECTLY", err)
		panic(err)
	}

	tkzr.initTempVariables(&text)

	finalScope := tk.InitScope()
	finalScope.SetType("File")
	tkzr.currentScope = &finalScope

	for tkzr.IndexInBound() {
		tkzr.skipIncrement = false
		foundString := tkzr.StringStartFunction(tkzr)
		foundComment := tkzr.CommentStartFunction(tkzr)
		foundStartScope := tkzr.ScopeStartFunction(tkzr)
		foundEndScope := tkzr.ScopeEndFunction(tkzr)
		if foundString || foundComment || foundStartScope || foundEndScope {
			tkzr.tempIgnoreChangesFromIncrement = true

			if tkzr.potentialKeyword != "" {
				tkzr.currentScope.Push(tkzr.createKeywordToken(tkzr.potentialKeyword))
				tkzr.potentialKeyword = ""
			}
			if foundString {
				resultingToken := tkzr.applyFunctionUntilFailureTokenCreation(tkzr.StringEndFunction, SYMBOLIC_NAME_STRING)
				if tkzr.IncludeStrings {
					tkzr.currentScope.Push(resultingToken)
				}
			} else if foundComment {
				resultingToken := tkzr.applyFunctionUntilFailureTokenCreation(tkzr.CommentEndFunction, SYMBOLIC_NAME_COMMENT)
				if tkzr.IncludeComments {
					tkzr.currentScope.Push(resultingToken)
				}
			} else if foundStartScope {
				preScopeToken := tkzr.createTokenType(tkzr.StartInfo)
				tkzr.currentScope.Push(preScopeToken)

				newScopeTkn := tk.InitScopeToken()
				tkzr.currentScope.Push(newScopeTkn)
				tkzr.currentScope = newScopeTkn.GetScopeToken()
			} else if foundEndScope {
				parentScope := tkzr.currentScope.GetScopeParent()
				if parentScope == nil {
					err := errors.New(fmt.Sprintf("Either malformed data attempted to be Tokenized or anonymous functions provided to tokenizers incorrectly defined when scopes being/end"))
					util.Error(err.Error(), err)
				} else {
					tkzr.currentScope = parentScope
				}
				if tkzr.EndInfo != "" {
					postScopeToken := tkzr.createTokenType(tkzr.EndInfo)
					tkzr.currentScope.Push(postScopeToken)
				}
			}
			for i := 0; i < len(tkzr.EndInfo)-1; i++ {
				if !tkzr.skipIncrement {
					tkzr.IncrementIndex() // TODO: Determine whether this is a good long term solution or should be left for anonymous functions to deal with
				}
			}
			tkzr.StartInfo = ""
			tkzr.EndInfo = ""
			tkzr.tempIgnoreChangesFromIncrement = false
		} else { // Not a scope identifier, not a comment, not a string
			char := tkzr.CurrentChar()
			if tkzr.IsKeywordCharacter(rune(char)) {
				tkzr.potentialKeyword += string(char)
			} else { // Found a symbol, which needs to be added and
				tkzr.AddPotentialKeyword()
				newSymbolToken := tkzr.createSymbolToken(string(char))

				if newSymbolToken.SymbolicName == SYMBOLIC_NAME_WHITESPACE {
					if !tkzr.IgnoreWhitespace {
						tkzr.currentScope.Push(newSymbolToken)
					}
				} else if newSymbolToken.SymbolicName == SYMBOLIC_NAME_NEWLINE {
					tkzr.dealWithNewline()
				} else {
					tkzr.currentScope.Push(newSymbolToken)
				}
			}
		}
		if !tkzr.skipIncrement {
			tkzr.currentIndex++
		}
	}

	if tkzr.potentialKeyword != "" { // TODO: SEE IF ONE LAST CHECK IS NECESSARY
		tkzr.currentScope.Push(tkzr.createKeywordToken(tkzr.potentialKeyword))
		tkzr.potentialKeyword = ""
	}

	return finalScope
}

func (tkzr *Tokenizer) AddPotentialKeyword() {
	if tkzr.potentialKeyword != "" {
		tkzr.currentScope.Push(tkzr.createKeywordToken(tkzr.potentialKeyword))
		tkzr.potentialKeyword = ""
	}
}

func (tkzr *Tokenizer) PrintCharIndices() {
	for i := 0; i < tkzr.TextSize(); i++ {
		char := tkzr.GetChar(i)
		fmt.Printf("%d\t:\t%d\n", i, int(char))
	}
}

func (tkzr *Tokenizer) SkipIncrement() {
	tkzr.skipIncrement = true
}

func (tkzr *Tokenizer) TextSize() int {
	// TODO: This will error if text is nil
	return len(*tkzr.Text)
}

func (tkzr *Tokenizer) TextRange(begin int, end int) (string, error) {
	if begin >= 0 && begin < tkzr.TextSize() && end > 0 && end < tkzr.TextSize() && begin != end {
		return (*tkzr.Text)[begin:end], nil
	}
	return "", errors.New("TextRange bounds were out of bounds or otherwise invalid")
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

func (tkzr *Tokenizer) GetCurrentTabLevel() int {
	return tkzr.currentTabLevel
}

func (tkzr *Tokenizer) applyFunctionUntilFailureTokenCreation(BooleanEndFunction func(tkzr *Tokenizer) bool, symbolicName string) *tk.Token {
	lineNumber := tkzr.currentLineNumber
	tempLineNumber := lineNumber
	tabLevel := tkzr.currentTabLevel
	tokenText := ""
	tkzr.IncrementIndex() // TODO: This should skip the char which initialed this function to be applied
	for !BooleanEndFunction(tkzr) && tkzr.IndexInBound() {
		if tkzr.currentLineNumber != tempLineNumber {
			tokenText += "\n"
			tempLineNumber = tkzr.currentLineNumber
		}
		tokenText += string(tkzr.CurrentChar())
		tkzr.IncrementIndex()
	}
	tokenText = tkzr.StartInfo + tokenText + tkzr.EndInfo

	if tkzr.IndexInBound() && len(tkzr.EndInfo) > 0 && tkzr.CurrentChar() != rune(tkzr.EndInfo[0]) {
		tkzr.SkipIncrement()
	}

	finalToken := tk.CreateUnidentifiedToken(tokenText, lineNumber, tabLevel)
	finalToken.SetValues(RULENAME_OTHER, symbolicName)

	return &finalToken
}

func (tkzr *Tokenizer) GetCurrentLineNumber() int {
	return tkzr.currentLineNumber
}

func (tkzr *Tokenizer) IncrementIndex() {
	tkzr.currentIndex++
	tkzr.dealWithNewline()
	for tkzr.IndexInBound() && tkzr.CurrentChar() == '\n' {
		tkzr.currentIndex++
		tkzr.dealWithNewline()
	}
}

func (tkzr *Tokenizer) dealWithNewline() {
	// Ensure we are on a newline
	if tkzr.IndexInBound() && tkzr.CurrentChar() == '\n' {
		if !tkzr.tempIgnoreChangesFromIncrement {
			tkzr.AddPotentialKeyword()
		}

		// Adds newline token, if applicable
		if !tkzr.IgnoreNewLines && !tkzr.tempIgnoreChangesFromIncrement {
			newToken := tk.CreateUnidentifiedToken("\n", tkzr.currentLineNumber, tkzr.currentTabLevel)
			newToken.SetValues(RULENAME_OTHER, SYMBOLIC_NAME_NEWLINE)
			tkzr.currentScope.Push(&newToken)
		}

		// Increments line number to keep track of line
		tkzr.currentLineNumber++

		// Gets the tab level for this line
		gatheredWhitespace := tkzr.gatherWhitespace(!tkzr.tempIgnoreChangesFromIncrement)
		numOfTabs := util.DetermineNumberOfTabs(gatheredWhitespace, tkzr.NumOfSpacesEquallyTab, true)
		tkzr.currentTabLevel = numOfTabs

		// Adds whitespace token, if applicable
		if !tkzr.IgnoreWhitespace && !tkzr.tempIgnoreChangesFromIncrement {
			newToken := tk.CreateUnidentifiedToken(gatheredWhitespace, tkzr.currentLineNumber, tkzr.currentTabLevel)
			newToken.SetValues(RULENAME_OTHER, SYMBOLIC_NAME_WHITESPACE)
			tkzr.currentScope.Push(&newToken)
		}
	}
}

func (tkzr *Tokenizer) PrintCurrentInfo(args ...string) {
	optionalInfo := ""
	for _, arg := range args {
		optionalInfo += arg
	}
	fmt.Println("------------------------------------------------------")
	fmt.Printf("CURRENT INFO (%s): ", optionalInfo)
	util.PrintTime()
	fmt.Println()
	fmt.Println("\tGENERAL: ")
	fmt.Printf("\t\tIndex: %d\tAsciiOfIndex: %d\tCharOfIndex: '%s'\n", tkzr.currentIndex, int(tkzr.CurrentChar()), string(tkzr.CurrentChar()))
	fmt.Printf("\t\tLine Num: %d\t Tab Level: %d\n", tkzr.currentLineNumber, tkzr.currentTabLevel)
	fmt.Println("\tBOOLS:")
	fmt.Printf("\t\ttempIgnoreChangesFromIncrement %t\tskipIncrement: %t\n", tkzr.tempIgnoreChangesFromIncrement, tkzr.skipIncrement)
	fmt.Println("\tFUNCTION INFO: ")
	fmt.Printf("\t\tStartInfo: %s\tEndInfo:%s\tFunctionSharedInfo: %s\n", tkzr.StartInfo, tkzr.EndInfo, tkzr.FunctionSharedInfo)
	fmt.Println("\tCURRENT SCOPE:")
	fmt.Printf("\t\tSize: %d\tType: %s\n", tkzr.currentScope.Size(), tkzr.currentScope.GetType())
	lastFewTokens := ""
	numberOfPreviousTokens := 5
	for i := tkzr.currentScope.Size() - 1; i > 0 && i > tkzr.currentScope.Size()-1-numberOfPreviousTokens; i-- {
		tkn, err := tkzr.currentScope.At(i)
		if err != nil {
			continue
		}
		lastFewTokens += tkn.ToString() + "\t"
	}
	fmt.Printf("\t\tLast few tokens: %s\n", lastFewTokens)
	fmt.Println("------------------------------------------------------")
}

func (tkzr *Tokenizer) gatherWhitespace(updateCurrentIndex bool) string {
	gatheredWhitespace := ""
	index := tkzr.currentIndex + 1
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
	if updateCurrentIndex && gatheredWhitespace != "" {
		tkzr.currentIndex = index - 1 // TODO: maybe we don't need to do this!
	}
	return gatheredWhitespace
}

func (tkzr *Tokenizer) CurrentChar() rune {
	// TODO: Doesn't check if it is out of bounds
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
	newToken.SetValues(RULENAME_KEYWORD, tkzr.identifyKeyword(keywordString))
	return &newToken
}

func (tkzr *Tokenizer) createSymbolToken(symbolString string) *tk.Token {
	newToken := tk.CreateUnidentifiedToken(symbolString, tkzr.currentLineNumber, tkzr.currentTabLevel)
	newToken.SetValues(RULENAME_SYMBOL, tkzr.identifySymbol(symbolString))
	if newToken.SymbolicName == SYMBOLIC_NAME_WHITESPACE || newToken.SymbolicName == SYMBOLIC_NAME_NEWLINE {
		newToken.RuleName = RULENAME_OTHER
	}
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
		symbolicName = SYMBOLIC_NAME_NON_KEYWORD
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
	if symbol == " " {
		symbolicName = SYMBOLIC_NAME_WHITESPACE
	} else if symbol == "\n" {
		symbolicName = SYMBOLIC_NAME_NEWLINE
	} else if symbolicName == "" {
		symbolicName = SYMBOLIC_NAME_UNKNOWN_SYMBOL
	}
	return symbolicName
}

func (tkzr *Tokenizer) createTokenType(tokenString string) *tk.Token {
	possibleKeyword := tkzr.identifyKeyword(tokenString)
	if possibleKeyword == SYMBOLIC_NAME_NON_KEYWORD {
		return tkzr.createSymbolToken(tokenString)
	}
	return tkzr.createKeywordToken(tokenString)
}
