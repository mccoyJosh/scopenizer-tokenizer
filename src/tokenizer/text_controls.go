package tokenizer

import (
	"errors"
	"tp/src/util"
)

// TextSize
// This returns the size of the text
// This results in an error and a panic
// if the text is currently nil
func (tkzr *Tokenizer) TextSize() int {
	if tkzr.Text == nil {
		err := errors.New("text is nil. Cannot get text size")
		util.Error(err.Error(), err)
		panic(err)
	}
	return len(*tkzr.Text)
}

// TextRange
// This will return a substring from the text. This substring
// will be the contents between the 'begin' and 'end' integer parameters.
// The being parameter is inclusive while the end parameter is exclusive.
// This will return an error if the indices are out of range.
func (tkzr *Tokenizer) TextRange(begin int, end int) (string, error) {
	if begin >= 0 && begin < tkzr.TextSize() && end > 0 && end < tkzr.TextSize() && begin != end {
		return (*tkzr.Text)[begin:end], nil
	}
	return "", errors.New("TextRange bounds were out of bounds or otherwise invalid")
}

// gatherWhitespace
// This start at the current index (+ 1) and get the whitespace found there.
// It will begin looking for the whitespace at this index and stop whenever
// it encounters a non-whitespace character.
//
// It accepts a boolean statement deciding whether or/not
// the index should be updated to wherever the whitespace stopped.
//
// It will return a string of this whitespace.
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
