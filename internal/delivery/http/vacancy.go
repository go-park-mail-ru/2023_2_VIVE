package http

import (
	"HnH/internal/delivery/http/middleware"
	"HnH/internal/domain"
	"HnH/internal/repository/psql"
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

const (
	SEARCH_QUERY_KEY           = "q"
	PAGE_NUM_QUERY_KEY         = "page_num"
	RESULTS_PER_PAGE_QUERY_KEY = "results_per_page"
)

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
			Count: metaVacancies.Vacancies.Count,
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
		responseTemplates.SendErrorMessage(w, getErr, http.StatusBadRequest)
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

	metaVacancies, getErr := vacancyHandler.vacancyUsecase.SearchVacancies(
		r.Context(),
		searchQuery,
		pageNum,
		resultsPerPage,
	)

	if getErr != nil {
		responseTemplates.SendErrorMessage(w, getErr, http.StatusBadRequest)
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

	vacancy, err := vacancyHandler.vacancyUsecase.GetVacancy(r.Context(), vacancyID)
	if err != nil {
		responseTemplates.SendErrorMessage(w, err, http.StatusBadRequest)
		return
	}

	sanitizedVacancy := vacancyHandler.sanitizeVacancies(*vacancy)

	responseTemplates.MarshalAndSend(w, sanitizedVacancy[0])
}

func (vacancyHandler *VacancyHandler) AddVacancy(w http.ResponseWriter, r *http.Request) {
	cookie, _ := r.Cookie("session")

	defer r.Body.Close()

	apiVac := new(domain.ApiVacancy)

	readErr := json.NewDecoder(r.Body).Decode(apiVac)
	if readErr != nil {
		responseTemplates.SendErrorMessage(w, readErr, http.StatusBadRequest)
		return
	}
	// fmt.Printf("apiVac: %v\n", apiVac)

	dbVac := apiVac.ToDb()
	// fmt.Printf("dbVac: %v\n", dbVac)

	vacID, addStatus := vacancyHandler.vacancyUsecase.AddVacancy(r.Context(), cookie.Value, dbVac)
	if addStatus != nil {
		responseTemplates.SendErrorMessage(w, addStatus, http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(fmt.Sprintf(`{"id":%d}`, vacID)))
}

func (vacancyHandler *VacancyHandler) UpdateVacancy(w http.ResponseWriter, r *http.Request) {
	cookie, _ := r.Cookie("session")

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

	updStatus := vacancyHandler.vacancyUsecase.UpdateVacancy(r.Context(), cookie.Value, vacancyID, updatedVac)
	if updStatus != nil {
		responseTemplates.SendErrorMessage(w, updStatus, http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (vacancyHandler *VacancyHandler) DeleteVacancy(w http.ResponseWriter, r *http.Request) {
	cookie, _ := r.Cookie("session")

	vars := mux.Vars(r)
	vacancyID, convErr := strconv.Atoi(vars["vacancyID"])
	if convErr != nil {
		responseTemplates.SendErrorMessage(w, convErr, http.StatusBadRequest)
		return
	}

	delStatus := vacancyHandler.vacancyUsecase.DeleteVacancy(r.Context(), cookie.Value, vacancyID)
	if delStatus != nil {
		responseTemplates.SendErrorMessage(w, delStatus, http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (vacancyHandler *VacancyHandler) GetUserVacancies(w http.ResponseWriter, r *http.Request) {
	cookie, _ := r.Cookie("session")

	vacanciesList, err := vacancyHandler.vacancyUsecase.GetUserVacancies(r.Context(), cookie.Value)
	if err != nil {
		responseTemplates.SendErrorMessage(w, err, http.StatusBadRequest)
		return
	}

	sanitizedList := vacancyHandler.sanitizeVacancies(vacanciesList...)

	responseTemplates.MarshalAndSend(w, sanitizedList)
}
