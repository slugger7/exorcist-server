package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func CheckError(err error) {
	if err != nil {
		fmt.Println("Found an error")
		panic(err)
	}
}

func main() {
	err := godotenv.Load()
	CheckError(err)

	host := os.Getenv("DATABASE_HOST")
	port := os.Getenv("DATABASE_PORT")
	user := os.Getenv("DATABASE_USER")
	password := os.Getenv("DATABASE_PASSWORD")
	dbname := os.Getenv("DATABASE_NAME")

	fmt.Printf("host=%s port=%s user=%s password=%s database=%s", host, port, user, password, dbname)
	psqlconn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	fmt.Println("Opening DB")
	db, err := sql.Open("postgres", psqlconn)
	CheckError(err)
	defer db.Close()

	err = db.Ping()
	CheckError(err)

	fmt.Println("Querying DB")
	rows, err := db.Query(`SELECT "name" FROM library`)
	CheckError(err)

	defer rows.Close()

	fmt.Println("Reading rows")
	for rows.Next() {
		var name string
		err = rows.Scan(&name)
		CheckError(err)

		fmt.Println(name)
	}
	CheckError(err)
}
