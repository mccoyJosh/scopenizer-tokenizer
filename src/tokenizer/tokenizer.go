package tokenizer

import (
	"errors"
	"fmt"
	tk "tp/src/tokenizer/tokens"
	"tp/src/util"
)

// Tokenizer
// Defines Tokenizer object.
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

	// Final Steps
	FinalSteps func(tkzr *Tokenizer, finalScope *tk.ScopeObj) error
}

// Tokenize
// Takes a string and tokenizes the contents of it into a ScopeObj object.
// Returns an error if the tokenizer is not configured correctly or an error results from the final steps.
func (tkzr *Tokenizer) Tokenize(text string) (tk.ScopeObj, error) {
	err := tkzr.IsConfigured()
	if err != nil {
		return tk.InitScope(), err
	}

	tkzr.initTempVariables(&text)

	finalScope := tk.InitScope()
	finalScope.SetType("File")
	tkzr.currentScope = &finalScope

	for tkzr.IndexInBound() {
		tkzr.skipIncrement = false

		if !tkzr.applyFunctions() {
			// Not a scope identifier, not a comment, not a string
			char := tkzr.CurrentChar()
			if tkzr.IsKeywordCharacter(char) {
				tkzr.potentialKeyword += string(char)
			} else { // Found a symbol
				// The previous keyword is over and needs to be added
				tkzr.addPotentialKeyword()

				// Current symbol needs to be added
				tkzr.addSymbol(char)
			}
		}

		if !tkzr.skipIncrement {
			tkzr.currentIndex++
		}
	}

	if tkzr.potentialKeyword != "" {
		tkzr.currentScope.Push(tkzr.createKeywordToken(tkzr.potentialKeyword))
		tkzr.potentialKeyword = ""
	}

	if tkzr.FinalSteps != nil {
		err = tkzr.FinalSteps(tkzr, &finalScope)
		if err != nil {
			return finalScope, err
		}
	}

	return finalScope, nil
}

// applyBeforeFunction
func (tkzr *Tokenizer) applyBeforeFunction() {
	tkzr.tempIgnoreChangesFromIncrement = true
	tkzr.addPotentialKeyword()
}

// applyAfterFunction
func (tkzr *Tokenizer) applyAfterFunction() {
	for i := 0; i < len(tkzr.EndInfo)-1; i++ {
		if !tkzr.skipIncrement {
			tkzr.IncrementIndex() // TODO: Determine whether this is a good long term solution or should be left for anonymous functions to deal with
		}
	}
	tkzr.StartInfo = ""
	tkzr.EndInfo = ""
	tkzr.tempIgnoreChangesFromIncrement = false
}

// applyFunctions
func (tkzr *Tokenizer) applyFunctions() bool {

	if tkzr.StringStartFunction(tkzr) {
		tkzr.applyBeforeFunction()
		// FOUND STRING
		resultingToken := tkzr.applyFunctionUntilFailureTokenCreation(tkzr.StringEndFunction, SYMBOLIC_NAME_STRING)
		if tkzr.IncludeStrings {
			tkzr.currentScope.Push(resultingToken)
		}
		tkzr.applyAfterFunction()
		return true
	}

	if tkzr.CommentStartFunction(tkzr) {
		tkzr.applyBeforeFunction()
		// FOUND COMMENT
		resultingToken := tkzr.applyFunctionUntilFailureTokenCreation(tkzr.CommentEndFunction, SYMBOLIC_NAME_COMMENT)
		if tkzr.IncludeComments {
			tkzr.currentScope.Push(resultingToken)
		}
		tkzr.applyAfterFunction()
		return true
	}

	if tkzr.ScopeStartFunction(tkzr) {
		tkzr.applyBeforeFunction()
		// FOUND SCOPE START
		preScopeToken := tkzr.createTokenType(tkzr.StartInfo)
		tkzr.currentScope.Push(preScopeToken)

		newScopeTkn := tk.InitScopeToken()
		tkzr.currentScope.Push(newScopeTkn)
		tkzr.currentScope = newScopeTkn.GetScopeToken()
		tkzr.applyAfterFunction()
		return true
	}

	if tkzr.ScopeEndFunction(tkzr) {
		tkzr.applyBeforeFunction()
		// FOUND SCOPE END
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
		tkzr.applyAfterFunction()
		return true
	}

	return false
}
