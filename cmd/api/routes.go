package main

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

func RegisterApiRoutes(e *echo.Echo, config *Config) {

	authGroup := e.Group("auth")

	authGroup.POST("/register", config.registerUser)
	authGroup.POST("/login", config.loginUser)

	homeGroup := e.Group("/home")
	homeGroup.Use(config.checkAuthToken)

	homeGroup.GET("/hello", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
}
