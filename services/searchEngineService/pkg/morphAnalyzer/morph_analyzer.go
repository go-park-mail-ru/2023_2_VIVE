package morphanalyzer

import "github.com/vbatushev/morph"

func InitMorphAnalyzer() error {
	if err := morph.Init(); err != nil {
		return err
	}
	return nil
}

// for correct work of this morph instal python dictionaries
// pip install --user pymorphy2-dicts-ru

func NormWord(word string) []string {
	_, norms, _ := morph.Parse(word)
	return norms
}
