package http

import (
	"HnH/internal/delivery/http/middleware"
	"HnH/internal/domain"
	"HnH/internal/usecase"
	"HnH/pkg/responseTemplates"

	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type CVHandler struct {
	cvUsecase usecase.ICVUsecase
}

func NewCVHandler(router *mux.Router, cvUCase usecase.ICVUsecase, sessionUCase usecase.ISessionUsecase) {
	handler := &CVHandler{
		cvUsecase: cvUCase,
	}

	router.Handle("/cv/{cvID}",
		middleware.AuthMiddleware(sessionUCase, http.HandlerFunc(handler.GetCV))).
		Methods("GET")

	router.Handle("/current_user/cvs",
		middleware.AuthMiddleware(sessionUCase, http.HandlerFunc(handler.GetCVList))).
		Methods("GET")

	router.Handle("/current_user/cvs",
		middleware.JSONBodyValidationMiddleware(middleware.AuthMiddleware(sessionUCase, http.HandlerFunc(handler.AddNewCV)))).
		Methods("POST")

	router.Handle("/current_user/cvs/{cvID}",
		middleware.AuthMiddleware(sessionUCase, http.HandlerFunc(handler.GetCVOfUser))).
		Methods("GET")

	router.Handle("/current_user/cvs/{cvID}",
		middleware.JSONBodyValidationMiddleware(middleware.AuthMiddleware(sessionUCase, http.HandlerFunc(handler.UpdateCVOfUser)))).
		Methods("PUT")

	router.Handle("/current_user/cvs/{cvID}",
		middleware.AuthMiddleware(sessionUCase, http.HandlerFunc(handler.DeleteCVOfUser))).
		Methods("DELETE")
}

func (cvHandler *CVHandler) GetCV(w http.ResponseWriter, r *http.Request) {
	cookie, _ := r.Cookie("session")

	vars := mux.Vars(r)
	cvID, convErr := strconv.Atoi(vars["cvID"])
	if convErr != nil {
		responseTemplates.SendErrorMessage(w, convErr, http.StatusBadRequest)
		return
	}

	cv, err := cvHandler.cvUsecase.GetCVById(cookie.Value, cvID)
	if err != nil {
		responseTemplates.SendErrorMessage(w, err, http.StatusForbidden)
		return
	}

	responseTemplates.MarshalAndSend(w, *cv)
}

func (cvHandler *CVHandler) GetCVList(w http.ResponseWriter, r *http.Request) {
	cookie, _ := r.Cookie("session")

	cvs, err := cvHandler.cvUsecase.GetCVList(cookie.Value)
	if err != nil {
		responseTemplates.SendErrorMessage(w, err, http.StatusBadRequest)
		return
	}

	responseTemplates.MarshalAndSend(w, cvs)
}

func (cvHandler *CVHandler) AddNewCV(w http.ResponseWriter, r *http.Request) {
	cookie, _ := r.Cookie("session")

	defer r.Body.Close()

	cv := new(domain.CV)

	readErr := json.NewDecoder(r.Body).Decode(cv)
	if readErr != nil {
		responseTemplates.SendErrorMessage(w, readErr, http.StatusBadRequest)
		return
	}

	newCVID, addErr := cvHandler.cvUsecase.AddNewCV(cookie.Value, cv)
	if addErr != nil {
		responseTemplates.SendErrorMessage(w, addErr, http.StatusUnauthorized)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(fmt.Sprintf(`{"id":%d}`, newCVID)))
}

func (cvHandler *CVHandler) GetCVOfUser(w http.ResponseWriter, r *http.Request) {
	cookie, _ := r.Cookie("session")

	vars := mux.Vars(r)
	cvID, convErr := strconv.Atoi(vars["cvID"])
	if convErr != nil {
		responseTemplates.SendErrorMessage(w, convErr, http.StatusBadRequest)
		return
	}

	cv, err := cvHandler.cvUsecase.GetCVOfUserById(cookie.Value, cvID)
	if err != nil {
		responseTemplates.SendErrorMessage(w, err, http.StatusBadRequest)
		return
	}

	responseTemplates.MarshalAndSend(w, *cv)
}

func (cvHandler *CVHandler) UpdateCVOfUser(w http.ResponseWriter, r *http.Request) {
	cookie, _ := r.Cookie("session")

	vars := mux.Vars(r)
	cvID, convErr := strconv.Atoi(vars["cvID"])
	if convErr != nil {
		responseTemplates.SendErrorMessage(w, convErr, http.StatusBadRequest)
		return
	}

	updateInfo := new(domain.CV)

	decodeErr := json.NewDecoder(r.Body).Decode(updateInfo)
	if decodeErr != nil {
		sendErrorMessage(w, decodeErr, http.StatusBadRequest)
		return
	}

	udpErr := cvHandler.cvUsecase.UpdateCVOfUserById(cookie.Value, cvID, updateInfo)
	if udpErr != nil {
		responseTemplates.SendErrorMessage(w, udpErr, http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (cvHandler *CVHandler) DeleteCVOfUser(w http.ResponseWriter, r *http.Request) {
	cookie, _ := r.Cookie("session")

	vars := mux.Vars(r)
	cvID, convErr := strconv.Atoi(vars["cvID"])
	if convErr != nil {
		responseTemplates.SendErrorMessage(w, convErr, http.StatusBadRequest)
		return
	}

	deleteErr := cvHandler.cvUsecase.DeleteCVOfUserById(cookie.Value, cvID)
	if deleteErr != nil {
		responseTemplates.SendErrorMessage(w, deleteErr, http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}
