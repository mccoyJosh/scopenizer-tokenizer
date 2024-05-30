package tokenizer

// GenerateDefaultTokenizerObject
// This creates a tokenizer object with most
// variables initialized with their default values
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

// initTempVariables
// This initializes various temporary variables
// needed for the tokenizer to function.
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

// initSpaceSizeString
// This takes the expected number of spaces in tabs and
// creates a string which represents the tab character.
func (tkzr *Tokenizer) initSpaceSizeString() {
	tkzr.spaceSizeString = ""
	for i := 0; i < tkzr.NumOfSpacesEquallyTab; i++ {
		tkzr.spaceSizeString += " "
	}
}
