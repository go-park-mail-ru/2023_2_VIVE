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

type redisConfig struct {
	protocol       string
	networkAddress string
	port           string
	password       string
}

func (rConf redisConfig) GetConnectionURL() string {
	return rConf.protocol + "://" + rConf.password + "@" + rConf.networkAddress + ":" + rConf.port
}

var HnHRedisConfig = redisConfig{
	protocol:       "redis",
	networkAddress: "212.233.90.231",
	port:           "8008",
	password:       "vive_password_redis",
}
