package models

import (
	"models/statuses"
	"net/http"
	"sync"
	"time"

	"github.com/google/uuid"
)

type cookieData struct {
	Name    string
	Value   string
	Expires time.Time
}

func getNewRecord(cookie *http.Cookie) cookieData {
	return cookieData{
		Name:    cookie.Name,
		Value:   cookie.Value,
		Expires: cookie.Expires,
	}
}

var sessions = &sync.Map{}

func AddSession() http.Cookie {
	uniqueID := uuid.NewString()

	cookie := http.Cookie{
		Name:     "authCookie",
		Value:    uniqueID,
		Expires:  time.Now().Add(10 * time.Hour),
		Path:     "/",
		Secure:   true,
		HttpOnly: true,
	}

	sessions.Store(getNewRecord(&cookie), uniqueID)

	return cookie
}

func ValidateSession(cookie *http.Cookie) statuses.Status {
	storedID, ok := sessions.Load(getNewRecord(cookie))

	if !ok {
		return statuses.UNAUTHORIZED
	}

	if cookie.Value != storedID {
		return statuses.UNAUTHORIZED
	}

	return 0
}
