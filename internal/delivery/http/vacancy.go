package http

import (
	"HnH/internal/appErrors"
	"HnH/internal/delivery/http/middleware"
	"HnH/internal/domain"
	"HnH/internal/repository/psql"
	"HnH/internal/usecase"
	"HnH/pkg/contextUtils"
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

	router.HandleFunc("/vacancies",
		handler.GetVacancies).
		Methods("GET")

	router.HandleFunc("/vacancies/search", handler.SearchVacancies).
		Methods("GET")

	router.Handle("/vacancies",
		middleware.JSONBodyValidationMiddleware(middleware.AuthMiddleware(sessionUCase, http.HandlerFunc(handler.AddVacancy)))).
		Methods("POST")

	router.Handle("/vacancies/current_user", middleware.AuthMiddleware(sessionUCase, http.HandlerFunc(handler.GetUserVacancies))).
		Methods("GET")

	router.HandleFunc("/vacancies/employer/{employerID}", handler.GetEmployerInfo).
		Methods("GET")

	router.HandleFunc("/vacancies/{vacancyID}",
		handler.GetVacancy).
		Methods("GET")

	router.Handle("/vacancies/{vacancyID}",
		middleware.JSONBodyValidationMiddleware(middleware.AuthMiddleware(sessionUCase, http.HandlerFunc(handler.UpdateVacancy)))).
		Methods("PUT")

	router.Handle("/vacancies/{vacancyID}",
		middleware.AuthMiddleware(sessionUCase, http.HandlerFunc(handler.DeleteVacancy))).
		Methods("DELETE")

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

	responseTemplates.MarshalAndSend(w, sanitizedVacancies)
}

func (vacancyHandler *VacancyHandler) SearchVacancies(w http.ResponseWriter, r *http.Request) {
	contextLogger := contextUtils.GetContextLogger(r.Context())
	query := r.URL.Query()
	contextLogger.WithFields(logrus.Fields{
		"query": query.Encode(),
	}).
		Debug("got search request with query")

	// searchEnginePB.SearchRequest
	// queryOptions := searchEnginePB.SearchOptions{}
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
		// option := searchEnginePB.SearchOption{
		// 	Name:   optionName,
		// 	Values: values,
		// }
		// options = append(options, &option)
	}
	options := searchEnginePB.SearchOptions{Options: queryOptions}

	// searchQuery := query.Get(SEARCH_QUERY_KEY)

	// pageNumStr := query.Get(PAGE_NUM_QUERY_KEY)
	// pageNum, convErr := strconv.ParseInt(pageNumStr, 10, 64)
	// if convErr != nil {
	// 	responseTemplates.SendErrorMessage(w, ErrWrongQueryParam, http.StatusBadRequest)
	// 	return
	// }

	// resultsPerPageStr := query.Get(RESULTS_PER_PAGE_QUERY_KEY)
	// resultsPerPage, convErr := strconv.ParseInt(resultsPerPageStr, 10, 64)
	// if convErr != nil {
	// 	responseTemplates.SendErrorMessage(w, ErrWrongQueryParam, http.StatusBadRequest)
	// 	return
	// }

	metaVacancies, getErr := vacancyHandler.vacancyUsecase.SearchVacancies(r.Context(), &options)

	if getErr != nil {
		errToSend, code := appErrors.GetErrAndCodeToSend(getErr)
		responseTemplates.SendErrorMessage(w, errToSend, code)
		return
	}

	sanitizedMetaVacancies := vacancyHandler.sanitizeMetaVacancies(metaVacancies)

	responseTemplates.MarshalAndSend(w, sanitizedMetaVacancies)

	w.WriteHeader(http.StatusOK)
}

func (vacancyHandler *VacancyHandler) GetVacancy(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	vacancyID, convErr := strconv.Atoi(vars["vacancyID"])
	if convErr != nil {
		responseTemplates.SendErrorMessage(w, convErr, http.StatusBadRequest)
		return
	}

	vacWithCompName, err := vacancyHandler.vacancyUsecase.GetVacancyWithCompanyName(r.Context(), vacancyID)
	if err != nil {
		errToSend, code := appErrors.GetErrAndCodeToSend(err)
		responseTemplates.SendErrorMessage(w, errToSend, code)
		return
	}

	vacWithCompName.Vacancy = vacancyHandler.sanitizeVacancies(vacWithCompName.Vacancy)[0]

	responseTemplates.MarshalAndSend(w, vacWithCompName)
}

func (vacancyHandler *VacancyHandler) AddVacancy(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	apiVac := new(domain.ApiVacancy)

	readErr := json.NewDecoder(r.Body).Decode(apiVac)
	if readErr != nil {
		responseTemplates.SendErrorMessage(w, readErr, http.StatusBadRequest)
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

	readErr := json.NewDecoder(r.Body).Decode(updatedVac)
	if readErr != nil {
		responseTemplates.SendErrorMessage(w, readErr, http.StatusBadRequest)
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

	responseTemplates.MarshalAndSend(w, sanitizedList)
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

	responseTemplates.MarshalAndSend(w, info)
}
