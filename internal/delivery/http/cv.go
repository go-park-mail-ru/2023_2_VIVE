package http

import (
	"HnH/internal/appErrors"
	"HnH/internal/domain"
	"HnH/internal/usecase"
	"HnH/pkg/contextUtils"
	"HnH/pkg/middleware"
	"HnH/pkg/responseTemplates"
	"HnH/pkg/sanitizer"
	"HnH/services/searchEngineService/searchEnginePB"

	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
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

	router.HandleFunc("/cvs/applicant/{applicantID}", handler.GetApplicantInfo).
		Methods("GET")

	router.HandleFunc("/cvs/search", handler.SearchCVs).
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

	router.Handle("/current_user/cvs/{cvID}/pdf",
		http.HandlerFunc(handler.GetCVsPDF)).
		Methods("GET")
}

func (cvHandler *CVHandler) sanitizeCVs(CVs ...domain.ApiCV) []domain.ApiCV {
	result := make([]domain.ApiCV, 0, len(CVs))

	for _, cv := range CVs {
		cv.ProfessionName = sanitizer.XSS.Sanitize(cv.ProfessionName)
		if cv.Description != nil {
			description := sanitizer.XSS.Sanitize(*cv.Description)
			cv.Description = &description
		}
		cv.FirstName = sanitizer.XSS.Sanitize(cv.FirstName)
		cv.LastName = sanitizer.XSS.Sanitize(cv.LastName)
		if cv.MiddleName != nil {
			middle_name := sanitizer.XSS.Sanitize(*cv.MiddleName)
			cv.MiddleName = &middle_name
		}
		if cv.Birthday != nil {
			birthday := sanitizer.XSS.Sanitize(*cv.Birthday)
			cv.Birthday = &birthday
		}
		if cv.Location != nil {
			location := sanitizer.XSS.Sanitize(*cv.Location)
			cv.Location = &location
		}

		result = append(result, cv)
	}

	return result
}

func (cvHandler *CVHandler) GetCV(w http.ResponseWriter, r *http.Request) {
	contextLogger := contextUtils.GetContextLogger(r.Context())
	vars := mux.Vars(r)
	cvID, convErr := strconv.Atoi(vars["cvID"])
	if convErr != nil {
		err := responseTemplates.SendErrorMessage(w, convErr, http.StatusBadRequest)
		if err != nil {
			contextLogger.WithFields(logrus.Fields{
				"error_msg": err,
			}).
				Error("could not send error message")
		}
		return
	}

	cv, err := cvHandler.cvUsecase.GetCVById(r.Context(), cvID)
	if err != nil {
		errToSend, code := appErrors.GetErrAndCodeToSend(err)
		err := responseTemplates.SendErrorMessage(w, errToSend, code)
		if err != nil {
			contextLogger.WithFields(logrus.Fields{
				"error_msg": err,
			}).
				Error("could not send error message")
		}
		return
	}

	sanitizedCV := cvHandler.sanitizeCVs(*cv)

	marshalAndSendErr := responseTemplates.MarshalAndSend(w, sanitizedCV[0])
	if marshalAndSendErr != nil {
		contextLogger.WithFields(
			logrus.Fields{
				"error_msg": marshalAndSendErr,
				"data":      sanitizedCV[0],
			},
		).
			Error("could not marshal and send data")
	}
}

func (cvHandler *CVHandler) sanitizeMetaCVs(metaCVs domain.ApiMetaCV) domain.ApiMetaCV {
	result := domain.ApiMetaCV{
		Filters: metaCVs.Filters,
		CVs: domain.ApiCVCount{
			Count: metaCVs.CVs.Count,
			CVs:   cvHandler.sanitizeCVs(metaCVs.CVs.CVs...),
		},
		// Count: metaCVs.Count,
		// CVs:   cvHandler.sanitizeCVs(metaCVs.CVs...),
	}
	return result
}

func (cvHandler *CVHandler) SearchCVs(w http.ResponseWriter, r *http.Request) {
	contextLogger := contextUtils.GetContextLogger(r.Context())
	query := r.URL.Query()
	contextLogger.WithFields(logrus.Fields{
		"query": query.Encode(),
	}).
		Debug("got search request with query")

	// options := searchEnginePB.SearchOptions{}
	queryOptions := make(map[string]*searchEnginePB.SearchOptionValues)
	for optionName, values := range query {
		contextLogger.WithFields(logrus.Fields{
			"option_name":   optionName,
			"option_values": values,
		}).
			Debug("parsing options")
		optionsValues := searchEnginePB.SearchOptionValues{
			Values: values,
		}
		queryOptions[optionName] = &optionsValues
	}
	options := searchEnginePB.SearchOptions{Options: queryOptions}

	metaCVs, getErr := cvHandler.cvUsecase.SearchCVs(r.Context(), &options)

	if getErr != nil {
		errToSend, code := appErrors.GetErrAndCodeToSend(getErr)
		err := responseTemplates.SendErrorMessage(w, errToSend, code)
		if err != nil {
			contextLogger.WithFields(logrus.Fields{
				"error_msg": err,
			}).
				Error("could not send error message")
		}
		return
	}

	sanitizedMetaVacancies := cvHandler.sanitizeMetaCVs(metaCVs)

	marshalAndSendErr := responseTemplates.MarshalAndSend(w, sanitizedMetaVacancies)
	if marshalAndSendErr != nil {
		contextLogger.WithFields(
			logrus.Fields{
				"error_msg": marshalAndSendErr,
				"data":      sanitizedMetaVacancies,
			},
		).
			Error("could not marshal and send data")
	}
}

