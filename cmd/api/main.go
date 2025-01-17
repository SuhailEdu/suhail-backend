package main

import (
	"context"
	"github.com/SuhailEdu/suhail-backend/internal/database/schema"
	_ "github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/lib/pq"
	"github.com/olahol/melody"
	"log"
	"os"
)

type Config struct {
	db     *schema.Queries
	logger *log.Logger
	melody *melody.Melody
}

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// create the db connection
	dbUrl := os.Getenv("DATABASE_URL")

	conn, err := pgxpool.New(context.Background(), dbUrl)
	if err != nil {
		log.Fatalf("Unable to connect to create pool: %v", err)
	}

	queries := schema.New(conn)

	m := melody.New()

	config := &Config{
		db:     queries,
		melody: m,
		logger: log.New(os.Stdout, "", log.LstdFlags),
	}
	e := echo.New()

	e.IPExtractor = echo.ExtractIPDirect()

	registerMelodyHandlers(e, config)

	e.Use(middleware.CORS())

	registerApiRoutes(e, config)

	appPort := os.Getenv("APP_PORT")

	e.Logger.Fatal(e.Start(":" + appPort))
}
