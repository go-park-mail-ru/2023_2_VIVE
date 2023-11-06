package app

import (
	"HnH/configs"
	deliveryHTTP "HnH/internal/delivery/http"
	"HnH/internal/delivery/http/middleware"
	"HnH/internal/repository"
	"HnH/internal/repository/mock"
	"HnH/internal/usecase"
	"HnH/pkg/logging"
	"HnH/pkg/serverErrors"

	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

func Run() error {
	logFile, err := os.OpenFile(configs.LOGFILE_NAME, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return err
	}

	defer logFile.Close()
	logger := logging.InitLogger(logFile)

	sessionRepo := repository.NewPsqlSessionRepository(&mock.SessionDB)
	userRepo := repository.NewPsqlUserRepository(&mock.UserDB)
	vacancyRepo := repository.NewPsqlVacancyRepository(&mock.VacancyDB)
	cvRepo := repository.NewPsqlCVRepository()
	responseRepo := repository.NewPsqlResponseRepository()

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
	loggedRouter := middleware.AccessLogMiddleware(logger, corsRouter)
	finalRouter := middleware.PanicRecoverMiddleware(logger, loggedRouter)

	http.Handle("/", finalRouter)

	fmt.Printf("\tstarting server at %s\n", configs.PORT)
	logger.Infof("starting server at %s", configs.PORT)

	err = http.ListenAndServe(configs.PORT, nil)
	if err != nil {
		logger.WithField("message", serverErrors.SERVER_IS_NOT_RUNNING).Error(err)
		return err
	}

	return nil
}
