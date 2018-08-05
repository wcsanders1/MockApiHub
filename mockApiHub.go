package main

import (
	"fmt"
	"os"

	"MockApiHub/config"
	"MockApiHub/api"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/BurntSushi/toml"
)

func main() {
	// *** Get configuration ****************************************************
	var config config.Config
	if _, err := toml.DecodeFile("config.toml", &config); err != nil {
		fmt.Println(err)
		return
	}

	// *** Instantiate server ***************************************************
	e := echo.New()
	e.Use(middleware.Logger())

	// *** Register APIs ********************************************************
	apis, err := api.GetAPIs()
	if err != nil {
		fmt.Println(err)
		return
	}

	for dir, api := range apis {
		api.Register(dir, e)
	}

	// *** Start server *********************************************************
	addr := fmt.Sprintf(":%d", config.HTTP.Port)
	if config.HTTP.UseTLS {
		startUsingTLS(e, &config.HTTP, addr)
	} else {
		e.Logger.Fatal(e.Start(addr))
	}
}

func startUsingTLS(e *echo.Echo, http *config.HTTP, addr string) {
	if _, err := os.Stat(http.CertFile); os.IsNotExist(err) {
		e.Logger.Fatal(fmt.Sprintf("%s cert file does not exist", http.CertFile))
	}

	if _, err := os.Stat(http.KeyFile); os.IsNotExist(err) {
		e.Logger.Fatal(fmt.Sprintf("%s key file does not exist", http.KeyFile))
	}

	e.Logger.Fatal(e.StartTLS(addr, http.CertFile, http.KeyFile))
}