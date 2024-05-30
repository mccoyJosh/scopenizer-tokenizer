package tokenizer

import (
	"errors"
	"fmt"
)

// ConfigureIgnores
// Sets up what kind of tokens should be ignored in the final output of tokens.
// This method is not necessary to be run since default values are provided, those defining that
// strings and comments are included, while whitespaces and newlines are ignored.
//
// ignoreString: determines whether string tokens should be ignored, i.e., not added to final output
//
// ignoreComments: determines whether comment tokens should be ignored, i.e., not added to final output
//
// ignoreWhitespaces: determines whether whitespace tokens should be ignored, i.e., not added to final output
//
// ignoreNewLines: determines whether new line tokens should be ignored, i.e., not added to final output
func (tkzr *Tokenizer) ConfigureIgnores(ignoreString bool, ignoreComments bool, ignoreWhitespaces bool, ignoreNewLines bool) {
	tkzr.IncludeStrings = !ignoreString
	tkzr.IncludeComments = !ignoreComments
	tkzr.IgnoreWhitespace = ignoreWhitespaces
	tkzr.IgnoreNewLines = ignoreNewLines
}

// ConfigureGeneral
// Sets up general functionality of a tokenizer.
// The variables this method sets up are necessary for the tokenizer to work.
//
// language: this parameter takes in a string which assigns the LanguageType variable of the tokenizer object
//
// symbols: this parameters will set the symbols which will be discovered by the tokenizer. It expects an array of arrays of strings,
// which have 2 items, that being the assigned name for the symbol and the symbol itself.
//
// keywords: this parameter will set the keywords for the tokenizer. It expects an array of strings of all the keywords this language expects.
//
// isKeywordCharacterFunction: this parameter expects a function that takes a rune and returns a bool. This should return true for characters
// which are considered valid characters for keywords in this language.
func (tkzr *Tokenizer) ConfigureGeneral(language string, symbols [][]string, keywords []string, isKeywordCharacterFunction func(c rune) bool) {
	tkzr.LanguageType = language
	tkzr.Symbols = symbols
	tkzr.Keywords = keywords
	tkzr.IsKeywordCharacter = isKeywordCharacterFunction
}

// ConfigureScope
// This function sets up the functions to control when a scope begins and ends
//
// startFunction: the parameter expects a function whose one parameter is a pointer to a tokenizer object,
// and it should return a boolean. It should return true if the index the tokenizer is currently on is the start of a
// scope object
//
// endFunction: the parameter expects a function whose one parameter is a pointer to a tokenizer object,
// and it should return a boolean. It should return true if the index the tokenizer is currently on is the end of a
// scope object
func (tkzr *Tokenizer) ConfigureScope(startFunction func(tkzr *Tokenizer) bool, endFunction func(tkzr *Tokenizer) bool) {
	tkzr.ScopeStartFunction = startFunction
	tkzr.ScopeEndFunction = endFunction
}

// ConfigureString
// This function sets up the functions to control when a string begins and ends
//
// startFunction: the parameter expects a function whose one parameter is a pointer to a tokenizer object,
// and it should return a boolean. It should return true if the index the tokenizer is currently on is the start of a
// string (or similar object)
//
// endFunction: the parameter expects a function whose one parameter is a pointer to a tokenizer object,
// and it should return a boolean. It should return true if the index the tokenizer is currently on is the end of a
// string (or similar object)
func (tkzr *Tokenizer) ConfigureString(startFunction func(tkzr *Tokenizer) bool, endFunction func(tkzr *Tokenizer) bool) {
	tkzr.StringStartFunction = startFunction
	tkzr.StringEndFunction = endFunction
}

// ConfigureComment
// This function sets up the functions to control when a comment begins and ends
//
// startFunction: the parameter expects a function whose one parameter is a pointer to a tokenizer object,
// and it should return a boolean. It should return true if the index the tokenizer is currently on is the start of a
// comment
//
// endFunction: the parameter expects a function whose one parameter is a pointer to a tokenizer object,
// and it should return a boolean. It should return true if the index the tokenizer is currently on is the end of a
// comment
func (tkzr *Tokenizer) ConfigureComment(startFunction func(tkzr *Tokenizer) bool, endFunction func(tkzr *Tokenizer) bool) {
	tkzr.CommentStartFunction = startFunction
	tkzr.CommentEndFunction = endFunction
}

// IsConfigured
// This determines whether the tokenizer is fully set up, as in, everything that needs to be configured is configured.
//
// If it finds that it is not configured correctly, it will return an error
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
