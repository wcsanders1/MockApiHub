package config

// Config is application configuration
type Config struct {
	HTTP http
	Apis map[string]api
}

type http struct {
	Port int
	SSL bool
}

type api struct {
	BaseURL string
	Endpoints []endpoint
}

type endpoint struct {
	Path string
	File string
}