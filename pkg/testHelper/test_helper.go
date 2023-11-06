package testHelper

import (
	"database/sql/driver"
	"fmt"
	"time"
)

var (
	ErrQuery = fmt.Errorf("some query error")

	location, _ = time.LoadLocation("Local")
	Created_at  = time.Date(2023, 11, 1, 0, 0, 0, 0, location)
	Updated_at  = time.Date(2023, 11, 2, 0, 0, 0, 0, location)
)

const (
	SELECT_QUERY = "SELECT(.|\n)+FROM(.|\n)+"
	INSERT_QUERY = "INSERT(.|\n)+INTO(.|\n)+RETURNING(.|\n)+"
	UPDATE_QUERY = "UPDATE(.|\n)+SET(.|\n)+FROM(.|\n)+WHERE(.|\n)+"
	DELETE_QUERY = "DELETE(.|\n)+FROM(.|\n)+"
)

// Converts given slice of ints into slice of driver.Vilue
func SliceIntToDriverValue(slice []int) []driver.Value {
	result := make([]driver.Value, len(slice))

	for i := 0; i < len(slice); i++ {
		result[i] = slice[i]
	}

	return result
}
