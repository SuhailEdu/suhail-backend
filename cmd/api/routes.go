package main

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

func RegisterApiRoutes(e *echo.Echo, config *Config) {

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	authGroup := e.Group("auth")

	authGroup.POST("/register", config.registerUser)

}
