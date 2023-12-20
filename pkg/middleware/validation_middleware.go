package middleware

import (
	"HnH/configs"
	"HnH/internal/appErrors"
	"HnH/internal/usecase"
	"HnH/pkg/contextUtils"
	"HnH/pkg/responseTemplates"
	"HnH/pkg/serverErrors"
	"context"

	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"mime"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

var authorizationURLs = map[string]string{
	"/users":   "POST",
	"/session": "POST",
}

func ifAuthURL(path string, method string) bool {
	authMethod, ok := authorizationURLs[path]
	if !ok {
		return false
	}

	return authMethod == method
}

func JSONBodyValidationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		contextLogger := contextUtils.GetContextLogger(r.Context())
		contentType := r.Header.Get("Content-Type")

		mt, _, err := mime.ParseMediaType(contentType)
		if err != nil {
			sendErr := responseTemplates.SendErrorMessage(w, MALFORMED_CONTENT_TYPE_HEADER, http.StatusBadRequest)
			if sendErr != nil {
				contextLogger.WithFields(logrus.Fields{
					"err_msg": sendErr,
					"error_to_send": MALFORMED_CONTENT_TYPE_HEADER,
				}).
				Error("could not send error")
			}
			return
		}

		if mt != "application/json" {
			sendErr := responseTemplates.SendErrorMessage(w, INCORRECT_CONTENT_TYPE_JSON, http.StatusUnsupportedMediaType)
			if sendErr != nil {
				contextLogger.WithFields(logrus.Fields{
					"err_msg": sendErr,
					"error_to_send": INCORRECT_CONTENT_TYPE_JSON,
				}).
				Error("could not send error")
			}
			return
		}

		next.ServeHTTP(w, r)
	})
}

func CSRFProtectionMiddleware(sessionUCase usecase.ISessionUsecase, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		contextLogger := contextUtils.GetContextLogger(r.Context())
		if r.Method != "GET" && !ifAuthURL(r.URL.Path, r.Method) {
			cookie, err := r.Cookie("session")
			if err != nil {
				sendErr := responseTemplates.SendErrorMessage(w, serverErrors.NO_COOKIE, http.StatusUnauthorized)
				if sendErr != nil {
					contextLogger.WithFields(logrus.Fields{
						"err_msg": sendErr,
						"error_to_send": serverErrors.NO_COOKIE,
					}).
					Error("could not send error")
				}
				return
			}

			ctxWithCookie := context.WithValue(r.Context(), contextUtils.SESSION_ID_KEY, cookie.Value)
			userID, err := sessionUCase.CheckLogin(ctxWithCookie)
			if err != nil {
				errToSend, code := appErrors.GetErrAndCodeToSend(err)
				sendErr := responseTemplates.SendErrorMessage(w, errToSend, code)
				if sendErr != nil {
					contextLogger.WithFields(logrus.Fields{
						"err_msg": sendErr,
						"error_to_send": errToSend,
					}).
					Error("could not send error")
				}
				return
			}

			csrfToken := r.Header.Get("X-CSRF-token")

			if csrfToken == "" {
				newToken := createToken(cookie.Value, userID, time.Now().Add(1*time.Hour).Unix())
				w.Header().Set("X-CSRF-token", newToken)

				sendErr := responseTemplates.SendErrorMessage(w, NO_TOKEN, http.StatusForbidden)
				if sendErr != nil {
					contextLogger.WithFields(logrus.Fields{
						"err_msg": sendErr,
						"error_to_send": NO_TOKEN,
					}).
					Error("could not send error")
				}
				return
			}

			ok, err := checkToken(cookie.Value, userID, csrfToken)
			if !ok {
				newToken := createToken(cookie.Value, userID, time.Now().Add(1*time.Hour).Unix())
				w.Header().Set("X-CSRF-token", newToken)

				if err == nil {
					err = BAD_TOKEN
				}

				sendErr := responseTemplates.SendErrorMessage(w, err, http.StatusForbidden)
				if sendErr != nil {
					contextLogger.WithFields(logrus.Fields{
						"err_msg": sendErr,
						"error_to_send": err,
					}).
					Error("could not send error")
				}
				return
			}
		}

		next.ServeHTTP(w, r)
	})
}

func createToken(sessionID string, userID int, tokenExpTimeSeconds int64) string {
	h := hmac.New(sha256.New, []byte(configs.SECRET_KEY))
	data := fmt.Sprintf("%s%d%d", sessionID, userID, tokenExpTimeSeconds)

	h.Write([]byte(data))
	token := hex.EncodeToString(h.Sum(nil)) + ":" + strconv.FormatInt(tokenExpTimeSeconds, 10)

	return token
}

func checkToken(sessionID string, userID int, inputToken string) (bool, error) {
	tokenData := strings.Split(inputToken, ":")
	if len(tokenData) != 2 {
		return false, BAD_TOKEN
	}

	tokenExpTime, err := strconv.ParseInt(tokenData[1], 10, 64)
	if err != nil {
		return false, BAD_TOKEN
	}

	if tokenExpTime < time.Now().Unix() {
		return false, EXPIRED_TOKEN
	}

	h := hmac.New(sha256.New, []byte(configs.SECRET_KEY))
	data := fmt.Sprintf("%s%d%d", sessionID, userID, tokenExpTime)
	h.Write([]byte(data))
	expectedMAC := h.Sum(nil)

	inputMAC, err := hex.DecodeString(tokenData[0])
	if err != nil {
		return false, BAD_TOKEN
	}

	return hmac.Equal(inputMAC, expectedMAC), nil
}
