package main

import (
	"fmt"
	"models/configs"
	"models/requestHandlers"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/session", requestHandlers.Login).Methods("POST")
	router.HandleFunc("/session", requestHandlers.Logout).Methods("DELETE")

	router.HandleFunc("/users", requestHandlers.SignUp).Methods("POST")
	router.HandleFunc("/current_user", requestHandlers.GetInfo).Methods("GET")

	router.HandleFunc("/vacancies", requestHandlers.GetVacancies).Methods("GET")

	http.Handle("/", router)

	fmt.Printf("starting server at %s\n", configs.PORT)
	err := http.ListenAndServe(configs.PORT, nil)

	if err != nil {
		fmt.Println("err: the server encountered a problem while starting")
		return
	}
}
