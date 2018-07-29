package main

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func main() {
	e := echo.New()

	e.Use(middleware.Logger())

	e.GET("/", hello)
	e.Logger.Fatal(e.Start(":5000"))
}

func hello(c echo.Context) error {
	return c.String(http.StatusOK, "Hello")
}
