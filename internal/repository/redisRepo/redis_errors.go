package redisRepo

import "fmt"

var (
	ENTITY_NOT_FOUND    = fmt.Errorf("The entity you requested is not found")
	ERROR_WHILE_WRITING = fmt.Errorf("An error occurred while writing to the database")
)
