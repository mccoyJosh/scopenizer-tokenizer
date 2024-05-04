package util

import (
	"encoding/json"
	"fmt"
	"strconv"
)

func FormatStringToSize(s string, size int) string {
	originalLength := len(s)
	if originalLength > size {
		return s[0:size-3] + "..."
	}
	spaces := ""
	neededSpaces := size - originalLength
	for i := 0; i < neededSpaces; i++ {
		spaces += " "
	}
	return s + spaces
}

func GetPercentage(numerator int, denominator int, round int) float64 {
	percentage := 0.0
	if denominator > 0 {
		percentage = float64(numerator) / float64(denominator)
		roundFormatting := "%." + fmt.Sprintf("%d", round+2) + "f"
		fmtString := fmt.Sprintf(roundFormatting, percentage)
		pP, _ := strconv.ParseFloat(fmtString, 64)
		percentage = pP * 100.0
	}
	return percentage
}

func ToJSONString(result any) string {
	b, err := json.MarshalIndent(result, "", "\t")
	if err != nil {
		return ""
	}
	return string(b)
}

func HideButAFewCharsOfInformation(str string, reveal int) string {
	stars := ""
	for i := 0; i < len(str)-reveal; i++ {
		stars += "*"
	}
	lastBit := str[len(str)-reveal : len(str)]
	return stars + lastBit
}
