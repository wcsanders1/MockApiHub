package api

import (
	"fmt"
	"net/http"

	"MockApiHub/utils"
)

// API contains information for an API
type API struct {
	BaseURL string
	Endpoints map[string]endpoint
	Server http.Server
	Port int
	Handlers map[string]func(http.ResponseWriter, *http.Request)
}

type endpoint struct {
	Path string 
	File string
}

type handler struct{}
var mux = make(map[string]func(http.ResponseWriter, *http.Request))
const apiDir = "./api/apis"



// Register registers an api with the server
func (api *API) Register(dir string) error {
	fmt.Println("Registering ", dir, " on port ", utils.GetPort(api.Port))
	api.Server = http.Server {
		Addr: utils.GetPort(api.Port),
		Handler: &handler{},
	}

	base := api.BaseURL
	for _, endpoint := range api.Endpoints {
		file := endpoint.File
		path := fmt.Sprintf("%s/%s", base, endpoint.Path)
		mux[path] = func(w http.ResponseWriter, r *http.Request) {
			json, err := utils.GetJSON(fmt.Sprintf("%s/%s/%s", apiDir, dir, file))
			if err != nil {
				fmt.Println(err)
			}
			w.Write(json)
		}
	}
	go api.Server.ListenAndServe()

	return nil
}

func (*handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if h, ok := mux[r.URL.String()[1:]]; ok {
		h(w, r)
		return
	}
}