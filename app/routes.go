package app

import (
	"HnH/configs"
	deliveryHTTP "HnH/internal/delivery/http"
	"HnH/internal/repository"
	"HnH/internal/repository/mock"
	"HnH/internal/usecase"

	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func Run() error {
	sessionRepo := repository.NewPsqlSessionRepository(&mock.SessionDB)
	userRepo := repository.NewPsqlUserRepository(&mock.UserDB)
	vacancyRepo := repository.NewPsqlVacancyRepository(&mock.VacancyDB)

	sessionUsecase := usecase.NewSessionUsecase(sessionRepo, userRepo)
	userUsecase := usecase.NewUserUsecase(userRepo, sessionRepo)
	vacancyUsecase := usecase.NewVacancyUsecase(vacancyRepo)

	router := mux.NewRouter()

	deliveryHTTP.NewSessionHandler(router, sessionUsecase)
	deliveryHTTP.NewUserHandler(router, userUsecase)
	deliveryHTTP.NewVacancyHandler(router, vacancyUsecase)

	corsRouter := configs.CORS.Handler(router)
	http.Handle("/", corsRouter)

	fmt.Printf("\tstarting server at %s\n", configs.PORT)
	err := http.ListenAndServe(configs.PORT, nil)
	if err != nil {
		return err
	}

	return nil
}
