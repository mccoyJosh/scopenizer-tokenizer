package tokenizer

import (
	tk "tp/src/tokenizer/tokens"
	"tp/src/util"
)

// Index
// Returns the index that the tokenizer is currently at.
func (tkzr *Tokenizer) Index() int {
	return tkzr.currentIndex
}

// IndexInBound
// Determine whether the current index is in bound of the text or not (returns false if not).
func (tkzr *Tokenizer) IndexInBound() bool {
	return tkzr.DetermineIfIndexInBound(tkzr.Index())
}

// DetermineIfIndexInBound
// Determines whether a provided integer index is within the
// bounds of the text (returns false if not).
func (tkzr *Tokenizer) DetermineIfIndexInBound(index int) bool {
	return index < len(*tkzr.Text)
}

// SkipIncrement
// Skips the next increment of the index
func (tkzr *Tokenizer) SkipIncrement() {
	tkzr.skipIncrement = true
}

// IncrementIndex
// This will move the index forward 1 (add one to it) and
// deal with any new line characters. This ensures that
// the current line count is correct.
func (tkzr *Tokenizer) IncrementIndex() {
	tkzr.currentIndex++
	tkzr.dealWithNewline()
	for tkzr.IndexInBound() && tkzr.CurrentChar() == '\n' {
		tkzr.currentIndex++
		tkzr.dealWithNewline()
	}
}

// dealWithNewline
// This checks whether the current char is a newline.
// If it is a newline, this will automatically update the current line counter and
// get the tab level of the new line.
func (tkzr *Tokenizer) dealWithNewline() {
	// Ensure we are on a newline
	if tkzr.IndexInBound() && tkzr.CurrentChar() == '\n' {
		if !tkzr.tempIgnoreChangesFromIncrement {
			tkzr.addPotentialKeyword()
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
