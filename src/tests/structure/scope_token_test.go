package structure

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"testing"
	"tp/src/tokenizer/tokens"
)

// CreateTestTokenArray
// This method will create a token array for the sake of testing.
// This will generate random tokens.
//
// idTest: this will be used for all the tokens' text value. This helps identify them when distinguishing between two different generated test token arrays
//
// size: this is the number of tokens which will be in the token array. This includes scope tokens
//
// numOfScopes: this the number of scope tokens which will be placed evenly throughout this test token array. Will not add any if < 0
//
// recursiveSteps (ONLY MATTERS IF YOU HAVE SCOPE TOKENS): This will determine the number of times it will attempt to recursively add scopes to the scopes within the returning tokens array
//
// recursiveStepNumOfScopeDecrements (ONLY MATTERS IF YOU HAVE SCOPE TOKENS): This will determine the decrement value of numOfScopes per recursive call. For example, if numOfScope is 3 and this parameter is set to 1, the first set of internal scopes will all have 2 scopes found within them.
func CreateTestTokenArray(idTest string, size int, numOfScopes int, recursiveSteps int, recursiveStepNumOfScopeDecrements int) []*tokens.Token {
	tokenArray := make([]*tokens.Token, 0)

	scopePosition := 0
	if numOfScopes > 0 && numOfScopes <= size {
		scopePosition = size/numOfScopes - 1
	} else {
		scopePosition = -2
	}
	scopeOffset := scopePosition + 1

	for i := 0; i < size; i++ {
		var tempToken tokens.Token

		if i == scopePosition { // Time to place a scope token
			var tempTokenArray []*tokens.Token
			if recursiveSteps > 0 {
				tempTokenArray = CreateTestTokenArray(idTest, size, numOfScopes-recursiveStepNumOfScopeDecrements, recursiveSteps-1, recursiveStepNumOfScopeDecrements)
			} else {
				tempTokenArray = CreateTestTokenArray(idTest, size, 0, 0, 0)
			}

			tempScope := tokens.InitScope(tempTokenArray)
			tempToken = *tokens.InitScopeToken(&tempScope)
			tempToken.Text = idTest
			scopePosition += scopeOffset
		} else {
			tempToken = tokens.CreateUnidentifiedToken(idTest, i, i*2, i*3)
			tempRuleName := "SYMBOL"
			tempSymbolicName := "PERIOD"
			rndNum := rand.Intn(100)
			if rndNum%2 == 0 { // Will be keyword
				tempRuleName = "KEYWORD"
				tempSymbolicName = "IF"
				if rndNum%3 == 0 {
					tempSymbolicName = "ELSE"
				}
			} else { // Will be symbol
				if rndNum%3 == 0 {
					tempSymbolicName = "BRACKET"
				}
			}
			tempToken.SetValues(tempRuleName, tempSymbolicName)
		}

		tokenArray = append(tokenArray, &tempToken)
	}

	return tokenArray
}

func Test_InitScope_Empty(t *testing.T) {
	exampleScope := tokens.InitScope()
	assert.Equal(t, 0, exampleScope.Size())
	assert.Equal(t, tokens.UNKNOWN_SCOPE_STRING, exampleScope.GetType())
}

func Test_InitScope_Start_One_List(t *testing.T) {
	tokenArray := CreateTestTokenArray("testToken", 10, 0, 0, 0)
	exampleScope := tokens.InitScope(tokenArray)
	assert.Equal(t, 10, exampleScope.Size())
	assert.Equal(t, 0, exampleScope.GetNumberOfScopes())
}

func Test_InitScope_Start_Two_List(t *testing.T) {
	tokenArray1 := CreateTestTokenArray("tkList1", 10, 0, 0, 0)
	tokenArray2 := CreateTestTokenArray("tkList2", 20, 0, 0, 0)
	exampleScope := tokens.InitScope(tokenArray1, tokenArray2)
	assert.Equal(t, 30, exampleScope.Size())
	assert.Equal(t, 0, exampleScope.GetNumberOfScopes())
}

func Test_InitScope_Start_One_List_With_Scope(t *testing.T) {
	tokenArray := CreateTestTokenArray("testToken", 10, 1, 0, 1)
	exampleScope := tokens.InitScope(tokenArray)
	assert.Equal(t, 10, exampleScope.Size())
	assert.Equal(t, 1, exampleScope.GetNumberOfScopes())
	assert.Equal(t, 19, exampleScope.TotalSize())
}

func Test_InitScope_Start_One_List_With_Scope_Rec(t *testing.T) {
	tokenArray := CreateTestTokenArray("testToken", 10, 2, 1, 1)

	exampleScope := tokens.InitScope(tokenArray)
	assert.Equal(t, 10, exampleScope.Size())
	assert.Equal(t, 2, exampleScope.GetNumberOfScopes())
	assert.Equal(t, 46, exampleScope.TotalSize())
}

