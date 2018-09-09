package log

import (
	"MockApiHub/config"

	"github.com/sirupsen/logrus"
	lumberjack "gopkg.in/natefinch/lumberjack.v2"
)

const (
	defaultFilename       = "default.log"
	defaultMaxFileSize    = 20
	defaultMaxFileBackups = 20
	defaultMaxFileDaysAge = 20

	// FuncField is the name of the log field denoting the name of the function doing the logging
	FuncField = "func"

	// PortField is the name of the log field denoting a server port number
	PortField = "port"

	// PkgField is the name of the log field denoting the name of the package doing the logging
	PkgField = "pkg"

	// MethodField is the name of the log field denoting the HTTP method name
	MethodField = "method"

	// PathField is the name of the log field denoting an HTTP path
	PathField = "path"
)

// NewLogger returns a new instance of a logger
func NewLogger(config *config.Log) *logrus.Logger {
	log := logrus.New()
	if config.FormatAsJSON {
		log.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: "2006-01-02 15:04:05",
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
	return log
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
