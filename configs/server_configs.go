package configs

import "github.com/rs/cors"

const (
	PORT = ":8081"
)

var CORS = cors.New(cors.Options{
	AllowedOrigins:   []string{"http://212.233.90.231:8082"},
	AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
	AllowCredentials: true,
})
