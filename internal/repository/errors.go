package repository

import "fmt"

var (
	ENTITY_NOT_FOUND = fmt.Errorf("The entity you requested is not found")
)
