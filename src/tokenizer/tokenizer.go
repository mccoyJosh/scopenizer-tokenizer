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
			return true
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
			return false
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

func (tkzr *Tokenizer) ConfigureIgnores(ignoreString bool, ignoreComments bool, ignoreWhitespaces bool, ignoreNewLines bool) {
	tkzr.IncludeStrings = !ignoreString
	tkzr.IncludeComments = !ignoreComments
	tkzr.IgnoreWhitespace = ignoreWhitespaces
	tkzr.IgnoreNewLines = ignoreNewLines
}

func (tkzr *Tokenizer) ConfigureGeneral(language string, symbols [][]string, keywords []string, isKeywordCharacterFunction func(c rune) bool) {
	tkzr.LanguageType = language
	tkzr.Symbols = symbols
	tkzr.Keywords = keywords
	tkzr.IsKeywordCharacter = isKeywordCharacterFunction
}

func (tkzr *Tokenizer) ConfigureScope(startFunction func(tkzr *Tokenizer) bool, endFunction func(tkzr *Tokenizer) bool) {
	tkzr.ScopeStartFunction = startFunction
	tkzr.ScopeEndFunction = endFunction
}

func (tkzr *Tokenizer) ConfigureString(startFunction func(tkzr *Tokenizer) bool, endFunction func(tkzr *Tokenizer) bool) {
	tkzr.StringStartFunction = startFunction
	tkzr.StringEndFunction = endFunction
}

func (tkzr *Tokenizer) ConfigureComment(startFunction func(tkzr *Tokenizer) bool, endFunction func(tkzr *Tokenizer) bool) {
	tkzr.CommentStartFunction = startFunction
	tkzr.CommentEndFunction = endFunction
}

func (tkzr *Tokenizer) IsConfigured() error {
	errorString := ""

	// General Configure
	if tkzr.LanguageType == "" {
		errorString += fmt.Sprintf("LanguageType not configured correctly (cannot be empty string)... USE .ConfigureGeneral to fix\n")
	}
	if tkzr.Symbols == nil {
		errorString += fmt.Sprintf("Symbols not configured correctly (cannot be nil)... USE .ConfigureGeneral to fix\n")
	}
	if tkzr.Keywords == nil {
		errorString += fmt.Sprintf("Keywords not configured correctly (cannot be nil)... USE .ConfigureGeneral to fix\n")
	}
	if tkzr.IsKeywordCharacter == nil {
		errorString += fmt.Sprintf("IsKeywordsCharacter not configured correctly (cannot be nil)... USE .ConfigureGeneral to fix\n")
	}

	// Scope Configure
	if tkzr.ScopeStartFunction == nil {
		errorString += fmt.Sprintf("ScopeStartFunction not configured correctly (cannot be nil)... USE .ConfigureScope to fix\n")
	}
	if tkzr.ScopeEndFunction == nil {
		errorString += fmt.Sprintf("ScopeEndFunction not configured correctly (cannot be nil)... USE .ConfigureScope to fix\n")
	}

	// String Configure
	if tkzr.StringStartFunction == nil {
		errorString += fmt.Sprintf("StringStartFunction not configured correctly (cannot be nil)... USE .ConfigureString to fix\n")
	}
	if tkzr.StringEndFunction == nil {
		errorString += fmt.Sprintf("StringEndFunction not configured correctly (cannot be nil)... USE .ConfigureString to fix\n")
	}

	// Comment Configure
	if tkzr.CommentStartFunction == nil {
		errorString += fmt.Sprintf("CommentStartFunction not configured correctly (cannot be nil)... USE .ConfigureComment to fix\n")
	}
	if tkzr.CommentEndFunction == nil {
		errorString += fmt.Sprintf("CommentEndFunction not configured correctly (cannot be nil)... USE .ConfigureComment to fix\n")
	}

	if errorString != "" {
		err := errors.New(errorString)
		return err
	}
	return nil
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
		foundString := tkzr.StringStartFunction(tkzr)
		foundComment := tkzr.CommentStartFunction(tkzr)
		foundStartScope := tkzr.ScopeStartFunction(tkzr)
		foundEndScope := tkzr.ScopeEndFunction(tkzr)
		//foundWhiteSpace := whiteSpaceStartFunction(tkzr)

		if foundString || foundComment || foundStartScope || foundEndScope {
			tkzr.tempIgnoreChangesFromIncrement = true

			if tkzr.potentialKeyword != "" {
				tkzr.currentScope.Push(tkzr.createKeywordToken(tkzr.potentialKeyword))
				tkzr.potentialKeyword = ""
			}

			//if foundWhiteSpace {
			//	resultingToken := tkzr.applyFunctionUntilFailureTokenCreation(whiteSpaceEndFunction, "WHITESPACE")
			//	if !tkzr.IgnoreWhitespace {
			//		tkzr.currentScope.Push(resultingToken)
			//	}
			//}
			if foundString {
				resultingToken := tkzr.applyFunctionUntilFailureTokenCreation(tkzr.StringEndFunction, "STRING")
				if tkzr.IncludeStrings {
					tkzr.currentScope.Push(resultingToken)
				}
			} else if foundComment {
				resultingToken := tkzr.applyFunctionUntilFailureTokenCreation(tkzr.CommentEndFunction, "COMMENT")
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
				postScopeToken := tkzr.createTokenType(tkzr.EndInfo)
				tkzr.currentScope.Push(postScopeToken)
			}
			for i := 0; i < len(tkzr.EndInfo)-1; i++ {
				tkzr.IncrementIndex() // TODO: Determine whether this is a good long term solution or should be left for anonymous functions to deal with
			}
			tkzr.tempIgnoreChangesFromIncrement = false
		} else { // Not a scope identifier, not a comment, not a string
			char := text[tkzr.Index()]
			if tkzr.IsKeywordCharacter(rune(char)) {
				tkzr.potentialKeyword += string(char)
			} else { // Found a symbol, which needs to be added and
				if tkzr.potentialKeyword != "" {
					tkzr.currentScope.Push(tkzr.createKeywordToken(tkzr.potentialKeyword))
					tkzr.potentialKeyword = ""
				}
				newSymbolToken := tkzr.createSymbolToken(string(char))

				if newSymbolToken.SymbolicName == "WHITESPACE" {
					if !tkzr.IgnoreWhitespace {
						newSymbolToken.RuleName = "OTHER"
						tkzr.currentScope.Push(newSymbolToken)
					}
				} else if newSymbolToken.SymbolicName == "NEWLINE" {
					if !tkzr.IgnoreWhitespace {
						newSymbolToken.RuleName = "OTHER"
						tkzr.currentScope.Push(newSymbolToken)
					}
				} else {
					tkzr.currentScope.Push(newSymbolToken)
				}
			}
		}
		tkzr.IncrementIndex()
	}

	if tkzr.potentialKeyword != "" { // TODO: SEE IF ONE LAST CHECK IS NECESSARY
		tkzr.currentScope.Push(tkzr.createKeywordToken(tkzr.potentialKeyword))
		tkzr.potentialKeyword = ""
	}

	return finalScope
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
	finalToken := tk.CreateUnidentifiedToken(tokenText, lineNumber, tabLevel)
	finalToken.SetValues("OTHER", symbolicName)
	return &finalToken
}

