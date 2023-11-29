package http

import (
	"HnH/internal/delivery/http/middleware"
	"HnH/internal/domain"
	"HnH/internal/usecase"
	"HnH/pkg/responseTemplates"

	"encoding/json"
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

	router.Handle("/session",
		middleware.JSONBodyValidationMiddleware(http.HandlerFunc(handler.Login))).
		Methods("POST")

	router.Handle("/session",
		middleware.AuthMiddleware(sessionUCase, http.HandlerFunc(handler.Logout))).
		Methods("DELETE")

	router.Handle("/session",
		middleware.AuthMiddleware(sessionUCase, http.HandlerFunc(handler.CheckLogin))).
		Methods("GET")
}

func (sessionHandler *SessionHandler) Login(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	user := new(domain.DbUser)

	err := json.NewDecoder(r.Body).Decode(user)
	if err != nil {
		responseTemplates.SendErrorMessage(w, err, http.StatusBadRequest)
		return
	}

	expiryTime := time.Now().Add(10 * time.Hour)

	sessionID, loginErr := sessionHandler.sessionUsecase.Login(r.Context(), user, expiryTime.Unix())
	if loginErr != nil {
		responseTemplates.SendErrorMessage(w, loginErr, http.StatusBadRequest)
		return
	}

	cookie := &http.Cookie{
		Name:     "session",
		Value:    sessionID,
		Expires:  expiryTime,
		Path:     "/",
		Secure:   false,
		HttpOnly: true,
	}

	http.SetCookie(w, cookie)
	w.WriteHeader(http.StatusOK)
}

func (sessionHandler *SessionHandler) Logout(w http.ResponseWriter, r *http.Request) {
	cookie, _ := r.Cookie("session")

	deleteErr := sessionHandler.sessionUsecase.Logout(r.Context(), cookie.Value)
	if deleteErr != nil {
		responseTemplates.SendErrorMessage(w, deleteErr, http.StatusUnauthorized)
		return
	}

	cookie.Expires = time.Now().AddDate(0, 0, -1)

	http.SetCookie(w, cookie)
	w.WriteHeader(http.StatusOK)
}

func (sessionHandler *SessionHandler) CheckLogin(w http.ResponseWriter, r *http.Request) {
	cookie, _ := r.Cookie("session")

	sessionErr := sessionHandler.sessionUsecase.CheckLogin(r.Context(), cookie.Value)
	if sessionErr != nil {
		responseTemplates.SendErrorMessage(w, sessionErr, http.StatusUnauthorized)
		return
	}

	w.WriteHeader(http.StatusOK)
}
