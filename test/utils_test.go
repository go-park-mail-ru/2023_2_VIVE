package requestHandlers_test

import (
	"net/http"
)

const (
	sessionUrl     = "http://hnh.ru/session"
	usersUrl       = "http://hnh.ru/users"
	currentUserUrl = "http://hnh.ru/current_user"
	vacanciesUrl   = "http://hnh.ru/vacancies"
)

const (
	INCORRECT_CREDENTIALS  = `{"message":"Incorrect credentials"}`
	NO_DATA_FOUND          = `{"message":"Account data not found"}`
	ACCOUNT_ALREADY_EXISTS = `{"message":"An account with given email already exists"}`
	INVALID_ROLE           = `{"message":"The entered role does not exist"}`
	INVALID_EMAIL          = `{"message":"The entered email-address is not a real one"}`
	INVALID_PASSWORD       = `{"message":"The entered password does not meet the requirements"}`
	INCORRECT_ROLE         = `{"message":"An account with chosen role does not exist"}`

	INVALID_COOKIE = `{"message":"The cookie provided is invalid"}`
	NO_COOKIE      = `{"message":"No cookie provided"}`
	AUTH_REQUIRED  = `{"message":"You need to be authenticated"}`

	MISSED_FIELD_JSON = `{"message":"invalid character ':' looking for beginning of object key string"}`
)

type JsonTestCase struct {
	requestBody  string
	statusCode   int
	responseBody string
}

type CookieTestCase struct {
	cookie             *http.Cookie
	expectedError      string
	expectedStatusCode int
}

type GetUserTestCase struct {
	authData           string
	expectedMessage    string
	expectedStatusCode int
}
