package http

import (
	"HnH/internal/domain"
	"HnH/internal/usecase"
	"HnH/pkg/serverErrors"

	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type VacancyHandler struct {
	vacancyUsecase usecase.IVacancyUsecase
}

func NewVacancyHandler(router *mux.Router, vacancyUCase usecase.IVacancyUsecase) {
	handler := &VacancyHandler{
		vacancyUsecase: vacancyUCase,
	}

	router.HandleFunc("/vacancies", handler.GetVacancies).Methods("GET")
	router.HandleFunc("/vacancies", handler.AddVacancy).Methods("POST")
	router.HandleFunc("/vacancies/{vacancyID}", handler.GetVacancy).Methods("GET")
	router.HandleFunc("/vacancies/{vacancyID}", handler.UpdateVacancy).Methods("PUT")
	router.HandleFunc("/vacancies/{vacancyID}", handler.DeleteVacancy).Methods("DELETE")
}

func (vacancyHandler *VacancyHandler) GetVacancy(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	vacancyID, convErr := strconv.Atoi(vars["vacancyID"])
	if convErr != nil {
		sendErrorMessage(w, convErr, http.StatusBadRequest)
		return
	}

	vacancy, err := vacancyHandler.vacancyUsecase.GetVacancy(vacancyID)
	if err != nil {
		sendErrorMessage(w, err, http.StatusBadRequest)
		return
	}

	js, err := json.Marshal(*vacancy)
	if err != nil {
		sendErrorMessage(w, err, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func (vacancyHandler *VacancyHandler) AddVacancy(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session")

	if errors.Is(err, http.ErrNoCookie) {
		sendErrorMessage(w, serverErrors.NO_COOKIE, http.StatusUnauthorized)
		return
	}

	defer r.Body.Close()

	vac := new(domain.Vacancy)

	readErr := json.NewDecoder(r.Body).Decode(vac)
	if readErr != nil {
		sendErrorMessage(w, readErr, http.StatusBadRequest)
		return
	}

	vacID, addStatus := vacancyHandler.vacancyUsecase.AddVacancy(cookie.Value, vac)
	if addStatus != nil {
		sendErrorMessage(w, addStatus, http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(fmt.Sprintf(`{"id":%d}`, vacID)))
}

func (vacancyHandler *VacancyHandler) UpdateVacancy(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session")

	if errors.Is(err, http.ErrNoCookie) {
		sendErrorMessage(w, serverErrors.NO_COOKIE, http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	vacancyID, convErr := strconv.Atoi(vars["vacancyID"])
	if convErr != nil {
		sendErrorMessage(w, convErr, http.StatusBadRequest)
		return
	}

	defer r.Body.Close()

	vac := new(domain.Vacancy)

	readErr := json.NewDecoder(r.Body).Decode(vac)
	if readErr != nil {
		sendErrorMessage(w, readErr, http.StatusBadRequest)
		return
	}

	updStatus := vacancyHandler.vacancyUsecase.UpdateVacancy(cookie.Value, vacancyID, vac)
	if updStatus != nil {
		sendErrorMessage(w, updStatus, http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (vacancyHandler *VacancyHandler) DeleteVacancy(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session")

	if errors.Is(err, http.ErrNoCookie) {
		sendErrorMessage(w, serverErrors.NO_COOKIE, http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	vacancyID, convErr := strconv.Atoi(vars["vacancyID"])
	if convErr != nil {
		sendErrorMessage(w, convErr, http.StatusBadRequest)
		return
	}

	delStatus := vacancyHandler.vacancyUsecase.DeleteVacancy(cookie.Value, vacancyID)
	if delStatus != nil {
		sendErrorMessage(w, delStatus, http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (vacancyHandler *VacancyHandler) GetVacancies(w http.ResponseWriter, r *http.Request) {
	vacancies, getErr := vacancyHandler.vacancyUsecase.GetAllVacancies()
	if getErr != nil {
		sendErrorMessage(w, getErr, http.StatusBadRequest)
		return
	}

	js, err := json.Marshal(vacancies)
	if err != nil {
		sendErrorMessage(w, err, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}
