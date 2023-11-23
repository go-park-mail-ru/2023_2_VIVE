package searchEngineUtils

import (
	morphanalyzer "HnH/services/searchEngineService/pkg/morphAnalyzer"
	"strings"
	"unicode"
)

func filterRunes(c rune) bool {
	return !unicode.IsLetter(c)
}

func getWords(query string) []string {
	query = strings.TrimSpace(query)
	if len(query) == 0 {
		return []string{}
	}
	query = strings.ToLower(query)
	return strings.FieldsFunc(query, filterRunes)
}

func filterUnique(words []string) []string {
	wordsMap := map[string]bool{}
	res := []string{}

	for _, word := range words {
		_, wordExists := wordsMap[word]
		if !wordExists {
			wordsMap[word] = true
			res = append(res, word)
		}
	}
	return res
}

func morphWords(queryWords []string) []string {
	res := []string{}
	for _, word := range queryWords {
		norms := morphanalyzer.NormWord(word)
		if len(norms) == 0 {
			norms = append(norms, word)
		}
		normsUnique := filterUnique(norms)

		res = append(res, normsUnique...)
	}
	return res
}

func getNorms(query string) []string {
	words := getWords(query)
	norms := morphWords(words)

	return norms
}

func ParseSearchQuery(query string) []string {
	norms := getNorms(query)
	searchWords := filterStopWords(norms)

	return searchWords
}
