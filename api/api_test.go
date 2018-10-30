package api

import (
	"net/http"
	"testing"

	"github.com/wcsanders1/MockApiHub/config"
	"github.com/wcsanders1/MockApiHub/route"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetPort_ReturnsPort_WhenCalled(t *testing.T) {
	port := 4000
	testAPI := API{
		httpConfig: config.HTTP{
			Port: port,
		},
	}

	result := testAPI.GetPort()

	assert := assert.New(t)
	assert.NotNil(result)
	assert.Equal(port, result)
}

func TestGetBaseURL_ReturnsBaseURL_WhenCalled(t *testing.T) {
	baseURL := "base/url"
	testAPI := API{
		baseURL: baseURL,
	}

	result := testAPI.GetBaseURL()

	assert := assert.New(t)
	assert.NotNil(result)
	assert.Equal(baseURL, result)
}

func TestGetEndpoints_ReturnsEndpoints_WhenCalled(t *testing.T) {
	endpointName := "testEndpoint"
	endpoints := map[string]config.Endpoint{
		endpointName: config.Endpoint{},
	}
	testAPI := API{
		endpoints: endpoints,
	}

	result := testAPI.GetEndpoints()

	assert := assert.New(t)
	assert.NotNil(result)
	assert.NotEmpty(result)
	assert.IsType(map[string]config.Endpoint{}, result)
	assert.Contains(result, endpointName)
	assert.Equal(1, len(result))
}

func TestEnsureRouteRegistered_AddsRoute_IfNotRegistered(t *testing.T) {
	url := "test/url"
	fakeRouteTree := &route.FakeTree{}
	fakeRouteTree.On("GetRoute", mock.AnythingOfType("string")).Return("", map[string]string{}, nil)
	fakeRouteTree.On("AddRoute", mock.AnythingOfType("string")).Return(url, nil)
	api := &API{
		routeTree: fakeRouteTree,
	}

	result := api.ensureRouteRegistered(url)

	assert := assert.New(t)
	assert.NotEmpty(result)
	assert.Equal(url, result)
	fakeRouteTree.AssertCalled(t, "GetRoute", url)
	fakeRouteTree.AssertCalled(t, "AddRoute", url)
}

func TestEnsureRouteRegistered_DoesNotAddRoute_IfRegistered(t *testing.T) {
	url := "test/url"
	fakeRouteTree := &route.FakeTree{}
	fakeRouteTree.On("GetRoute", mock.AnythingOfType("string")).Return(url, map[string]string{}, nil)
	api := &API{
		routeTree: fakeRouteTree,
	}

	result := api.ensureRouteRegistered(url)

	assert := assert.New(t)
	assert.NotEmpty(result)
	assert.Equal(url, result)
	fakeRouteTree.AssertCalled(t, "GetRoute", url)
	fakeRouteTree.AssertNotCalled(t, "AddRoute", mock.Anything)
}

func TestCreateAPIServer_ReturnsServer_WhenPortProvided(t *testing.T) {
	httpConfig := &config.HTTP{
		Port: 4000,
	}
	api, _ := NewAPI(&config.APIConfig{})

	result, err := createAPIServer(httpConfig, api)

	assert := assert.New(t)
	assert.NoError(err)
	assert.NotNil(result)
	assert.IsType(&http.Server{}, result)
}

func TestCreateAPIServer_ReturnsError_WhenPortNotProvided(t *testing.T) {
	httpConfig := &config.HTTP{
		Port: 0,
	}
	api, _ := NewAPI(&config.APIConfig{})

	result, err := createAPIServer(httpConfig, api)

	assert := assert.New(t)
	assert.Nil(result)
	assert.Error(err)
}
