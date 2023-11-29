package app

import (
	"HnH/configs"
	"HnH/configs/metrics"
	deliveryHTTP "HnH/internal/delivery/http"
	"HnH/internal/delivery/http/middleware"
	grpcRepo "HnH/internal/repository/grpc"
	"HnH/internal/repository/psql"
	"HnH/internal/usecase"
	"HnH/pkg/logging"
	"HnH/services/auth/authPB"
	authConfig "HnH/services/auth/config"
	csatConfig "HnH/services/csat/config"
	"HnH/services/csat/csatPB"
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

/*var pingCounter = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "ping_request_count",
		Help: "No of request handled by Ping handler",
	},
	[]string{"path"},
)*/

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

	csatClient, err := initCsatClient(csatConfig.CsatServiceConfig)
	if err != nil {
		return err
	}

	authClient, err := initAuthClient(authConfig.AuthServiceConfig)
	if err != nil {
		return err
	}

	//sessionRepo := redisRepo.NewRedisSessionRepository(redisDB)
	authRepo := grpcRepo.NewGrpcAuthRepository(authClient)
	userRepo := psql.NewPsqlUserRepository(db)
	vacancyRepo := psql.NewPsqlVacancyRepository(db)
	cvRepo := psql.NewPsqlCVRepository(db)
	responseRepo := psql.NewPsqlResponseRepository(db)
	experienceRepo := psql.NewPsqlExperienceRepository(db)
	institutionRepo := psql.NewPsqlEducationInstitutionRepository(db)
	csatRepo := grpcRepo.NewGrpcCsatRepository(csatClient)
	skillRepo := psql.NewPsqlSkillRepository(db)

	searchEngineClient, err := initSearchEngineClient(searchConfig.SearchEngineServiceConfig)
	if err != nil {
		return err
	}
	searchEngineClientRepo := grpcRepo.NewGrpcSearchEngineRepository(searchEngineClient)

	sessionUsecase := usecase.NewSessionUsecase(authRepo, userRepo)
	userUsecase := usecase.NewUserUsecase(userRepo, authRepo)
	vacancyUsecase := usecase.NewVacancyUsecase(vacancyRepo, authRepo, userRepo, searchEngineClientRepo, skillRepo)
	cvUsecase := usecase.NewCVUsecase(cvRepo, experienceRepo, institutionRepo, authRepo, userRepo, responseRepo, vacancyRepo, searchEngineClientRepo, skillRepo)
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

	/*pinger := func(w http.ResponseWriter, r *http.Request) {
		pingCounter.WithLabelValues("GET").Inc()
		w.WriteHeader(200)
		w.Write([]byte("pong"))
	}*/

	prometheus.MustRegister(metrics.HitCounter, metrics.ErrorCounter)
	router.Handle("/metrics", promhttp.Handler())

	corsRouter := configs.CORS.Handler(router)
	loggedRouter := middleware.AccessLogMiddleware( /* logging.Logger,  */ corsRouter)
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
