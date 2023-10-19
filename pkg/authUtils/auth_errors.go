package authUtils

import "fmt"

var (
	INCORRECT_CREDENTIALS = fmt.Errorf("Incorrect credentials")

	INVALID_EMAIL = fmt.Errorf("The entered email-address is not a real one")
	EMPTY_EMAIL   = fmt.Errorf("You haven't passed any email-address")

	EMPTY_PASSWORD   = fmt.Errorf("You haven't passed a password")
	INVALID_PASSWORD = fmt.Errorf("The entered password does not meet the requirements")
)
