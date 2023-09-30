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
		errResp := serverErrors.ServerError{Message: err.Error()}
		errJs, _ := json.Marshal(errResp)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write(errJs)
		return
	}

	addStatus := modelHandlers.AddUser(newUser)
	if addStatus != nil {
		errResp := serverErrors.ServerError{Message: addStatus.Error()}
		errJs, _ := json.Marshal(errResp)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusConflict)
		w.Write(errJs)
		return
	}

	cookie := modelHandlers.AddSession(newUser)
	http.SetCookie(w, &cookie)
	w.WriteHeader(http.StatusOK)
}

func GetInfo(w http.ResponseWriter, r *http.Request) {
	session, err := r.Cookie("session")

	if errors.Is(err, http.ErrNoCookie) {
		errResp := serverErrors.ServerError{Message: serverErrors.NO_COOKIE.Error()}
		errJs, _ := json.Marshal(errResp)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		w.Write(errJs)
		return
	}

	validStatus := modelHandlers.ValidateSession(session)
	if validStatus != nil {
		errResp := serverErrors.ServerError{Message: validStatus.Error()}
		errJs, _ := json.Marshal(errResp)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		w.Write(errJs)
		return
	}

	user := modelHandlers.GetUserInfo(session)

	js, err := json.Marshal(*user)
	if err != nil {
		errResp := serverErrors.ServerError{Message: serverErrors.INTERNAL_SERVER_ERROR.Error()}
		errJs, _ := json.Marshal(errResp)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(errJs)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}
