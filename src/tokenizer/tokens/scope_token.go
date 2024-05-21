package tokens

import (
	"errors"
	"fmt"
	"tp/src/util"
)

// ScopeObj
// Defines the scope object
//
// scopeType: Defines what type of scope this is; will be an identifier for the info contained within. For example: if this scope is a method, the info within will be of type MethodInfo, or whatever info is needed for this type of object
//
// info: Custom info pertaining to this scope
//
// tokenList: List of all tokens in this given scope
//
// scopeIndices: This array give quick access to where ScopeObj objects are found within the tokenList by storing the indices of said ScopeObj's
//
// size: This is the number of tokens within this scope object. DOES NOT INCLUDE INNER SCOPE SIZES
type ScopeObj struct {
	scopeType    string
	info         *any
	tokenList    []*Token
	scopeIndices []int
	size         int
}

// InitScope
// This will construct a scope object.
// This method can take an array of tokens as parameters
// as to have some initial tokens in the scope
func InitScope(lists ...[]*Token) ScopeObj {
	sc := ScopeObj{
		scopeType:    UNKNOWN_SCOPE_STRING,
		info:         nil,
		tokenList:    make([]*Token, 0),
		scopeIndices: make([]int, 0),
		size:         0,
	}

	for _, list := range lists {
		for _, token := range list {
			sc.Push(token)
		}
	}

	if len(lists) != 0 {
		sc.FixScopeIndices()
	}

	return sc
}

func InitScopeToken(scopes ...*ScopeObj) *Token {
	// If a SINGLE scope is provided, it will use it as the scope within the scope token
	// If not, it will create an empty scope object for this token
	var providedScope ScopeObj
	if len(scopes) == 1 {
		providedScope = *scopes[0]
	} else {
		providedScope = InitScope()
	}

	return &Token{
		LineNumber:    0,
		TabNumber:     0,
		BracketNumber: 0,
		SymbolicName:  "",
		RuleName:      SCOPE_TOKEN_STIRNG,
		Text:          "",
		scopeToken:    &providedScope,
	}
}

// Size
// This returns the number of tokens in this list as an integer
// This size does not include the number of tokens inside of inner scopes
func (tl *ScopeObj) Size() int {
	return tl.size
}

func (tl *ScopeObj) GetType() string {
	return tl.scopeType
}

func (tl *ScopeObj) SetType(typeString string) {
	tl.scopeType = typeString
}

func (tl *ScopeObj) GetNumberOfScopes() int {
	return len(tl.scopeIndices)
}

// GetScope
// THIS ASSUMES FixScopeIndices HAS BEEN RAN
// If you have been using the Insert and Delete methods, this method SHOULD HAVE BEEN run already
// The goal of storing the indices of the token is to save time from searching for them every time, and calling FixScopeIndices everytime in this method would defeat that purpose
//
// index: the index of the scope token in the scope indices array. If you use 0 as index, it would be essentially asking for the first scope
func (tl *ScopeObj) GetScope(index int) *ScopeObj {
	if index < tl.GetNumberOfScopes() {
		return tl.tokenList[tl.scopeIndices[index]].scopeToken
	}
	err := errors.New(fmt.Sprintf("Invalid (out of bounds) index provided to scope indicies list GETSCOPE: %d", index))
	util.Error(err.Error(), err)
	return nil
}

// TotalSize
// This returns the number of tokens in the entire list.
// This size DOES include the number of tokens in inner scopes.
// This does not include the number of scopes.
// WARNING: This does re-evaluate the total size everytime it is called due to
// the fact one could change the token list at any point!
func (tl *ScopeObj) TotalSize() int {
	return tl.setTotalSize()
}

// setTotalSize
// This will go through the scopes found within this list
// and get the total number of tokens found within.
// This should get the total number of tokens within this and all inner scopes.
func (tl *ScopeObj) setTotalSize() int {
	// If there are no inner scopes, we return
	if len(tl.scopeIndices) <= 0 {
		return tl.size
	}

	innerTotals := 0
	for _, index := range tl.scopeIndices {
		token := tl.tokenList[index] //
		if token.ValidScopeToken() {
			scopeToken := token.scopeToken
			innerTotals += scopeToken.setTotalSize()
		} else {
			err := errors.New(fmt.Sprintf("Non-Scope Token Found At Supposed Scope Index: %d", index))
			util.Error(err.Error(), err)
		}
	}

	// Removes the scope tokens for the count, as we only want 'real' tokens
	total := innerTotals + tl.size - len(tl.scopeIndices)
	return total
}

