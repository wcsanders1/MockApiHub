package config

// Config is application configuration
type Config struct {
	HTTP http
}

type http struct {
	Port int
	SSL bool
}