package psql

import "fmt"

var (
	ErrNoLastUpdate   = fmt.Errorf("no last update of given user found")
	ErrEntityNotFound = fmt.Errorf("the entity you requested is not found")
)
