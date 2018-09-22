package log

import (
	"io/ioutil"
	"strings"

	"github.com/wcsanders1/MockApiHub/config"

	"github.com/sirupsen/logrus"
	lumberjack "gopkg.in/natefinch/lumberjack.v2"
)

const (
	defaultFilename       = "default.log"
	defaultMaxFileSize    = 20
	defaultMaxFileBackups = 20
	defaultMaxFileDaysAge = 20
	pkgField              = "pkg"

	// FuncField is the name of the log field denoting the name of the function doing the logging
	FuncField = "func"

	// PortField is the name of the log field denoting a server port number
	PortField = "port"

	// MethodField is the name of the log field denoting the HTTP method name
	MethodField = "method"

	// PathField is the name of the log field denoting an HTTP path
	PathField = "path"

	// UseTLSField is the name of the log field denoting whether a server is configured to use TLS
	UseTLSField = "useTLS"

	// CertFileField is the name of the log field denoting the configured path to a server's TLS certificate
	CertFileField = "certFile"

	// KeyFileField is the name of the log field denoting the configured path to a server's TLS key
	KeyFileField = "keyFile"

	// DefaultCertFileField is the name of the log field denoting the default cert file
	DefaultCertFileField = "defaultCertFile"

	// DefaultKeyFileField is the name of the log field denoting the default key file
	DefaultKeyFileField = "defaultKeyFile"

	// BaseURLField is the name of the log field denoting the base URL of an API
	BaseURLField = "baseURL"

	// RouteField is the name of the log field denoting a registered route in an API
	RouteField = "route"

	// FileField is the name of the log field denoting a file that a mock API route will serve
	FileField = "file"

	// APIDirField is the name of the log field denoting the directory a mock API is in
	APIDirField = "apiDir"

	// APINameField is the name of the log field denoting the name of a mock API
	APINameField = "mockAPIName"

	// EndpointNameField is the name of the log field denoting the name of an endpoint
	EndpointNameField = "endpointName"
)

// NewLogger returns a new instance of a logger
func NewLogger(config *config.Log, pkgName string) *logrus.Entry {
	log := logrus.New()

	if !config.LoggingEnabled {
		log.Out = ioutil.Discard
		return log.WithField("", "")
	}

	if config.FormatAsJSON {
		log.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: "2006-01-02 15:04:05",
			PrettyPrint:     config.PrettyJSON,
		})
	} else {
		log.SetFormatter(&logrus.TextFormatter{
			TimestampFormat: "2006-01-02 15:04:05",
			FullTimestamp:   true,
		})
	}

	rotate := &lumberjack.Logger{
		Filename:   getLogFilename(config.Filename),
		MaxSize:    getMaxFileSize(config.MaxFileSize),
		MaxBackups: getMaxFileBackups(config.MaxFileBackups),
		MaxAge:     getMaxFileDaysAge(config.MaxFileDaysAge),
	}

	log.SetOutput(rotate)
	log.SetLevel(getLogLevel(config.Level))
	return log.WithField(pkgField, pkgName)
}

func getLogFilename(filename string) string {
	if len(filename) == 0 {
		return defaultFilename
	}
	return filename
}

func getMaxFileSize(maxSize int) int {
	if maxSize < 1 {
		return defaultMaxFileSize
	}
	return maxSize
}

func getMaxFileBackups(maxFileBackups int) int {
	if maxFileBackups < 1 {
		return defaultMaxFileBackups
	}
	return maxFileBackups
}

func getMaxFileDaysAge(maxFileDaysAge int) int {
	if maxFileDaysAge < 1 {
		return defaultMaxFileDaysAge
	}
	return maxFileDaysAge
}

func getLogLevel(level string) logrus.Level {
	lowerCaseLevel := strings.ToLower(level)
	switch lowerCaseLevel {
	case "debug":
		return logrus.DebugLevel
	case "info":
		return logrus.InfoLevel
	case "warn":
		return logrus.WarnLevel
	case "error":
		return logrus.ErrorLevel
	case "fatal":
		return logrus.FatalLevel
	case "panic":
		return logrus.PanicLevel
	default:
		return logrus.DebugLevel
	}
}
