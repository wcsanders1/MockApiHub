package manager

import (
	"errors"
	"net/http"
	"os"
	"testing"

	"github.com/wcsanders1/MockApiHub/api"
	"github.com/wcsanders1/MockApiHub/config"
	"github.com/wcsanders1/MockApiHub/fake"
	"github.com/wcsanders1/MockApiHub/helper"
	"github.com/wcsanders1/MockApiHub/log"
	"github.com/wcsanders1/MockApiHub/wrapper"

	"github.com/stretchr/testify/mock"
)

func TestShowRegisteredMockAPIs_WritesJSON_WhenCalled(t *testing.T) {
	endpoints := map[string]config.Endpoint{
		"testEndpoint": config.Endpoint{
			Path:             "test",
			File:             "test",
			Method:           "GET",
			EnforceValidJSON: false,
		},
	}
	fakeAPI := new(api.FakeAPI)
	fakeAPI.On("GetBaseURL").Return("testURL")
	fakeAPI.On("GetPort").Return(4000)
	fakeAPI.On("GetEndpoints").Return(endpoints)
	fakeAPIs := make(map[string]api.IAPI)
	fakeAPIs["fakeAPI"] = fakeAPI
	mgr := Manager{
		apis: fakeAPIs,
		log:  log.GetFakeLogger(),
	}
	w := new(fake.ResponseWriter)
	w.On("Write", mock.AnythingOfType("[]uint8")).Return(1, nil)
	request, _ := http.NewRequest("GET", "test/url", nil)

	mgr.showRegisteredMockAPIs(w, request)

	fakeAPI.AssertCalled(t, "GetBaseURL")
	fakeAPI.AssertCalled(t, "GetPort")
	fakeAPI.AssertCalled(t, "GetEndpoints")
	w.AssertCalled(t, "Write", mock.AnythingOfType("[]uint8"))
}

func TestRefreshMockAPIs_WritesToResponse_WhenSuccessful(t *testing.T) {
	fileInfo := new(fake.FileInfo)
	fileInfo.On("Name").Return("testName")
	fileInfoCollection := []os.FileInfo{fileInfo}

	fileOps := new(wrapper.FakeFileOps)
	fileOps.On("ReadDir", mock.AnythingOfType("string")).Return(fileInfoCollection, nil)
	apiConfig := helper.GetFakeAPIConfig(4000)
	configManager := new(config.FakeManager)
	configManager.On("GetAPIConfig", mock.AnythingOfType("*fake.FileInfo")).Return(apiConfig, nil)
	mgr := Manager{
		log:           log.GetFakeLogger(),
		file:          fileOps,
		configManager: configManager,
		config:        helper.GetFakeAppConfig("certFile", "keyFile"),
	}
	w := new(fake.ResponseWriter)
	w.On("Write", mock.AnythingOfType("[]uint8")).Return(1, nil)
	request, _ := http.NewRequest("GET", "test/url", nil)

	mgr.refreshMockAPIs(w, request)

	w.AssertCalled(t, "Write", mock.AnythingOfType("[]uint8"))
}

func TestRefreshMockAPIs_DoesNotWriteToResponse_WhenFails(t *testing.T) {
	fileInfo := new(fake.FileInfo)
	fileInfo.On("Name").Return("testName")
	fileInfoCollection := []os.FileInfo{fileInfo}

	fileOps := new(wrapper.FakeFileOps)
	fileOps.On("ReadDir", mock.AnythingOfType("string")).Return(fileInfoCollection, errors.New(""))
	apiConfig := helper.GetFakeAPIConfig(4000)
	configManager := new(config.FakeManager)
	configManager.On("GetAPIConfig", mock.AnythingOfType("*fake.FileInfo")).Return(apiConfig, nil)
	mgr := Manager{
		log:           log.GetFakeLogger(),
		file:          fileOps,
		configManager: configManager,
		config:        helper.GetFakeAppConfig("certFile", "keyFile"),
	}
	w := new(fake.ResponseWriter)
	w.On("Write", mock.AnythingOfType("[]uint8")).Return(1, nil)
	request, _ := http.NewRequest("GET", "test/url", nil)

	mgr.refreshMockAPIs(w, request)

	w.AssertNotCalled(t, "Write", mock.Anything)
}
