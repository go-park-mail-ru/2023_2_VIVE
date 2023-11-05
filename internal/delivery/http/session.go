package http

import (
	"HnH/internal/domain"
	"HnH/internal/usecase"
	"HnH/pkg/serverErrors"

	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type SessionHandler struct {
	sessionUsecase usecase.ISessionUsecase
}

func NewSessionHandler(router *mux.Router, sessionUCase usecase.ISessionUsecase) {
	handler := &SessionHandler{
		sessionUsecase: sessionUCase,
	}

	router.HandleFunc("/session", handler.Login).Methods("POST")
	router.HandleFunc("/session", handler.Logout).Methods("DELETE")
	router.HandleFunc("/session", handler.CheckLogin).Methods("GET")
}

func (sessionHandler *SessionHandler) Login(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	user := new(domain.User)

	err := json.NewDecoder(r.Body).Decode(user)
	if err != nil {
		sendErrorMessage(w, err, http.StatusBadRequest)
		return
	}

	sessionID, loginErr := sessionHandler.sessionUsecase.Login(user)
	if loginErr != nil {
		sendErrorMessage(w, loginErr, http.StatusBadRequest)
		return
	}

	cookie := &http.Cookie{
		Name:     "session",
		Value:    sessionID,
		Expires:  time.Now().Add(10 * time.Hour),
		Path:     "/",
		Secure:   false,
		HttpOnly: true,
	}

	http.SetCookie(w, cookie)
	w.WriteHeader(http.StatusOK)
}

func (sessionHandler *SessionHandler) Logout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session")

	if errors.Is(err, http.ErrNoCookie) {
		sendErrorMessage(w, serverErrors.NO_COOKIE, http.StatusUnauthorized)
		return
	}

	deleteErr := sessionHandler.sessionUsecase.Logout(cookie.Value)
	if deleteErr != nil {
		sendErrorMessage(w, deleteErr, http.StatusUnauthorized)
		return
	}

	cookie.Expires = time.Now().AddDate(0, 0, -1)

	http.SetCookie(w, cookie)
	w.WriteHeader(http.StatusOK)
}

func (sessionHandler *SessionHandler) CheckLogin(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session")

	if errors.Is(err, http.ErrNoCookie) {
		sendErrorMessage(w, serverErrors.NO_COOKIE, http.StatusUnauthorized)
		return
	}

	sessionErr := sessionHandler.sessionUsecase.CheckLogin(cookie.Value)
	if sessionErr != nil {
		sendErrorMessage(w, sessionErr, http.StatusUnauthorized)
		return
	}

	w.WriteHeader(http.StatusOK)
}
