package configs

import "github.com/rs/cors"

const (
	PORT         = ":8081"
	LOGFILE_NAME = "server.log"
)

var CORS = cors.New(cors.Options{
	AllowedOrigins:   []string{"http://212.233.90.231:8082", "http://212.233.90.231:8083", "http://212.233.90.231:8084", "http://212.233.90.231:8085"},
	AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
	AllowCredentials: true,
})
