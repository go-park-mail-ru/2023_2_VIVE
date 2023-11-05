package middleware

import "fmt"

var (
	MALFORMED_CONTENT_TYPE_HEADER = fmt.Errorf("Malformed Content-Type header")
	INCORRECT_CONTENT_TYPE_JSON   = fmt.Errorf("Content-Type header must be application/json")
)
