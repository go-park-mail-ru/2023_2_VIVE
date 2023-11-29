package testHelper

import (
	"database/sql/driver"
	"fmt"
	"io"
	"time"

	"github.com/sirupsen/logrus"
)

var (
	ErrQuery = fmt.Errorf("some query error")

	location, _ = time.LoadLocation("Local")
	Created_at  = time.Date(2023, 11, 1, 0, 0, 0, 0, location)
	Updated_at  = time.Date(2023, 11, 2, 0, 0, 0, 0, location)
)

const (
	SelectQuery      = "SELECT(.|\n)+FROM(.|\n)+"
	SelectExistQuery = "SELECT EXISTS(.|\n)+"
	InsertQuery      = "INSERT(.|\n)+INTO(.|\n)+"
	UpdateQuery      = "UPDATE(.|\n)+SET(.|\n)+WHERE(.|\n)+"
	DeleteQuery      = "DELETE(.|\n)+FROM(.|\n)+"
)

// Converts given slice of ints into slice of driver.Vilue
func SliceIntToDriverValue(slice []int) []driver.Value {
	result := make([]driver.Value, len(slice))

	for i := 0; i < len(slice); i++ {
		result[i] = slice[i]
	}

	return result
}

func InitCtxLogger() *logrus.Entry {
	logger := &logrus.Entry{
		Logger: &logrus.Logger{
			Out: io.Discard,
		},
	}
	return logger
}

func ErrNotEqual(expected, actual any) string {
	return fmt.Sprintf(
		"actual does not match expected:\n\tactual: %v\n\texpected: %v\n",
		actual,
		expected,
	)
}
