package requestHandlers

import (
	"encoding/json"
	"models/errors"
	"models/modelHandlers"
	"models/models"
	"net/http"
)

func SignUp(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	newUser := new(models.User)

	err := json.NewDecoder(r.Body).Decode(newUser)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if len(newUser.Email) == 0 || len(newUser.Password) == 0 {
		http.Error(w, errors.INCORRECT_CREDENTIALS.Error(), http.StatusUnauthorized)
		return
	}

	addStatus := modelHandlers.AddUser(*newUser)
	if addStatus != nil {
		http.Error(w, addStatus.Error(), http.StatusConflict)
	}

	cookie := modelHandlers.AddSession(*newUser)
	http.SetCookie(w, &cookie)
	w.WriteHeader(http.StatusOK)
}

func GetInfo(w http.ResponseWriter, r *http.Request) {
	session, err := r.Cookie("session")

	if err == http.ErrNoCookie {
		http.Error(w, errors.NO_COOKIE.Error(), http.StatusUnauthorized)
		return
	}

	validStatus := modelHandlers.ValidateSession(session)
	if validStatus != nil {
		http.Error(w, validStatus.Error(), http.StatusUnauthorized)
		return
	}

	user := modelHandlers.GetUserInfo(session)

	js, err := json.Marshal(user)
	if err != nil {
		http.Error(w, errors.INTERNAL_SERVER_ERROR.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}
