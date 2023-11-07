package app

import (
	"HnH/configs"
	deliveryHTTP "HnH/internal/delivery/http"
	"HnH/internal/delivery/http/middleware"
	"HnH/internal/repository/psql"
	"HnH/internal/repository/redisRepo"
	"HnH/internal/usecase"
	"HnH/pkg/logging"
	"os"

	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func Run() error {
	logFile, err := os.OpenFile(configs.LOGFILE_NAME, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return err
	}

	defer logFile.Close()
	logger := logging.InitLogger(logFile)

	db, err := getPostgres()
	if err != nil {
		return err
	}
	defer db.Close()

	redisDB, err := getRedis()
	if err != nil {
		return err
	}
	defer redisDB.Close()

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
	// router.Use(func(h http.Handler) http.Handler {
	// 	return middleware.CSRFProtectionMiddleware(sessionRepo, h)
	// })


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
		return err
	}

	return nil
}
