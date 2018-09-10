package config

type (
	// AppConfig is application configuration
	AppConfig struct {
		HTTP HTTP
		Log  Log
	}

	// APIConfig is configuration for an individual mock api
	APIConfig struct {
		HTTP      HTTP
		BaseURL   string
		Endpoints map[string]Endpoint
	}

	// Log is configuration for logging
	Log struct {
		Filename       string
		MaxFileSize    int
		MaxFileBackups int
		MaxFileDaysAge int
		FormatAsJSON   bool
		Level          string
	}

	// HTTP contains information regarding server setup
	HTTP struct {
		Port     int
		UseTLS   bool
		CertFile string
		KeyFile  string
	}

	// Endpoint contains information regarding an endpoint
	Endpoint struct {
		Path   string
		File   string
		Method string
	}
)
