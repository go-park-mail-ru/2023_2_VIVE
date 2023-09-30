package requestHandlers

import (
	"encoding/json"
	"models/errors"
	"models/modelHandlers"
	"models/models"
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

	if len(user.Email) == 0 || len(user.Password) == 0 {
		http.Error(w, errors.INCORRECT_CREDENTIALS.Error(), http.StatusUnauthorized)
		return
	}

	passwordStatus := modelHandlers.CheckPassword(*user)
	if passwordStatus != nil {
		http.Error(w, errors.INCORRECT_CREDENTIALS.Error(), http.StatusUnauthorized)
		return
	}

	cookie := modelHandlers.AddSession(*user)
	http.SetCookie(w, &cookie)
	w.WriteHeader(http.StatusOK)
}

func Logout(w http.ResponseWriter, r *http.Request) {
	session, err := r.Cookie("session")

	if err == http.ErrNoCookie {
		http.Error(w, errors.NO_COOKIE.Error(), http.StatusUnauthorized)
		return
	}

	modelHandlers.DeleteSession(session)

	session.Expires = time.Now().AddDate(0, 0, -1)
	http.SetCookie(w, session)
	w.WriteHeader(http.StatusOK)
}

func CheckLogin(w http.ResponseWriter, r *http.Request) {
	session, err := r.Cookie("session")

	if err == http.ErrNoCookie {
		http.Error(w, errors.NO_COOKIE.Error(), http.StatusUnauthorized)
		return
	}

	sessionErr := modelHandlers.ValidateSession(session)
	if sessionErr != nil {
		http.Error(w, sessionErr.Error(), http.StatusUnauthorized)
		return
	}

	w.WriteHeader(http.StatusOK)
}
