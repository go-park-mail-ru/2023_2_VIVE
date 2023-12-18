package http

import (
	"HnH/internal/appErrors"
	"HnH/internal/domain"
	"HnH/internal/usecase"
	"HnH/pkg/contextUtils"
	"HnH/pkg/middleware"
	"HnH/pkg/responseTemplates"
	"HnH/pkg/sanitizer"
	"HnH/pkg/serverErrors"
	"errors"

	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
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
		middleware.AuthMiddleware(sessionUCase, http.HandlerFunc(handler.GetUserAvatar))).
		Methods("GET")

	fileServer := http.FileServer(http.Dir("./"))
	router.PathPrefix("/image").Handler(http.StripPrefix("/image", fileServer)).
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
	contextLogger := contextUtils.GetContextLogger(r.Context())
	defer r.Body.Close()

	newUser := new(domain.ApiUser)

	err := json.NewDecoder(r.Body).Decode(newUser)
	if err != nil {
		sendErr := responseTemplates.SendErrorMessage(w, err, http.StatusBadRequest)
		if sendErr != nil {
			contextLogger.WithFields(logrus.Fields{
				"error_msg":     sendErr,
				"error_to_send": err,
			}).
				Error("could not send error message")
		}
		return
	}

	expiryTime := time.Now().Add(10 * time.Hour)

	sessionID, err := userHandler.userUsecase.SignUp(r.Context(), newUser, expiryTime.Unix())
	if err != nil {
		errToSend, code := appErrors.GetErrAndCodeToSend(err)
		sendErr := responseTemplates.SendErrorMessage(w, errToSend, code)
		if sendErr != nil {
			contextLogger.WithFields(logrus.Fields{
				"error_msg":     sendErr,
				"error_to_send": errToSend,
			}).
				Error("could not send error message")
		}
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
	contextLogger := contextUtils.GetContextLogger(r.Context())
	// logger := r.Context().Value(middleware.LOGGER_KEY).(*logrus.Entry)
	// logger.Info("GOT REQUEST")

	// requesID, requestID := middleware.GetRequestIDCtx(r.Context())
	// logging.Logger.WithField(requesID, requestID).
	// 	Info("got request")

	// logging.Logger.Info("got request")

	user, err := userHandler.userUsecase.GetInfo(r.Context())
	if err != nil {
		errToSend, code := appErrors.GetErrAndCodeToSend(err)
		sendErr := responseTemplates.SendErrorMessage(w, errToSend, code)
		if sendErr != nil {
			contextLogger.WithFields(logrus.Fields{
				"error_msg":     sendErr,
				"error_to_send": errToSend,
			}).
				Error("could not send error message")
		}
		return
	}

	userHandler.sanitizeUser(user)

	marshalErr := responseTemplates.MarshalAndSend(w, *user)
	if marshalErr != nil {
		contextLogger.WithFields(logrus.Fields{
			"err_msg": marshalErr,
			"data":    user,
		}).
			Error("could not marshal and send data")
	}
}

func (userHandler *UserHandler) UpdateInfo(w http.ResponseWriter, r *http.Request) {
	contextLogger := contextUtils.GetContextLogger(r.Context())
	defer r.Body.Close()

	updateInfo := new(domain.UserUpdate)

	decodeErr := json.NewDecoder(r.Body).Decode(updateInfo)
	if decodeErr != nil {
		sendErr := responseTemplates.SendErrorMessage(w, decodeErr, http.StatusBadRequest)
		if sendErr != nil {
			contextLogger.WithFields(logrus.Fields{
				"error_msg":     sendErr,
				"error_to_send": decodeErr,
			}).
				Error("could not send error message")
		}
		return
	}

	updStatus := userHandler.userUsecase.UpdateInfo(r.Context(), updateInfo)
	if updStatus != nil {
		errToSend, code := appErrors.GetErrAndCodeToSend(updStatus)
		sendErr := responseTemplates.SendErrorMessage(w, errToSend, code)
		if sendErr != nil {
			contextLogger.WithFields(logrus.Fields{
				"error_msg":     sendErr,
				"error_to_send": errToSend,
			}).
				Error("could not send error message")
		}
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (userHandler *UserHandler) UploadAvatar(w http.ResponseWriter, r *http.Request) {
	contextLogger := contextUtils.GetContextLogger(r.Context())
	uploadedData, header, err := r.FormFile("avatar")
	if err != nil {
		sendErr := responseTemplates.SendErrorMessage(w, err, http.StatusBadRequest)
		if sendErr != nil {
			contextLogger.WithFields(logrus.Fields{
				"error_msg":     sendErr,
				"error_to_send": err,
			}).
				Error("could not send error message")
		}
		return
	}
	defer uploadedData.Close()

	uplErr := userHandler.userUsecase.UploadAvatar(r.Context(), uploadedData, header)
	if errors.Is(uplErr, usecase.BadAvatarSize) {
		sendErr := responseTemplates.SendErrorMessage(w, usecase.BadAvatarSize, http.StatusBadRequest)
		if sendErr != nil {
			contextLogger.WithFields(logrus.Fields{
				"error_msg":     sendErr,
				"error_to_send": usecase.BadAvatarSize,
			}).
				Error("could not send error message")
		}
		return
	} else if errors.Is(uplErr, usecase.BadAvatarType) {
		sendErr := responseTemplates.SendErrorMessage(w, usecase.BadAvatarType, http.StatusBadRequest)
		if sendErr != nil {
			contextLogger.WithFields(logrus.Fields{
				"error_msg":     sendErr,
				"error_to_send": usecase.BadAvatarType,
			}).
				Error("could not send error message")
		}
		return
	} else if uplErr != nil {
		errToSend, code := appErrors.GetErrAndCodeToSend(uplErr)
		sendErr := responseTemplates.SendErrorMessage(w, errToSend, code)
		if sendErr != nil {
			contextLogger.WithFields(logrus.Fields{
				"error_msg":     sendErr,
				"error_to_send": errToSend,
			}).
				Error("could not send error message")
		}
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (userHandler *UserHandler) GetUserAvatar(w http.ResponseWriter, r *http.Request) {
	contextLogger := contextUtils.GetContextLogger(r.Context())
	file, err := userHandler.userUsecase.GetUserAvatar(r.Context())
	if file == nil && err == nil {
		sendErr := responseTemplates.SendErrorMessage(w, serverErrors.NO_DATA_FOUND, http.StatusNotFound)
		if sendErr != nil {
			contextLogger.WithFields(logrus.Fields{
				"error_msg":     sendErr,
				"error_to_send": serverErrors.NO_DATA_FOUND,
			}).
				Error("could not send error message")
		}
		return
	} else if err != nil {
		errToSend, code := appErrors.GetErrAndCodeToSend(err)
		sendErr := responseTemplates.SendErrorMessage(w, errToSend, code)
		if sendErr != nil {
			contextLogger.WithFields(logrus.Fields{
				"error_msg":     sendErr,
				"error_to_send": errToSend,
			}).
				Error("could not send error message")
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/octet-stream")
	_, wErr := w.Write(file)
	if wErr != nil {
		contextLogger.WithFields(logrus.Fields{
			"error_msg": wErr,
		}).
			Error("could not send avatar file")
	}
}
