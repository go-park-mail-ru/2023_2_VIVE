package app

import (
	"HnH/configs"
	deliveryHTTP "HnH/internal/delivery/http"
	"HnH/internal/repository/psql"
	"HnH/internal/repository/redisRepo"
	"HnH/internal/usecase"

	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func Run() error {
	db, err := getPostgres()
	if err != nil {
		return err
	}

	redisDB, err := getRedis()
	if err != nil {
		return err
	}

	sessionRepo := redisRepo.NewPsqlSessionRepository(redisDB)
	userRepo := psql.NewPsqlUserRepository(db)
	vacancyRepo := psql.NewPsqlVacancyRepository(db)
	cvRepo := psql.NewPsqlCVRepository(db)
	responseRepo := psql.NewPsqlResponseRepository(db)

	sessionUsecase := usecase.NewSessionUsecase(sessionRepo, userRepo)
	userUsecase := usecase.NewUserUsecase(userRepo, sessionRepo)
	vacancyUsecase := usecase.NewVacancyUsecase(vacancyRepo, sessionRepo, userRepo)
	cvUsecase := usecase.NewCVUsecase(cvRepo, sessionRepo, userRepo, responseRepo, vacancyRepo)
	responseUsecase := usecase.NewResponseUsecase(responseRepo, sessionRepo, userRepo, vacancyRepo, cvRepo)

	router := mux.NewRouter()

	deliveryHTTP.NewSessionHandler(router, sessionUsecase)
	deliveryHTTP.NewUserHandler(router, userUsecase, sessionUsecase)
	deliveryHTTP.NewVacancyHandler(router, vacancyUsecase, sessionUsecase)
	deliveryHTTP.NewCVHandler(router, cvUsecase, sessionUsecase)
	deliveryHTTP.NewResponseHandler(router, responseUsecase, sessionUsecase)

	corsRouter := configs.CORS.Handler(router)
	http.Handle("/", corsRouter)

	fmt.Printf("\tstarting server at %s\n", configs.PORT)
	err = http.ListenAndServe(configs.PORT, nil)
	if err != nil {
		return err
	}

	return nil
}
