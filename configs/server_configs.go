package configs

import "github.com/rs/cors"

const (
	PORT = ":8081"
)

var CORS = cors.New(cors.Options{
	AllowedOrigins:   []string{"*"},
	AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
	AllowCredentials: true,
})
