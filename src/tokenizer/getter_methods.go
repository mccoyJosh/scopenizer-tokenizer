package tokenizer

import (
	"errors"
	"fmt"
	"tp/src/util"
)

// GetChar
// Gets a character from the text as long as the provided
// integer index is within the bounds of the Text. If
// it out of bounds, it panics.
func (tkzr *Tokenizer) GetChar(index int) rune {
	if !tkzr.DetermineIfIndexInBound(index) {
		err := errors.New(fmt.Sprintf("index %d is out of bounds of text size %d", index, tkzr.TextSize()))
		util.Error(err.Error(), err)
		panic(err)
	}
	return rune((*tkzr.Text)[index])
}

// GetCurrentTabLevel
// Returns the current tab level the tokenizer is at
func (tkzr *Tokenizer) GetCurrentTabLevel() int {
	return tkzr.currentTabLevel
}

// GetCurrentLineNumber
// Returns the current line number the tokenizer is at
func (tkzr *Tokenizer) GetCurrentLineNumber() int {
	return tkzr.currentLineNumber
}

// CurrentChar
// Returns the current char the tokenizer is at.
// This is determined by wherever the index is at
func (tkzr *Tokenizer) CurrentChar() rune {
	return tkzr.GetChar(tkzr.currentIndex)
}
