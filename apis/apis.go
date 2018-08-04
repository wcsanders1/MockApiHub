package apis

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"os"

	"github.com/labstack/echo"
	"github.com/BurntSushi/toml"
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

// Register registers an api with the ECHO server
func (api *API) Register(dir string, e *echo.Echo) error {
	fmt.Println("Registering ", dir)
	for _, endpoint := range api.Endpoints {
		file := endpoint.File
		path := endpoint.Path
		e.GET(path, func(c echo.Context) (err error) {
			return getJSON(c, fmt.Sprintf("apis/%s/%s", dir, file))
		})
	}

	return nil
}

// GetAPIs gets all the apis in the apis directory
func GetAPIs() (map[string]*API, error) {
	apiDir, err := ioutil.ReadDir("./apis")
	if err != nil {
		return nil, err
	}

	apis := make(map[string]*API)
	for _, file := range apiDir {
		api, err := extractAPIFromFile(file)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}

		if api != nil {
			apis[file.Name()] = api
		}
	}

	return apis, nil
}

func getAPI(dir string, files []os.FileInfo) (*API, error) {
	for _, file := range files {
		if (isAPIConfig(file)) {
			return decodeAPIConfig(dir, file)
		}
	}
	return nil, nil
}

func extractAPIFromFile(file os.FileInfo) (*API, error) {
	if (!file.IsDir() || !isAPI(file.Name())) {
		return nil, nil
	}
	
	dir := file.Name()
	fmt.Println("Found the following mock api: ", dir)
	files, _ := ioutil.ReadDir(fmt.Sprintf("./apis/%s", dir))
	api, err := getAPI(dir, files)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return api, nil	
}

func decodeAPIConfig(dir string, file os.FileInfo) (*API, error) {
	path := fmt.Sprintf("./apis/%s/%s", dir, file.Name())
	var api API
	if _, err := toml.DecodeFile(path, &api); err != nil {
		fmt.Println(err)
		return nil, err
	}
	return &api, nil
}

func isAPI(dir string) bool {
	return len(dir) > 3 && dir[len(dir)-3:] == "Api"
}

func isAPIConfig(file os.FileInfo) bool {
	ext := filepath.Ext(file.Name())
	return ext == ".toml"
}