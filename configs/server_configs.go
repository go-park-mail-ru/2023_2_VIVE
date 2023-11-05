package configs

import "github.com/rs/cors"

const (
	PORT = ":8081"
)

var CORS = cors.New(cors.Options{
	AllowedOrigins:   []string{"http://212.233.90.231:8082", "http://212.233.90.231:8083", "http://212.233.90.231:8084", "http://212.233.90.231:8085"},
	AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
	AllowCredentials: true,
})

type postgresConfig struct {
	user     string
	password string
	dbname   string
	host     string
	port     string
	sslmode  string
}

func (pConf postgresConfig) GetConnectionString() string {
	return "user=" + pConf.user + " password=" + pConf.password + " dbname=" + pConf.dbname +
		" host=" + pConf.host + " port=" + pConf.port +
		" sslmode=" + pConf.sslmode
}

var HnHPostgresConfig = postgresConfig{
	user:     "vive_admin",
	password: "vive_password",
	dbname:   "hnh",
	host:     "212.233.90.231",
	port:     "8054",
	sslmode:  "disable",
}
