package log

import (
	"MockApiHub/config"
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
	lumberjack "gopkg.in/natefinch/lumberjack.v2"
)

const (
	defaultFilename       = "default"
	defaultMaxFileSize    = 20
	defaultMaxFileBackups = 20
	defaultMaxFileDaysAge = 20
)

// NewLogger returns a new instance of a logger
func NewLogger(config *config.LogConfig) (*logrus.Logger, error) {
	logFile := getLogFilename(config.Filename)

	file, err := os.OpenFile(logFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

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
		Filename:   file.Name(),
		MaxSize:    getMaxFileSize(config.MaxFileSize),
		MaxBackups: getMaxFileBackups(config.MaxFileBackups),
		MaxAge:     getMaxFileDaysAge(config.MaxFileDaysAge),
	}

	log.SetOutput(rotate)

	return log, nil
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
