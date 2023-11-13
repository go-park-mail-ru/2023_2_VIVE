package queryUtils

import (
	"reflect"
	"testing"
)

var testIntToAnySliceCases = []struct {
	input    []int
	expected []any
}{
	{
		input:    []int{1, 2, 3},
		expected: []any{1, 2, 3},
	},
	{
		input:    []int{1},
		expected: []any{1},
	},
	{
		input:    []int{},
		expected: []any{},
	},
	{
		input:    []int{0},
		expected: []any{0},
	},
	{
		input:    []int{0, 0, 0, 0},
		expected: []any{0, 0, 0, 0},
	},
}

func TestIntToAnySlice(t *testing.T) {
	for _, testCase := range testIntToAnySliceCases {
		actual := IntToAnySlice(testCase.input)
		if !reflect.DeepEqual(*actual, testCase.expected) {
			t.Errorf("Two slices must be equal\nExpected: %v\nActual: %v", testCase.expected, actual)
		}
	}
}

var testQueryPlaceHoldersCases = []struct {
	inputStartIndex int
	inputElemNum    int
	expected        string
}{
	{
		inputStartIndex: 1,
		inputElemNum:    3,
		expected:        "$1, $2, $3",
	},
	{
		inputStartIndex: 1,
		inputElemNum:    4,
		expected:        "$1, $2, $3, $4",
	},
	{
		inputStartIndex: 1,
		inputElemNum:    1,
		expected:        "$1",
	},
	{
		inputStartIndex: 1,
		inputElemNum:    0,
		expected:        "",
	},
	{
		inputStartIndex: 3,
		inputElemNum:    4,
		expected:        "$3, $4, $5, $6",
	},
}

func TestQueryPlaceHolders(t *testing.T) {
	for _, testCase := range testQueryPlaceHoldersCases {
		actual := QueryPlaceHolders(testCase.inputStartIndex, testCase.inputElemNum)
		if !reflect.DeepEqual(actual, testCase.expected) {
			t.Errorf("Two strings must be equal\nExpected: %v\nActual: %v", testCase.expected, actual)
		}
	}
}

var testQueryPlaceHoldersMultipleRows = []struct {
	inputStartIndex int
	inputElemPerRow int
	inputRowNum     int
	expected        string
}{
	{
		inputStartIndex: 1,
		inputElemPerRow: 2,
		inputRowNum:     2,
		expected:        "($1, $2), ($3, $4)",
	},
	{
		inputStartIndex: 1,
		inputElemPerRow: 3,
		inputRowNum:     3,
		expected:        "($1, $2, $3), ($4, $5, $6), ($7, $8, $9)",
	},
	{
		inputStartIndex: 3,
		inputElemPerRow: 2,
		inputRowNum:     3,
		expected:        "($3, $4), ($5, $6), ($7, $8)",
	},
	{
		inputStartIndex: 1,
		inputElemPerRow: 0,
		inputRowNum:     3,
		expected:        "(), (), ()",
	},
	{
		inputStartIndex: 1,
		inputElemPerRow: 4,
		inputRowNum:     0,
		expected:        "",
	},
	{
		inputStartIndex: 1,
		inputElemPerRow: 0,
		inputRowNum:     0,
		expected:        "",
	},
}

func TestQueryPlaceHoldersMultipleRows(t *testing.T) {
	for _, testCase := range testQueryPlaceHoldersMultipleRows {
		actual := QueryPlaceHoldersMultipleRows(testCase.inputStartIndex, testCase.inputElemPerRow, testCase.inputRowNum)
		if !reflect.DeepEqual(actual, testCase.expected) {
			t.Errorf("Two strings must be equal\nExpected: %v\nActual: %v", testCase.expected, actual)
		}
	}
}

