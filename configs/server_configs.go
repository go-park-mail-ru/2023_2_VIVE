package configs

import (
	"os"

	"github.com/rs/cors"
)

var CURRENT_DIR, _ = os.Getwd()

const UPLOADS_DIR = "/assets/avatars/"

const (
	PORT         = ":8081"
	LOGFILE_NAME = "server.log"
)

var CORS = cors.New(cors.Options{
	AllowedOrigins: []string{
		"http://localhost:8082",
		"http://localhost:8083",
		"http://localhost:8084",
		"http://localhost:8085",
		"http://212.233.90.231:8082",
		"http://212.233.90.231:8083",
		"http://212.233.90.231:8084",
		"http://212.233.90.231:8085",
		"http://212.233.90.231:8086",
	},
	AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
	AllowCredentials: true,
})

var HnHRedisConfig = redisConfig{
	protocol:       "redis",
	networkAddress: "localhost",
	port:           "8008",
	password:       "vive_password_redis",
}

var HnHPostgresConfig = postgresConfig{
	user:     "vive_admin",
	password: "vive_password",
	dbname:   "hnh",
	host:     "localhost",
	port:     "8054",
	sslmode:  "disable",
}

type redisConfig struct {
	protocol       string
	networkAddress string
	port           string
	password       string
}

type postgresConfig struct {
	user     string
	password string
	dbname   string
	host     string
	port     string
	sslmode  string
}

func (rConf redisConfig) GetConnectionURL() string {
	return rConf.protocol + "://" + rConf.password + "@" + rConf.networkAddress + ":" + rConf.port
}

func (pConf postgresConfig) GetConnectionString() string {
	return "user=" + pConf.user + " password=" + pConf.password + " dbname=" + pConf.dbname +
		" host=" + pConf.host + " port=" + pConf.port +
		" sslmode=" + pConf.sslmode
}
