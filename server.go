package main

import (
	"github.com/SuhailEdu/suhail-backend/routes"
	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()

	routes.ApiRoutes(e)

	e.Logger.Fatal(e.Start(":1323"))
}
