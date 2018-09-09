package api

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"MockApiHub/config"
	"MockApiHub/json"
	"MockApiHub/route"
	"MockApiHub/str"
)

// API contains information for an API
type API struct {
	baseURL    string
	endpoints  map[string]config.Endpoint
	server     *http.Server
	handlers   map[string]map[string]func(http.ResponseWriter, *http.Request)
	routeTree  *route.Tree
	httpConfig config.HTTP
}

const apiDir = "./api/apis"

// NewAPI returns a new API
func NewAPI(config *config.APIConfig) (*API, error) {
	api := &API{}
	server, err := createAPIServer(&config.HTTP, api)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	api.baseURL = config.BaseURL
	api.endpoints = config.Endpoints
	api.server = server
	api.handlers = make(map[string]map[string]func(http.ResponseWriter, *http.Request))
	api.routeTree = route.NewRouteTree()
	api.httpConfig = config.HTTP

	return api, nil
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

// Register registers an api with the server
func (api *API) Register(dir, defaultCert, defaultKey string) error {
	fmt.Println("Registering ", dir, " on port ", str.GetPort(api.httpConfig.Port))

	base := api.baseURL
	for _, endpoint := range api.endpoints {
		path := fmt.Sprintf("%s/%s", base, endpoint.Path)
		registeredRoute := api.ensureRouteRegistered(path)
		file := endpoint.File
		method := strings.ToUpper(endpoint.Method)

		if route, methodExists := api.handlers[method]; !methodExists {
			api.handlers[method] = make(map[string]func(http.ResponseWriter, *http.Request))
		} else {
			if _, routeExists := route[registeredRoute]; routeExists {
				fmt.Println(fmt.Sprintf("WARNING: the following route with method %s already exists: %s",
					method, registeredRoute))

				continue
			}
		}

		api.handlers[method][registeredRoute] = func(w http.ResponseWriter, r *http.Request) {
			json, err := json.GetJSON(fmt.Sprintf("%s/%s/%s", apiDir, dir, file))
			if err != nil {
				fmt.Println(err)
			}
			w.Write(json)
		}
	}

	if api.httpConfig.UseTLS {
		cert, key, err := api.getCertAndKeyFile(defaultCert, defaultKey)
		if err != nil {
			fmt.Println(err)
			return err
		}

		go func() {
			api.server.ListenAndServeTLS(cert, key)
			defer api.Shutdown()
		}()
	}

	go func() {
		api.server.ListenAndServe()
		defer api.Shutdown()
	}()

	return nil
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
	registeredRoute, _ := api.routeTree.GetRoute(url)
	if len(registeredRoute) == 0 {
		registeredRoute, _ = api.routeTree.AddRoute(url)
	}

	return registeredRoute
}

// Shutdown shutsdown the server
func (api *API) Shutdown() error {
	if err := api.server.Shutdown(context.Background()); err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println(fmt.Sprintf("shut down server on port %d", api.httpConfig.Port))
	return nil
}

func (api *API) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path, err := api.routeTree.GetRoute(str.CleanURL(r.URL.String()))
	if err != nil {
		switch err.(type) {
		case *route.HTTPError:
			httpError, _ := err.(*route.HTTPError)
			w.WriteHeader(httpError.Status)
			w.Write([]byte(httpError.Msg))
			return
		default:
			fmt.Println(err)
			return
		}
	}

	method := strings.ToUpper(r.Method)
	if handler, exists := api.handlers[method][path]; exists {
		handler(w, r)
		return
	}

	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("endpoint not found"))
}
