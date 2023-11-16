package http

import (
	"HnH/internal/delivery/http/middleware"
	"HnH/internal/domain"
	"HnH/internal/usecase"
	"HnH/pkg/responseTemplates"
	"HnH/pkg/sanitizer"

	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type ResponseHandler struct {
	responseUsecase usecase.IResponseUsecase
}

func NewResponseHandler(router *mux.Router, responseUCase usecase.IResponseUsecase, sessionUCase usecase.ISessionUsecase) {
	handler := &ResponseHandler{
		responseUsecase: responseUCase,
	}

	router.Handle("/vacancies/{vacancyID}/respond/{cvID}",
		middleware.AuthMiddleware(sessionUCase, http.HandlerFunc(handler.CreateResponse))).
		Methods("POST")

	router.Handle("/vacancies/{vacancyID}/applicants",
		middleware.AuthMiddleware(sessionUCase, http.HandlerFunc(handler.GetApplicants))).
		Methods("GET")
}

func (responseHandler *ResponseHandler) sanitizeApplicants(applicants ...domain.ApplicantInfo) []domain.ApplicantInfo {
	result := make([]domain.ApplicantInfo, 0, len(applicants))

	for _, app := range applicants {
		app.FirstName = sanitizer.XSS.Sanitize(app.FirstName)
		app.LastName = sanitizer.XSS.Sanitize(app.LastName)

		for i, skill := range app.Skills {
			app.Skills[i] = sanitizer.XSS.Sanitize(skill)
		}

		result = append(result, app)
	}

	return result
}

func (responseHandler *ResponseHandler) CreateResponse(w http.ResponseWriter, r *http.Request) {
	cookie, _ := r.Cookie("session")

	vars := mux.Vars(r)

	vacancyID, convErr := strconv.Atoi(vars["vacancyID"])
	if convErr != nil {
		responseTemplates.SendErrorMessage(w, convErr, http.StatusBadRequest)
		return
	}

	cvID, convErr := strconv.Atoi(vars["cvID"])
	if convErr != nil {
		responseTemplates.SendErrorMessage(w, convErr, http.StatusBadRequest)
		return
	}

	createStatus := responseHandler.responseUsecase.RespondToVacancy(cookie.Value, vacancyID, cvID)
	if createStatus != nil {
		responseTemplates.SendErrorMessage(w, createStatus, http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (responseHandler *ResponseHandler) GetApplicants(w http.ResponseWriter, r *http.Request) {
	cookie, _ := r.Cookie("session")

	vars := mux.Vars(r)

	vacancyID, convErr := strconv.Atoi(vars["vacancyID"])
	if convErr != nil {
		responseTemplates.SendErrorMessage(w, convErr, http.StatusBadRequest)
		return
	}

	applicantsList, err := responseHandler.responseUsecase.GetApplicantsList(cookie.Value, vacancyID)
	if err != nil {
		responseTemplates.SendErrorMessage(w, err, http.StatusForbidden)
		return
	}

	sanitizedApplicants := responseHandler.sanitizeApplicants(applicantsList...)

	responseTemplates.MarshalAndSend(w, sanitizedApplicants)
}
