package http

import (
	"HnH/internal/delivery/http/middleware"
	"HnH/internal/domain"
	"HnH/internal/usecase"
	"HnH/pkg/responseTemplates"
	"HnH/pkg/serverErrors"

	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type UserHandler struct {
	userUsecase usecase.IUserUsecase
}

func NewUserHandler(router *mux.Router, userUCase usecase.IUserUsecase, sessionUCase usecase.ISessionUsecase) {
	handler := &UserHandler{
		userUsecase: userUCase,
	}

	router.Handle("/users",
		middleware.JSONBodyValidationMiddleware(http.HandlerFunc(handler.SignUp))).
		Methods("POST")

	router.Handle("/current_user",
		middleware.AuthMiddleware(sessionUCase, http.HandlerFunc(handler.GetInfo))).
		Methods("GET")

	router.Handle("/current_user",
		middleware.JSONBodyValidationMiddleware(middleware.AuthMiddleware(sessionUCase, http.HandlerFunc(handler.UpdateInfo)))).
		Methods("PUT")
}

func (userHandler *UserHandler) SignUp(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	newUser := new(domain.User)

	err := json.NewDecoder(r.Body).Decode(newUser)
	if err != nil {
		responseTemplates.SendErrorMessage(w, err, http.StatusBadRequest)
		return
	}

	expiryTime := time.Now().Add(10 * time.Hour)

	sessionID, err := userHandler.userUsecase.SignUp(newUser, expiryTime.Unix())
	if err != nil {
		responseTemplates.SendErrorMessage(w, err, http.StatusBadRequest)
		return
	}

	cookie := &http.Cookie{
		Name:     "session",
		Value:    sessionID,
		Expires:  expiryTime,
		Path:     "/",
		Secure:   false,
		HttpOnly: true,
	}

	http.SetCookie(w, cookie)
	w.WriteHeader(http.StatusOK)
}

func (userHandler *UserHandler) GetInfo(w http.ResponseWriter, r *http.Request) {
	cookie, _ := r.Cookie("session")

	user, err := userHandler.userUsecase.GetInfo(cookie.Value)
	if err != nil {
		responseTemplates.SendErrorMessage(w, serverErrors.AUTH_REQUIRED, http.StatusUnauthorized)
		return
	}

	responseTemplates.MarshalAndSend(w, *user)
}

func (userHandler *UserHandler) UpdateInfo(w http.ResponseWriter, r *http.Request) {
	cookie, _ := r.Cookie("session")

	defer r.Body.Close()

	updateInfo := new(domain.UserUpdate)

	decodeErr := json.NewDecoder(r.Body).Decode(updateInfo)
	if decodeErr != nil {
		responseTemplates.SendErrorMessage(w, decodeErr, http.StatusBadRequest)
		return
	}

	updStatus := userHandler.userUsecase.UpdateInfo(cookie.Value, updateInfo)
	if updStatus != nil {
		responseTemplates.SendErrorMessage(w, updStatus, http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}
