package http

import (
	"HnH/internal/appErrors"
	"HnH/internal/domain"
	"HnH/internal/repository/psql"
	"HnH/internal/usecase"
	"HnH/pkg/contextUtils"
	"HnH/pkg/middleware"
	"HnH/pkg/responseTemplates"
	"HnH/pkg/sanitizer"
	"HnH/services/searchEngineService/searchEnginePB"

	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/mailru/easyjson"
	"github.com/sirupsen/logrus"
)

// const (
// 	SEARCH_QUERY_KEY           = "q"
// 	PAGE_NUM_QUERY_KEY         = "page_num"
// 	RESULTS_PER_PAGE_QUERY_KEY = "results_per_page"
// )

type VacancyHandler struct {
	vacancyUsecase usecase.IVacancyUsecase
}

func (vacancyHandler *VacancyHandler) sanitizeVacancies(vacancies ...domain.ApiVacancy) []domain.ApiVacancy {
	result := make([]domain.ApiVacancy, 0, len(vacancies))

	for _, vac := range vacancies {
		vac.VacancyName = sanitizer.XSS.Sanitize(vac.VacancyName)
		vac.Description = sanitizer.XSS.Sanitize(vac.Description)
		vac.Employment = domain.EmploymentType((sanitizer.XSS.Sanitize(string(vac.Employment))))

		if vac.Location != nil {
			*vac.Location = sanitizer.XSS.Sanitize(*vac.Location)
		}

		result = append(result, vac)
	}

	return result
}

func (vacancyHandler *VacancyHandler) sanitizeMetaVacancies(metaVacancies domain.ApiMetaVacancy) domain.ApiMetaVacancy {
	result := domain.ApiMetaVacancy{
		Filters: metaVacancies.Filters,
		Vacancies: domain.ApiVacancyCount{
			Count:     metaVacancies.Vacancies.Count,
			Vacancies: vacancyHandler.sanitizeVacancies(metaVacancies.Vacancies.Vacancies...),
		},
	}
	return result
}

func NewVacancyHandler(router *mux.Router, vacancyUCase usecase.IVacancyUsecase, sessionUCase usecase.ISessionUsecase) {
	handler := &VacancyHandler{
		vacancyUsecase: vacancyUCase,
	}

	router.Handle("/vacancies",
		middleware.SetSessionIDIfExists(sessionUCase, http.HandlerFunc(handler.GetVacancies))).
		Methods("GET")

	router.Handle("/vacancies/search",
		middleware.SetSessionIDIfExists(sessionUCase, http.HandlerFunc(handler.SearchVacancies))).
		Methods("GET")

	router.Handle("/vacancies",
		middleware.JSONBodyValidationMiddleware(middleware.AuthMiddleware(sessionUCase, http.HandlerFunc(handler.AddVacancy)))).
		Methods("POST")

	router.Handle("/vacancies/current_user", middleware.AuthMiddleware(sessionUCase, http.HandlerFunc(handler.GetUserVacancies))).
		Methods("GET")

	router.Handle("/vacancies/employer/{employerID}",
		middleware.SetSessionIDIfExists(sessionUCase, http.HandlerFunc(handler.GetEmployerInfo))).
		Methods("GET")

	router.Handle("/vacancies/{vacancyID}",
		middleware.SetSessionIDIfExists(sessionUCase, http.HandlerFunc(handler.GetVacancy))).
		Methods("GET")

	router.Handle("/vacancies/{vacancyID}",
		middleware.JSONBodyValidationMiddleware(middleware.AuthMiddleware(sessionUCase, http.HandlerFunc(handler.UpdateVacancy)))).
		Methods("PUT")

	router.Handle("/vacancies/{vacancyID}",
		middleware.AuthMiddleware(sessionUCase, http.HandlerFunc(handler.DeleteVacancy))).
		Methods("DELETE")

	router.Handle("/vacancies/favourite/{vacancyID}",
		middleware.AuthMiddleware(sessionUCase, http.HandlerFunc(handler.AddToFavourite))).
		Methods("POST")

	router.Handle("/vacancies/favourite/{vacancyID}",
		middleware.AuthMiddleware(sessionUCase, http.HandlerFunc(handler.DeleteFromFavourite))).
		Methods("DELETE")

	router.Handle("/vacancies/favourite",
		middleware.AuthMiddleware(sessionUCase, http.HandlerFunc(handler.GetFavourite))).
		Methods("GET")
}

