package tokenizer

import (
	"fmt"
	"tp/src/util"
)

// PrintCharIndices
// Prints out all the characters and their indices of the text for the sake of debugging
func (tkzr *Tokenizer) PrintCharIndices() {
	for i := 0; i < tkzr.TextSize(); i++ {
		char := tkzr.GetChar(i)
		fmt.Printf("%d\t:\t%d\n", i, int(char))
	}
}

// PrintCurrentInfo
// Prints out the temporary information of the tokenizer's current char for the sake of debugging
func (tkzr *Tokenizer) PrintCurrentInfo(args ...string) {
	optionalInfo := ""
	for _, arg := range args {
		optionalInfo += arg
	}
	fmt.Println("------------------------------------------------------")
	fmt.Printf("CURRENT INFO (%s): ", optionalInfo)
	util.PrintTime()
	fmt.Println()
	fmt.Println("\tGENERAL: ")
	fmt.Printf("\t\tIndex: %d\tAsciiOfIndex: %d\tCharOfIndex: '%s'\n", tkzr.currentIndex, int(tkzr.CurrentChar()), string(tkzr.CurrentChar()))
	fmt.Printf("\t\tLine Num: %d\t Tab Level: %d\n", tkzr.currentLineNumber, tkzr.currentTabLevel)
	fmt.Println("\tBOOLS:")
	fmt.Printf("\t\ttempIgnoreChangesFromIncrement %t\tskipIncrement: %t\n", tkzr.tempIgnoreChangesFromIncrement, tkzr.skipIncrement)
	fmt.Println("\tFUNCTION INFO: ")
	fmt.Printf("\t\tStartInfo: %s\tEndInfo:%s\tFunctionSharedInfo: %s\n", tkzr.StartInfo, tkzr.EndInfo, tkzr.FunctionSharedInfo)
	fmt.Println("\tCURRENT SCOPE:")
	fmt.Printf("\t\tSize: %d\tType: %s\n", tkzr.currentScope.Size(), tkzr.currentScope.GetType())
	lastFewTokens := ""
	numberOfPreviousTokens := 5
	for i := tkzr.currentScope.Size() - 1; i > 0 && i > tkzr.currentScope.Size()-1-numberOfPreviousTokens; i-- {
		tkn, err := tkzr.currentScope.At(i)
		if err != nil {
			continue
		}
		lastFewTokens += tkn.ToString() + "\t"
	}
	fmt.Printf("\t\tLast few tokens: %s\n", lastFewTokens)
	fmt.Println("------------------------------------------------------")
}
