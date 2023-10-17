package serverErrors

import "fmt"

var (
	INCORRECT_CREDENTIALS  = fmt.Errorf("Incorrect credentials")
	INVALID_EMAIL          = fmt.Errorf("The entered email-address is not a real one")
	INVALID_PASSWORD       = fmt.Errorf("The entered password does not meet the requirements")
	INVALID_ROLE           = fmt.Errorf("The entered role does not exist")
	INCORRECT_ROLE         = fmt.Errorf("An account with chosen role does not exist")
	NO_DATA_FOUND          = fmt.Errorf("Account data not found")
	ACCOUNT_ALREADY_EXISTS = fmt.Errorf("An account with given email already exists")
	SESSION_ALREADY_EXISTS = fmt.Errorf("Session has already started")
	INVALID_COOKIE         = fmt.Errorf("The cookie provided is invalid")
	NO_COOKIE              = fmt.Errorf("No cookie provided")
	AUTH_REQUIRED          = fmt.Errorf("You need to be authenticated")
	INTERNAL_SERVER_ERROR  = fmt.Errorf("The server encountered a problem and could not process your request")

	SERVER_IS_NOT_RUNNUNG = fmt.Errorf("the server encountered a problem while starting")
)
