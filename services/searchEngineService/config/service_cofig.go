package config

type SearchEngineConfig struct {
	ServiceName string
	Host        string
	Port        int
	LogFile     string
}

var SearchEngineServiceConfig = SearchEngineConfig{
	ServiceName: "SearchEngine",
	Host:        "212.233.90.231",
	Port:        8063,
	LogFile:     "search_engine_service.log",
}
