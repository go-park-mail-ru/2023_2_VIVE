package http

import (
	"HnH/internal/appErrors"
	"HnH/internal/domain"
	"HnH/internal/usecase"
	"HnH/pkg/contextUtils"
	"HnH/pkg/middleware"
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
		errToSend, code := appErrors.GetErrAndCodeToSend(loginErr)
		responseTemplates.SendErrorMessage(w, errToSend, code)
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
	deleteErr := sessionHandler.sessionUsecase.Logout(r.Context())
	if deleteErr != nil {
		errToSend, code := appErrors.GetErrAndCodeToSend(deleteErr)
		responseTemplates.SendErrorMessage(w, errToSend, code)
		return
	}

	cookie := &http.Cookie{
		Name:     "session",
		Value:    contextUtils.GetSessionIDFromCtx(r.Context()),
		Expires:  time.Now().AddDate(0, 0, -1),
		Path:     "/",
		Secure:   false,
		HttpOnly: true,
	}

	http.SetCookie(w, cookie)
	w.WriteHeader(http.StatusOK)
}

func (sessionHandler *SessionHandler) CheckLogin(w http.ResponseWriter, r *http.Request) {
	_, sessionErr := sessionHandler.sessionUsecase.CheckLogin(r.Context())
	if sessionErr != nil {
		errToSend, code := appErrors.GetErrAndCodeToSend(sessionErr)
		responseTemplates.SendErrorMessage(w, errToSend, code)
		return
	}

	w.WriteHeader(http.StatusOK)
}
