package main

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

func (config *Config) RegisterUser(c echo.Context) error {

	return c.JSON(http.StatusNotImplemented, "not implemented")

	// how to use the db here?
}