func (cvHandler *CVHandler) GetCVList(w http.ResponseWriter, r *http.Request) {
	contextLogger := contextUtils.GetContextLogger(r.Context())
	cvs, err := cvHandler.cvUsecase.GetCVList(r.Context())
	if err != nil {
		errToSend, code := appErrors.GetErrAndCodeToSend(err)
		err := responseTemplates.SendErrorMessage(w, errToSend, code)
		if err != nil {
			contextLogger.WithFields(logrus.Fields{
				"error_msg": err,
			}).
				Error("could not send error message")
		}
		return
	}

	sanitizedCVs := cvHandler.sanitizeCVs(cvs...)

	marshalAndSendErr := responseTemplates.MarshalAndSend(w, sanitizedCVs)
	if marshalAndSendErr != nil {
		contextLogger.WithFields(
			logrus.Fields{
				"error_msg": marshalAndSendErr,
				"data":      sanitizedCVs,
			},
		).
			Error("could not marshal and send data")
	}
}

func (cvHandler *CVHandler) AddNewCV(w http.ResponseWriter, r *http.Request) {
	contextLogger := contextUtils.GetContextLogger(r.Context())
	defer r.Body.Close()

	apiCV := new(domain.ApiCV)

	readErr := json.NewDecoder(r.Body).Decode(apiCV)
	if readErr != nil {
		err := responseTemplates.SendErrorMessage(w, readErr, http.StatusBadRequest)
		if err != nil {
			contextLogger.WithFields(logrus.Fields{
				"error_msg": err,
			}).
				Error("could not send error message")
		}
		return
	}
	// fmt.Println(cv)
	// bdCV := apiCV.ToDb()

	newCVID, addErr := cvHandler.cvUsecase.AddNewCV(r.Context(), apiCV)
	if addErr != nil {
		errToSend, code := appErrors.GetErrAndCodeToSend(addErr)
		err := responseTemplates.SendErrorMessage(w, errToSend, code)
		if err != nil {
			contextLogger.WithFields(logrus.Fields{
				"error_msg": err,
			}).
				Error("could not send error message")
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	_, wErr := w.Write([]byte(fmt.Sprintf(`{"id":%d}`, newCVID)))
	contextLogger.WithFields(logrus.Fields{
		"err_msg": wErr,
		"data":    fmt.Sprintf(`{"id":%d}`, newCVID),
	}).
		Error("could not send data")
}

func (cvHandler *CVHandler) GetCVOfUser(w http.ResponseWriter, r *http.Request) {
	contextLogger := contextUtils.GetContextLogger(r.Context())
	vars := mux.Vars(r)
	cvID, convErr := strconv.Atoi(vars["cvID"])
	if convErr != nil {
		err := responseTemplates.SendErrorMessage(w, convErr, http.StatusBadRequest)
		if err != nil {
			contextLogger.WithFields(logrus.Fields{
				"error_msg": err,
			}).
				Error("could not send error message")
		}
		return
	}

	cv, err := cvHandler.cvUsecase.GetCVOfUserById(r.Context(), cvID)
	if err != nil {
		errToSend, code := appErrors.GetErrAndCodeToSend(err)
		err := responseTemplates.SendErrorMessage(w, errToSend, code)
		if err != nil {
			contextLogger.WithFields(logrus.Fields{
				"error_msg": err,
			}).
				Error("could not send error message")
		}
		return
	}

	sanitizedCV := cvHandler.sanitizeCVs(*cv)

	marshalAndSendErr := responseTemplates.MarshalAndSend(w, sanitizedCV[0])
	if marshalAndSendErr != nil {
		contextLogger.WithFields(
			logrus.Fields{
				"error_msg": marshalAndSendErr,
				"data":      sanitizedCV[0],
			},
		).
			Error("could not marshal and send data")
	}
}

func (cvHandler *CVHandler) UpdateCVOfUser(w http.ResponseWriter, r *http.Request) {
	contextLogger := contextUtils.GetContextLogger(r.Context())
	vars := mux.Vars(r)
	cvID, convErr := strconv.Atoi(vars["cvID"])
	if convErr != nil {
		err := responseTemplates.SendErrorMessage(w, convErr, http.StatusBadRequest)
		if err != nil {
			contextLogger.WithFields(logrus.Fields{
				"error_msg": err,
			}).
				Error("could not send error message")
		}
		return
	}

	cv := new(domain.ApiCV)

	decodeErr := json.NewDecoder(r.Body).Decode(cv)
	if decodeErr != nil {
		err := responseTemplates.SendErrorMessage(w, decodeErr, http.StatusBadRequest)
		if err != nil {
			contextLogger.WithFields(logrus.Fields{
				"error_msg": err,
			}).
				Error("could not send error message")
		}
		return
	}

	updErr := cvHandler.cvUsecase.UpdateCVOfUserById(r.Context(), cvID, cv)
	if updErr != nil {
		errToSend, code := appErrors.GetErrAndCodeToSend(updErr)
		err := responseTemplates.SendErrorMessage(w, errToSend, code)
		if err != nil {
			contextLogger.WithFields(logrus.Fields{
				"error_msg": err,
			}).
				Error("could not send error message")
		}
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (cvHandler *CVHandler) DeleteCVOfUser(w http.ResponseWriter, r *http.Request) {
	contextLogger := contextUtils.GetContextLogger(r.Context())
	vars := mux.Vars(r)
	cvID, convErr := strconv.Atoi(vars["cvID"])
	if convErr != nil {
		err := responseTemplates.SendErrorMessage(w, convErr, http.StatusBadRequest)
		if err != nil {
			contextLogger.WithFields(logrus.Fields{
				"error_msg": err,
			}).
				Error("could not send error message")
		}
		return
	}

	deleteErr := cvHandler.cvUsecase.DeleteCVOfUserById(r.Context(), cvID)
	if deleteErr != nil {
		errToSend, code := appErrors.GetErrAndCodeToSend(deleteErr)
		err := responseTemplates.SendErrorMessage(w, errToSend, code)
		if err != nil {
			contextLogger.WithFields(logrus.Fields{
				"error_msg": err,
			}).
				Error("could not send error message")
		}
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (cvHandler *CVHandler) GetApplicantInfo(w http.ResponseWriter, r *http.Request) {
	contextLogger := contextUtils.GetContextLogger(r.Context())
	vars := mux.Vars(r)
	appID, convErr := strconv.Atoi(vars["applicantID"])
	if convErr != nil {
		err := responseTemplates.SendErrorMessage(w, convErr, http.StatusBadRequest)
		if err != nil {
			contextLogger.WithFields(logrus.Fields{
				"error_msg": err,
			}).
				Error("could not send error message")
		}
		return
	}

	info, err := cvHandler.cvUsecase.GetApplicantInfo(r.Context(), appID)
	if err != nil {
		errToSend, code := appErrors.GetErrAndCodeToSend(err)
		err := responseTemplates.SendErrorMessage(w, errToSend, code)
		if err != nil {
			contextLogger.WithFields(logrus.Fields{
				"error_msg": err,
			}).
				Error("could not send error message")
		}
		return
	}

	info.CVs = cvHandler.sanitizeCVs(info.CVs...)

	marshalAndSendErr := responseTemplates.MarshalAndSend(w, info)
	if marshalAndSendErr != nil {
		contextLogger.WithFields(
			logrus.Fields{
				"error_msg": marshalAndSendErr,
				"data":      info,
			},
		).
			Error("could not marshal and send data")
	}
}

func (cvHandler *CVHandler) GetCVsPDF(w http.ResponseWriter, r *http.Request) {
	contextLogger := contextUtils.GetContextLogger(r.Context())
	vars := mux.Vars(r)
	cvID, convErr := strconv.Atoi(vars["cvID"])
	if convErr != nil {
		err := responseTemplates.SendErrorMessage(w, convErr, http.StatusBadRequest)
		if err != nil {
			contextLogger.WithFields(logrus.Fields{
				"error_msg": err,
			}).
				Error("could not send error message")
		}
		return
	}

	pdf, err := cvHandler.cvUsecase.GenerateCVsPDF(r.Context(), cvID)
	if err != nil {
		errToSend, code := appErrors.GetErrAndCodeToSend(err)
		err := responseTemplates.SendErrorMessage(w, errToSend, code)
		if err != nil {
			contextLogger.WithFields(logrus.Fields{
				"error_msg": err,
			}).
				Error("could not send error message")
		}
		return
	}

	pdfErr := responseTemplates.SendPDF(w, pdf, "file") // TODO: rename file
	if pdfErr != nil {
		contextLogger.WithFields(logrus.Fields{
			"error_msg": pdfErr,
		}).
			Error("could not send pdf")
	}
}
