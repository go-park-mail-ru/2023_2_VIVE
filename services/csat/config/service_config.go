package config

type CsatConfig struct {
	ServiceName string
	Host        string
	Port        int
	LogFile     string
}

var CsatServiceConfig = CsatConfig{
	ServiceName: "csat service",
	Host:        "84.23.53.171",
	Port:        8061,
	LogFile:     "csat_service.log",
}

type postgresConfig struct {
	user     string
	password string
	dbname   string
	host     string
	port     string
	sslmode  string
}

var CsatPostgresConfig = postgresConfig{
	user:     "vive_admin",
	password: "vive_password",
	dbname:   "hnh_csat",
	host:     "84.23.53.171",
	port:     "8055",
	sslmode:  "disable",
}

func (pConf postgresConfig) GetConnectionString() string {
	return "user=" + pConf.user + " password=" + pConf.password + " dbname=" + pConf.dbname +
		" host=" + pConf.host + " port=" + pConf.port +
		" sslmode=" + pConf.sslmode
}
