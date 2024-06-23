package main

import (
	"database/sql"
	"github.com/labstack/echo/v4"
	"log"
)

type Config struct {
	db     *sql.DB
	logger *log.Logger
}

func main() {

	// create the db connection
	db, err := sql.Open("postgres", "postgres://peter:password@localhost:5432/suhail")
	if err != nil {
		log.Fatal(err)
	}

	config := &Config{
		db: db,
	}

	e := echo.New()

	registerApiRoutes(e, config)

	//config(e)

	e.Logger.Fatal(e.Start(":1323"))
}
