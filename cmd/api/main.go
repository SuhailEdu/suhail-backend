package main

import (
	"context"
	"github.com/SuhailEdu/suhail-backend/internal/database/schema"
	_ "github.com/SuhailEdu/suhail-backend/internal/database/schema"
	"github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5"
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
	//connection, err := sql.Open("postgres", dbUrl)
	conn, err := pgx.Connect(context.Background(), dbUrl)
	if err != nil {
		log.Fatal(err)
	}

	queries := schema.New(conn)

	config := &Config{
		db: queries,
	}

	e := echo.New()

	registerApiRoutes(e, config)

	e.Logger.Fatal(e.Start(":1323"))
}
