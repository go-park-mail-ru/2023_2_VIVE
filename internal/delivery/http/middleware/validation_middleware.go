package middleware

import (
	"HnH/configs"
	"HnH/internal/repository/redisRepo"
	"HnH/pkg/responseTemplates"
	"HnH/pkg/serverErrors"

	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"mime"
	"net/http"
	"strconv"
	"strings"
	"time"
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
		contentType := r.Header.Get("Content-Type")

		mt, _, err := mime.ParseMediaType(contentType)
		if err != nil {
			responseTemplates.SendErrorMessage(w, MALFORMED_CONTENT_TYPE_HEADER, http.StatusBadRequest)
			return
		}

		if mt != "application/json" {
			responseTemplates.SendErrorMessage(w, INCORRECT_CONTENT_TYPE_JSON, http.StatusUnsupportedMediaType)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func CSRFProtectionMiddleware(sessionRepo redisRepo.ISessionRepository, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" && !ifAuthURL(r.URL.Path, r.Method) {
			cookie, err := r.Cookie("session")
			if err != nil {
				responseTemplates.SendErrorMessage(w, serverErrors.NO_COOKIE, http.StatusUnauthorized)
				return
			}

			sessionID := cookie.Value
			userID, err := sessionRepo.GetUserIdBySession(sessionID)
			if err != nil {
				responseTemplates.SendErrorMessage(w, err, http.StatusUnauthorized)
				return
			}

			csrfToken := r.Header.Get("X-CSRF-token")

			if csrfToken == "" {
				newToken := createToken(sessionID, userID, time.Now().Add(1*time.Hour).Unix())
				w.Header().Set("X-CSRF-token", newToken)

				responseTemplates.SendErrorMessage(w, NO_TOKEN, http.StatusForbidden)
				return
			}

			ok, err := checkToken(sessionID, userID, csrfToken)
			if !ok {
				newToken := createToken(sessionID, userID, time.Now().Add(1*time.Hour).Unix())
				w.Header().Set("X-CSRF-token", newToken)

				if err == nil {
					err = BAD_TOKEN
				}

				responseTemplates.SendErrorMessage(w, err, http.StatusForbidden)
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
