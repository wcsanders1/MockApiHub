package manager

import (
	"errors"
	"os"
	"testing"

	"github.com/wcsanders1/MockApiHub/api"
	"github.com/wcsanders1/MockApiHub/config"
	"github.com/wcsanders1/MockApiHub/fake"
	"github.com/wcsanders1/MockApiHub/log"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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
		apis: make(map[string]api.IAPI),
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

func TestLoadMockAPIs(t *testing.T) {
	fileInfoCollection := []os.FileInfo{}
	fileInfoInner := new(fake.FileInfo)
	fileInfoInner.On("Name").Return("testconfig.toml")
	fileInfoCollection = append(fileInfoCollection, fileInfoInner)

	basicOpsIsAPI := new(fake.BasicOps)
	basicOpsIsAPI.On("ReadDir", mock.AnythingOfType("string")).Return(fileInfoCollection, nil)

	testAPIConfig := &config.APIConfig{
		HTTP: config.HTTP{
			Port: 4000,
		},
	}
	configManager := new(config.FakeManager)
	configManager.On("GetAPIConfig", mock.AnythingOfType("*fake.FileInfo")).Return(testAPIConfig, nil)

	mgr := Manager{
		file:          basicOpsIsAPI,
		configManager: configManager,
		log:           log.GetFakeLogger(),
		apis:          make(map[string]api.IAPI),
	}

	err := mgr.loadMockAPIs()
	assert := assert.New(t)
	assert.Nil(err)

	configManagerErr := new(config.FakeManager)
	configManagerErr.On("GetAPIConfig", mock.AnythingOfType("*fake.FileInfo")).Return(testAPIConfig, errors.New(""))

	mgrNoConfig := Manager{
		file:          basicOpsIsAPI,
		configManager: configManagerErr,
		log:           log.GetFakeLogger(),
		apis:          make(map[string]api.IAPI),
	}

	errNoConfig := mgrNoConfig.loadMockAPIs()
	assert.Nil(errNoConfig)
}
