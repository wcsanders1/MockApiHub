package manager

import (
	"testing"

	"github.com/wcsanders1/MockApiHub/config"

	"github.com/stretchr/testify/assert"
)

func TestNewManager(t *testing.T) {
	port := 4000
	cfg := &config.AppConfig{
		HTTP: config.HTTP{
			Port: port,
		},
	}

	result, err := NewManager(cfg)

	assert := assert.New(t)
	assert.Nil(err)
	assert.NotNil(result)

	badCfg := &config.AppConfig{
		HTTP: config.HTTP{
			Port: 0,
		},
	}

	result, err = NewManager(badCfg)

	assert.Nil(result)
	assert.Error(err)
}