func (vacancyHandler *VacancyHandler) GetVacancies(w http.ResponseWriter, r *http.Request) {
	vacancies, getErr := vacancyHandler.vacancyUsecase.GetAllVacancies(r.Context())

	if getErr == psql.ErrEntityNotFound {
		vacancies = []domain.ApiVacancy{}
	} else if getErr != nil {
		errToSend, code := appErrors.GetErrAndCodeToSend(getErr)
		responseTemplates.SendErrorMessage(w, errToSend, code)
		return
	}

	sanitizedVacancies := vacancyHandler.sanitizeVacancies(vacancies...)
	toSend := domain.ApiVacancySlice(sanitizedVacancies)

	responseTemplates.MarshalAndSend(w, toSend)
}

func (vacancyHandler *VacancyHandler) SearchVacancies(w http.ResponseWriter, r *http.Request) {
	contextLogger := contextUtils.GetContextLogger(r.Context())
	query := r.URL.Query()
	contextLogger.WithFields(logrus.Fields{
		"query": query.Encode(),
	}).
		Debug("got search request with query")

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

	metaVacancies, getErr := vacancyHandler.vacancyUsecase.SearchVacancies(r.Context(), &options)

	if getErr != nil {
		errToSend, code := appErrors.GetErrAndCodeToSend(getErr)
		responseTemplates.SendErrorMessage(w, errToSend, code)
		return
	}

	sanitizedMetaVacancies := vacancyHandler.sanitizeMetaVacancies(metaVacancies)

	responseTemplates.MarshalAndSend(w, sanitizedMetaVacancies)
}

func (vacancyHandler *VacancyHandler) GetVacancy(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	vacancyID, convErr := strconv.Atoi(vars["vacancyID"])
	if convErr != nil {
		responseTemplates.SendErrorMessage(w, convErr, http.StatusBadRequest)
		return
	}

	vacancy, err := vacancyHandler.vacancyUsecase.GetVacancy(r.Context(), vacancyID)
	if err != nil {
		errToSend, code := appErrors.GetErrAndCodeToSend(err)
		responseTemplates.SendErrorMessage(w, errToSend, code)
		return
	}

	*vacancy = vacancyHandler.sanitizeVacancies(*vacancy)[0]

	responseTemplates.MarshalAndSend(w, *vacancy)
}

func (vacancyHandler *VacancyHandler) AddVacancy(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	apiVac := new(domain.ApiVacancy)
	err := easyjson.UnmarshalFromReader(r.Body, apiVac)
	if err != nil {
		responseTemplates.SendErrorMessage(w, ErrWrongBodyParam, http.StatusBadRequest)
		return
	}

	vacID, addStatus := vacancyHandler.vacancyUsecase.AddVacancy(r.Context(), apiVac)
	if addStatus != nil {
		errToSend, code := appErrors.GetErrAndCodeToSend(addStatus)
		responseTemplates.SendErrorMessage(w, errToSend, code)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(fmt.Sprintf(`{"id":%d}`, vacID)))
}

