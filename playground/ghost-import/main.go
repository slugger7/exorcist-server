package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	"github.com/slugger7/exorcist/internal/db/exorcist/public/model"
	"github.com/slugger7/exorcist/internal/db/exorcist/public/table"
	gmodel "github.com/slugger7/exorcist/internal/db/ghost/model"
	gtable "github.com/slugger7/exorcist/internal/db/ghost/table"
	errs "github.com/slugger7/exorcist/internal/errors"
)

func createPostgresDb() *sql.DB {
	host := os.Getenv("DATABASE_HOST")
	port := os.Getenv("DATABASE_PORT")
	user := os.Getenv("DATABASE_USER")
	password := os.Getenv("DATABASE_PASSWORD")
	dbname := os.Getenv("DATABASE_NAME")

	psqlconn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	fmt.Printf("Opening DB: %v\n", psqlconn)

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

type UserMap struct {
	ExorcistUser string `json:"exorcist_user"`
	GhostUser    string `json:"ghost_user"`
}

func parseUserMap(filePath string) ([]UserMap, error) {
	userMapFile, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}

	userMapBytes, err := io.ReadAll(userMapFile)
	if err != nil {
		return nil, err
	}

	var users []UserMap

	err = json.Unmarshal(userMapBytes, &users)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func transferLibraries(dbs *DBs) error {
	log.Println("Transferring libraries")

	ghostStatement := gtable.Libraries.SELECT(gtable.Libraries.AllColumns)

	var ghostLibraries []gmodel.Libraries

	if err := ghostStatement.Query(dbs.GhostDb, &ghostLibraries); err != nil {
		return err
	}

	log.Printf("Found %v libraries in ghost", len(ghostLibraries))

	exorcistLibraries := make([]model.Library, len(ghostLibraries))

	for i, l := range ghostLibraries {
		exorcistLibraries[i] = model.Library{
			Name:    l.Name,
			GhostID: &l.ID,
		}
	}

	insertStmnt := table.Library.INSERT(table.Library.GhostID, table.Library.Name).
		MODELS(exorcistLibraries).
		ON_CONFLICT(table.Library.GhostID).DO_NOTHING()

	res, err := insertStmnt.Exec(dbs.ExorcistDb)
	if err != nil {
		return err
	}

	rows, _ := res.RowsAffected()
	log.Printf("Altered %v rows in exorcist libraries", rows)

	return nil
}

type DBs struct {
	ExorcistDb *sql.DB
	GhostDb    *sql.DB
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

	dbs := DBs{
		ExorcistDb: pgDb,
		GhostDb:    sqlLiteDb,
	}

	err = transferLibraries(&dbs)
	if err != nil {
		errs.PanicError(err)
	}

	userMap, err := parseUserMap("./playground/ghost-import/user_map.json")
	if err != nil {
		errs.PanicError(err)
	}

	_ = userMap
}
