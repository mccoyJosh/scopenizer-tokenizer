package tokens

import (
	"errors"
	"fmt"
	"tp/src/util"
)

type ScopeToken struct {
	scopeType    string
	tokenList    []*interface{}
	scopeIndices []int
	size         int
	totalSize    int
}

// InitTokenList
// This will construct a token list object
func InitTokenList() ScopeToken {
	return ScopeToken{
		scopeType:    UNKNOWN_SCOPE_STRING,
		tokenList:    make([]*interface{}, 0),
		scopeIndices: make([]int, 0),
		size:         0,
		// totalSize is initially -1 as it NEEDS be updated once a scope is finished. -1 should indicate this was not done
		totalSize: -1,
	}
}

// Size
// This returns the number of tokens in this list as an integer
// This size does not include the number of tokens inside of inner scopes
func (tl *ScopeToken) Size() int {
	return tl.size
}

// TotalSize
// This returns the number of tokens in the entire list.
// This size DOES include the number of tokens in inner scopes.
// This does not include the number of scopes.
// WARNING: This does re-evaluate the total size everytime it is called due to
// the fact one could change the token list at any point!
func (tl *ScopeToken) TotalSize() int {
	tl.setTotalSize()
	return tl.totalSize
}

// setTotalSize
// This will go through the scopes found within this list
// and get the total number of tokens found within.
// This should get the total number of tokens within this and all inner scopes.
func (tl *ScopeToken) setTotalSize() int {
	// If there are no inner scopes, we return
	if len(tl.scopeIndices) <= 0 {
		tl.totalSize = tl.size
		return tl.totalSize
	}

	innerTotals := 0
	for _, index := range tl.scopeIndices {
		token := tl.tokenList[index] //
		if ValidScopeToken(token) {
			scopeToken := (*token).(ScopeToken)
			innerTotals += scopeToken.setTotalSize()
		} else {
			err := errors.New(fmt.Sprintf("Non-Scope Token Found At Supposed Scope Index: %d", index))
			util.Error(err.Error(), err)
		}
	}

	// Removes the scope tokens for the count, as we only want 'real' tokens
	total := innerTotals + tl.size - len(tl.scopeIndices)
	tl.totalSize = total
	return total
}

// FixScopeIndices
// This will assign the correct values of scope tokens in the event where
// they may need to be fixed
func (tl *ScopeToken) FixScopeIndices() {
	tl.scopeIndices = make([]int, 0)
	for index, token := range tl.tokenList {
		if ValidScopeToken(token) {
			tl.scopeIndices = append(tl.scopeIndices, index)
		}
	}
}

// Push
// This adds a token to the token list at the end of the list,
// much like one would push an item to the top of a stack.
func (tl *ScopeToken) Push(i *interface{}) {
	tl.Insert(i, tl.size)
}

// Insert
// Inserts a token into the token list. Returns nothing.
// If this token is not inserted at the end of the list,
// it may need to fix the scope index list
func (tl *ScopeToken) Insert(i *interface{}, index int) {
	if ValidTokenType(i) && index >= 0 && index <= tl.size {
		tl.tokenList = append(tl.tokenList[:index+1], tl.tokenList[index:]...)
		tl.tokenList[index] = i
		tl.size++

		if index != tl.size-1 {
			tl.FixScopeIndices()
		} else if ValidScopeToken(i) {
			tl.scopeIndices = append(tl.scopeIndices, index)
		}
	} else {
		err := errors.New(fmt.Sprintf("Invalid (out of bounds) index provided to token list INSERT: %d", index))
		util.Error(err.Error(), err)
	}
}

// GetIndex
// This will return the token at the given index.
// Returns a pointer to an interface object. This object
// Will need to be converted to either a ScopeToken or Token object
func (tl *ScopeToken) GetIndex(index int) *interface{} {
	if index >= 0 && index < tl.size {
		return tl.tokenList[index]
	}
	err := errors.New(fmt.Sprintf("Invalid (out of bounds) index provided to token list GET INDEX: %d", index))
	util.Error(err.Error(), err)
	return nil
}

// Pop
// This will get and remove the item at the end of the token list,
// much like one would pop the top item off the top of a stack
func (tl *ScopeToken) Pop() *interface{} {
	i := tl.Front()
	tl.Delete(tl.size - 1)
	return i
}

// Front
// Returns the items at the front of the list of tokens,
// i.e. the value at the index of size - 1
func (tl *ScopeToken) Front() *interface{} {
	return tl.GetIndex(tl.size - 1)
}

// Delete
// This removes a token from the token list given its index.
// If this index is out of bounds, nothing is done and an error is printed
func (tl *ScopeToken) Delete(index int) {
	if index >= 0 && index < tl.size {
		tl.tokenList = append(tl.tokenList[:index], tl.tokenList[index+1:]...)
		tl.size--
		tl.FixScopeIndices()
	} else {
		err := errors.New(fmt.Sprintf("Invalid (out of bounds) index provided to token list DELETE: %d", index))
		util.Error(err.Error(), err)
	}
}

// ScopifyRange
// This method, given a start index and an end index, will
// add the given range of tokens to a scope token and place the new
// scope token at the start index. The tokens in the range will only
// be in the new scope token and will be removed from the token it is called from.
func (tl *ScopeToken) ScopifyRange(start int, end int) {
	if start < end && start >= 0 && start < tl.size && end >= 0 && end < tl.size {
		// TODO: actually 'scopify' a range of tokens
	} else {
		err := errors.New(fmt.Sprintf("Invalid (out of bounds) index provided to token list SCOPIFY RANGE: %d, %d", start, end))
		util.Error(err.Error(), err)
	}
}
