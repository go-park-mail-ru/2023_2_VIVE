package searchOptions

import "fmt"

var (
	ErrNoOption         = fmt.Errorf("no such option provided")
	ErrWrongValueFormat = fmt.Errorf("wrong value format of given option")
)
