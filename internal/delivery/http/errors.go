package http

import "fmt"

var (
	ErrWrongQueryParam = fmt.Errorf("invalid query parameters")
	ErrWrongBodyParam  = fmt.Errorf("incorrect JSON parameters")
)
