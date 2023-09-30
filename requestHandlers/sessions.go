package requestHandlers

import (
	"encoding/json"
	"errors"
	"models/modelHandlers"
	"models/models"
	"models/serverErrors"
	"net/http"
	"time"
)

func Login(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	user := new(models.User)

	err := json.NewDecoder(r.Body).Decode(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	loginErr := modelHandlers.CheckUser(user)
	if loginErr != nil {
		http.Error(w, loginErr.Error(), http.StatusUnauthorized)
		return
	}

	cookie := modelHandlers.AddSession(user)
	http.SetCookie(w, &cookie)
	w.WriteHeader(http.StatusOK)
}

func Logout(w http.ResponseWriter, r *http.Request) {
	session, err := r.Cookie("session")

	if errors.Is(err, http.ErrNoCookie) {
		http.Error(w, serverErrors.NO_COOKIE.Error(), http.StatusUnauthorized)
		return
	}

	deleteErr := modelHandlers.DeleteSession(session)
	if deleteErr != nil {
		http.Error(w, deleteErr.Error(), http.StatusUnauthorized)
		return
	}

	session.Expires = time.Now().AddDate(0, 0, -1)
	http.SetCookie(w, session)
	w.WriteHeader(http.StatusOK)
}

func CheckLogin(w http.ResponseWriter, r *http.Request) {
	session, err := r.Cookie("session")

	if errors.Is(err, http.ErrNoCookie) {
		http.Error(w, serverErrors.NO_COOKIE.Error(), http.StatusUnauthorized)
		return
	}

	sessionErr := modelHandlers.ValidateSession(session)
	if sessionErr != nil {
		http.Error(w, sessionErr.Error(), http.StatusUnauthorized)
		return
	}

	w.WriteHeader(http.StatusOK)
}
