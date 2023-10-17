package app

import (
	"HnH/configs"
	deliveryHTTP "HnH/internal/delivery/http"

	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func Run() error {
	router := mux.NewRouter()

	router.HandleFunc("/session", deliveryHTTP.Login).Methods("POST")
	router.HandleFunc("/session", deliveryHTTP.Logout).Methods("DELETE")
	router.HandleFunc("/session", deliveryHTTP.CheckLogin).Methods("GET")

	router.HandleFunc("/users", deliveryHTTP.SignUp).Methods("POST")
	router.HandleFunc("/current_user", deliveryHTTP.GetInfo).Methods("GET")

	router.HandleFunc("/vacancies", deliveryHTTP.GetVacancies).Methods("GET")

	corsRouter := configs.CORS.Handler(router)
	http.Handle("/", corsRouter)

	fmt.Printf("\tstarting server at %s\n", configs.PORT)
	err := http.ListenAndServe(configs.PORT, nil)
	if err != nil {
		return err
	}

	return nil
}
