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
	"github.com/wcsanders1/MockApiHub/json"
	"github.com/wcsanders1/MockApiHub/log"
	"github.com/wcsanders1/MockApiHub/ref"
	"github.com/wcsanders1/MockApiHub/route"
	"github.com/wcsanders1/MockApiHub/str"

	"github.com/sirupsen/logrus"
)

// API contains information for an API
type API struct {
	baseURL    string
	endpoints  map[string]config.Endpoint
	server     *http.Server
	handlers   map[string]map[string]func(http.ResponseWriter, *http.Request)
	routeTree  *route.Tree
	httpConfig config.HTTP
	log        *logrus.Entry
}

const apiDir = "./api/apis"

// NewAPI returns a new API
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

	api.server = server
	api.baseURL = config.BaseURL
	api.endpoints = config.Endpoints
	api.handlers = make(map[string]map[string]func(http.ResponseWriter, *http.Request))
	api.routeTree = route.NewRouteTree()
	api.httpConfig = config.HTTP

	contextLogger.Info("successfully created mock API")
	return api, nil
}

// Register registers an api with the server
func (api *API) Register(dir, defaultCert, defaultKey string) error {
	contextLogger := api.log.WithFields(logrus.Fields{
		log.FuncField:            ref.GetFuncName(),
		log.DefaultCertFileField: defaultCert,
		log.DefaultKeyFileField:  defaultKey,
		log.APIDirField:          dir,
	})
	contextLogger.Debug("registering API")

	base := api.baseURL
	for endpointName, endpoint := range api.endpoints {
		path := fmt.Sprintf("%s/%s", base, endpoint.Path)
		registeredRoute := api.ensureRouteRegistered(path)
		file := endpoint.File
		method := strings.ToUpper(endpoint.Method)
		contextLoggerEndpoint := api.log.WithFields(logrus.Fields{
			log.PathField:         path,
			log.RouteField:        registeredRoute,
			log.FileField:         file,
			log.MethodField:       method,
			log.EndpointNameField: endpointName,
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

		api.handlers[method][registeredRoute] = func(w http.ResponseWriter, r *http.Request) {
			json, err := json.GetJSON(fmt.Sprintf("%s/%s/%s", apiDir, dir, file))
			if err != nil {
				contextLoggerEndpoint.WithError(err).Error("error serving from this endpoint")
			}
			contextLoggerEndpoint.Debug("successfully retrieved JSON; serving it")
			w.Write(json)
		}
	}

	var err error
	if api.httpConfig.UseTLS {
		cert, key, err := api.getCertAndKeyFile(defaultCert, defaultKey)
		if err != nil {
			contextLogger.WithError(err).Error("error getting TLS cert and key")
			return err
		}

		go func() error {
			if err = api.server.ListenAndServeTLS(cert, key); err != nil {
				contextLogger.WithError(err).Error("error starting mock API with TLS")
				return err
			}
			return nil
		}()
		return nil
	}

	go func() error {
		if err = api.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			contextLogger.WithError(err).Error("mock API server error")
			return err
		}
		return nil
	}()

	return err
}

// Shutdown shutsdown the server
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
	path, params, err := api.routeTree.GetRoute(str.CleanURL(r.URL.String()))
	contextLogger := api.log.WithFields(logrus.Fields{
		log.FuncField: ref.GetFuncName(),
		log.PathField: path,
		"params":      params,
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

// GetPort returns the API's port number
func (api *API) GetPort() int {
	return api.httpConfig.Port
}

// GetBaseURL returns the API's base URL
func (api *API) GetBaseURL() string {
	return api.baseURL
}

// GetEndpoints returns the API's endpoints
func (api *API) GetEndpoints() map[string]config.Endpoint {
	return api.endpoints
}

func (api *API) getCertAndKeyFile(defaultCert, defaultKey string) (string, string, error) {
	if len(api.httpConfig.CertFile) > 0 && len(api.httpConfig.KeyFile) > 0 {
		return api.httpConfig.CertFile, api.httpConfig.KeyFile, nil
	}

	if len(api.httpConfig.CertFile) == 0 && len(api.httpConfig.KeyFile) > 0 {
		return "", "", errors.New("key provided without cert")
	}

	if len(api.httpConfig.KeyFile) == 0 && len(api.httpConfig.CertFile) > 0 {
		return "", "", errors.New("cert provided without key")
	}

	return defaultCert, defaultKey, nil
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
