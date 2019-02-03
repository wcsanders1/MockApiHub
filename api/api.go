/*
Package api creates a mock API based on configuration placed in the mockApis directory.

Example configuration:

	baseUrl = "studentsApi/:districtNumber"

	[log]
	loggingEnabled = true
	fileName = "testLogs/studentsApi/default.log"
	maxFileDaysAge = 3
	formatAsJSON = true
	prettyJSON = true

	[http]
	port = 5002
	useTLS = false
	certFile = ""
	keyFile = ""

	[endpoints]

		[endpoints.getAllAccounts]
		path = "accounts"
		file = "accounts.json"
		method = "GET"

		[endpoints.getCustomerBalances]
		path = "customers/:id/balances"
		file = "customers.json"
		method = "GET"

*/
package api

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"path"
	"strings"
	"time"

	"github.com/wcsanders1/MockApiHub/config"
	"github.com/wcsanders1/MockApiHub/log"
	"github.com/wcsanders1/MockApiHub/ref"
	"github.com/wcsanders1/MockApiHub/route"
	"github.com/wcsanders1/MockApiHub/str"
	"github.com/wcsanders1/MockApiHub/wrapper"

	"github.com/sirupsen/logrus"
)

type (
	// IAPI is an interface providing functionality to manage an API.
	IAPI interface {
		Start(dir, defaultCert, defaultKey string) error
		Shutdown() error
		ServeHTTP(w http.ResponseWriter, r *http.Request)
		GetPort() int
		GetBaseURL() string
		GetEndpoints() map[string]config.Endpoint
	}

	// API contains information for an API.
	API struct {
		baseURL    string
		endpoints  map[string]config.Endpoint
		server     wrapper.IServerOps
		handlers   map[string]map[string]func(http.ResponseWriter, *http.Request)
		routeTree  route.ITree
		httpConfig config.HTTP
		log        *logrus.Entry
		file       wrapper.IFileOps
		creator    iCreator
	}
)

// NewAPI returns a new API.
func NewAPI(config *config.APIConfig) (*API, error) {
	api := &API{}
	api.log = log.NewLogger(&config.Log, "api").WithFields(logrus.Fields{
		log.BaseURLField:  config.BaseURL,
		log.PortField:     config.HTTP.Port,
		log.UseTLSField:   config.HTTP.UseTLS,
		log.CertFileField: config.HTTP.CertFile,
		log.KeyFileField:  config.HTTP.KeyFile,
	})
	contextLogger := api.log.WithField(log.FuncField, ref.GetFuncName())

	server, err := createAPIServer(&config.HTTP, api)
	if err != nil {
		contextLogger.WithError(err).Error("error creating mock API")
		return nil, err
	}

	api.server = wrapper.NewServerOps(server)
	api.baseURL = config.BaseURL
	api.endpoints = config.Endpoints
	api.handlers = make(map[string]map[string]func(http.ResponseWriter, *http.Request))
	api.routeTree = route.NewRouteTree()
	api.httpConfig = config.HTTP
	api.file = &wrapper.FileOps{}
	api.creator = newCreator(api.log)

	contextLogger.Info("successfully created mock API")
	return api, nil
}

// Start starts an api server.
func (api *API) Start(dir, defaultCert, defaultKey string) error {
	contextLogger := api.log.WithFields(logrus.Fields{
		log.FuncField:            ref.GetFuncName(),
		log.DefaultCertFileField: defaultCert,
		log.DefaultKeyFileField:  defaultKey,
		log.APIDirField:          dir,
	})
	contextLogger.Debug("starting API")

	for endpointName, endpoint := range api.endpoints {
		var path string
		if len(api.baseURL) > 0 {
			path = fmt.Sprintf("%s/%s", api.baseURL, endpoint.Path)
		} else {
			path = endpoint.Path
		}
		registeredRoute := api.ensureRouteRegistered(path)
		file := endpoint.File
		method := strings.ToUpper(endpoint.Method)
		contextLoggerEndpoint := api.log.WithFields(logrus.Fields{
			log.PathField:            path,
			log.RouteField:           registeredRoute,
			log.FileField:            file,
			log.MethodField:          method,
			log.EndpointNameField:    endpointName,
			log.ResponseHeadersField: endpoint.Headers,
		})

		contextLoggerEndpoint.Debug("registering endpoint")
		if route, methodExists := api.handlers[method]; !methodExists {
			api.handlers[method] = make(map[string]func(http.ResponseWriter, *http.Request))
		} else {
			if _, routeExists := route[registeredRoute]; routeExists {
				contextLoggerEndpoint.Warn("endpoint already exists, moving on to next endpoint...")
				delete(api.endpoints, endpointName)
				continue
			}
		}

		contextLoggerEndpoint.Debug("registered endpoint; now assigning handler")
		api.handlers[method][registeredRoute] = api.creator.getHandler(endpoint.EnforceValidJSON, dir, file, api.file)
	}

	return api.creator.startAPI(defaultCert, defaultKey, api.server, api.httpConfig)
}

// Shutdown shutsdown the server.
func (api *API) Shutdown() error {
	contextLogger := api.log.WithField(log.FuncField, ref.GetFuncName())
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := api.server.Shutdown(ctx); err != nil {
		contextLogger.WithError(err).Error("error shutting down mock API")
		return err
	}

	contextLogger.Info("successfully shut down mock API")
	return nil
}

func (api *API) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path, params, err := api.routeTree.GetRoute(str.CleanURL(r.URL.Path))
	contextLogger := api.log.WithFields(logrus.Fields{
		log.FuncField: ref.GetFuncName(),
		log.PathField: path,
		"params":      params,
		"query":       r.URL.Query(),
	})

	if err != nil {
		switch err.(type) {
		case *route.HTTPError:
			httpError, _ := err.(*route.HTTPError)
			contextLogger.WithError(httpError).Error("server error")
			w.WriteHeader(httpError.Status)
			w.Write([]byte(httpError.Msg))
			return
		default:
			contextLogger.WithError(err).Error("server error")
			return
		}
	}

	method := strings.ToUpper(r.Method)
	if handler, exists := api.handlers[method][path]; exists {
		contextLogger.Debug("handler exists for this path")
		handler(w, r)
		return
	}

	contextLogger.Warn("endpoint not found")
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("endpoint not found"))
}

// GetPort returns the API's port number.
func (api *API) GetPort() int {
	return api.httpConfig.Port
}

// GetBaseURL returns the API's base URL.
func (api *API) GetBaseURL() string {
	return api.baseURL
}

// GetEndpoints returns the API's endpoints.
func (api *API) GetEndpoints() map[string]config.Endpoint {
	return api.endpoints
}

func (api *API) ensureRouteRegistered(url string) string {
	url = path.Clean(url)
	registeredRoute, _, _ := api.routeTree.GetRoute(url)
	if len(registeredRoute) == 0 {
		registeredRoute, _ = api.routeTree.AddRoute(url)
	}

	return registeredRoute
}

func createAPIServer(config *config.HTTP, api *API) (*http.Server, error) {
	if config.Port == 0 {
		return nil, errors.New("no port provided")
	}

	server := &http.Server{
		Addr:    str.GetPort(config.Port),
		Handler: api,
	}

	return server, nil
}
