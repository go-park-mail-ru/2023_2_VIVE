package http

import (
	"HnH/internal/appErrors"
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

	router.Handle("/users/{userID}/responses",
		middleware.AuthMiddleware(sessionUCase, http.HandlerFunc(handler.GetUserResponses))).
		Methods("GET")
}

func (responseHandler *ResponseHandler) sanitizeApplicants(applicants ...domain.ApiApplicant) []domain.ApiApplicant {
	result := make([]domain.ApiApplicant, 0, len(applicants))

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

func (responseHandler *ResponseHandler) sanitizeResponses(responses ...domain.ApiResponse) []domain.ApiResponse {
	result := make([]domain.ApiResponse, 0, len(responses))

	for _, resp := range responses {
		resp.OrganizationName = sanitizer.XSS.Sanitize(resp.OrganizationName)
		resp.VacancyName = sanitizer.XSS.Sanitize(resp.VacancyName)

		result = append(result, resp)
	}

	return result
}

func (responseHandler *ResponseHandler) CreateResponse(w http.ResponseWriter, r *http.Request) {
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

	createStatus := responseHandler.responseUsecase.RespondToVacancy(r.Context(), vacancyID, cvID)
	if createStatus != nil {
		errToSend, code := appErrors.GetErrAndCodeToSend(createStatus)
		responseTemplates.SendErrorMessage(w, errToSend, code)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (responseHandler *ResponseHandler) GetApplicants(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	vacancyID, convErr := strconv.Atoi(vars["vacancyID"])
	if convErr != nil {
		responseTemplates.SendErrorMessage(w, convErr, http.StatusBadRequest)
		return
	}

	applicantsList, err := responseHandler.responseUsecase.GetApplicantsList(r.Context(), vacancyID)
	if err != nil {
		errToSend, code := appErrors.GetErrAndCodeToSend(err)
		responseTemplates.SendErrorMessage(w, errToSend, code)
		return
	}

	sanitizedApplicants := responseHandler.sanitizeApplicants(applicantsList...)

	responseTemplates.MarshalAndSend(w, sanitizedApplicants)
}

func (responseHandler *ResponseHandler) GetUserResponses(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	userID, convErr := strconv.Atoi(vars["userID"])
	if convErr != nil {
		responseTemplates.SendErrorMessage(w, convErr, http.StatusBadRequest)
		return
	}

	responses, err := responseHandler.responseUsecase.GetUserResponses(r.Context(), userID)
	if err != nil {
		errToSend, code := appErrors.GetErrAndCodeToSend(err)
		responseTemplates.SendErrorMessage(w, errToSend, code)
	}

	sanitizedResponses := responseHandler.sanitizeResponses(responses...)
	responseTemplates.MarshalAndSend(w, sanitizedResponses)
}
