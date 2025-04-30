package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/go-jet/jet/v2/postgres"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	. "github.com/slugger7/exorcist/internal/db/exorcist/public/table"
	. "github.com/slugger7/exorcist/internal/errors"
)

func main() {
	err := godotenv.Load()
	PanicError(err)

	host := os.Getenv("DATABASE_HOST")
	port := os.Getenv("DATABASE_PORT")
	user := os.Getenv("DATABASE_USER")
	password := os.Getenv("DATABASE_PASSWORD")
	dbname := os.Getenv("DATABASE_NAME")

	fmt.Printf("host=%s port=%s user=%s password=%s database=%s", host, port, user, password, dbname)
	psqlconn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	fmt.Println("Opening DB")
	db, err := sql.Open("postgres", psqlconn)
	PanicError(err)
	defer db.Close()

	err = db.Ping()
	PanicError(err)

	stmnt := Video.SELECT(postgres.COUNT(Video.ID).AS("total")).FROM(Video)

	sql, _ := stmnt.Sql()
	fmt.Printf("SQL: %v\n", sql)
	var dest struct {
		Total int
	}
_:
	stmnt.Query(db, &dest)

	fmt.Println(dest)
	// rows, _ := db.Query(sql)

	// defer rows.Close()
	// for rows.Next() {
	// 	var total int

	// 	_ = rows.Scan(&total)

	// 	fmt.Println(total)
	// }
}
