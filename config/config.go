package config

// Config is application configuration
type Config struct {
	HTTP HTTP
}

// HTTP contains information regardin server setup
type HTTP struct {
	Port int
	UseTLS bool
	CertFile string
	KeyFile string
}