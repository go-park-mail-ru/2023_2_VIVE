package queryUtils

import (
	"fmt"
	"strings"
)

// Retruns []any that were converted from []int
func IntToAnySlice(ints []int) *[]any {
	sliceToReturn := make([]any, len(ints))
	for i, val := range ints {
		sliceToReturn[i] = val
	}
	return &sliceToReturn
}

// Returns string containing query placeholders separated by comma
//
// Example: []int{3, 8, 9} -> "$1, $2, $3"
func QueryPlaceHolders(startIndex int, elementsNum int) string {
	queryPlaceHolders := make([]string, elementsNum)
	for i := 0; i < elementsNum; i++ {
		queryPlaceHolders[i] = fmt.Sprintf("$%d", i+startIndex)
	}
	return strings.Join(queryPlaceHolders, ", ")
}

func QueryPlaceHoldersMultipleRows(startIndex, elementsPerRow, rowsNum int) string {
	rows := make([]string, rowsNum)
	// currPHIdx := startIndex

	for i := 0; i < rowsNum; i++ {
		rows[i] = fmt.Sprintf("(%s)", QueryPlaceHolders(startIndex+i*elementsPerRow, elementsPerRow))
	}
	return strings.Join(rows, ", ")
}
