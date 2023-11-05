package repository

import "fmt"

var (
	ErrEntityNotFound = fmt.Errorf("the entity you requested is not found")
	ErrNotInserted    = fmt.Errorf("could not insert data into db")
	ErrNoRowsUpdated  = fmt.Errorf("after your query no rows were updated")
	ErrNoRowsDeleted  = fmt.Errorf("after your query no rows were deleted")
)
