package config

// AppConfig is application configuration
type AppConfig struct {
	HTTP HTTP
}

// HTTP contains information regardin server setup
type HTTP struct {
	Port int
	UseTLS bool
	CertFile string
	KeyFile string
}