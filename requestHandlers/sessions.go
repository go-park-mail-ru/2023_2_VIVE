package requestHandlers

import (
	"HnH/modelHandlers"
	"HnH/models"
	"HnH/serverErrors"
	"encoding/json"
	"errors"
	"net/http"
	"time"
)

func Login(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	user := new(models.User)

	err := json.NewDecoder(r.Body).Decode(user)
	if err != nil {
		sendErrorMessage(w, err, http.StatusBadRequest)
		return
	}

	loginErr := modelHandlers.CheckUser(user)
	if loginErr != nil {
		sendErrorMessage(w, loginErr, http.StatusUnauthorized)
		return
	}

	cookie := modelHandlers.AddSession(user)
	http.SetCookie(w, cookie)
	w.WriteHeader(http.StatusOK)
}

func Logout(w http.ResponseWriter, r *http.Request) {
	session, err := r.Cookie("session")

	if errors.Is(err, http.ErrNoCookie) {
		sendErrorMessage(w, serverErrors.NO_COOKIE, http.StatusUnauthorized)
		return
	}

	deleteErr := modelHandlers.DeleteSession(session)
	if deleteErr != nil {
		sendErrorMessage(w, deleteErr, http.StatusUnauthorized)
		return
	}

	session.Expires = time.Now().AddDate(0, 0, -1)
	http.SetCookie(w, session)
	w.WriteHeader(http.StatusOK)
}

func CheckLogin(w http.ResponseWriter, r *http.Request) {
	session, err := r.Cookie("session")

	if errors.Is(err, http.ErrNoCookie) {
		sendErrorMessage(w, serverErrors.NO_COOKIE, http.StatusUnauthorized)
	}

	sessionErr := modelHandlers.ValidateSession(session)
	if sessionErr != nil {
		sendErrorMessage(w, sessionErr, http.StatusUnauthorized)
		return
	}

	w.WriteHeader(http.StatusOK)
}
