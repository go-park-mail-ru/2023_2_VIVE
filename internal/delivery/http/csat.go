package http

import (
	"HnH/internal/appErrors"
	"HnH/internal/usecase"
	"HnH/pkg/contextUtils"
	"HnH/pkg/middleware"
	"HnH/pkg/responseTemplates"
	"HnH/services/csat/csatPB"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
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
	contextLogger := contextUtils.GetContextLogger(r.Context())
	questionList, err := handler.csatUsecase.GetQuestions(r.Context())
	if err != nil {
		errToSend, code := appErrors.GetErrAndCodeToSend(err)
		sendErr := responseTemplates.SendErrorMessage(w, errToSend, code)
		if sendErr != nil {
			contextLogger.WithFields(logrus.Fields{
				"err_msg":       sendErr,
				"error_to_send": errToSend,
			}).
				Error("could not send error")
		}
		return
	}

	marshalErr := responseTemplates.MarshalAndSend(w, questionList)
	if marshalErr != nil {
		contextLogger.WithFields(logrus.Fields{
			"err_msg": marshalErr,
			"data":    questionList,
		}).
			Error("could not marshal and send data")
	}
}

func (handler *CsatHandler) RegisterAnswer(w http.ResponseWriter, r *http.Request) {
	contextLogger := contextUtils.GetContextLogger(r.Context())
	answer := &csatPB.Answer{}
	readErr := json.NewDecoder(r.Body).Decode(answer)
	if readErr != nil {
		sendErr := responseTemplates.SendErrorMessage(w, readErr, http.StatusBadRequest)
		if sendErr != nil {
			contextLogger.WithFields(logrus.Fields{
				"err_msg":       sendErr,
				"error_to_send": readErr,
			}).
				Error("could not send error")
		}
		return
	}

	err := handler.csatUsecase.RegisterAnswer(r.Context(), answer)
	if err != nil {
		errToSend, code := appErrors.GetErrAndCodeToSend(err)
		sendErr := responseTemplates.SendErrorMessage(w, errToSend, code)
		if sendErr != nil {
			contextLogger.WithFields(logrus.Fields{
				"err_msg":       sendErr,
				"error_to_send": errToSend,
			}).
				Error("could not send error")
		}
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (handler *CsatHandler) GetStatistics(w http.ResponseWriter, r *http.Request) {
	contextLogger := contextUtils.GetContextLogger(r.Context())
	statistics, err := handler.csatUsecase.GetStatistic(r.Context())
	if err != nil {
		errToSend, code := appErrors.GetErrAndCodeToSend(err)
		sendErr := responseTemplates.SendErrorMessage(w, errToSend, code)
		if sendErr != nil {
			contextLogger.WithFields(logrus.Fields{
				"err_msg":       sendErr,
				"error_to_send": errToSend,
			}).
				Error("could not send error")
		}
		return
	}

	marshalErr := responseTemplates.MarshalAndSend(w, statistics)
	if marshalErr != nil {
		contextLogger.WithFields(logrus.Fields{
			"err_msg": marshalErr,
			"data":    statistics,
		}).
			Error("could not marshal and send data")
	}
}
