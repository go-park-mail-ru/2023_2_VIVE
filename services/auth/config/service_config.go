package config

const (
	LOGS_DIR = "logs/"
)

type AuthConfig struct {
	ServiceName string
	Host        string
	Port        int
	LogFile     string
}

var AuthServiceConfig = AuthConfig{
	ServiceName: "auth service",
	Host:        "hnh_auth",
	Port:        8062,
	LogFile:     "auth_service.log",
}

type redisConfig struct {
	protocol       string
	networkAddress string
	port           string
	password       string
}

var AuthRedisConfig = redisConfig{
	protocol:       "redis",
	networkAddress: "sessions_hnh",
	port:           "6379",
	password:       "vive_password_redis",
}

func (rConf redisConfig) GetConnectionURL() string {
	return rConf.protocol + "://" + rConf.password + "@" + rConf.networkAddress + ":" + rConf.port
}
