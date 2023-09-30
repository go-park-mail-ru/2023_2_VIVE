package app

import (
	"models/configs"
	"models/requestHandlers"
	"net/http"

	"github.com/gorilla/mux"
)

func Run() error {
	router := mux.NewRouter()

	router.HandleFunc("/session", requestHandlers.Login).Methods("POST")
	router.HandleFunc("/session", requestHandlers.Logout).Methods("DELETE")

	router.HandleFunc("/users", requestHandlers.SignUp).Methods("POST")
	router.HandleFunc("/current_user", requestHandlers.GetInfo).Methods("GET")

	router.HandleFunc("/vacancies", requestHandlers.GetVacancies).Methods("GET")

	http.Handle("/", router)

	err := http.ListenAndServe(configs.PORT, nil)
	if err != nil {
		return err
	}

	return nil
}
