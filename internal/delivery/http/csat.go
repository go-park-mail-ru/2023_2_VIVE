package http

import (
	"HnH/internal/appErrors"
	"HnH/internal/delivery/http/middleware"
	"HnH/internal/usecase"
	"HnH/pkg/responseTemplates"
	"HnH/services/csat/csatPB"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

type CsatHandler struct {
	csatUsecase usecase.ICsatUsecase
}

func NewCsatHandler(router *mux.Router, csatUCase usecase.ICsatUsecase, sessionUCase usecase.ISessionUsecase) {
	handler := &CsatHandler{
		csatUsecase: csatUCase,
	}

	router.Handle("/statistics/questions",
		middleware.AuthMiddleware(sessionUCase, http.HandlerFunc(handler.GetQuestions))).
		Methods("GET")

	router.Handle("/statistics/questions",
		middleware.JSONBodyValidationMiddleware(middleware.AuthMiddleware(sessionUCase, http.HandlerFunc(handler.RegisterAnswer)))).
		Methods("POST")

	router.Handle("/statistics",
		middleware.AuthMiddleware(sessionUCase, http.HandlerFunc(handler.GetStatistics))).
		Methods("GET")
}

func (handler *CsatHandler) GetQuestions(w http.ResponseWriter, r *http.Request) {
	questionList, err := handler.csatUsecase.GetQuestions(r.Context())
	if err != nil {
		errToSend, code := appErrors.GetErrAndCodeToSend(err)
		responseTemplates.SendErrorMessage(w, errToSend, code)
		return
	}

	responseTemplates.MarshalAndSend(w, questionList)
}

func (handler *CsatHandler) RegisterAnswer(w http.ResponseWriter, r *http.Request) {
	answer := &csatPB.Answer{}
	readErr := json.NewDecoder(r.Body).Decode(answer)
	if readErr != nil {
		responseTemplates.SendErrorMessage(w, readErr, http.StatusBadRequest)
		return
	}

	err := handler.csatUsecase.RegisterAnswer(r.Context(), answer)
	if err != nil {
		errToSend, code := appErrors.GetErrAndCodeToSend(err)
		responseTemplates.SendErrorMessage(w, errToSend, code)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (handler *CsatHandler) GetStatistics(w http.ResponseWriter, r *http.Request) {
	statistics, err := handler.csatUsecase.GetStatistic(r.Context())
	if err != nil {
		errToSend, code := appErrors.GetErrAndCodeToSend(err)
		responseTemplates.SendErrorMessage(w, errToSend, code)
		return
	}

	responseTemplates.MarshalAndSend(w, statistics)
}
