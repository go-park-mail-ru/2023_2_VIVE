package usecase

import "fmt"

var (
	INAPPROPRIATE_ROLE = fmt.Errorf("Such a request is not possible with your role")
)
