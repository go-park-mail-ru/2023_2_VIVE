package searchEngineUtils

import (
	morphanalyzer "HnH/services/searchEngineService/pkg/morphAnalyzer"
	"reflect"
	"testing"

	"github.com/vbatushev/morph"
)

var testGetWordsCases = []struct {
	input    string
	expected []string
}{
	{
		input:    "hello world",
		expected: []string{"hello", "world"},
	},
	{
		input:    "HellO WoRld",
		expected: []string{"hello", "world"},
	},
	{
		input:    "",
		expected: []string{},
	},
	{
		input:    "hello   world",
		expected: []string{"hello", "world"},
	},
	{
		input:    "    hello   world    ",
		expected: []string{"hello", "world"},
	},
	{
		input:    "        ",
		expected: []string{},
	},
	{
		input:    "Привет Мир",
		expected: []string{"привет", "мир"},
	},
	{
		input:    "привет, мир",
		expected: []string{"привет", "мир"},
	},
	{
		input:    "привет, мир1",
		expected: []string{"привет", "мир"},
	},
	{
		input:    "привет, мир1 1145",
		expected: []string{"привет", "мир"},
	},
	{
		input:    "привет, мир!",
		expected: []string{"привет", "мир"},
	},
	{
		input:    "привет.,;!\"№;%:?*()_+=-@#$^'[]{}<>`~ мир",
		expected: []string{"привет", "мир"},
	},
}

func TestGetWords(t *testing.T) {
	for _, testCase := range testGetWordsCases {
		actual := getWords(testCase.input)
		if !reflect.DeepEqual(actual, testCase.expected) {
			t.Errorf("wrong result,\n\tactual: %v\n\texpected: %v\n", actual, testCase.expected)
		}
	}
}

var testMorphWordsCases = []struct {
	input    []string
	expected []string
}{
	{
		input:    []string{"все"},
		expected: []string{"всё", "весь"},
	},
	{
		input:    []string{"поиска"},
		expected: []string{"поиск"},
	},
	{
		input:    []string{"криком"},
		expected: []string{"крик"},
	},
	{
		input:    []string{"потому", "что"},
		expected: []string{"потому", "что"},
	},
	{
		input:    []string{"VK"},
		expected: []string{"VK"},
	},
}

func TestMorphWords(t *testing.T) {
	err := morphanalyzer.InitMorphAnalyzer()
	if err != nil && err != morph.ErrAlreadyInitialized {
		t.Errorf("could not initialize morph analyzer: %s\n", err)
	}

	for _, testCase := range testMorphWordsCases {
		actual := morphWords(testCase.input)
		if !reflect.DeepEqual(actual, testCase.expected) {
			t.Errorf("wrong result,\n\tactual: %v\n\texpected: %v\n", actual, testCase.expected)
		}
	}
}

var testGetNormsCases = []struct {
	input    string
	expected []string
}{
	{
		input:    "ищу работника на постоянную работу",
		expected: []string{"искать", "работник", "на", "постоянный", "работа"},
	},
	{
		input:    "   ищу    работника    на   постоянную   работу   ",
		expected: []string{"искать", "работник", "на", "постоянный", "работа"},
	},
	{
		input:    "Ищу работникА на постоянную РАБОТУ",
		expected: []string{"искать", "работник", "на", "постоянный", "работа"},
	},
}

func TestGetNorms(t *testing.T) {
	err := morphanalyzer.InitMorphAnalyzer()
	if err != nil && err != morph.ErrAlreadyInitialized {
		t.Errorf("could not initialize morph analyzer: %s\n", err)
	}

	for _, testCase := range testGetNormsCases {
		actual := getNorms(testCase.input)
		if !reflect.DeepEqual(actual, testCase.expected) {
			t.Errorf("wrong result,\n\tactual: %v\n\texpected: %v\n", actual, testCase.expected)
		}
	}
}

var testParseSearchQueryCases = []struct {
	input    string
	expected []string
}{
	{
		input:    "ищу работника на постоянную работу",
		expected: []string{"искать", "работник", "постоянный", "работа"},
	},
	{
		input:    "   ищу    работника    на   постоянную   работу   ",
		expected: []string{"искать", "работник", "постоянный", "работа"},
	},
	{
		input:    "Ищу работникА на постоянную РАБОТУ",
		expected: []string{"искать", "работник", "постоянный", "работа"},
	},
	{
		input:    "ищу работника на постоянную работу, потому что большая текучка",
		expected: []string{"искать", "работник", "постоянный", "работа", "больший", "большой", "текучка"},
	},
}

func TestParseSearchQuery(t *testing.T) {
	err := morphanalyzer.InitMorphAnalyzer()
	if err != nil && err != morph.ErrAlreadyInitialized {
		t.Errorf("could not initialize morph analyzer: %s\n", err)
	}

	for _, testCase := range testParseSearchQueryCases {
		actual := ParseSearchQuery(testCase.input)
		if !reflect.DeepEqual(actual, testCase.expected) {
			t.Errorf("wrong result,\n\tactual: %v\n\texpected: %v\n", actual, testCase.expected)
		}
	}
}
