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
	vars := mux.Vars(r)
	cvID, convErr := strconv.Atoi(vars["cvID"])
	if convErr != nil {
		responseTemplates.SendErrorMessage(w, convErr, http.StatusBadRequest)
		return
	}

	cv, err := cvHandler.cvUsecase.GetCVById(r.Context(), cvID)
	if err != nil {
		errToSend, code := appErrors.GetErrAndCodeToSend(err)
		responseTemplates.SendErrorMessage(w, errToSend, code)
		return
	}

	sanitizedCV := cvHandler.sanitizeCVs(*cv)

	responseTemplates.MarshalAndSend(w, sanitizedCV[0])
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
		responseTemplates.SendErrorMessage(w, errToSend, code)
		return
	}

	sanitizedMetaVacancies := cvHandler.sanitizeMetaCVs(metaCVs)

	responseTemplates.MarshalAndSend(w, sanitizedMetaVacancies)

	w.WriteHeader(http.StatusOK)
}

func (cvHandler *CVHandler) GetCVList(w http.ResponseWriter, r *http.Request) {
	cvs, err := cvHandler.cvUsecase.GetCVList(r.Context())
	if err != nil {
		errToSend, code := appErrors.GetErrAndCodeToSend(err)
		responseTemplates.SendErrorMessage(w, errToSend, code)
		return
	}

	sanitizedCVs := cvHandler.sanitizeCVs(cvs...)

	responseTemplates.MarshalAndSend(w, sanitizedCVs)
}

func (cvHandler *CVHandler) AddNewCV(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	apiCV := new(domain.ApiCV)

	readErr := json.NewDecoder(r.Body).Decode(apiCV)
	if readErr != nil {
		responseTemplates.SendErrorMessage(w, readErr, http.StatusBadRequest)
		return
	}
	// fmt.Println(cv)
	// bdCV := apiCV.ToDb()

	newCVID, addErr := cvHandler.cvUsecase.AddNewCV(r.Context(), apiCV)
	if addErr != nil {
		errToSend, code := appErrors.GetErrAndCodeToSend(addErr)
		responseTemplates.SendErrorMessage(w, errToSend, code)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(fmt.Sprintf(`{"id":%d}`, newCVID)))
}

func (cvHandler *CVHandler) GetCVOfUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	cvID, convErr := strconv.Atoi(vars["cvID"])
	if convErr != nil {
		responseTemplates.SendErrorMessage(w, convErr, http.StatusBadRequest)
		return
	}

	cv, err := cvHandler.cvUsecase.GetCVOfUserById(r.Context(), cvID)
	if err != nil {
		errToSend, code := appErrors.GetErrAndCodeToSend(err)
		responseTemplates.SendErrorMessage(w, errToSend, code)
		return
	}

	sanitizedCV := cvHandler.sanitizeCVs(*cv)

	responseTemplates.MarshalAndSend(w, sanitizedCV[0])
}

func (cvHandler *CVHandler) UpdateCVOfUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	cvID, convErr := strconv.Atoi(vars["cvID"])
	if convErr != nil {
		responseTemplates.SendErrorMessage(w, convErr, http.StatusBadRequest)
		return
	}

	cv := new(domain.ApiCV)

	decodeErr := json.NewDecoder(r.Body).Decode(cv)
	if decodeErr != nil {
		responseTemplates.SendErrorMessage(w, decodeErr, http.StatusBadRequest)
		return
	}

	updErr := cvHandler.cvUsecase.UpdateCVOfUserById(r.Context(), cvID, cv)
	if updErr != nil {
		errToSend, code := appErrors.GetErrAndCodeToSend(updErr)
		responseTemplates.SendErrorMessage(w, errToSend, code)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (cvHandler *CVHandler) DeleteCVOfUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	cvID, convErr := strconv.Atoi(vars["cvID"])
	if convErr != nil {
		responseTemplates.SendErrorMessage(w, convErr, http.StatusBadRequest)
		return
	}

	deleteErr := cvHandler.cvUsecase.DeleteCVOfUserById(r.Context(), cvID)
	if deleteErr != nil {
		errToSend, code := appErrors.GetErrAndCodeToSend(deleteErr)
		responseTemplates.SendErrorMessage(w, errToSend, code)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (cvHandler *CVHandler) GetApplicantInfo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	appID, convErr := strconv.Atoi(vars["applicantID"])
	if convErr != nil {
		responseTemplates.SendErrorMessage(w, convErr, http.StatusBadRequest)
		return
	}

	info, err := cvHandler.cvUsecase.GetApplicantInfo(r.Context(), appID)
	if err != nil {
		errToSend, code := appErrors.GetErrAndCodeToSend(err)
		responseTemplates.SendErrorMessage(w, errToSend, code)
		return
	}

	info.CVs = cvHandler.sanitizeCVs(info.CVs...)

	responseTemplates.MarshalAndSend(w, info)
}
