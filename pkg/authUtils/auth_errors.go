package authUtils

import "fmt"

var (
	INCORRECT_CREDENTIALS = fmt.Errorf("Incorrect credentials")

	INVALID_EMAIL = fmt.Errorf("The entered email-address is not a real one")
	EMPTY_EMAIL   = fmt.Errorf("You haven't passed any email-address")

	EMPTY_PASSWORD   = fmt.Errorf("You haven't passed a password")
	INVALID_PASSWORD = fmt.Errorf("The entered password does not meet the requirements")

	ENTITY_NOT_FOUND     = fmt.Errorf("The entity you requested is not found")
	ERROR_WHILE_WRITING  = fmt.Errorf("An error occurred while writing to the database")
	ERROR_WHILE_DELETING = fmt.Errorf("Error while deleting session")
)