func Test_InitScope_Start_Two_Lists_With_Scopes_Rec(t *testing.T) {
	tokenArray := CreateTestTokenArray("testToken", 10, 2, 1, 1)
	tokenArray2 := CreateTestTokenArray("testToken", 10, 2, 1, 1)

	exampleScope := tokens.InitScope(tokenArray, tokenArray2)
	assert.Equal(t, 20, exampleScope.Size())
	assert.Equal(t, 4, exampleScope.GetNumberOfScopes())
	assert.Equal(t, 92, exampleScope.TotalSize())
}

func Test_SetType(t *testing.T) {
	exampleScope := tokens.InitScope()
	exampleScope.SetType("Method")
	assert.Equal(t, "Method", exampleScope.GetType())
}

func Test_Get_Scope(t *testing.T) {
	tokenArray := CreateTestTokenArray("testToken", 10, 2, 1, 1)
	exampleScope := tokens.InitScope(tokenArray)
	exampleInnerScope, err := exampleScope.GetScope(0)

	assert.Nil(t, err)
	assert.Equal(t, 10, exampleInnerScope.Size())
	assert.Equal(t, 1, exampleInnerScope.GetNumberOfScopes())
	assert.Equal(t, 19, exampleInnerScope.TotalSize())
}

func Test_Get_Scope_Out_Of_Bounds(t *testing.T) {
	tokenArray := CreateTestTokenArray("testToken", 10, 2, 1, 1)
	exampleScope := tokens.InitScope(tokenArray)
	exampleInnerScope, err := exampleScope.GetScope(10)

	assert.Error(t, err)
	var expectedResult *tokens.ScopeObj // is nil
	assert.Equal(t, expectedResult, exampleInnerScope)
}

func Test_Push(t *testing.T) {
	tokenArray := CreateTestTokenArray("testToken", 10, 0, 0, 1)
	exampleScope := tokens.InitScope(tokenArray)

	pushToken := tokens.CreateUnidentifiedToken("if", 1, 1, 1)
	pushToken.SetValues("KEYWORD", "IF")
	exampleScope.Push(&pushToken)

	token, err := exampleScope.At(10)
	assert.Nil(t, err)
	assert.Equal(t, &pushToken, token)
}

func Test_Insert(t *testing.T) {
	tokenArray := CreateTestTokenArray("testToken", 10, 0, 0, 1)
	exampleScope := tokens.InitScope(tokenArray)

	pushToken := tokens.CreateUnidentifiedToken("if", 1, 1, 1)
	pushToken.SetValues("KEYWORD", "IF")
	err := exampleScope.Insert(&pushToken, 0)

	assert.Nil(t, err)

	pushToken2 := tokens.CreateUnidentifiedToken("else", 1, 1, 1)
	pushToken2.SetValues("KEYWORD", "ELSE")
	err = exampleScope.Insert(&pushToken2, 0)

	assert.Nil(t, err)

	token1, err1 := exampleScope.At(0)
	token2, err2 := exampleScope.At(1)

	assert.Nil(t, err1)
	assert.Nil(t, err2)
	assert.Equal(t, &pushToken2, token1)
	assert.Equal(t, &pushToken, token2)
}

func Test_At_Nothing_Found(t *testing.T) {
	exampleScope := tokens.InitScope()
	var nullToken *tokens.Token

	token, err := exampleScope.At(0)
	assert.Error(t, err)
	assert.Equal(t, nullToken, token)
}

func Test_Pop(t *testing.T) {
	tokenArray := CreateTestTokenArray("testToken", 4, 0, 0, 1)
	exampleScope := tokens.InitScope(tokenArray)

	pushToken := tokens.CreateUnidentifiedToken("if", 1, 1, 1)
	pushToken.SetValues("KEYWORD", "IF")
	exampleScope.Push(&pushToken)

	receivedToken := exampleScope.Pop()

	assert.Equal(t, 4, exampleScope.Size())
	assert.Equal(t, pushToken, *receivedToken)
}

func Test_Pop_Nothing(t *testing.T) {
	exampleScope := tokens.InitScope()

	receivedToken := exampleScope.Pop()

	assert.Equal(t, 0, exampleScope.Size())
	assert.Nil(t, receivedToken)
}

func Test_Front(t *testing.T) {
	tokenArray := CreateTestTokenArray("testToken", 4, 0, 0, 1)
	exampleScope := tokens.InitScope(tokenArray)

	pushToken := tokens.CreateUnidentifiedToken("if", 1, 1, 1)
	pushToken.SetValues("KEYWORD", "IF")
	exampleScope.Push(&pushToken)

	receivedToken := exampleScope.Front()

	assert.Equal(t, 5, exampleScope.Size())
	assert.Equal(t, pushToken, *receivedToken)
}

func Test_Front_Nothing(t *testing.T) {
	exampleScope := tokens.InitScope()

	receivedToken := exampleScope.Front()

	assert.Equal(t, 0, exampleScope.Size())
	assert.Nil(t, receivedToken)
}

