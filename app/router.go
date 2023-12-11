package app

import (
	"HnH/configs"
	"HnH/configs/metrics"
	deliveryHTTP "HnH/internal/delivery/http"
	grpcRepo "HnH/internal/repository/grpc"
	"HnH/internal/repository/psql"
	"HnH/internal/usecase"
	"HnH/pkg/logging"
	"HnH/pkg/middleware"
	"HnH/services/auth/authPB"
	authConfig "HnH/services/auth/config"
	csatConfig "HnH/services/csat/config"
	"HnH/services/csat/csatPB"
	notificationsPB "HnH/services/notifications/api/proto"
	notificationsConfig "HnH/services/notifications/config"
	searchConfig "HnH/services/searchEngineService/config"
	"HnH/services/searchEngineService/searchEnginePB"
	"os"

	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func initCsatClient(config csatConfig.CsatConfig) (csatPB.CsatClient, error) {
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

func initAuthClient(config authConfig.AuthConfig) (authPB.AuthClient, error) {
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

func initNotificationsClient(config notificationsConfig.NotificationsGRPCConfig) (notificationsPB.NotificationServiceClient, error) {
	connAddr := fmt.Sprintf("%s:%d", config.Host, config.Port)

	opts := []grpc.DialOption{}
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))

	conn, err := grpc.Dial(connAddr, opts...)
	if err != nil {
		return nil, err
	}

	client := notificationsPB.NewNotificationServiceClient(conn)
	return client, nil
}

func initSearchEngineClient(config searchConfig.SearchEngineConfig) (searchEnginePB.SearchEngineClient, error) {
	connAddr := fmt.Sprintf("%s:%d", config.Host, config.Port)

	opts := []grpc.DialOption{}
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))

	conn, err := grpc.Dial(connAddr, opts...)
	if err != nil {
		return nil, err
	}

	client := searchEnginePB.NewSearchEngineClient(conn)
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

	csatClient, err := initCsatClient(csatConfig.CsatServiceConfig)
	if err != nil {
		return err
	}

	authClient, err := initAuthClient(authConfig.AuthServiceConfig)
	if err != nil {
		return err
	}

	notificationsClient, err := initNotificationsClient(notificationsConfig.NotificationGRPCServiceConfig)
	if err != nil {
		return err
	}

	authRepo := grpcRepo.NewGrpcAuthRepository(authClient)
	userRepo := psql.NewPsqlUserRepository(db)
	vacancyRepo := psql.NewPsqlVacancyRepository(db)
	cvRepo := psql.NewPsqlCVRepository(db)
	responseRepo := psql.NewPsqlResponseRepository(db)
	experienceRepo := psql.NewPsqlExperienceRepository(db)
	institutionRepo := psql.NewPsqlEducationInstitutionRepository(db)
	csatRepo := grpcRepo.NewGrpcCsatRepository(csatClient)
	skillRepo := psql.NewPsqlSkillRepository(db)
	notificationsRepo := grpcRepo.NewGrpcNotificationRepository(notificationsClient)

	searchEngineClient, err := initSearchEngineClient(searchConfig.SearchEngineServiceConfig)
	if err != nil {
		return err
	}
	searchEngineClientRepo := grpcRepo.NewGrpcSearchEngineRepository(searchEngineClient)

	sessionUsecase := usecase.NewSessionUsecase(authRepo, userRepo)
	userUsecase := usecase.NewUserUsecase(userRepo, authRepo)
	vacancyUsecase := usecase.NewVacancyUsecase(vacancyRepo, authRepo, userRepo, searchEngineClientRepo, skillRepo)
	cvUsecase := usecase.NewCVUsecase(cvRepo, experienceRepo, institutionRepo, authRepo, userRepo, responseRepo, vacancyRepo, searchEngineClientRepo, skillRepo)
	responseUsecase := usecase.NewResponseUsecase(responseRepo, authRepo, userRepo, vacancyRepo, cvRepo, notificationsRepo)
	csatUsecase := usecase.NewCsatUsecase(csatRepo, authRepo)
	notificationUsecase := usecase.NewNotificationUsecase(notificationsRepo)

	router := mux.NewRouter()
	//router.Use(func(h http.Handler) http.Handler {
	//return middleware.CSRFProtectionMiddleware(sessionUsecase, h)
	//})

	deliveryHTTP.NewSessionHandler(router, sessionUsecase)
	deliveryHTTP.NewUserHandler(router, userUsecase, sessionUsecase)
	deliveryHTTP.NewVacancyHandler(router, vacancyUsecase, sessionUsecase)
	deliveryHTTP.NewCVHandler(router, cvUsecase, sessionUsecase)
	deliveryHTTP.NewCsatHandler(router, csatUsecase, sessionUsecase)
	deliveryHTTP.NewResponseHandler(router, responseUsecase, sessionUsecase)
	deliveryHTTP.NewNotificationHandler(router, notificationUsecase, sessionUsecase)

	prometheus.MustRegister(metrics.HitCounter, metrics.ErrorCounter)
	router.Handle("/metrics", promhttp.Handler())

	corsRouter := configs.CORS.Handler(router)
	loggedRouter := middleware.AccessLogMiddleware(corsRouter)
	requestIDRouter := middleware.RequestID(loggedRouter)
	recoverRouter := middleware.PanicRecoverMiddleware(logging.Logger, requestIDRouter)

	hitCountRouter := metrics.HitCounterMiddleware(recoverRouter)
	timingCountRouter := metrics.TimingHistogramMiddleware(hitCountRouter)
	finalRouter := metrics.ErrorCounterMiddleware(timingCountRouter)

	http.Handle("/", finalRouter)

	fmt.Printf("\tstarting server at %s\n", configs.PORT)
	logging.Logger.Infof("starting server at %s", configs.PORT)

	err = http.ListenAndServe(configs.PORT, nil)
	if err != nil {
		return err
	}

	return nil
}
