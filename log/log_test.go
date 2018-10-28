package log

import (
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestGetLogLevel_ReturnsDebugLevel_ByDefault(t *testing.T) {
	assert.Equal(t, logrus.DebugLevel, getLogLevel("noLevel"))
}

func TestGetLogLevel_ReturnsPanicLevel_WhenLevelPanic(t *testing.T) {
	assert.Equal(t, logrus.PanicLevel, getLogLevel("panic"))
}

func TestGetLogLevel_ReturnsFatalLevel_WhenLevelFatal(t *testing.T) {
	assert.Equal(t, logrus.FatalLevel, getLogLevel("fatal"))
}

func TestGetLogLevel_ReturnsErrorLevel_WhenLevelFatal(t *testing.T) {
	assert.Equal(t, logrus.ErrorLevel, getLogLevel("error"))
}

func TestGetLogLevel_ReturnsWarnLevel_WhenLevelWarn(t *testing.T) {
	assert.Equal(t, logrus.WarnLevel, getLogLevel("warn"))
}

func TestGetLogLevel_ReturnsInfoLevel_WhenLevelInfo(t *testing.T) {
	assert.Equal(t, logrus.InfoLevel, getLogLevel("info"))
}

func TestGetLogLevel_ReturnsDebug_WhenLevelDebug(t *testing.T) {
	assert.Equal(t, logrus.DebugLevel, getLogLevel("debug"))
}
