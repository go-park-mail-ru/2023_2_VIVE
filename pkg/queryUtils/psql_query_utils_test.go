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
