package domain

import "net/http"

type SessionRepository interface {
	AddSession(user *User) (*http.Cookie, error)
	DeleteSession(cookie *http.Cookie) error
	ValidateSession(cookie *http.Cookie) error
}
