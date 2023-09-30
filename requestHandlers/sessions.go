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
		errResp := serverErrors.ServerError{Message: err.Error()}
		errJs, _ := json.Marshal(errResp)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write(errJs)
		return
	}

	loginErr := modelHandlers.CheckUser(user)
	if loginErr != nil {
		errResp := serverErrors.ServerError{Message: loginErr.Error()}
		errJs, _ := json.Marshal(errResp)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		w.Write(errJs)
		return
	}

	cookie := modelHandlers.AddSession(user)
	http.SetCookie(w, &cookie)
	w.WriteHeader(http.StatusOK)
}

func Logout(w http.ResponseWriter, r *http.Request) {
	session, err := r.Cookie("session")

	if errors.Is(err, http.ErrNoCookie) {
		errResp := serverErrors.ServerError{Message: serverErrors.NO_COOKIE.Error()}
		errJs, _ := json.Marshal(errResp)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		w.Write(errJs)
		return
	}

	deleteErr := modelHandlers.DeleteSession(session)
	if deleteErr != nil {
		errResp := serverErrors.ServerError{Message: deleteErr.Error()}
		errJs, _ := json.Marshal(errResp)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		w.Write(errJs)
		return
	}

	session.Expires = time.Now().AddDate(0, 0, -1)
	http.SetCookie(w, session)
	w.WriteHeader(http.StatusOK)
}

func CheckLogin(w http.ResponseWriter, r *http.Request) {
	session, err := r.Cookie("session")

	if errors.Is(err, http.ErrNoCookie) {
		errResp := serverErrors.ServerError{Message: serverErrors.NO_COOKIE.Error()}
		errJs, _ := json.Marshal(errResp)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		w.Write(errJs)
		return
	}

	sessionErr := modelHandlers.ValidateSession(session)
	if sessionErr != nil {
		errResp := serverErrors.ServerError{Message: sessionErr.Error()}
		errJs, _ := json.Marshal(errResp)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		w.Write(errJs)
		return
	}

	w.WriteHeader(http.StatusOK)
}
