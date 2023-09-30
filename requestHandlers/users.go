package requestHandlers

import (
	"encoding/json"
	"errors"
	"models/modelHandlers"
	"models/models"
	"models/serverErrors"
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

	if errors.Is(err, http.ErrNoCookie) {
		http.Error(w, serverErrors.NO_COOKIE.Error(), http.StatusUnauthorized)
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
		http.Error(w, serverErrors.INTERNAL_SERVER_ERROR.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}
