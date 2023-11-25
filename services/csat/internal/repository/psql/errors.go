package psql

import "fmt"

var (
	ErrNoLastUpdate = fmt.Errorf("no last update of given user found")
)
