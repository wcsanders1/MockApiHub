package manager

import (
	"errors"
	"os"
	"testing"

	"github.com/wcsanders1/MockApiHub/api"
	"github.com/wcsanders1/MockApiHub/config"
	"github.com/wcsanders1/MockApiHub/fake"
	"github.com/wcsanders1/MockApiHub/log"
	"github.com/wcsanders1/MockApiHub/wrapper"

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

	basicOpsIsAPI := new(wrapper.FakeFileOps)
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
	configManager.AssertCalled(t, "GetAPIConfig", mock.AnythingOfType("*fake.FileInfo"))

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

	basicOpsReadDirErr := new(wrapper.FakeFileOps)
	basicOpsReadDirErr.On("ReadDir", mock.AnythingOfType("string")).Return([]os.FileInfo{}, errors.New(""))
	configMgrReadDirErr := new(config.FakeManager)

	mgrReadDirErr := Manager{
		file:          basicOpsReadDirErr,
		configManager: configMgrReadDirErr,
		log:           log.GetFakeLogger(),
	}

	errReadDir := mgrReadDirErr.loadMockAPIs()
	assert.Error(errReadDir)
	configMgrReadDirErr.AssertNotCalled(t, "GetAPIConfig", mock.Anything)

	configManagerDupPort := new(config.FakeManager)
	configManagerDupPort.On("GetAPIConfig", mock.AnythingOfType("*fake.FileInfo")).Return(testAPIConfig, nil)
	fileInfoInner2 := new(fake.FileInfo)
	fileInfoInner2.On("Name").Return("testconfig2.toml")
	fileInfoCollection = append(fileInfoCollection, fileInfoInner2)

	basicOpsDupPort := new(wrapper.FakeFileOps)
	basicOpsDupPort.On("ReadDir", mock.AnythingOfType("string")).Return(fileInfoCollection, nil)

	mgrDupPort := Manager{
		file:          basicOpsDupPort,
		configManager: configManagerDupPort,
		log:           log.GetFakeLogger(),
		apis:          make(map[string]api.IAPI),
	}

	errDupPort := mgrDupPort.loadMockAPIs()
	assert.Nil(errDupPort)
	assert.Equal(1, len(mgrDupPort.apis))

	testAPIConfigPortZero := &config.APIConfig{
		HTTP: config.HTTP{
			Port: 0,
		},
	}

	configMgrPortZero := new(config.FakeManager)
	configMgrPortZero.On("GetAPIConfig", mock.AnythingOfType("*fake.FileInfo")).Return(testAPIConfigPortZero, nil)
	basicOpsPortZero := new(wrapper.FakeFileOps)
	basicOpsPortZero.On("ReadDir", mock.AnythingOfType("string")).Return(fileInfoCollection, nil)

	mgrPortZero := Manager{
		file:          basicOpsPortZero,
		configManager: configMgrPortZero,
		log:           log.GetFakeLogger(),
		apis:          make(map[string]api.IAPI),
	}

	errPortZero := mgrPortZero.loadMockAPIs()
	assert.Nil(errPortZero)
	assert.Empty(mgrPortZero.apis)
}

func TestRegisterMockAPIs(t *testing.T) {
	baseURL := "baseURL"
	port := 4000
	certFile := "testCert"
	keyFile := "testKey"
	dir := "fakeAPI"
	fakeAPI := new(api.FakeAPI)
	fakeAPI.On("GetBaseURL").Return(baseURL)
	fakeAPI.On("GetPort").Return(port)
	fakeAPI.On("Register", mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(nil)

	apis := map[string]api.IAPI{
		dir: fakeAPI,
	}

	fakeConfig := &config.AppConfig{
		HTTP: config.HTTP{
			CertFile: certFile,
			KeyFile:  keyFile,
		},
	}

	mgrNoErr := Manager{
		apis:   apis,
		config: fakeConfig,
		log:    log.GetFakeLogger(),
	}

	mgrNoErr.registerMockAPIs()
	fakeAPI.AssertCalled(t, "Register", dir, certFile, keyFile)

	fakeAPIErr := new(api.FakeAPI)
	fakeAPIErr.On("GetBaseURL").Return(baseURL)
	fakeAPIErr.On("GetPort").Return(port)
	fakeAPIErr.On("Register", mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(errors.New(""))

	apisErr := map[string]api.IAPI{
		dir: fakeAPIErr,
	}

	mgrErr := Manager{
		apis:   apisErr,
		config: fakeConfig,
		log:    log.GetFakeLogger(),
	}

	mgrErr.registerMockAPIs()
	fakeAPIErr.AssertCalled(t, "Register", dir, certFile, keyFile)
}
