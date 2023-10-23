package http

import (
	"HnH/internal/usecase"
	"HnH/pkg/serverErrors"

	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type ResponseHandler struct {
	responseUsecase usecase.IResponseUsecase
}

func NewResponseHandler(router *mux.Router, responseUCase usecase.IResponseUsecase) {
	handler := &ResponseHandler{
		responseUsecase: responseUCase,
	}

	router.HandleFunc("/vacancies/{vacancyID}/respond/{cvID}", handler.CreateResponse).Methods("POST")
	router.HandleFunc("/vacancies/{vacancyID}/applicants", handler.GetApplicants).Methods("GET")
}

func (responseHandler *ResponseHandler) CreateResponse(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session")

	if errors.Is(err, http.ErrNoCookie) {
		sendErrorMessage(w, serverErrors.NO_COOKIE, http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)

	vacancyID, convErr := strconv.Atoi(vars["vacancyID"])
	if convErr != nil {
		sendErrorMessage(w, convErr, http.StatusBadRequest)
		return
	}

	cvID, convErr := strconv.Atoi(vars["cvID"])
	if convErr != nil {
		sendErrorMessage(w, convErr, http.StatusBadRequest)
		return
	}

	createStatus := responseHandler.responseUsecase.RespondToVacancy(cookie.Value, vacancyID, cvID)
	if createStatus != nil {
		sendErrorMessage(w, createStatus, http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (responseHandler *ResponseHandler) GetApplicants(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session")

	if errors.Is(err, http.ErrNoCookie) {
		sendErrorMessage(w, serverErrors.NO_COOKIE, http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)

	vacancyID, convErr := strconv.Atoi(vars["vacancyID"])
	if convErr != nil {
		sendErrorMessage(w, convErr, http.StatusBadRequest)
		return
	}

	applicantsList, err := responseHandler.responseUsecase.GetApplicantsList(cookie.Value, vacancyID)
	if err != nil {
		sendErrorMessage(w, err, http.StatusForbidden)
		return
	}

	js, err := json.Marshal(applicantsList)
	if err != nil {
		sendErrorMessage(w, serverErrors.INTERNAL_SERVER_ERROR, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}
