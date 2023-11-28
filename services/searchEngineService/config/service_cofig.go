package config

type SearchEngineConfig struct {
	ServiceName string
	Host        string
	Port        int
	LogFile     string
}

var SearchEngineServiceConfig = SearchEngineConfig{
	ServiceName: "SearchEngine",
	Host:        "localhost",
	Port:        8063,
	LogFile:     "search_engine_service.log",
}
