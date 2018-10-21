package manager

import (
	"errors"
	"testing"

	"github.com/wcsanders1/MockApiHub/api"
	"github.com/wcsanders1/MockApiHub/config"
	"github.com/wcsanders1/MockApiHub/fake"
	"github.com/wcsanders1/MockApiHub/helper"
	"github.com/wcsanders1/MockApiHub/log"
	"github.com/wcsanders1/MockApiHub/wrapper"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNewManager_ReturnsManager_WhenConfigValid(t *testing.T) {
	cfg := &config.AppConfig{
		HTTP: config.HTTP{
			Port: 4000,
		},
	}

	result, err := NewManager(cfg)

	assert := assert.New(t)
	assert.Nil(err)
	assert.NotNil(result)
	assert.IsType(&Manager{}, result)
}

func TestNewManager_ReturnsError_WhenPortNotProvided(t *testing.T) {
	cfg := &config.AppConfig{
		HTTP: config.HTTP{
			Port: 0,
		},
	}

	result, err := NewManager(cfg)

	assert := assert.New(t)
	assert.Nil(result)
	assert.Error(err)
}

func TestAPIByPortExists_ReturnsFalse_WhenProvidedUnregisteredPort(t *testing.T) {
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
	mgr.apis["firstAPI"], _ = api.NewAPI(&firstAPIConfig)
	mgr.apis["secondAPI"], _ = api.NewAPI(&secondAPIConfig)

	result := mgr.apiByPortExists(4000)

	assert.False(t, result)
}

func TestAPIPortExists_ReturnsTrue_WhenProvidedRegisteredPort(t *testing.T) {
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
	mgr.apis["firstAPI"], _ = api.NewAPI(&firstAPIConfig)
	mgr.apis["secondAPI"], _ = api.NewAPI(&secondAPIConfig)

	result := mgr.apiByPortExists(4001)

	assert.True(t, result)

}

func TestLoadMockAPIs_ReturnsNil_WhenProvidedValidAPI(t *testing.T) {
	_, fileCollection := helper.GetFakeFileInfoAndCollection("", "testconfig.toml")
	fileOps := new(wrapper.FakeFileOps)
	fileOps.On("ReadDir", mock.AnythingOfType("string")).Return(fileCollection, nil)
	apiConfig := helper.GetFakeAPIConfig(4000)
	configManager := new(config.FakeManager)
	configManager.On("GetAPIConfig", mock.AnythingOfType("*fake.FileInfo")).Return(apiConfig, nil)
	mgr := Manager{
		file:          fileOps,
		configManager: configManager,
		log:           log.GetFakeLogger(),
		apis:          make(map[string]api.IAPI),
	}

	err := mgr.loadMockAPIs()

	assert := assert.New(t)
	assert.Nil(err)
	configManager.AssertCalled(t, "GetAPIConfig", mock.AnythingOfType("*fake.FileInfo"))
}

func TestLoadMockAPIs_ReturnsNil_WhenGetConfigFails(t *testing.T) {
	_, fileCollection := helper.GetFakeFileInfoAndCollection("", "testconfig.toml")
	fileOps := new(wrapper.FakeFileOps)
	fileOps.On("ReadDir", mock.AnythingOfType("string")).Return(fileCollection, nil)
	apiConfig := helper.GetFakeAPIConfig(4000)
	configManager := new(config.FakeManager)
	configManager.On("GetAPIConfig", mock.AnythingOfType("*fake.FileInfo")).Return(apiConfig, errors.New(""))
	mgr := Manager{
		file:          fileOps,
		configManager: configManager,
		log:           log.GetFakeLogger(),
		apis:          make(map[string]api.IAPI),
	}

	err := mgr.loadMockAPIs()

	assert := assert.New(t)
	assert.Nil(err)
	configManager.AssertCalled(t, "GetAPIConfig", mock.AnythingOfType("*fake.FileInfo"))
}

func TestLoadMockAPIs_ReturnsError_WhenReadDirFails(t *testing.T) {
	_, fileCollection := helper.GetFakeFileInfoAndCollection("", "testconfig.toml")
	fileOps := new(wrapper.FakeFileOps)
	fileOps.On("ReadDir", mock.AnythingOfType("string")).Return(fileCollection, errors.New(""))
	configManager := new(config.FakeManager)
	mgr := Manager{
		file:          fileOps,
		configManager: configManager,
		log:           log.GetFakeLogger(),
		apis:          make(map[string]api.IAPI),
	}

	err := mgr.loadMockAPIs()

	assert := assert.New(t)
	assert.Error(err)
	configManager.AssertNotCalled(t, "GetAPIConfig", mock.Anything)
}

func TestLoadMockAPIs_LoadsOneAPI_WhenProvidedTwoWithSamePort(t *testing.T) {
	_, fileCollection := helper.GetFakeFileInfoAndCollection("", "testconfig.toml")
	fileInfo, _ := helper.GetFakeFileInfoAndCollection("", "testconfig2.toml")
	fileCollection = append(fileCollection, fileInfo)
	fileOps := new(wrapper.FakeFileOps)
	fileOps.On("ReadDir", mock.AnythingOfType("string")).Return(fileCollection, nil)
	apiConfig := helper.GetFakeAPIConfig(4000)
	configManager := new(config.FakeManager)
	configManager.On("GetAPIConfig", mock.AnythingOfType("*fake.FileInfo")).Return(apiConfig, nil)
	mgr := Manager{
		file:          fileOps,
		configManager: configManager,
		log:           log.GetFakeLogger(),
		apis:          make(map[string]api.IAPI),
	}

	err := mgr.loadMockAPIs()

	assert := assert.New(t)
	assert.Nil(err)
	assert.Equal(1, len(mgr.apis))
	configManager.AssertCalled(t, "GetAPIConfig", mock.AnythingOfType("*fake.FileInfo"))
}

func TestLoadMockAPIs_DoesNotLoadAPI_WhenPortIsZero(t *testing.T) {
	_, fileCollection := helper.GetFakeFileInfoAndCollection("", "testconfig.toml")
	fileOps := new(wrapper.FakeFileOps)
	fileOps.On("ReadDir", mock.AnythingOfType("string")).Return(fileCollection, nil)
	apiConfig := helper.GetFakeAPIConfig(0)
	configManager := new(config.FakeManager)
	configManager.On("GetAPIConfig", mock.AnythingOfType("*fake.FileInfo")).Return(apiConfig, nil)
	mgr := Manager{
		file:          fileOps,
		configManager: configManager,
		log:           log.GetFakeLogger(),
		apis:          make(map[string]api.IAPI),
	}

	err := mgr.loadMockAPIs()

	assert := assert.New(t)
	assert.Nil(err)
	assert.Empty(mgr.apis)
}

func TestRegisterMockAPIs_RegistersAPI_WithCertAndKey(t *testing.T) {
	certFile := "testCert"
	keyFile := "testKey"
	dir := "fakeAPI"
	fakeAPI := new(api.FakeAPI)
	fakeAPI.On("GetBaseURL").Return("baseURL")
	fakeAPI.On("GetPort").Return(4000)
	fakeAPI.On("Register", mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(nil)
	apis := map[string]api.IAPI{
		dir: fakeAPI,
	}
	mgr := Manager{
		apis:   apis,
		config: helper.GetFakeAppConfig(certFile, keyFile),
		log:    log.GetFakeLogger(),
	}

	mgr.registerMockAPIs()

	fakeAPI.AssertCalled(t, "Register", dir, certFile, keyFile)
}

func TestRegisterMockAPIs_DoesNotPanic_WhenRegisterFails(t *testing.T) {
	certFile := "testCert"
	keyFile := "testKey"
	dir := "fakeAPI"
	fakeAPI := new(api.FakeAPI)
	fakeAPI.On("GetBaseURL").Return("baseURL")
	fakeAPI.On("GetPort").Return(4000)
	fakeAPI.On("Register", mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(errors.New(""))
	apis := map[string]api.IAPI{
		dir: fakeAPI,
	}
	mgr := Manager{
		apis:   apis,
		config: helper.GetFakeAppConfig(certFile, keyFile),
		log:    log.GetFakeLogger(),
	}

	assert.NotPanics(t, func() { mgr.registerMockAPIs() })
	fakeAPI.AssertCalled(t, "Register", dir, certFile, keyFile)
}

func TestStartHubServerUsingTLS(t *testing.T) {
	certFile := "testCertFile"
	keyFile := "testKeyFile"
	fakeLogger := log.GetFakeLogger()
	assert := assert.New(t)

	basicFileOpsNoErr := new(wrapper.FakeFileOps)
	basicFileOpsNoErr.On("Stat", mock.AnythingOfType("string")).Return(new(fake.FileInfo), nil)

	fakeConfig := &config.AppConfig{
		HTTP: config.HTTP{
			CertFile: certFile,
			KeyFile:  keyFile,
		},
	}

	fakeServerNoErr := new(wrapper.FakeServerOps)
	fakeServerNoErr.On("ListenAndServeTLS", mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(nil)

	mgrNoErr := Manager{
		config: fakeConfig,
		log:    fakeLogger,
		file:   basicFileOpsNoErr,
		server: fakeServerNoErr,
	}

	resultNoErr := mgrNoErr.startHubServerUsingTLS()

	assert.Nil(resultNoErr)
	basicFileOpsNoErr.AssertCalled(t, "Stat", certFile)
	basicFileOpsNoErr.AssertCalled(t, "Stat", keyFile)
	fakeServerNoErr.AssertCalled(t, "ListenAndServeTLS", certFile, keyFile)

	basicFileOpsCertErr := new(wrapper.FakeFileOps)
	basicFileOpsCertErr.On("Stat", certFile).Return(new(fake.FileInfo), errors.New(""))
	fakeServerCertErr := new(wrapper.FakeServerOps)

	mgrCertErr := Manager{
		config: fakeConfig,
		log:    fakeLogger,
		file:   basicFileOpsCertErr,
		server: fakeServerCertErr,
	}

	resultCertErr := mgrCertErr.startHubServerUsingTLS()
	assert.Error(resultCertErr)
	basicFileOpsCertErr.AssertCalled(t, "Stat", certFile)
	basicFileOpsCertErr.AssertNotCalled(t, "Stat", keyFile)
	fakeServerCertErr.AssertNotCalled(t, "ListenAndServeTLS", mock.Anything, mock.Anything)

	basicFileOpsKeyErr := new(wrapper.FakeFileOps)
	basicFileOpsKeyErr.On("Stat", certFile).Return(new(fake.FileInfo), nil)
	basicFileOpsKeyErr.On("Stat", keyFile).Return(new(fake.FileInfo), errors.New(""))
	fakeServerKeyErr := new(wrapper.FakeServerOps)

	mgrKeyErr := Manager{
		config: fakeConfig,
		log:    fakeLogger,
		file:   basicFileOpsKeyErr,
		server: fakeServerKeyErr,
	}

	resultKeyErr := mgrKeyErr.startHubServerUsingTLS()
	assert.Error(resultKeyErr)
	basicFileOpsKeyErr.AssertCalled(t, "Stat", certFile)
	basicFileOpsKeyErr.AssertCalled(t, "Stat", keyFile)
	fakeServerCertErr.AssertNotCalled(t, "ListenAndServeTLS", mock.Anything, mock.Anything)

	basicFileOpsServerErr := new(wrapper.FakeFileOps)
	basicFileOpsServerErr.On("Stat", mock.AnythingOfType("string")).Return(new(fake.FileInfo), nil)
	fakeServerErr := new(wrapper.FakeServerOps)
	fakeServerErr.On("ListenAndServeTLS", mock.Anything, mock.Anything).Return(errors.New(""))

	mgrSrvErr := Manager{
		config: fakeConfig,
		log:    fakeLogger,
		file:   basicFileOpsServerErr,
		server: fakeServerErr,
	}

	resultSrvErr := mgrSrvErr.startHubServerUsingTLS()
	assert.Error(resultSrvErr)
	basicFileOpsServerErr.AssertCalled(t, "Stat", certFile)
	basicFileOpsServerErr.AssertCalled(t, "Stat", keyFile)
	fakeServerErr.AssertCalled(t, "ListenAndServeTLS", certFile, keyFile)
}
