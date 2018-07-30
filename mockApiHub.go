package main

import (
	"encoding/json"
	"net/http"
	"fmt"
	"os"
	"io/ioutil"

	configurator "MockApiHub/config"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/BurntSushi/toml"
)

func main() {
	var config configurator.Config

	// *** Get configuration ****************************************************
	if _, err := toml.DecodeFile("config.toml", &config); err != nil {
		fmt.Println(err)
		return
	}



	// *** Start server *********************************************************
	e := echo.New()

	e.Use(middleware.Logger())

	for _, api := range config.Apis {
		for _, endpoint := range api.Endpoints {
			file := endpoint.File
			path := endpoint.Path

			e.GET(path, func(c echo.Context) (err error) {
				return getJSON(c, file)
			})
		}
	}

	configurator.GetApis()

	// e.GET("/", getJSON)``
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", config.HTTP.Port)))
}

func getJSON(c echo.Context, filePath string) error {
	fmt.Println(filePath)

	jsonFile, err := os.Open(filePath)
	if err != nil {
		fmt.Println(err)
	}

	defer jsonFile.Close()
	
	bytes, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		fmt.Println(err)
	}

	if (!isValidJSON(bytes)) {
		return c.String(http.StatusInternalServerError, "bad json")
	}

	return c.JSONBlob(http.StatusOK, bytes)
}

func isValidJSON(bytes []byte) bool {
	var js json.RawMessage
	return json.Unmarshal(bytes, &js) == nil
}