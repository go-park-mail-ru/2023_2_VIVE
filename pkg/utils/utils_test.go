package utils

import (
	"reflect"
	"testing"
)

var testContainsCases = []struct {
	inputElem  int
	inputElems []int
	expected   bool
}{
	{
		inputElem:  1,
		inputElems: []int{1, 2, 3},
		expected:   true,
	},
	{
		inputElem:  24,
		inputElems: []int{1, 2, 3},
		expected:   false,
	},
	{
		inputElem:  1,
		inputElems: []int{},
		expected:   false,
	},
	{
		inputElem:  0,
		inputElems: []int{0, 0, 1, 1},
		expected:   true,
	},
}

func TestContains(t *testing.T) {
	for _, testCase := range testContainsCases {
		actual := Contains(testCase.inputElem, testCase.inputElems)
		if actual != testCase.expected {
			t.Errorf("wrong answer:\n\texpected: %v\n\tactual: %v\n", testCase.expected, actual)
		}
	}
}

var testDifferenceCases = []struct {
	inputSlice1 []int
	inputSlice2 []int
	expected    []int
}{
	{
		inputSlice1: []int{1, 2, 3},
		inputSlice2: []int{1, 2},
		expected:    []int{3},
	},
	{
		inputSlice1: []int{1, 2},
		inputSlice2: []int{1, 2},
		expected:    []int{},
	},
	{
		inputSlice1: []int{1, 2, 3},
		inputSlice2: []int{},
		expected:    []int{1, 2, 3},
	},
	{
		inputSlice1: []int{},
		inputSlice2: []int{1, 2, 3},
		expected:    []int{},
	},
}

func TestDifference(t *testing.T) {
	for _, testCase := range testDifferenceCases {
		actual := Difference(testCase.inputSlice1, testCase.inputSlice2)
		if !reflect.DeepEqual(actual, testCase.expected) {
			t.Errorf("wrong answer:\n\texpected: %v\n\tactual: %v\n", testCase.expected, actual)
		}
	}
}
