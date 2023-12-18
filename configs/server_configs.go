package configs

import (
	"os"

	"github.com/rs/cors"
)

const (
	SERVER_ADDRESS = "http://84.23.53.171:8081"
	SERVER_DOMAIN  = "https://hunt-n-hire/api"
)

var CURRENT_DIR, _ = os.Getwd()

const (
	PORT         = ":8081"
	LOGS_DIR     = "/logs/"
	LOGFILE_NAME = "server.log"
	UPLOADS_DIR  = "/assets/avatars/"
)

var CORS = cors.New(cors.Options{
	AllowedOrigins: []string{
		"http://localhost:8082",
		"http://localhost:8083",
		"http://localhost:8084",
		"http://localhost:8085",
		"http://84.23.53.171:8082",
		"http://84.23.53.171:8083",
		"http://84.23.53.171:8084",
		"http://84.23.53.171:8085",
		"http://84.23.53.171:8086",
	},
	AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
	AllowCredentials: true,
	ExposedHeaders:   []string{"Content-Disposition"},
})

// var HnHRedisConfig = redisConfig{
// 	protocol:       "redis",
// 	networkAddress: "localhost",
// 	port:           "8008",
// 	password:       "vive_password_redis",
// }

var HnHPostgresConfig = postgresConfig{
	user:     "vive_admin",
	password: os.Getenv("POSTGRES_PASSWORD"),
	dbname:   "hnh",
	host:     "db_hnh",
	port:     "5432",
	sslmode:  "disable",
}

// type redisConfig struct {
// 	protocol       string
// 	networkAddress string
// 	port           string
// 	password       string
// }

type postgresConfig struct {
	user     string
	password string
	dbname   string
	host     string
	port     string
	sslmode  string
}

// func (rConf redisConfig) GetConnectionURL() string {
// 	return rConf.protocol + "://" + rConf.password + "@" + rConf.networkAddress + ":" + rConf.port
// }

func (pConf postgresConfig) GetConnectionString() string {
	return "user=" + pConf.user + " password=" + pConf.password + " dbname=" + pConf.dbname +
		" host=" + pConf.host + " port=" + pConf.port +
		" sslmode=" + pConf.sslmode
}