var testQueryCase = []struct {
	inputStartIndex          int
	inputColumnName          string
	inputCaseCondition       string
	inputCaseConditionValues []string
	expected                 string
}{
	{
		inputStartIndex:          1,
		inputColumnName:          "name",
		inputCaseCondition:       "id",
		inputCaseConditionValues: []string{"1", "2", "3"},
		expected:                 "name = CASE id\nWHEN 1 THEN $1\nWHEN 2 THEN $2\nWHEN 3 THEN $3\nELSE name\nEND",
	},
	{
		inputStartIndex:          1,
		inputColumnName:          "name",
		inputCaseCondition:       "id",
		inputCaseConditionValues: []string{"1", "2"},
		expected:                 "name = CASE id\nWHEN 1 THEN $1\nWHEN 2 THEN $2\nELSE name\nEND",
	},
	{
		inputStartIndex:          2,
		inputColumnName:          "name",
		inputCaseCondition:       "id",
		inputCaseConditionValues: []string{"1", "2"},
		expected:                 "name = CASE id\nWHEN 1 THEN $2\nWHEN 2 THEN $3\nELSE name\nEND",
	},
}

func TestQueryCase(t *testing.T) {
	for _, testCase := range testQueryCase {
		actual := QueryCase(
			testCase.inputStartIndex,
			testCase.inputColumnName,
			testCase.inputCaseCondition,
			testCase.inputCaseConditionValues,
		)
		if !reflect.DeepEqual(actual, testCase.expected) {
			t.Errorf("Two strings must be equal\nExpected: %v\nActual: %v", testCase.expected, actual)
		}
	}
}

var testQueryCases = []struct {
	inputStartIndex          int
	inputColumnNames         []string
	inputCaseConditionValues []string
	inputCaseCondition       string
	expected                 string
}{
	{
		inputStartIndex:          1,
		inputColumnNames:         []string{"name1", "name2", "name3"},
		inputCaseCondition:       "id",
		inputCaseConditionValues: []string{"1", "2", "3"},
		expected: `name1 = CASE id
WHEN 1 THEN $1
WHEN 2 THEN $2
WHEN 3 THEN $3
ELSE name1
END,
name2 = CASE id
WHEN 1 THEN $4
WHEN 2 THEN $5
WHEN 3 THEN $6
ELSE name2
END,
name3 = CASE id
WHEN 1 THEN $7
WHEN 2 THEN $8
WHEN 3 THEN $9
ELSE name3
END`,
	},
	{
		inputStartIndex:          1,
		inputColumnNames:         []string{"name1", "name2"},
		inputCaseCondition:       "id",
		inputCaseConditionValues: []string{"1", "2", "3"},
		expected: `name1 = CASE id
WHEN 1 THEN $1
WHEN 2 THEN $2
WHEN 3 THEN $3
ELSE name1
END,
name2 = CASE id
WHEN 1 THEN $4
WHEN 2 THEN $5
WHEN 3 THEN $6
ELSE name2
END`,
	},
	{
		inputStartIndex:          5,
		inputColumnNames:         []string{"name1", "name2"},
		inputCaseCondition:       "id",
		inputCaseConditionValues: []string{"1", "2", "3"},
		expected: `name1 = CASE id
WHEN 1 THEN $5
WHEN 2 THEN $6
WHEN 3 THEN $7
ELSE name1
END,
name2 = CASE id
WHEN 1 THEN $8
WHEN 2 THEN $9
WHEN 3 THEN $10
ELSE name2
END`,
	},
}

func TestQueryCases(t *testing.T) {
	for _, testCase := range testQueryCases {
		actual := QueryCases(
			testCase.inputStartIndex,
			testCase.inputColumnNames,
			testCase.inputCaseConditionValues,
			testCase.inputCaseCondition,
		)
		if !reflect.DeepEqual(actual, testCase.expected) {
			t.Errorf("Two strings must be equal\nExpected: %v\nActual: %v", testCase.expected, actual)
		}
	}
}

var testGetColumnNames = []struct {
	inputColumnNames []string
	inputExcept      []string
	expected         []string
}{
	{
		inputColumnNames: []string{"name1", "name2", "name3"},
		inputExcept:      []string{"name1"},
		expected:         []string{"name2", "name3"},
	},
	{
		inputColumnNames: []string{"name1", `"name2"`, "name3"},
		inputExcept:      []string{},
		expected:         []string{"name1", `"name2"`, "name3"},
	},
}

func TestGetColumnNames(t *testing.T) {
	for _, testCase := range testGetColumnNames {
		actual := GetColumnNames(testCase.inputColumnNames, testCase.inputExcept...)
		if !reflect.DeepEqual(actual, testCase.expected) {
			t.Errorf("Two strings must be equal\nExpected: %v\nActual: %v", testCase.expected, actual)
		}
	}
}
