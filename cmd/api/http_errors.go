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

	errors := map[string]interface{}{"validationError": err}
	return c.JSON(http.StatusUnprocessableEntity, errors)
}
func unAuthorizedError(c echo.Context, err interface{}) error {

	errors := map[string]interface{}{"error": err}
	return c.JSON(http.StatusUnauthorized, errors)
}
func forbiddenError(c echo.Context, err interface{}) error {

	errors := map[string]interface{}{"error": err}
	return c.JSON(http.StatusForbidden, errors)
}

func dataResponse(c echo.Context, data interface{}) error {

	//if len(data) == 0 {
	//
	//}

	return c.JSON(http.StatusOK, map[string]interface{}{"data": data})
}
