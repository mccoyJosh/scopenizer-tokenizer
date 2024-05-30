package tokenizer

import "unicode"

// CreateDullTokenizer
// Creates a 'dull' tokenizer object.
// This is a mostly empty example of what a tokenizer object should look like.
// This will essentially split up words and symbols from each other and tokenize those.
// The anonymous functions all just return either true or false and not much else...
// you obviously would fill these out in the real deal.
//
// This returns a pointer to a newly created tokenizer object
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
