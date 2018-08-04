package main

import (
	"fmt"

	configurator "MockApiHub/config"
	"MockApiHub/apis"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/BurntSushi/toml"
)

func main() {
	// *** Get configuration ****************************************************
	var config configurator.Config
	if _, err := toml.DecodeFile("config.toml", &config); err != nil {
		fmt.Println(err)
		return
	}

	// *** Instantiate server ***************************************************
	e := echo.New()
	e.Use(middleware.Logger())

	// *** Register APIs ********************************************************
	apis, err :=apis.GetAPIs()
	if err != nil {
		fmt.Println(err)
		return
	}

	for dir, api := range apis {
		api.Register(dir, e)
	}

	// *** Start server *********************************************************
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", config.HTTP.Port)))
}