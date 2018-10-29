package api

import (
	"net/http"
	"testing"

	"github.com/wcsanders1/MockApiHub/config"
	"github.com/wcsanders1/MockApiHub/route"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

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
