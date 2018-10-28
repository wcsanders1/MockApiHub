package manager

import (
	"net/http"
	"testing"

	"github.com/wcsanders1/MockApiHub/api"
	"github.com/wcsanders1/MockApiHub/config"
	"github.com/wcsanders1/MockApiHub/fake"
	"github.com/wcsanders1/MockApiHub/log"

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
