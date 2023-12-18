package config

const (
	LOGS_DIR = "logs/"
)

type SearchEngineConfig struct {
	ServiceName string
	Host        string
	Port        int
	LogFile     string
}

var SearchEngineServiceConfig = SearchEngineConfig{
	ServiceName: "SearchEngine",
	Host:        "hnh_search",
	Port:        8063,
	LogFile:     "search_engine_service.log",
}
