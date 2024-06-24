package main

import (
	"database/sql"
	"github.com/SuhailEdu/suhail-backend/internal/database/schema"
	_ "github.com/SuhailEdu/suhail-backend/internal/database/schema"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
	"log"
	"os"
)

type Config struct {
	db     *schema.Queries
	logger *log.Logger
}

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// create the db connection
	dbUrl := os.Getenv("DATABASE_URL")
	connection, err := sql.Open("postgres", dbUrl)
	if err != nil {
		log.Fatal(err)
	}

	queries := schema.New(connection)

	config := &Config{
		db: queries,
	}

	e := echo.New()

	RegisterApiRoutes(e, config)

	e.Logger.Fatal(e.Start(":1323"))
}
