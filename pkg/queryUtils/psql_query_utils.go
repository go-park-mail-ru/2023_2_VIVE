package queryUtils

import (
	"fmt"
	"strings"
)

func IntToAnySlice(ints []int) *[]any {
	sliceToReturn := make([]any, len(ints))
	for i, val := range ints {
		sliceToReturn[i] = val
	}
	return &sliceToReturn
	// copy(sliceToReturn, ints)
}

// Returns string containing query placeholders separated by comma
//
// Example: []int{3, 8, 9} -> []any{3, 8, 9}, "$1, $2, $3"
func QueryPlaceHolders(values ...any) string {
	elementsNum := len(values)
	queryPlaceHolders := make([]string, elementsNum)
	for i := 0; i < elementsNum; i++ {
		queryPlaceHolders[i] = fmt.Sprintf("$%d", i+1)
	}
	return strings.Join(queryPlaceHolders, ", ")

	// copy(sliceToReturn, values)
	// for i, val := range values {
	// 	sliceToReturn[i] = val
	// }

}
