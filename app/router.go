package app

import (
	"HnH/configs"
	deliveryHTTP "HnH/internal/delivery/http"
	"HnH/internal/delivery/http/middleware"
	grpcRepo "HnH/internal/repository/grpc"
	"HnH/internal/repository/psql"
	"HnH/internal/usecase"
	"HnH/pkg/logging"
	authPB "HnH/services/auth/authPB"
	auth "HnH/services/auth/config"
	csat "HnH/services/csat/config"
	csatPB "HnH/services/csat/csatPB"
	"os"

	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func initCsatClient(config csat.CsatConfig) (csatPB.CsatClient, error) {
	connAddr := fmt.Sprintf("%s:%d", config.Host, config.Port)

	opts := []grpc.DialOption{}
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))

	conn, err := grpc.Dial(connAddr, opts...)
	if err != nil {
		return nil, err
	}

	client := csatPB.NewCsatClient(conn)
	return client, nil
}

func initAuthClient(config auth.AuthConfig) (authPB.AuthClient, error) {
	connAddr := fmt.Sprintf("%s:%d", config.Host, config.Port)

	opts := []grpc.DialOption{}
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))

	conn, err := grpc.Dial(connAddr, opts...)
	if err != nil {
		return nil, err
	}

	client := authPB.NewAuthClient(conn)
	return client, nil
}

func Run() error {
	logFile, err := os.OpenFile(configs.LOGFILE_NAME, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	defer logFile.Close()

	logging.InitLogger(logFile)

	db, err := GetPostgres()
	if err != nil {
		return err
	}
	defer db.Close()

	/*redisDB := getRedis()
	if err != nil {
		return err
	}
	defer redisDB.Close()*/

	csatClient, err := initCsatClient(csat.CsatServiceConfig)
	if err != nil {
		return err
	}

	authClient, err := initAuthClient(auth.AuthServiceConfig)
	if err != nil {
		return err
	}

	//sessionRepo := redisRepo.NewRedisSessionRepository(redisDB)
	authRepo := grpcRepo.NewGrpcAuthRepository(authClient)
	userRepo := psql.NewPsqlUserRepository(db)
	organizationRepo := psql.NewPsqlOrganizationRepository(db)
	vacancyRepo := psql.NewPsqlVacancyRepository(db)
	cvRepo := psql.NewPsqlCVRepository(db)
	responseRepo := psql.NewPsqlResponseRepository(db)
	experienceRepo := psql.NewPsqlExperienceRepository(db)
	institutionRepo := psql.NewPsqlEducationInstitutionRepository(db)
	csatRepo := grpcRepo.NewGrpcSearchEngineRepository(csatClient)

	sessionUsecase := usecase.NewSessionUsecase(authRepo, userRepo)
	userUsecase := usecase.NewUserUsecase(userRepo, authRepo, organizationRepo)
	vacancyUsecase := usecase.NewVacancyUsecase(vacancyRepo, authRepo, userRepo)
	cvUsecase := usecase.NewCVUsecase(cvRepo, experienceRepo, institutionRepo, authRepo, userRepo, responseRepo, vacancyRepo)
	responseUsecase := usecase.NewResponseUsecase(responseRepo, authRepo, userRepo, vacancyRepo, cvRepo)
	csatUsecase := usecase.NewCsatUsecase(csatRepo, authRepo)

	router := mux.NewRouter()
	// router.Use(func(h http.Handler) http.Handler {
	// 	return middleware.CSRFProtectionMiddleware(sessionRepo, h)
	// })

	deliveryHTTP.NewSessionHandler(router, sessionUsecase)
	deliveryHTTP.NewUserHandler(router, userUsecase, sessionUsecase)
	deliveryHTTP.NewVacancyHandler(router, vacancyUsecase, sessionUsecase)
	deliveryHTTP.NewCVHandler(router, cvUsecase, sessionUsecase)
	deliveryHTTP.NewCsatHandler(router, csatUsecase, sessionUsecase)
	deliveryHTTP.NewResponseHandler(router, responseUsecase, sessionUsecase)

	corsRouter := configs.CORS.Handler(router)
	loggedRouter := middleware.AccessLogMiddleware( /* logging.Logger,  */ corsRouter)
	requestIDRouter := middleware.RequestID(loggedRouter)
	finalRouter := middleware.PanicRecoverMiddleware(logging.Logger, requestIDRouter)

	http.Handle("/", finalRouter)

	fmt.Printf("\tstarting server at %s\n", configs.PORT)
	logging.Logger.Infof("starting server at %s", configs.PORT)

	err = http.ListenAndServe(configs.PORT, nil)
	if err != nil {
		return err
	}

	return nil
}
