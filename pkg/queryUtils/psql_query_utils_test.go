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
	input    []int
	expected string
}{
	{
		input:    []int{1, 2, 3},
		expected: "$1, $2, $3",
	},
	{
		input:    []int{5, 6, 1, 8},
		expected: "$1, $2, $3, $4",
	},
	{
		input:    []int{5},
		expected: "$1",
	},
	{
		input:    []int{},
		expected: "",
	},
}

func TestQueryPlaceHolders(t *testing.T) {
	for _, testCase := range testQueryPlaceHoldersCases {
		args := *IntToAnySlice(testCase.input)
		actual := QueryPlaceHolders(args...)
		if !reflect.DeepEqual(actual, testCase.expected) {
			t.Errorf("Two strings must be equal\nExpected: %v\nActual: %v", testCase.expected, actual)
		}
	}
}
