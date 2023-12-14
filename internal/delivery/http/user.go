package http

import (
	"HnH/internal/appErrors"
	"HnH/internal/domain"
	"HnH/internal/usecase"
	"HnH/pkg/middleware"
	"HnH/pkg/responseTemplates"
	"HnH/pkg/sanitizer"
	"HnH/pkg/serverErrors"
	"errors"

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

	router.Handle("/upload_avatar",
		middleware.AuthMiddleware(sessionUCase, http.HandlerFunc(handler.UploadAvatar))).
		Methods("POST")

	router.Handle("/get_avatar",
		middleware.AuthMiddleware(sessionUCase, http.HandlerFunc(handler.GetAvatar))).
		Methods("GET")
}

func (userHandler *UserHandler) sanitizeUser(user *domain.ApiUser) {
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

	newUser := new(domain.ApiUser)

	err := json.NewDecoder(r.Body).Decode(newUser)
	if err != nil {
		responseTemplates.SendErrorMessage(w, ErrWrongBodyParam, http.StatusBadRequest)
		return
	}

	expiryTime := time.Now().Add(10 * time.Hour)

	sessionID, err := userHandler.userUsecase.SignUp(r.Context(), newUser, expiryTime.Unix())
	if err != nil {
		errToSend, code := appErrors.GetErrAndCodeToSend(err)
		responseTemplates.SendErrorMessage(w, errToSend, code)
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
	// logger := r.Context().Value(middleware.LOGGER_KEY).(*logrus.Entry)
	// logger.Info("GOT REQUEST")

	// requesID, requestID := middleware.GetRequestIDCtx(r.Context())
	// logging.Logger.WithField(requesID, requestID).
	// 	Info("got request")

	// logging.Logger.Info("got request")

	user, err := userHandler.userUsecase.GetInfo(r.Context())
	if err != nil {
		errToSend, code := appErrors.GetErrAndCodeToSend(err)
		responseTemplates.SendErrorMessage(w, errToSend, code)
		return
	}

	userHandler.sanitizeUser(user)

	responseTemplates.MarshalAndSend(w, *user)
}

func (userHandler *UserHandler) UpdateInfo(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	updateInfo := new(domain.UserUpdate)

	decodeErr := json.NewDecoder(r.Body).Decode(updateInfo)
	if decodeErr != nil {
		responseTemplates.SendErrorMessage(w, decodeErr, http.StatusBadRequest)
		return
	}

	updStatus := userHandler.userUsecase.UpdateInfo(r.Context(), updateInfo)
	if updStatus != nil {
		errToSend, code := appErrors.GetErrAndCodeToSend(updStatus)
		responseTemplates.SendErrorMessage(w, errToSend, code)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (userHandler *UserHandler) UploadAvatar(w http.ResponseWriter, r *http.Request) {
	uploadedData, header, err := r.FormFile("avatar")
	if err != nil {
		responseTemplates.SendErrorMessage(w, err, http.StatusBadRequest)
		return
	}
	defer uploadedData.Close()

	uplErr := userHandler.userUsecase.UploadAvatar(r.Context(), uploadedData, header)
	if errors.Is(uplErr, usecase.BadAvatarSize) {
		responseTemplates.SendErrorMessage(w, usecase.BadAvatarSize, http.StatusBadRequest)
		return
	} else if errors.Is(uplErr, usecase.BadAvatarType) {
		responseTemplates.SendErrorMessage(w, usecase.BadAvatarType, http.StatusBadRequest)
		return
	} else if uplErr != nil {
		errToSend, code := appErrors.GetErrAndCodeToSend(uplErr)
		responseTemplates.SendErrorMessage(w, errToSend, code)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (userHandler *UserHandler) GetAvatar(w http.ResponseWriter, r *http.Request) {
	file, err := userHandler.userUsecase.GetAvatar(r.Context())
	if file == nil && err == nil {
		responseTemplates.SendErrorMessage(w, serverErrors.NO_DATA_FOUND, http.StatusNotFound)
		return
	} else if err != nil {
		errToSend, code := appErrors.GetErrAndCodeToSend(err)
		responseTemplates.SendErrorMessage(w, errToSend, code)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Write(file)
}
