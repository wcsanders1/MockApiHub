package manager

import (
	"testing"

	"github.com/wcsanders1/MockApiHub/api"
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

func TestApiByPortExists(t *testing.T) {
	testPort := 4000

	mgr := Manager{
		apis: make(map[string]*api.API),
	}

	firstAPIConfig := config.APIConfig{
		HTTP: config.HTTP{
			Port: 3999,
		},
	}
	secondAPIConfig := config.APIConfig{
		HTTP: config.HTTP{
			Port: 4001,
		},
	}

	firstAPI, err1 := api.NewAPI(&firstAPIConfig)
	secondAPI, err2 := api.NewAPI(&secondAPIConfig)

	mgr.apis["firstAPI"] = firstAPI
	mgr.apis["secondAPI"] = secondAPI

	result := mgr.apiByPortExists(testPort)
	assert := assert.New(t)
	assert.False(result)
	assert.Nil(err1)
	assert.Nil(err2)

	thirdAPIConfig := config.APIConfig{
		HTTP: config.HTTP{
			Port: testPort,
		},
	}

	thirdAPI, err3 := api.NewAPI(&thirdAPIConfig)
	mgr.apis["thirdAPI"] = thirdAPI

	result2 := mgr.apiByPortExists(testPort)
	assert.True(result2)
	assert.Nil(err3)
}
