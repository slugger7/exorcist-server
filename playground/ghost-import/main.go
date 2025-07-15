package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	errs "github.com/slugger7/exorcist/internal/errors"
)

func createPostgresDb() *sql.DB {
	host := os.Getenv("DATABASE_HOST")
	port := os.Getenv("DATABASE_PORT")
	user := os.Getenv("DATABASE_USER")
	password := os.Getenv("DATABASE_PASSWORD")
	dbname := os.Getenv("DATABASE_NAME")

	fmt.Printf("host=%s port=%s user=%s password=%s database=%s", host, port, user, password, dbname)
	psqlconn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	fmt.Println("Opening DB")
	db, err := sql.Open("postgres", psqlconn)
	errs.PanicError(err)

	return db
}

func createSqlLiteDb() *sql.DB {
	dbPath := os.Getenv("GHOST_DB_PATH")

	db, err := sql.Open("sqlite3", dbPath)
	errs.PanicError(err)

	return db
}

func main() {
	err := godotenv.Load()
	errs.PanicError(err)

	pgDb := createPostgresDb()
	defer pgDb.Close()

	err = pgDb.Ping()
	errs.PanicError(err)

	sqlLiteDb := createSqlLiteDb()
	defer sqlLiteDb.Close()

	err = sqlLiteDb.Ping()
	errs.PanicError(err)

}
