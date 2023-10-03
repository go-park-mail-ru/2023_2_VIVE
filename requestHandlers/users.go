package requestHandlers

import (
	"HnH/modelHandlers"
	"HnH/models"
	"HnH/serverErrors"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

func SignUp(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	fmt.Printf("Signing Up\n\n")
	newUser := new(models.User)

	err := json.NewDecoder(r.Body).Decode(newUser)
	fmt.Printf("%v\n\n", *newUser)
	if err != nil {
		sendErrorMessage(w, err, http.StatusBadRequest)
		return
	}

	addStatus := modelHandlers.AddUser(newUser)
	fmt.Printf("Status: %v\n\n", addStatus)
	if addStatus != nil {
		sendErrorMessage(w, addStatus, http.StatusUnauthorized)
		return
	}

	cookie := modelHandlers.AddSession(newUser)
	http.SetCookie(w, cookie)
	w.WriteHeader(http.StatusOK)
}

func GetInfo(w http.ResponseWriter, r *http.Request) {
	session, err := r.Cookie("session")

	if errors.Is(err, http.ErrNoCookie) {
		sendErrorMessage(w, serverErrors.NO_COOKIE, http.StatusUnauthorized)
		return
	}

	validStatus := modelHandlers.ValidateSession(session)
	if validStatus != nil {
		sendErrorMessage(w, validStatus, http.StatusUnauthorized)
		return
	}

	user := modelHandlers.GetUserInfo(session)

	js, err := json.Marshal(*user)
	if err != nil {
		sendErrorMessage(w, serverErrors.INTERNAL_SERVER_ERROR, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}
