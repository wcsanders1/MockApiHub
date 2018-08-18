package api

import (
	"fmt"
	"net/http"

	// "github.com/labstack/echo"
)

// API contains information for an API
type API struct {
	BaseURL string
	Endpoints map[string]endpoint
	Server http.Server
}

type endpoint struct {
	Path string 
	File string
}

type handler struct{}
var mux map[string]func(http.ResponseWriter, *http.Request)
const apiDir = "./api/apis"

// Register registers an api with the server
func (api *API) Register(dir string) error {
	fmt.Println("Registering ", dir)

	mux = make(map[string]func(http.ResponseWriter, *http.Request))
	api.Server = http.Server {
		Addr: ":5000",
		Handler: &handler{},
	}

	base := api.BaseURL
	for _, endpoint := range api.Endpoints {
		file := endpoint.File
		path := fmt.Sprintf("%s/%s", base, endpoint.Path)
		// e.GET(path, func(c echo.Context) (err error) {
		// 	return getJSON(c, fmt.Sprintf("%s/%s/%s", apiDir, dir, file))
		// })
		fmt.Println(path)
		mux[path] = func(w http.ResponseWriter, r *http.Request) {
			json, err := getJSON(fmt.Sprintf("%s/%s/%s", apiDir, dir, file))
			if err != nil {
				fmt.Println(err)
			}
			w.Write(json)
		}
	}
	api.Server.ListenAndServe()

	return nil
}

func (*handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.URL.String()[1:])
	if h, ok := mux[r.URL.String()[1:]]; ok {
		fmt.Println("here it is")
		h(w, r)
		return
	}
}