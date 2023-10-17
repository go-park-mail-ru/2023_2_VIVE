package http

import (
	"HnH/internal/domain"
	"HnH/internal/usecase"
	"HnH/pkg/serverErrors"

	"encoding/json"
	"errors"
	"net/http"
)

func SignUp(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	newUser := new(domain.User)

	err := json.NewDecoder(r.Body).Decode(newUser)
	if err != nil {
		sendErrorMessage(w, err, http.StatusBadRequest)
		return
	}

	cookie, err := usecase.SignUp(newUser)
	if err != nil {
		sendErrorMessage(w, err, http.StatusBadRequest)
		return
	}

	http.SetCookie(w, cookie)
	w.WriteHeader(http.StatusOK)
}

func GetInfo(w http.ResponseWriter, r *http.Request) {
	session, err := r.Cookie("session")

	if errors.Is(err, http.ErrNoCookie) {
		sendErrorMessage(w, serverErrors.NO_COOKIE, http.StatusUnauthorized)
		return
	}

	user, err := usecase.GetInfo(session)
	if err != nil {
		sendErrorMessage(w, serverErrors.AUTH_REQUIRED, http.StatusUnauthorized)
		return
	}

	js, err := json.Marshal(*user)
	if err != nil {
		sendErrorMessage(w, serverErrors.INTERNAL_SERVER_ERROR, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}
