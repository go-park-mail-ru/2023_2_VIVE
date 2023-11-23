package searchEngineUtils

var (
	StopWords = map[string]bool{
		"и": true,
		"но": true,
		"а": true,
		"в": true,
		"на": true,
		"под": true,
		"о": true,
		"об": true,
		"при": true,
		"к": true,
		"с": true,
		"у": true,
		"за": true,
		"до": true,
		"для": true,
		"по": true,
		"из": true,
		"что": true,
		"это": true,
		"тот": true,
		"эта": true,
		"те": true,
		"так": true,
		"такой": true,
		"такие": true,
		"там": true,
		"тут": true,
		"где": true,
		"когда": true,
		"если": true,
		"потому": true,
		"потому что": true, 
		"как": true,
		"чтобы": true,
		"который": true,
		"которая": true,
		"которые": true,
		"чем": true,
		"без": true,
		"над": true,
		"через": true,
		"между": true,
	}
)


func filterStopWords(words []string) []string {
	res := []string{}
	for _, word := range words {
		if _, isStopWord := StopWords[word]; !isStopWord {
			res = append(res, word)
		}
	}
	return res
}