// FixScopeIndices
// This will assign the correct values of scope tokens in the event where
// they may need to be fixed
func (tl *ScopeObj) FixScopeIndices() {
	tl.scopeIndices = make([]int, 0)
	for index, token := range tl.tokenList {
		if token.ValidScopeToken() {
			tl.scopeIndices = append(tl.scopeIndices, index)
		}
	}
}

// Push
// This adds a token to the token list at the end of the list,
// much like one would push an item to the top of a stack.
func (tl *ScopeObj) Push(tt *Token) {
	tl.Insert(tt, tl.size)
}

// Insert
// Inserts a token into the token list. Returns nothing.
// If this token is not inserted at the end of the list,
// it may need to fix the scope index list
func (tl *ScopeObj) Insert(tt *Token, index int) {
	if index >= 0 && index <= tl.size {
		if tl.size == index {
			tl.tokenList = append(tl.tokenList, tt)
		} else {
			tl.tokenList = append(tl.tokenList[:index+1], tl.tokenList[index:]...)
			tl.tokenList[index] = tt
		}
		tl.size++

		if index != tl.size-1 {
			tl.FixScopeIndices()
		} else if tt.ValidScopeToken() {
			tl.scopeIndices = append(tl.scopeIndices, index)
		}
	} else {
		err := errors.New(fmt.Sprintf("Invalid (out of bounds) index provided to token list INSERT: %d", index))
		util.Error(err.Error(), err)
	}
}

// At
// This will return the token at the given index.
// Returns a pointer to an interface object. This object
// Will need to be converted to either a ScopeObj or Token object
func (tl *ScopeObj) At(index int) *Token {
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
func (tl *ScopeObj) Pop() *Token {
	i := tl.Front()
	tl.Delete(tl.size - 1)
	return i
}

// Front
// Returns the items at the front of the list of tokens,
// i.e. the value at the index of size - 1
func (tl *ScopeObj) Front() *Token {
	return tl.At(tl.size - 1)
}

// Delete
// This removes a token from the token list given its index.
// If this index is out of bounds, nothing is done and an error is printed
func (tl *ScopeObj) Delete(index int) {
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
func (tl *ScopeObj) ScopifyRange(start int, end int) {
	if start < end && start >= 0 && start < tl.size && end >= 0 && end < tl.size {
		tokensSubset := tl.tokenList[start : end+1] // TODO: ensure this includes said element
		tokensSubsetLength := end - start + 1
		newScopeObj := InitScope(tokensSubset)
		for i := 0; i < tokensSubsetLength; i++ {
			tl.Delete(start)
		}
		tl.Insert(InitScopeToken(&newScopeObj), start)
	} else {
		err := errors.New(fmt.Sprintf("Invalid (out of bounds) index provided to token list SCOPIFY RANGE: %d, %d", start, end))
		util.Error(err.Error(), err)
	}
}

// ConvertToArray
// This will convert the scope object to an array
// of all of its tokens. This reduces the tree type structure
// of the scopes into a single array. This does eliminate
// The purpose of this method is to provide just the tokens alone without
// any extra data.
func (tl *ScopeObj) ConvertToArray() []*Token {
	rtnArray := make([]*Token, 0)
	initialIndex := 0
	endIndex := len(tl.tokenList)

	for i := 0; i < tl.GetNumberOfScopes(); i++ {
		endIndex = tl.scopeIndices[i]
		rtnArray = append(rtnArray, tl.tokenList[initialIndex:endIndex]...)
		rtnArray = append(rtnArray, tl.tokenList[endIndex].GetScopeToken().ConvertToArray()...) // Makes ecursive Call
		initialIndex = endIndex + 1
	}
	rtnArray = append(rtnArray, tl.tokenList[initialIndex:len(tl.tokenList)]...)

	return rtnArray
}

func (tl *ScopeObj) PrintSimpleText() {
	for i := 0; i < tl.Size(); i++ {
		fmt.Print(tl.tokenList[i].SymbolicName + " ")
	}
	fmt.Println()
}
