package http

import (
	"HnH/internal/delivery/http/middleware"
	"HnH/internal/domain"
	"HnH/internal/usecase"
	"HnH/pkg/contextUtils"
	"HnH/pkg/responseTemplates"
	"HnH/pkg/sanitizer"

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
	cookie, _ := r.Cookie("session")

	vars := mux.Vars(r)
	cvID, convErr := strconv.Atoi(vars["cvID"])
	if convErr != nil {
		responseTemplates.SendErrorMessage(w, convErr, http.StatusBadRequest)
		return
	}

	cv, err := cvHandler.cvUsecase.GetCVById(r.Context(), cookie.Value, cvID)
	if err != nil {
		responseTemplates.SendErrorMessage(w, err, http.StatusForbidden)
		return
	}

	sanitizedCV := cvHandler.sanitizeCVs(*cv)

	responseTemplates.MarshalAndSend(w, sanitizedCV[0])
}

func (cvHandler *CVHandler) sanitizeMetaCVs(metaCVs domain.ApiMetaCV) domain.ApiMetaCV {
	result := domain.ApiMetaCV{
		Count: metaCVs.Count,
		CVs:   cvHandler.sanitizeCVs(metaCVs.CVs...),
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
	searchQuery := query.Get(SEARCH_QUERY_KEY)

	pageNumStr := query.Get(PAGE_NUM_QUERY_KEY)
	pageNum, convErr := strconv.ParseInt(pageNumStr, 10, 64)
	if convErr != nil {
		responseTemplates.SendErrorMessage(w, ErrWrongQueryParam, http.StatusBadRequest)
		return
	}

	resultsPerPageStr := query.Get(RESULTS_PER_PAGE_QUERY_KEY)
	resultsPerPage, convErr := strconv.ParseInt(resultsPerPageStr, 10, 64)
	if convErr != nil {
		responseTemplates.SendErrorMessage(w, ErrWrongQueryParam, http.StatusBadRequest)
		return
	}

	metaCVs, getErr := cvHandler.cvUsecase.SearchCVs(
		r.Context(),
		searchQuery,
		pageNum,
		resultsPerPage,
	)

	if getErr != nil {
		responseTemplates.SendErrorMessage(w, getErr, http.StatusBadRequest)
		return
	}

	sanitizedMetaVacancies := cvHandler.sanitizeMetaCVs(metaCVs)

	responseTemplates.MarshalAndSend(w, sanitizedMetaVacancies)

	w.WriteHeader(http.StatusOK)
}

func (cvHandler *CVHandler) GetCVList(w http.ResponseWriter, r *http.Request) {
	cookie, _ := r.Cookie("session")

	cvs, err := cvHandler.cvUsecase.GetCVList(r.Context(), cookie.Value)
	if err != nil {
		responseTemplates.SendErrorMessage(w, err, http.StatusBadRequest)
		return
	}

	sanitizedCVs := cvHandler.sanitizeCVs(cvs...)

	responseTemplates.MarshalAndSend(w, sanitizedCVs)
}

func (cvHandler *CVHandler) AddNewCV(w http.ResponseWriter, r *http.Request) {
	cookie, _ := r.Cookie("session")

	defer r.Body.Close()

	apiCV := new(domain.ApiCV)

	readErr := json.NewDecoder(r.Body).Decode(apiCV)
	if readErr != nil {
		responseTemplates.SendErrorMessage(w, readErr, http.StatusBadRequest)
		return
	}
	// fmt.Println(cv)
	// bdCV := apiCV.ToDb()

	newCVID, addErr := cvHandler.cvUsecase.AddNewCV(r.Context(), cookie.Value, apiCV)
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

	cv, err := cvHandler.cvUsecase.GetCVOfUserById(r.Context(), cookie.Value, cvID)
	if err != nil {
		responseTemplates.SendErrorMessage(w, err, http.StatusBadRequest)
		return
	}

	sanitizedCV := cvHandler.sanitizeCVs(*cv)

	responseTemplates.MarshalAndSend(w, sanitizedCV[0])
}

func (cvHandler *CVHandler) UpdateCVOfUser(w http.ResponseWriter, r *http.Request) {
	cookie, _ := r.Cookie("session")

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

	udpErr := cvHandler.cvUsecase.UpdateCVOfUserById(r.Context(), cookie.Value, cvID, cv)
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

	deleteErr := cvHandler.cvUsecase.DeleteCVOfUserById(r.Context(), cookie.Value, cvID)
	if deleteErr != nil {
		responseTemplates.SendErrorMessage(w, deleteErr, http.StatusBadRequest)
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
		responseTemplates.SendErrorMessage(w, err, http.StatusBadRequest)
		return
	}

	info.CVs = cvHandler.sanitizeCVs(info.CVs...)

	responseTemplates.MarshalAndSend(w, info)
}
