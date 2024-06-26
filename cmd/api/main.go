package main

import (
	"database/sql"
	_ "github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
	"github.com/volatiletech/sqlboiler/v4/boil"
	_ "github.com/volatiletech/sqlboiler/v4/queries/qm"
	"log"
	"os"
)

//go:generate sqlboiler --wipe psql

type Config struct {
	db     *sql.DB
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

	boil.SetDB(connection)

	//queries := schema.New(conn)
	//
	config := &Config{
		db: connection,
	}

	e := echo.New()

	registerApiRoutes(e, config)

	e.Logger.Fatal(e.Start(":1323"))
}
