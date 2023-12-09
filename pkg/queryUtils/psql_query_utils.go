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

func QueryCase(startIndex int, columnName, caseCondition string, caseConditionValues []string) string {
	lines := []string{}
	lines = append(lines, fmt.Sprintf("%s = CASE %s", columnName, caseCondition))
	for i, caseCondition := range caseConditionValues {
		lines = append(lines, fmt.Sprintf("WHEN %s THEN $%d", caseCondition, startIndex+i))
	}
	lines = append(lines, fmt.Sprintf("ELSE %s", columnName))
	lines = append(lines, "END")

	return strings.Join(lines, "\n")
}

func QueryCases(startIndex int, columnNames, caseConditionValues []string, caseCondition string) string {
	blocks := []string{}
	for i, columnName := range columnNames {
		blocks = append(blocks, QueryCase(i*len(caseConditionValues)+startIndex, columnName, caseCondition, caseConditionValues))
	}
	return strings.Join(blocks, ",\n")
}

func GetColumnNames(columnNames []string, except ...string) []string {
	res := []string{}
	for _, name := range columnNames {
		excluded := false
		for _, exceptName := range except {
			if exceptName == name {
				excluded = true
				break
			}
		}
		if !excluded {
			res = append(res, name)
		}
	}
	return res
}

// returns slice of psql place holders according to given parametres. For example:
//
// 	GetPlaceholdersSliceFrom(1, 2) -> []string{"$1", "$2"}
// 	GetPlaceholdersSliceFrom(3, 4) -> []string{"$3", "$4", "$5", "$6"}
func GetPlaceholdersSliceFrom(indexFrom, placeHoldersNum int) []string {
	res := make([]string, placeHoldersNum)
	
	for i := range res {
		res[i] = fmt.Sprintf("$%d", indexFrom + i)
	}
	return res
}