func Test_Delete_Front(t *testing.T) {
	exampleScope := tokens.InitScope()

	pushToken := tokens.CreateUnidentifiedToken("if", 1, 1, 1)
	pushToken.SetValues("KEYWORD", "IF")
	exampleScope.Push(&pushToken)

	pushToken2 := tokens.CreateUnidentifiedToken("else", 1, 1, 1)
	pushToken2.SetValues("KEYWORD", "ELSE")
	exampleScope.Push(&pushToken2)

	pushToken3 := tokens.CreateUnidentifiedToken("{", 1, 1, 1)
	pushToken3.SetValues("SYMBOL", "RBRACKET")
	exampleScope.Push(&pushToken3)

	assert.Equal(t, 3, exampleScope.Size())

	err := exampleScope.Delete(0)
	assert.Nil(t, err)

	tokenA, err1 := exampleScope.At(0)
	tokenB, err2 := exampleScope.At(1)

	assert.Nil(t, err1)
	assert.Nil(t, err2)

	assert.Equal(t, pushToken2, *tokenA)
	assert.Equal(t, pushToken3, *tokenB)
}

func Test_Delete_Middle(t *testing.T) {
	exampleScope := tokens.InitScope()

	pushToken := tokens.CreateUnidentifiedToken("if", 1, 1, 1)
	pushToken.SetValues("KEYWORD", "IF")
	exampleScope.Push(&pushToken)

	pushToken2 := tokens.CreateUnidentifiedToken("else", 1, 1, 1)
	pushToken2.SetValues("KEYWORD", "ELSE")
	exampleScope.Push(&pushToken2)

	pushToken3 := tokens.CreateUnidentifiedToken("{", 1, 1, 1)
	pushToken3.SetValues("SYMBOL", "RBRACKET")
	exampleScope.Push(&pushToken3)

	assert.Equal(t, 3, exampleScope.Size())

	err := exampleScope.Delete(1)
	assert.Nil(t, err)
	assert.Equal(t, 2, exampleScope.Size())

	tokenA, err1 := exampleScope.At(0)
	tokenB, err2 := exampleScope.At(1)

	assert.Nil(t, err1)
	assert.Nil(t, err2)

	assert.Equal(t, pushToken, *tokenA)
	assert.Equal(t, pushToken3, *tokenB)
}

func Test_Delete_End(t *testing.T) {
	exampleScope := tokens.InitScope()

	pushToken := tokens.CreateUnidentifiedToken("if", 1, 1, 1)
	pushToken.SetValues("KEYWORD", "IF")
	exampleScope.Push(&pushToken)

	pushToken2 := tokens.CreateUnidentifiedToken("else", 1, 1, 1)
	pushToken2.SetValues("KEYWORD", "ELSE")
	exampleScope.Push(&pushToken2)

	pushToken3 := tokens.CreateUnidentifiedToken("{", 1, 1, 1)
	pushToken3.SetValues("SYMBOL", "RBRACKET")
	exampleScope.Push(&pushToken3)

	assert.Equal(t, 3, exampleScope.Size())

	err := exampleScope.Delete(2)
	assert.Nil(t, err)
	assert.Equal(t, 2, exampleScope.Size())

	tokenA, err1 := exampleScope.At(0)
	tokenB, err2 := exampleScope.At(1)

	assert.Nil(t, err1)
	assert.Nil(t, err2)

	assert.Equal(t, pushToken, *tokenA)
	assert.Equal(t, pushToken2, *tokenB)
}

func Test_ScopifyRange_OneSection(t *testing.T) {
	tokenArray := CreateTestTokenArray("testToken", 10, 0, 0, 1)
	exampleScope := tokens.InitScope(tokenArray)
	err := exampleScope.ScopifyRange(0, 9)
	assert.Nil(t, err)
	assert.Equal(t, 1, exampleScope.Size())
	assert.Equal(t, 10, exampleScope.TotalSize())
	assert.Equal(t, 1, exampleScope.GetNumberOfScopes())

	scopifiedRangeFound, err := exampleScope.GetScope(0)
	assert.Nil(t, err)
	for i := 0; i < scopifiedRangeFound.Size(); i++ {
		tokenFound, err := scopifiedRangeFound.At(i)
		assert.Nil(t, err, fmt.Sprintf("Expected no error at index %d in Test_ScopifyRange_OneSection", i))
		assert.Equal(t, tokenArray[i], tokenFound, fmt.Sprintf("Expected same token at index %d in Test_ScopifyRange_OneSection", i))
	}
}

func Test_ConvertToArray(t *testing.T) {
	tokenArray := CreateTestTokenArray("testToken", 4, 4, 0, 1)
	exampleScope := tokens.InitScope(tokenArray)
	resultingArray := exampleScope.ConvertToArray()
	assert.Equal(t, 16, len(resultingArray))
}
