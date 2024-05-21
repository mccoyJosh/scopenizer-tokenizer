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
		sc.fixScopeIndices()
	}

	return sc
}

// InitScopeToken
// This will initialize a token object which is of scope type
// (i.e. a scope token) and returns a pointer to it.
//
// This method takes an optional number of scope objects as a parameter.
// If no scopes are provided, it will create an empty scope token.
// If scopes are provided, it will combine them all together into this token's scope object
func InitScopeToken(scopes ...*ScopeObj) *Token {
	var providedScope ScopeObj

	if len(scopes) == 1 {
		providedScope = *scopes[0]
	} else {
		providedScope = InitScope()
		if len(scopes) > 0 {
			for _, scope := range scopes {
				providedScope.Concatenate(scope)
			}
		}
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
func (so *ScopeObj) Size() int {
	return so.size
}

// GetType
// Returns the scope object's type
func (so *ScopeObj) GetType() string {
	return so.scopeType
}

// SetType
// Given a string parameter, this will set the scope's type
func (so *ScopeObj) SetType(typeString string) {
	so.scopeType = typeString
}

// GetNumberOfScopes
// Returns the number of scopes found directly in this scope obj
func (so *ScopeObj) GetNumberOfScopes() int {
	return len(so.scopeIndices)
}

// Concatenate
// Takes a scope object as a parameters
// Will add all the tokens from the parameter scope object to
// THIS scope object
func (so *ScopeObj) Concatenate(additionalScope *ScopeObj) {
	for _, token := range additionalScope.tokenList {
		so.Push(token)
	}
}

// GetScope
// Returns the i-th scope in this scope object
// index: the index of the scope token in the scope indices array. If you use 0 as index, it would be requesting for the first scope
//
// This will return an error if the provided index is out of bounds
func (so *ScopeObj) GetScope(index int) (*ScopeObj, error) {
	if index < so.GetNumberOfScopes() {
		return so.tokenList[so.scopeIndices[index]].scopeToken, nil
	}
	err := errors.New(fmt.Sprintf("Invalid (out of bounds) index provided to scope indicies list GETSCOPE: %d", index))
	return nil, err
}

// TotalSize
// This returns the number of tokens in this and all inner scopes.
// This size DOES include the number of tokens in inner scopes.
// This does not include the number of scopes.
func (so *ScopeObj) TotalSize() int {
	return so.setTotalSize()
}

// setTotalSize
// This will go through the scopes found within this list
// and get the total number of tokens found within.
// This should get the total number of tokens within this and all inner scopes.
func (so *ScopeObj) setTotalSize() int {
	// If there are no inner scopes, we return
	if len(so.scopeIndices) <= 0 {
		return so.size
	}

	innerTotals := 0
	for _, index := range so.scopeIndices {
		token := so.tokenList[index] //
		if token.ValidScopeToken() {
			scopeToken := token.scopeToken
			innerTotals += scopeToken.setTotalSize()
		} else {
			err := errors.New(fmt.Sprintf("Non-Scope Token Found At Supposed Scope Index: %d", index))
			util.Error(err.Error(), err)
		}
	}

	// Removes the scope tokens for the count, as we only want 'real' tokens
	total := innerTotals + so.size - len(so.scopeIndices)
	return total
}

// fixScopeIndices
// This will go through the tokens and find the indices
// of all the scope tokens and add them to the scopeIndices array
//
// This will go through the entirety of the token list
// everytime this is called to ensure all scopes are accounted for.
// If there is a significant number of tokens,
// this method may become an inefficiency
func (so *ScopeObj) fixScopeIndices() {
	so.scopeIndices = make([]int, 0)
	for index, token := range so.tokenList {
		if token.ValidScopeToken() {
			so.scopeIndices = append(so.scopeIndices, index)
		}
	}
}

// Push
// This adds a token to the token list at the end of the list,
// much like one would push an item to the top of a stack.
func (so *ScopeObj) Push(tt *Token) {
	_ = so.Insert(tt, so.size)
}

// Insert
// Inserts a token into the token list. Returns nothing.
// If this token is not inserted at the end of the list,
// it may need to fix the scope index list
//
// Will return an error if the provided index is not possible to be inserted (negative value or > Size())
func (so *ScopeObj) Insert(tt *Token, index int) error {
	if index >= 0 && index <= so.size {
		if so.size == index {
			so.tokenList = append(so.tokenList, tt)
		} else {
			so.tokenList = append(so.tokenList[:index+1], so.tokenList[index:]...)
			so.tokenList[index] = tt
		}
		so.size++

		if index != so.size-1 {
			so.fixScopeIndices()
		} else if tt.ValidScopeToken() {
			so.scopeIndices = append(so.scopeIndices, index)
		}
		return nil
	} else {
		err := errors.New(fmt.Sprintf("Invalid (out of bounds) index provided to token list INSERT: %d", index))
		return err
	}
}

// At
// This will return the pointer to the token at the given index.
//
// If the provided index is out of range, it produces an error and returns nil
func (so *ScopeObj) At(index int) (*Token, error) {
	if index >= 0 && index < so.size {
		return so.tokenList[index], nil
	}
	err := errors.New(fmt.Sprintf("Invalid (out of bounds) index provided to token list GET INDEX: %d", index))
	return nil, err
}

// Pop
// This will get and remove the item at the end of the token list,
// much like one would pop the top item off the top of a stack
//
// Returns nil if the scope has no tokens
func (so *ScopeObj) Pop() *Token {
	if so.Size() > 0 {
		i := so.Front()
		err := so.Delete(so.size - 1)
		if err != nil {
			return nil
		}
		return i
	}
	return nil
}

// Front
// Returns the items at the front of the list of tokens,
// i.e. the value at the index of size - 1
//
// Return nil if the scope has no tokens
func (so *ScopeObj) Front() *Token {
	if so.Size() > 0 {
		token, _ := so.At(so.size - 1)
		return token
	}
	return nil
}

// Delete
// This removes a token from the token list given its index.
//
// If this index is out of bounds, an error is returned
func (so *ScopeObj) Delete(index int) error {
	if index >= 0 && index < so.size {
		so.tokenList = append(so.tokenList[:index], so.tokenList[index+1:]...)
		so.size--
		so.fixScopeIndices()
		return nil
	} else {
		err := errors.New(fmt.Sprintf("Invalid (out of bounds) index provided to token list DELETE: %d", index))
		return err
	}
}

// ScopifyRange
// This method, given a start index and an end index, will
// add the given range of tokens to a scope token and place the new
// scope token at the start index. The tokens in the range will only
// be in the new scope token and will be removed from the token it is called from.
//
// The start and end indices are INCLUSIVE meaning all values between the two indices,
// including the values at the indices, are included in the scope token
//
// Returns an error if the start and/or end indices are out of range
func (so *ScopeObj) ScopifyRange(start int, end int) error {
	if start <= end && start >= 0 && start < so.size && end >= 0 && end < so.size {
		tokensSubset := so.tokenList[start : end+1]
		tokensSubsetLength := end - start + 1
		newScopeObj := InitScope(tokensSubset)
		for i := 0; i < tokensSubsetLength; i++ {
			_ = so.Delete(start)
		}
		_ = so.Insert(InitScopeToken(&newScopeObj), start)
		return nil
	} else {
		err := errors.New(fmt.Sprintf("Invalid (out of bounds) index provided to token list SCOPIFY RANGE: %d, %d", start, end))
		return err
	}
}

// ConvertToArray
// This will convert the scope object to an array
// of all of its tokens. This reduces the tree type structure
// of the scopes into a single array.
// The purpose of this method is to provide just the tokens alone without
// any extra data.
func (so *ScopeObj) ConvertToArray() []*Token {
	rtnArray := make([]*Token, 0)
	initialIndex := 0
	endIndex := len(so.tokenList)

	for i := 0; i < so.GetNumberOfScopes(); i++ {
		endIndex = so.scopeIndices[i]
		rtnArray = append(rtnArray, so.tokenList[initialIndex:endIndex]...)
		rtnArray = append(rtnArray, so.tokenList[endIndex].GetScopeToken().ConvertToArray()...) // Makes ecursive Call
		initialIndex = endIndex + 1
	}
	rtnArray = append(rtnArray, so.tokenList[initialIndex:len(so.tokenList)]...)

	return rtnArray
}

// PrintSymbolicNames
// This will print out the symbolic names of all the tokens
// in this scope object on a single line.
//
// Created and used for testing purposes
func (so *ScopeObj) PrintSymbolicNames() {
	for i := 0; i < so.Size(); i++ {
		fmt.Print(so.tokenList[i].SymbolicName + " ")
	}
	fmt.Println()
}

// PrintTexts
// This will print out the text of all the tokens
// in this scope object on a single line
//
// Created and used for testing purposes
func (so *ScopeObj) PrintTexts() {
	for i := 0; i < so.Size(); i++ {
		fmt.Print(so.tokenList[i].Text + " ")
	}
	fmt.Println()
}
