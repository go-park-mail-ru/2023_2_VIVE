package errors

import "fmt"

var INCORRECT_CREDENTIALS error = fmt.Errorf("Incorrect credentials")
var NO_DATA_FOUND error = fmt.Errorf("Account data not found")
var ACCOUNT_ALREADY_EXISTS error = fmt.Errorf("An account with given email already exists")
var SESSION_ALREADY_EXISTS error = fmt.Errorf("Session has already started")
var COOKIE_HAS_EXPIRED error = fmt.Errorf("The cookie provided has expired")
var INVALID_COOKIE error = fmt.Errorf("The cookie provided is invalid")
var NO_COOKIE error = fmt.Errorf("No cookie provided")
var AUTH_REQUIRED error = fmt.Errorf("You need to be authenticated")
var INTERNAL_SERVER_ERROR error = fmt.Errorf("The server encountered a problem and could not process your request")

var SERVER_IS_NOT_RUNNUNG error = fmt.Errorf("the server encountered a problem while starting")