func (vacancyHandler *VacancyHandler) UpdateVacancy(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	vacancyID, convErr := strconv.Atoi(vars["vacancyID"])
	if convErr != nil {
		responseTemplates.SendErrorMessage(w, convErr, http.StatusBadRequest)
		return
	}

	defer r.Body.Close()

	updatedVac := new(domain.ApiVacancy)
	err := easyjson.UnmarshalFromReader(r.Body, updatedVac)
	if err != nil {
		responseTemplates.SendErrorMessage(w, ErrWrongBodyParam, http.StatusBadRequest)
		return
	}

	updStatus := vacancyHandler.vacancyUsecase.UpdateVacancy(r.Context(), vacancyID, updatedVac)
	if updStatus != nil {
		errToSend, code := appErrors.GetErrAndCodeToSend(updStatus)
		responseTemplates.SendErrorMessage(w, errToSend, code)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (vacancyHandler *VacancyHandler) DeleteVacancy(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	vacancyID, convErr := strconv.Atoi(vars["vacancyID"])
	if convErr != nil {
		responseTemplates.SendErrorMessage(w, convErr, http.StatusBadRequest)
		return
	}

	delStatus := vacancyHandler.vacancyUsecase.DeleteVacancy(r.Context(), vacancyID)
	if delStatus != nil {
		errToSend, code := appErrors.GetErrAndCodeToSend(delStatus)
		responseTemplates.SendErrorMessage(w, errToSend, code)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (vacancyHandler *VacancyHandler) GetUserVacancies(w http.ResponseWriter, r *http.Request) {
	vacanciesList, err := vacancyHandler.vacancyUsecase.GetUserVacancies(r.Context())
	if err != nil {
		errToSend, code := appErrors.GetErrAndCodeToSend(err)
		responseTemplates.SendErrorMessage(w, errToSend, code)
		return
	}

	sanitizedList := vacancyHandler.sanitizeVacancies(vacanciesList...)
	toSend := domain.ApiVacancySlice(sanitizedList)

	responseTemplates.MarshalAndSend(w, toSend)
}

func (vacancyHandler *VacancyHandler) GetEmployerInfo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	empID, convErr := strconv.Atoi(vars["employerID"])
	if convErr != nil {
		responseTemplates.SendErrorMessage(w, convErr, http.StatusBadRequest)
		return
	}

	info, err := vacancyHandler.vacancyUsecase.GetEmployerInfo(r.Context(), empID)
	if err != nil {
		errToSend, code := appErrors.GetErrAndCodeToSend(err)
		responseTemplates.SendErrorMessage(w, errToSend, code)
		return
	}

	info.Vacancies = vacancyHandler.sanitizeVacancies(info.Vacancies...)

	responseTemplates.MarshalAndSend(w, *info)
}

func (vacancyHandler *VacancyHandler) AddToFavourite(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	vacID, convErr := strconv.Atoi(vars["vacancyID"])
	if convErr != nil {
		responseTemplates.SendErrorMessage(w, convErr, http.StatusBadRequest)
		return
	}

	err := vacancyHandler.vacancyUsecase.AddToFavourite(r.Context(), vacID)
	if err != nil {
		errToSend, code := appErrors.GetErrAndCodeToSend(err)
		responseTemplates.SendErrorMessage(w, errToSend, code)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (vacancyHandler *VacancyHandler) DeleteFromFavourite(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	vacID, convErr := strconv.Atoi(vars["vacancyID"])
	if convErr != nil {
		responseTemplates.SendErrorMessage(w, convErr, http.StatusBadRequest)
		return
	}

	err := vacancyHandler.vacancyUsecase.DeleteFromFavourite(r.Context(), vacID)
	if err != nil {
		errToSend, code := appErrors.GetErrAndCodeToSend(err)
		responseTemplates.SendErrorMessage(w, errToSend, code)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (vacancyHandler *VacancyHandler) GetFavourite(w http.ResponseWriter, r *http.Request) {
	vacsList, err := vacancyHandler.vacancyUsecase.GetFavourite(r.Context())
	if err != nil {
		errToSend, code := appErrors.GetErrAndCodeToSend(err)
		responseTemplates.SendErrorMessage(w, errToSend, code)
		return
	}

	sanitizedVacancies := vacancyHandler.sanitizeVacancies(vacsList...)
	toSend := domain.ApiVacancySlice(sanitizedVacancies)

	responseTemplates.MarshalAndSend(w, toSend)
}
