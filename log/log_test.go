package log

import (
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestGetLogFilename_ReturnsDefault_IfNoNameProvided(t *testing.T) {
	assert.Equal(t, defaultFilename, getLogFilename(""))
}

func TestGetLogFilename_ReturnsFilename_WhenFilenameProvided(t *testing.T) {
	assert.Equal(t, "fileName", getLogFilename("fileName"))
}

func TestGetMaxFileSize_ReturnsDefault_IfNoSizeProvided(t *testing.T) {
	assert.Equal(t, defaultMaxFileSize, getMaxFileSize(0))
}

func TestGetMaxFileSize_ReturnsFileSize_WhenProvided(t *testing.T) {
	assert.Equal(t, 50, getMaxFileSize(50))
}

func TestGetMaxFileBackups_ReturnsDefault_WhenNoAmountProvided(t *testing.T) {
	assert.Equal(t, defaultMaxFileBackups, getMaxFileBackups(0))
}

func TestGetMaxFileBackups_ReturnsMaxBackups_WhenProvided(t *testing.T) {
	assert.Equal(t, 50, getMaxFileBackups(50))
}

func TestGetMaxFileDaysAge_ReturnsDefault_WhenDaysNotProvided(t *testing.T) {
	assert.Equal(t, defaultMaxFileDaysAge, getMaxFileDaysAge(0))
}

func TestGetMaxFileDaysAge_ReturnsDays_WhenDaysProvided(t *testing.T) {
	assert.Equal(t, 50, getMaxFileDaysAge(50))
}

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
