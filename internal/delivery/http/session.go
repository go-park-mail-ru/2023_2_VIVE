package http

import (
	"HnH/internal/appErrors"
	"HnH/internal/domain"
	"HnH/internal/usecase"
	"HnH/pkg/contextUtils"
	"HnH/pkg/middleware"
	"HnH/pkg/responseTemplates"

	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/mailru/easyjson"
	"github.com/sirupsen/logrus"
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
	contextLogger := contextUtils.GetContextLogger(r.Context())
	defer r.Body.Close()

	user := new(domain.DbUser)
	err := easyjson.UnmarshalFromReader(r.Body, user)
	if err != nil {
		sendErr := responseTemplates.SendErrorMessage(w, err, http.StatusBadRequest)
		if sendErr != nil {
			contextLogger.WithFields(logrus.Fields{
				"error_msg":     sendErr,
				"error_to_send": err,
			}).
				Error("could not send error message")
		}
		return
	}

	expiryTime := time.Now().Add(10 * time.Hour)

	sessionID, loginErr := sessionHandler.sessionUsecase.Login(r.Context(), user, expiryTime.Unix())
	if loginErr != nil {
		errToSend, code := appErrors.GetErrAndCodeToSend(loginErr)
		sendErr := responseTemplates.SendErrorMessage(w, errToSend, code)
		if sendErr != nil {
			contextLogger.WithFields(logrus.Fields{
				"error_msg":     sendErr,
				"error_to_send": errToSend,
			}).
				Error("could not send error message")
		}
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
	contextLogger := contextUtils.GetContextLogger(r.Context())
	deleteErr := sessionHandler.sessionUsecase.Logout(r.Context())
	if deleteErr != nil {
		errToSend, code := appErrors.GetErrAndCodeToSend(deleteErr)
		sendErr := responseTemplates.SendErrorMessage(w, errToSend, code)
		if sendErr != nil {
			contextLogger.WithFields(logrus.Fields{
				"error_msg":     sendErr,
				"error_to_send": errToSend,
			}).
				Error("could not send error message")
		}
		return
	}

	sessionID, err := contextUtils.GetSessionIDFromCtx(r.Context())
	if err != nil {
		errToSend, code := appErrors.GetErrAndCodeToSend(err)
		sendErr := responseTemplates.SendErrorMessage(w, errToSend, code)
		if sendErr != nil {
			contextLogger.WithFields(logrus.Fields{
				"error_msg":     sendErr,
				"error_to_send": errToSend,
			}).
				Error("could not send error message")
		}
		return 
	}

	cookie := &http.Cookie{
		Name:     "session",
		Value:    sessionID,
		Expires:  time.Now().AddDate(0, 0, -1),
		Path:     "/",
		Secure:   false,
		HttpOnly: true,
	}

	http.SetCookie(w, cookie)
	w.WriteHeader(http.StatusOK)
}

func (sessionHandler *SessionHandler) CheckLogin(w http.ResponseWriter, r *http.Request) {
	contextLogger := contextUtils.GetContextLogger(r.Context())
	_, sessionErr := sessionHandler.sessionUsecase.CheckLogin(r.Context())
	if sessionErr != nil {
		errToSend, code := appErrors.GetErrAndCodeToSend(sessionErr)
		sendErr := responseTemplates.SendErrorMessage(w, errToSend, code)
		if sendErr != nil {
			contextLogger.WithFields(logrus.Fields{
				"error_msg":     sendErr,
				"error_to_send": errToSend,
			}).
				Error("could not send error message")
		}
		return
	}

	w.WriteHeader(http.StatusOK)
}
