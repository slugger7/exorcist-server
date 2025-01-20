package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/slugger7/exorcist/internal/db/exorcist/public/model"
	. "github.com/slugger7/exorcist/internal/db/exorcist/public/table"
	. "github.com/slugger7/exorcist/internal/errors"
	"github.com/slugger7/exorcist/internal/media"
)

func main() {
	path := ""
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

	libraryPathId := createLibWithPath(db, path)

	values, err := media.GetFilesByExtensions(path, []string{".mp4", ".m4v", ".mkv", ".avi", ".wmv", ".flv", ".webm", ".f4v", ".mpg", ".m2ts", ".mov"})
	CheckError(err)

	fmt.Println("Printing out results")
	videoModels := []model.Video{}
	for _, v := range values {
		fmt.Println(v)
		checksum := "lol"
		videoModels = append(videoModels, model.Video{
			LibraryPathID: libraryPathId,
			RelativePath:  v,
			Title:         "",
			FileName:      "",
			Height:        666,
			Width:         666,
			Runtime:       666,
			Size:          666,
			Checksum:      &checksum,
		})
	}

	insertStatement := Video.INSERT(
		Video.LibraryPathID,
		Video.RelativePath,
		Video.Title,
		Video.FileName,
		Video.Height,
		Video.Width,
		Video.Runtime,
		Video.Size,
		Video.Checksum,
	).
		MODELS(videoModels).
		RETURNING(Video.ID)

	var newVideos []struct {
		model.Video
	}
	err = insertStatement.Query(db, &newVideos)
}

func createLibWithPath(db *sql.DB, path string) uuid.UUID {
	newLib := model.Library{
		Name: "New Lib",
	}

	insertStatement := Library.INSERT(Library.Name).MODEL(newLib).RETURNING(Library.ID)

	var library []struct {
		model.Library
	}

	err := insertStatement.Query(db, &library)
	CheckError(err)

	newLibPath := model.LibraryPath{
		LibraryID: library[0].ID,
		Path:      path,
	}

	insertStatement = LibraryPath.INSERT(LibraryPath.LibraryID, LibraryPath.Path).MODEL(newLibPath).RETURNING(LibraryPath.ID)

	var libraryPath []struct {
		model.LibraryPath
	}

	err = insertStatement.Query(db, &libraryPath)

	return libraryPath[0].ID
}
