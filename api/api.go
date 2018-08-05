package api

import (
	"fmt"

	"github.com/labstack/echo"
)

// API contains information for an API
type API struct {
	BaseURL string
	Endpoints map[string]endpoint
}

type endpoint struct {
	Path string 
	File string
}

const apiDir = "./api/apis"

// Register registers an api with the server
func (api *API) Register(dir string, e *echo.Echo) error {
	fmt.Println("Registering ", dir)
	base := api.BaseURL
	for _, endpoint := range api.Endpoints {
		file := endpoint.File
		path := fmt.Sprintf("%s/%s", base, endpoint.Path)
		e.GET(path, func(c echo.Context) (err error) {
			return getJSON(c, fmt.Sprintf("%s/%s/%s", apiDir, dir, file))
		})
	}

	return nil
}