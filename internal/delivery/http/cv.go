package http

import (
	"HnH/internal/domain"
	"HnH/internal/usecase"
	"HnH/pkg/serverErrors"
	"fmt"

	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type CVHandler struct {
	cvUsecase usecase.ICVUsecase
}

func NewCVHandler(router *mux.Router, cvUCase usecase.ICVUsecase) {
	handler := &CVHandler{
		cvUsecase: cvUCase,
	}

	router.HandleFunc("/cv/{cvID}", handler.GetCV).Methods("GET")
	router.HandleFunc("/current_user/cvs", handler.GetCVList).Methods("GET")
	router.HandleFunc("/current_user/cvs", handler.AddNewCV).Methods("POST")
	router.HandleFunc("/current_user/cvs/{cvID}", handler.GetCVOfUser).Methods("GET")
	router.HandleFunc("/current_user/cvs/{cvID}", handler.UpdateCVOfUser).Methods("PUT")
	router.HandleFunc("/current_user/cvs/{cvID}", handler.DeleteCVOfUser).Methods("DELETE")
}

func (cvHandler *CVHandler) GetCV(w http.ResponseWriter, r *http.Request) {

}

func (cvHandler *CVHandler) GetCVList(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session")

	if errors.Is(err, http.ErrNoCookie) {
		sendErrorMessage(w, serverErrors.NO_COOKIE, http.StatusUnauthorized)
		return
	}

	cvs, err := cvHandler.cvUsecase.GetCVList(cookie.Value)
	if err != nil {
		sendErrorMessage(w, err, http.StatusBadRequest)
		return
	}

	js, err := json.Marshal(cvs)
	if err != nil {
		sendErrorMessage(w, serverErrors.INTERNAL_SERVER_ERROR, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func (cvHandler *CVHandler) AddNewCV(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session")

	if errors.Is(err, http.ErrNoCookie) {
		sendErrorMessage(w, serverErrors.NO_COOKIE, http.StatusUnauthorized)
		return
	}

	defer r.Body.Close()

	cv := new(domain.CV)

	readErr := json.NewDecoder(r.Body).Decode(cv)
	if readErr != nil {
		sendErrorMessage(w, readErr, http.StatusBadRequest)
		return
	}

	newCVID, addErr := cvHandler.cvUsecase.AddNewCV(cookie.Value, cv)
	if addErr != nil {
		sendErrorMessage(w, addErr, http.StatusUnauthorized)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(fmt.Sprintf(`{"id":%d}`, newCVID)))
}

func (cvHandler *CVHandler) GetCVOfUser(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session")

	if errors.Is(err, http.ErrNoCookie) {
		sendErrorMessage(w, serverErrors.NO_COOKIE, http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	cvID, convErr := strconv.Atoi(vars["cvID"])
	if convErr != nil {
		sendErrorMessage(w, convErr, http.StatusBadRequest)
		return
	}

	cv, err := cvHandler.cvUsecase.GetCVOfUserById(cookie.Value, cvID)
	if err != nil {
		sendErrorMessage(w, err, http.StatusBadRequest)
		return
	}

	js, err := json.Marshal(*cv)
	if err != nil {
		sendErrorMessage(w, serverErrors.INTERNAL_SERVER_ERROR, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func (cvHandler *CVHandler) UpdateCVOfUser(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session")

	if errors.Is(err, http.ErrNoCookie) {
		sendErrorMessage(w, serverErrors.NO_COOKIE, http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	cvID, convErr := strconv.Atoi(vars["cvID"])
	if convErr != nil {
		sendErrorMessage(w, convErr, http.StatusBadRequest)
		return
	}

	udpErr := cvHandler.cvUsecase.UpdateCVOfUserById(cookie.Value, cvID)
	if udpErr != nil {
		sendErrorMessage(w, err, http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (cvHandler *CVHandler) DeleteCVOfUser(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session")

	if errors.Is(err, http.ErrNoCookie) {
		sendErrorMessage(w, serverErrors.NO_COOKIE, http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	cvID, convErr := strconv.Atoi(vars["cvID"])
	if convErr != nil {
		sendErrorMessage(w, convErr, http.StatusBadRequest)
		return
	}

	deleteErr := cvHandler.cvUsecase.DeleteCVOfUserById(cookie.Value, cvID)
	if deleteErr != nil {
		sendErrorMessage(w, err, http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}
