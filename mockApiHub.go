package main

import (
	"net/http"
	"fmt"

	"MockApiHub/config"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/BurntSushi/toml"
)

func main() {
	var config config.Config

	// Configure
	if _, err := toml.DecodeFile("config.toml", &config); err != nil {
		fmt.Println(err)
		return;
	}

	e := echo.New()

	e.Use(middleware.Logger())

	e.GET("/", hello)
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", config.HTTP.Port)))

	
}

func hello(c echo.Context) error {
	return c.String(http.StatusOK, "Hello")
}
