package log

import (
	"github.com/sirupsen/logrus"
	"github.com/wcsanders1/MockApiHub/config"
)

// GetFakeLogger returns a fake logger that logs nothing.
func GetFakeLogger() *logrus.Entry {
	config := config.Log{
		LoggingEnabled: false,
	}
	return NewLogger(&config, "fake")
}