func (tkzr *Tokenizer) GetCurrentLineNumber() int {
	return tkzr.currentLineNumber
}

func whiteSpaceStartFunction(tkzr *Tokenizer) bool {
	return util.IsWhitespaceCharacter(tkzr.CurrentChar())
}

func whiteSpaceEndFunction(tkzr *Tokenizer) bool {
	return util.IsWhitespaceCharacter(tkzr.CurrentChar())
}

func (tkzr *Tokenizer) IncrementIndex() {
	tkzr.currentIndex++
	if tkzr.currentIndex < tkzr.TextSize() && tkzr.CurrentChar() == '\n' { // Encounters a newline
		if !tkzr.IgnoreNewLines && !tkzr.tempIgnoreChangesFromIncrement { // (potentially) adding the new line character token
			newToken := tk.CreateUnidentifiedToken("\n", tkzr.currentLineNumber, tkzr.currentTabLevel)
			newToken.SetValues("OTHER", "NEWLINE")
			tkzr.currentScope.Push(&newToken)
		}
		tkzr.currentLineNumber++

		// Gets the tab level for this line
		gatheredWhitespace := tkzr.GatherWhitespace(!tkzr.tempIgnoreChangesFromIncrement)
		numOfTabs := util.DetermineNumberOfTabs(gatheredWhitespace, tkzr.NumOfSpacesEquallyTab, true)
		tkzr.currentTabLevel = numOfTabs
		if !tkzr.IgnoreWhitespace && !tkzr.tempIgnoreChangesFromIncrement {
			newToken := tk.CreateUnidentifiedToken(gatheredWhitespace, tkzr.currentLineNumber, tkzr.currentTabLevel)
			newToken.SetValues("OTHER", "WHITESPACE")
			tkzr.currentScope.Push(&newToken)
		}
		tkzr.IncrementIndex()
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
	if symbol == " " {
		symbolicName = "WHITESPACE"
	} else if symbol == "\n" {
		symbolicName = "NEWLINE"
	} else if symbolicName == "" {
		symbolicName = "UNKNOWN"
	}
	return symbolicName
}

func (tkzr *Tokenizer) createTokenType(tokenString string) *tk.Token {
	possibleKeyword := tkzr.identifyKeyword(tokenString)
	if possibleKeyword == "IDENTIFIER" {
		return tkzr.createSymbolToken(tokenString)
	}
	return tkzr.createKeywordToken(tokenString)
}
