package tokenizer

import (
	"strings"
	tk "tp/src/tokenizer/tokens"
)

// addPotentialKeyword
// This will add the potential keyword to the current scope as a keyword token.
// If potentialKeyword is an empty string, this method does nothing
func (tkzr *Tokenizer) addPotentialKeyword() {
	if tkzr.potentialKeyword != "" {
		tkzr.currentScope.Push(tkzr.createKeywordToken(tkzr.potentialKeyword))
		tkzr.potentialKeyword = ""
	}
}

// addSymbol
// This method accepts a single character rune as a parameter
// and will identify and add the char to the current scope as
// a symbol token
func (tkzr *Tokenizer) addSymbol(char rune) {
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

// createKeywordToken
// This will take a keyword string, identify its type,
// create a keyword token, and return a pointer to the newly created token
func (tkzr *Tokenizer) createKeywordToken(keywordString string) *tk.Token {
	newToken := tk.CreateUnidentifiedToken(keywordString, tkzr.currentLineNumber, tkzr.currentTabLevel)
	newToken.SetValues(RULENAME_KEYWORD, tkzr.identifyKeyword(keywordString))
	return &newToken
}

// createSymbolToken
// This will take a symbol string, identify its type,
// create a symbol token, and return a pointer to the newly created token
func (tkzr *Tokenizer) createSymbolToken(symbolString string) *tk.Token {
	newToken := tk.CreateUnidentifiedToken(symbolString, tkzr.currentLineNumber, tkzr.currentTabLevel)
	newToken.SetValues(RULENAME_SYMBOL, tkzr.identifySymbol(symbolString))
	if newToken.SymbolicName == SYMBOLIC_NAME_WHITESPACE || newToken.SymbolicName == SYMBOLIC_NAME_NEWLINE {
		newToken.RuleName = RULENAME_OTHER
	}
	return &newToken
}

// identifyKeyword
// This takes a keyword as a string, and it will go through the keywords array and
// see if this is an identifiable keyword. If it is not an identifiable keyword,
// it will be labeled as SYMBOLIC_NAME_NON_KEYWORD. The label is returned as a string.
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

// identifySymbol
// This takes a symbol as a parameter and goes through the Symbols array
// to see if this symbol is present. If the symbol is found, its identified name
// if returned. If the symbol cannot be identified, it is returned as
// SYMBOLIC_NAME_UNKNOWN_SYMBOL.
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

// createTokenType
// This method takes a token string, and will find out whether it is a keyword,
// or symbol, and fully identify it. This will then create a token using the identified string and return it.
func (tkzr *Tokenizer) createTokenType(tokenString string) *tk.Token {
	possibleKeyword := tkzr.identifyKeyword(tokenString)
	if possibleKeyword == SYMBOLIC_NAME_NON_KEYWORD {
		return tkzr.createSymbolToken(tokenString)
	}
	return tkzr.createKeywordToken(tokenString)
}

// applyFunctionUntilFailureTokenCreation
// This method takes in a method to be checked every index until it fails. Until it fails,
// it will be accumulating the characters and create a token which it will return.
// The symbolic name will be used to set the values for the tokens.
// The rule nam for the returning token will be set to RULENAME_OTHER.
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
