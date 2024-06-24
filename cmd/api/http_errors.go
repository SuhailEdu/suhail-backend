package main

import (
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
)

func serverError(c echo.Context, err error) error {

	log.Println(err)

	return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal Server Error"})

}

func badRequestError(c echo.Context, err error) error {

	log.Println(err)

	return c.JSON(http.StatusBadRequest, map[string]string{"error": "Your request is invalid"})
}

func validationError(c echo.Context, err interface{}) error {

	log.Println(err)

	return c.JSON(http.StatusUnprocessableEntity, err)
}
