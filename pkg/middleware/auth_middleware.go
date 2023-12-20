package middleware

import (
	"HnH/internal/appErrors"
	"HnH/internal/usecase"
	"HnH/pkg/contextUtils"
	"HnH/pkg/responseTemplates"
	"HnH/pkg/serverErrors"

	"context"
	"errors"
	"net/http"

	"github.com/sirupsen/logrus"
)

func AuthMiddleware(sessionUsecase usecase.ISessionUsecase, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		contextLogger := contextUtils.GetContextLogger(r.Context())
		cookie, err := r.Cookie("session")
		if errors.Is(err, http.ErrNoCookie) {
			sendErr := responseTemplates.SendErrorMessage(w, serverErrors.NO_COOKIE, http.StatusBadRequest)
			if sendErr != nil {
				contextLogger.WithFields(logrus.Fields{
					"err_msg":       sendErr,
					"error_to_send": serverErrors.NO_COOKIE,
				}).
					Error("could not send error")
			}
			return
		} else if err != nil {
			sendErr := responseTemplates.SendErrorMessage(w, err, http.StatusBadRequest)
			if sendErr != nil {
				contextLogger.WithFields(logrus.Fields{
					"err_msg":       sendErr,
					"error_to_send": err,
				}).
					Error("could not send error")
			}
			return
		}

		ctxWithCookie := context.WithValue(r.Context(), contextUtils.SESSION_ID_KEY, cookie.Value)
		userID, authErr := sessionUsecase.CheckLogin(ctxWithCookie)
		if authErr != nil {
			errToSend, code := appErrors.GetErrAndCodeToSend(authErr)
			sendErr := responseTemplates.SendErrorMessage(w, errToSend, code)
			if sendErr != nil {
				contextLogger.WithFields(logrus.Fields{
					"err_msg":       sendErr,
					"error_to_send": errToSend,
				}).
					Error("could not send error")
			}
			return
		}

		ctxWithUID := context.WithValue(ctxWithCookie, contextUtils.USER_ID_KEY, userID)

		next.ServeHTTP(w, r.WithContext(ctxWithUID))
	})
}

func SetSessionIDIfExists(sessionUsecase usecase.ISessionUsecase, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		contextLogger := contextUtils.GetContextLogger(r.Context())
		cookie, err := r.Cookie("session")
		if errors.Is(err, http.ErrNoCookie) {
			next.ServeHTTP(w, r)
			return
		} else if err != nil {
			sendErr := responseTemplates.SendErrorMessage(w, err, http.StatusBadRequest)
			if sendErr != nil {
				contextLogger.WithFields(logrus.Fields{
					"err_msg":       sendErr,
					"error_to_send": err,
				}).
					Error("could not send error")
			}
			return
		}

		ctxWithCookie := context.WithValue(r.Context(), contextUtils.SESSION_ID_KEY, cookie.Value)
		userID, authErr := sessionUsecase.CheckLogin(ctxWithCookie)
		if authErr != nil {
			next.ServeHTTP(w, r)
			return
		}

		ctxWithUID := context.WithValue(ctxWithCookie, contextUtils.USER_ID_KEY, userID)

		next.ServeHTTP(w, r.WithContext(ctxWithUID))
	})
}
