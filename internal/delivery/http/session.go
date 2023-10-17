package http

import (
	"HnH/internal/domain"
	"HnH/internal/usecase"
	"HnH/pkg/serverErrors"

	"encoding/json"
	"errors"
	"net/http"
)

func Login(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	user := new(domain.User)

	err := json.NewDecoder(r.Body).Decode(user)
	if err != nil {
		sendErrorMessage(w, err, http.StatusBadRequest)
		return
	}

	cookie, loginErr := usecase.Login(user)
	if loginErr != nil {
		sendErrorMessage(w, err, http.StatusBadRequest)
		return
	}

	http.SetCookie(w, cookie)
	w.WriteHeader(http.StatusOK)
}

func Logout(w http.ResponseWriter, r *http.Request) {
	session, err := r.Cookie("session")

	if errors.Is(err, http.ErrNoCookie) {
		sendErrorMessage(w, serverErrors.NO_COOKIE, http.StatusUnauthorized)
		return
	}

	deleteErr := usecase.Logout(session)
	if deleteErr != nil {
		sendErrorMessage(w, deleteErr, http.StatusUnauthorized)
		return
	}

	http.SetCookie(w, session)
	w.WriteHeader(http.StatusOK)
}

func CheckLogin(w http.ResponseWriter, r *http.Request) {
	session, err := r.Cookie("session")

	if errors.Is(err, http.ErrNoCookie) {
		sendErrorMessage(w, serverErrors.NO_COOKIE, http.StatusUnauthorized)
		return
	}

	sessionErr := usecase.CheckLogin(session)
	if sessionErr != nil {
		sendErrorMessage(w, sessionErr, http.StatusUnauthorized)
		return
	}

	w.WriteHeader(http.StatusOK)
}
