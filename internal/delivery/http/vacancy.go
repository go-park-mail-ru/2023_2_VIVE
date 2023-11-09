package http

import (
	"HnH/internal/delivery/http/middleware"
	"HnH/internal/domain"
	"HnH/internal/usecase"
	"HnH/pkg/responseTemplates"
	"HnH/pkg/sanitizer"

	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type VacancyHandler struct {
	vacancyUsecase usecase.IVacancyUsecase
}

func (vacancyHandler *VacancyHandler) sanitizeVacancies(vacancies ...domain.DbVacancy) []domain.DbVacancy {
	result := make([]domain.DbVacancy, 0, len(vacancies))

	for _, vac := range vacancies {
		vac.VacancyName = sanitizer.XSS.Sanitize(vac.VacancyName)
		vac.Description = sanitizer.XSS.Sanitize(vac.Description)

		if vac.Employment != nil {
			*vac.Employment = sanitizer.XSS.Sanitize(*vac.Employment)
		}
		if vac.Location != nil {
			*vac.Location = sanitizer.XSS.Sanitize(*vac.Location)
		}

		result = append(result, vac)
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
	vacancies, getErr := vacancyHandler.vacancyUsecase.GetAllVacancies()

	if getErr != nil {
		responseTemplates.SendErrorMessage(w, getErr, http.StatusBadRequest)
		return
	}

	sanitizedVacancies := vacancyHandler.sanitizeVacancies(vacancies...)

	responseTemplates.MarshalAndSend(w, sanitizedVacancies)
}

func (vacancyHandler *VacancyHandler) GetVacancy(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	vacancyID, convErr := strconv.Atoi(vars["vacancyID"])
	if convErr != nil {
		responseTemplates.SendErrorMessage(w, convErr, http.StatusBadRequest)
		return
	}

	vacancy, err := vacancyHandler.vacancyUsecase.GetVacancy(vacancyID)
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

	vac := new(domain.DbVacancy)

	readErr := json.NewDecoder(r.Body).Decode(vac)
	if readErr != nil {
		responseTemplates.SendErrorMessage(w, readErr, http.StatusBadRequest)
		return
	}

	vacID, addStatus := vacancyHandler.vacancyUsecase.AddVacancy(cookie.Value, vac)
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

	vac := new(domain.DbVacancy)

	readErr := json.NewDecoder(r.Body).Decode(vac)
	if readErr != nil {
		responseTemplates.SendErrorMessage(w, readErr, http.StatusBadRequest)
		return
	}

	updStatus := vacancyHandler.vacancyUsecase.UpdateVacancy(cookie.Value, vacancyID, vac)
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

	delStatus := vacancyHandler.vacancyUsecase.DeleteVacancy(cookie.Value, vacancyID)
	if delStatus != nil {
		responseTemplates.SendErrorMessage(w, delStatus, http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (vacancyHandler *VacancyHandler) GetUserVacancies(w http.ResponseWriter, r *http.Request) {
	cookie, _ := r.Cookie("session")

	vacanciesList, err := vacancyHandler.vacancyUsecase.GetUserVacancies(cookie.Value)
	if err != nil {
		responseTemplates.SendErrorMessage(w, err, http.StatusBadRequest)
		return
	}

	sanitizedList := vacancyHandler.sanitizeVacancies(vacanciesList...)

	responseTemplates.MarshalAndSend(w, sanitizedList)
}
