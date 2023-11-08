package http

import (
	"HnH/configs"
	"HnH/internal/delivery/http/middleware"
	"HnH/internal/domain"
	"HnH/internal/usecase"
	"HnH/pkg/responseTemplates"
	"HnH/pkg/sanitizer"
	"HnH/pkg/serverErrors"

	"encoding/json"
	"io"
	"net/http"
	"os"
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

	router.Handle("/upload_avatar",
		middleware.AuthMiddleware(sessionUCase, http.HandlerFunc(handler.UploadAvatar))).
		Methods("PUT")
}

func (userHandler *UserHandler) sanitizeUser(user *domain.User) {
	user.Email = sanitizer.XSS.Sanitize(user.Email)
	user.FirstName = sanitizer.XSS.Sanitize(user.FirstName)
	user.LastName = sanitizer.XSS.Sanitize(user.LastName)

	if user.Birthday != nil {
		*user.Birthday = sanitizer.XSS.Sanitize(*user.Birthday)
	}

	if user.PhoneNumber != nil {
		*user.PhoneNumber = sanitizer.XSS.Sanitize(*user.PhoneNumber)
	}

	if user.Location != nil {
		*user.Location = sanitizer.XSS.Sanitize(*user.Location)
	}
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

	userHandler.sanitizeUser(user)

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

func (userHandler *UserHandler) UploadAvatar(w http.ResponseWriter, r *http.Request) {
	uploadedData, handler, err := r.FormFile("avatar")
	if err != nil {
		responseTemplates.SendErrorMessage(w, err, http.StatusBadRequest)
		return
	}
	defer uploadedData.Close()

	f, err := os.OpenFile(configs.CURRENT_DIR+configs.UPLOADS_DIR+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		responseTemplates.SendErrorMessage(w, err, http.StatusInternalServerError)
		return
	}

	io.Copy(f, uploadedData)

	f.Sync()
	f.Close()

	cookie, _ := r.Cookie("session")

	uplErr := userHandler.userUsecase.UploadAvatar(cookie.Value, configs.UPLOADS_DIR+handler.Filename)
	if uplErr != nil {
		responseTemplates.SendErrorMessage(w, uplErr, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

/*func (userHandler *UserHandler) GetAvatar(w http.ResponseWriter, r *http.Request) {
	cookie, _ := r.Cookie("session")


}*/
