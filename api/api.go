package api

import (
	"fmt"
	"net/http"
	"errors"

	"MockApiHub/str"
	"MockApiHub/config"
	"MockApiHub/json"
)

// API contains information for an API
type API struct {
	baseURL string
	endpoints map[string]config.Endpoint
	server *http.Server
	port int
	handlers map[string]func(http.ResponseWriter, *http.Request)
}

const apiDir = "./api/apis"

// NewAPI returns a new API
func NewAPI (config *config.APIConfig) (*API, error) {
	api := &API{}
	server, err := createAPIServer(&config.HTTP, api)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	api.baseURL = config.BaseURL
	api.endpoints = config.Endpoints
	api.port = config.HTTP.Port
	api.server = server
	api.handlers = make(map[string]func(http.ResponseWriter, *http.Request))

	return api, nil
}

func createAPIServer(config *config.HTTP, api *API) (*http.Server, error) {
	if config.Port == 0 {
		return nil, errors.New("no port provided")
	}
	
	server := &http.Server {
		Addr: str.GetPort(config.Port),
		Handler: api,
	}

	return server, nil
}

// Register registers an api with the server
func (api *API) Register(dir string) error {
	fmt.Println("Registering ", dir, " on port ", str.GetPort(api.port))

	base := api.baseURL
	for _, endpoint := range api.endpoints {
		file := endpoint.File
		path := fmt.Sprintf("%s/%s", base, endpoint.Path)
		fmt.Println(path)
		api.handlers[path] = func(w http.ResponseWriter, r *http.Request) {
			json, err := json.GetJSON(fmt.Sprintf("%s/%s/%s", apiDir, dir, file))
			if err != nil {
				fmt.Println(err)
			}
			w.Write(json)
		}
	}
	go api.server.ListenAndServe()

	return nil
}

func (api *API) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if handler, ok := api.handlers[r.URL.String()[1:]]; ok {
		handler(w, r)
		return
	}
}