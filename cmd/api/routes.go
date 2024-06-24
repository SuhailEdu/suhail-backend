package main

import (
	"github.com/labstack/echo/v4"
)

func RegisterApiRoutes(e *echo.Echo, config *Config) {

	authGroup := e.Group("auth")

	authGroup.POST("/register", config.registerUser)
	authGroup.POST("/login", config.loginUser)

}
