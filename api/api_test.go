package api

import (
	"errors"
	"net/http"
	"testing"

	"github.com/wcsanders1/MockApiHub/config"
	"github.com/wcsanders1/MockApiHub/fake"
	"github.com/wcsanders1/MockApiHub/log"
	"github.com/wcsanders1/MockApiHub/route"
	"github.com/wcsanders1/MockApiHub/str"
	"github.com/wcsanders1/MockApiHub/wrapper"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestStart_DeletesEndpoint_WhenEndpointIsDuplicate(t *testing.T) {
	endpointName := "testEndpoint"
	baseURL := "baseURL"
	path := "test/endpoint"
	file := "testFile"
	cert := "testCert"
	key := "testKey"
	dir := "testDir"
	creator := fakeAPICreator{}
	creator.On("getHandler", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(func(w http.ResponseWriter, r *http.Request) {})
	creator.On("startAPI", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
	routeTree := route.FakeTree{}
	routeTree.On("GetRoute", mock.AnythingOfType("string")).Return(path, map[string]string{}, nil)
	endpoints := map[string]config.Endpoint{
		endpointName: config.Endpoint{
			Path:   path,
			File:   file,
			Method: "GET",
		},
		"somethingElse": config.Endpoint{
			Path:   path,
			File:   file,
			Method: "GET",
		},
	}
	testAPI := API{
		server:     &wrapper.FakeServerOps{},
		httpConfig: config.HTTP{},
		endpoints:  endpoints,
		baseURL:    baseURL,
		routeTree:  &routeTree,
		log:        log.GetFakeLogger(),
		creator:    &creator,
		handlers:   make(map[string]map[string]func(http.ResponseWriter, *http.Request)),
	}

	err := testAPI.Start(dir, cert, key)

	assert := assert.New(t)
	assert.NoError(err)
	assert.Equal(1, len(testAPI.endpoints))
}

func TestStart_ReturnsNil_WhenStartSuccessful(t *testing.T) {
	endpointName := "testEndpoint"
	baseURL := "baseURL"
	path := "test/endpoint"
	file := "testFile"
	cert := "testCert"
	key := "testKey"
	dir := "testDir"
	creator := fakeAPICreator{}
	creator.On("getHandler", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(func(w http.ResponseWriter, r *http.Request) {})
	creator.On("startAPI", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
	routeTree := route.FakeTree{}
	routeTree.On("GetRoute", mock.AnythingOfType("string")).Return(path, map[string]string{}, nil)
	endpoints := map[string]config.Endpoint{
		endpointName: config.Endpoint{
			Path: path,
			File: file,
		},
	}
	testAPI := API{
		server:     &wrapper.FakeServerOps{},
		httpConfig: config.HTTP{},
		endpoints:  endpoints,
		baseURL:    baseURL,
		routeTree:  &routeTree,
		log:        log.GetFakeLogger(),
		creator:    &creator,
		handlers:   make(map[string]map[string]func(http.ResponseWriter, *http.Request)),
	}

	err := testAPI.Start(dir, cert, key)

	assert.NoError(t, err)
	creator.AssertCalled(t, "startAPI", cert, key, mock.Anything, mock.Anything)
}

func TestShutdown_ReturnsNil_WhenShutdownSuccessful(t *testing.T) {
	fakeServer := wrapper.FakeServerOps{}
	fakeServer.On("Shutdown", mock.Anything).Return(nil)
	api := API{
		log:    log.GetFakeLogger(),
		server: &fakeServer,
	}

	err := api.Shutdown()

	assert.NoError(t, err)
	fakeServer.AssertCalled(t, "Shutdown", mock.Anything)
}

func TestShutdown_ReturnsError_WhenShutdownFails(t *testing.T) {
	fakeServer := wrapper.FakeServerOps{}
	fakeServer.On("Shutdown", mock.Anything).Return(errors.New(""))
	api := API{
		log:    log.GetFakeLogger(),
		server: &fakeServer,
	}

	err := api.Shutdown()

	assert.Error(t, err)
	fakeServer.AssertCalled(t, "Shutdown", mock.Anything)
}

func TestServeHTTP_WritesStatusNotFound_WhenNoHandlerForRoute(t *testing.T) {
	path := "/test/path"
	method := "GET"
	routeTree := route.FakeTree{}
	routeTree.On("GetRoute", mock.AnythingOfType("string")).Return(path, map[string]string{}, nil)
	writer := fake.ResponseWriter{}
	writer.On("WriteHeader", mock.AnythingOfType("int")).Return(1)
	writer.On("Write", mock.AnythingOfType("[]uint8")).Return(1, nil)
	request, _ := http.NewRequest(method, path, nil)
	testAPI := API{
		log:       log.GetFakeLogger(),
		routeTree: &routeTree,
	}

	testAPI.ServeHTTP(&writer, request)

	routeTree.AssertCalled(t, "GetRoute", str.CleanURL(path))
	writer.AssertCalled(t, "WriteHeader", http.StatusNotFound)
}

func TestServeHTTP_DoesNotWriteError_WhenRouteTreeReturnsNonHTTPError(t *testing.T) {
	path := "/test/path"
	routeTree := route.FakeTree{}
	routeTree.On("GetRoute", mock.AnythingOfType("string")).Return(path, map[string]string{}, errors.New(""))
	writer := fake.ResponseWriter{}
	writer.On("WriteHeader", mock.AnythingOfType("int")).Return(1)
	writer.On("Write", mock.AnythingOfType("[]uint8")).Return(1, nil)
	request, _ := http.NewRequest("GET", path, nil)
	testAPI := API{
		log:       log.GetFakeLogger(),
		routeTree: &routeTree,
	}

	testAPI.ServeHTTP(&writer, request)

	routeTree.AssertCalled(t, "GetRoute", str.CleanURL(path))
	writer.AssertNotCalled(t, "WriteHeader", mock.Anything)
}

func TestServeHTTP_WritesError_WhenRouteTreeReturnsHTTPError(t *testing.T) {
	path := "/test/path"
	status := http.StatusBadRequest
	httpError := &route.HTTPError{
		Status: status,
		Msg:    "test",
	}
	routeTree := route.FakeTree{}
	routeTree.On("GetRoute", mock.AnythingOfType("string")).Return(path, map[string]string{}, httpError)
	writer := fake.ResponseWriter{}
	writer.On("WriteHeader", mock.AnythingOfType("int")).Return(1)
	writer.On("Write", mock.AnythingOfType("[]uint8")).Return(1, nil)
	request, _ := http.NewRequest("GET", path, nil)
	testAPI := API{
		log:       log.GetFakeLogger(),
		routeTree: &routeTree,
	}

	testAPI.ServeHTTP(&writer, request)

	routeTree.AssertCalled(t, "GetRoute", str.CleanURL(path))
	writer.AssertCalled(t, "WriteHeader", status)
}

func TestServeHTTP_DoesNotWriteErrorStatus_WhenHandlerExists(t *testing.T) {
	path := "/test/path"
	method := "GET"
	routeTree := route.FakeTree{}
	routeTree.On("GetRoute", mock.AnythingOfType("string")).Return(path, map[string]string{}, nil)
	writer := fake.ResponseWriter{}
	request, _ := http.NewRequest(method, path, nil)
	handlers := make(map[string]map[string]func(http.ResponseWriter, *http.Request))
	handlers[method] = make(map[string]func(http.ResponseWriter, *http.Request))
	handlers[method][path] = func(http.ResponseWriter, *http.Request) {}
	testAPI := API{
		handlers:  handlers,
		log:       log.GetFakeLogger(),
		routeTree: &routeTree,
	}

	testAPI.ServeHTTP(&writer, request)

	routeTree.AssertCalled(t, "GetRoute", str.CleanURL(path))
	writer.AssertNotCalled(t, "WriteHeader", http.StatusNotFound)
}

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
