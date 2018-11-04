package api

import (
	"net/http"

	"github.com/stretchr/testify/mock"
	"github.com/wcsanders1/MockApiHub/config"
)

// FakeAPI is a mockable API.
type FakeAPI struct {
	mock.Mock
}

// Start is a mockable api.Start().
func (api *FakeAPI) Start(dir, defaultCert, defaultKey string) error {
	args := api.Called(dir, defaultCert, defaultKey)
	return args.Error(0)
}

// Shutdown is a mockable api.Shutdown().
func (api *FakeAPI) Shutdown() error {
	args := api.Called()
	return args.Error(0)
}

// ServeHTTP is a mockable api.ServeHTTP().
func (api *FakeAPI) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	return
}

// GetPort is a mockable api.GetPort().
func (api *FakeAPI) GetPort() int {
	args := api.Called()
	return args.Int(0)
}

// GetBaseURL is a mockable api.GetBaseURL().
func (api *FakeAPI) GetBaseURL() string {
	args := api.Called()
	return args.String(0)
}

// GetEndpoints is a mockable api.GetEndpoints().
func (api *FakeAPI) GetEndpoints() map[string]config.Endpoint {
	args := api.Called()
	return args.Get(0).(map[string]config.Endpoint)
}
