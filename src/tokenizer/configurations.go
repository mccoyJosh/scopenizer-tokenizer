package tokenizer

import (
	"errors"
	"fmt"
)

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
