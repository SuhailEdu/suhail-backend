package main

import (
	"github.com/labstack/echo/v4"
	"github.com/thedevsaddam/govalidator"
	"net/http"
)

func (config *Config) registerUser(c echo.Context) error {

	rules := govalidator.MapData{
		"email":            []string{"required", "email"},
		"first_name":       []string{"required", "between:3,8"},
		"last_name":        []string{"required", "between:3,8"},
		"password":         []string{"required", "between:8,20"},
		"password_confirm": []string{"required", "between:8,20"},
	}

	validationOptions := govalidator.Options{
		Rules:           rules,
		RequiredDefault: true,
	}

	v := govalidator.New(validationOptions)

	e := v.Validate()

	err := map[string]interface{}{"validationError": e}

	return c.JSON(http.StatusNotImplemented, err)

	// how to use the db here?
}
