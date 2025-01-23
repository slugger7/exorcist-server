package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/slugger7/exorcist/internal/db/exorcist/public/model"
	"github.com/slugger7/exorcist/internal/db/exorcist/public/table"
	er "github.com/slugger7/exorcist/internal/errors"
	ff "github.com/slugger7/exorcist/internal/ffmpeg"
	"github.com/slugger7/exorcist/internal/media"
	ffmpeg "github.com/u2takey/ffmpeg-go"
)

func main() {
	path := "."
	err := godotenv.Load()
	er.CheckError(err)

	db := setupDB()

	libraryPathId := getOrCreateLibraryPathID(db, path)
	fmt.Printf("Library path id %v\n", libraryPathId)

	values, err := media.GetFilesByExtensions(path, []string{".mp4", ".m4v", ".mkv", ".avi", ".wmv", ".flv", ".webm", ".f4v", ".mpg", ".m2ts", ".mov"})
	er.CheckError(err)

	fmt.Println("Printing out results")
	videoModels := []model.Video{}
	for _, v := range values {
		fmt.Println(v)
		checksum := "lol"

		probeData, err := ffmpeg.Probe(v.Path)
		er.CheckError(err)

		fmt.Println(probeData)

		var data *ff.Probe
		err = json.Unmarshal([]byte(probeData), &data)
		er.CheckError(err)

		width, height, err := ff.GetDimensions(data.Streams)
		er.CheckError(err)

		runtime, err := strconv.Atoi(data.Format.Duration)
		if err != nil {
			fmt.Printf("Could not convert duration from string (%v) to int for video %v. Setting runtime to 0\n", data.Format.Duration, v)
			runtime = 0
		}
		size, err := strconv.Atoi(data.Format.Size)
		if err != nil {
			fmt.Printf("Could not convert size from string (%v) to int for video %v. Setting size to 0\n", data.Format.Size, v)
			size = 0
		}

		videoModels = append(videoModels, model.Video{
			LibraryPathID: libraryPathId,
			RelativePath:  v.Path,
			Title:         v.Name,
			FileName:      v.FileName,
			Height:        int32(height),
			Width:         int32(width),
			Runtime:       int64(runtime),
			Size:          int64(size),
			Checksum:      &checksum,
		})
	}

	insertStatement := table.Video.INSERT(
		table.Video.LibraryPathID,
		table.Video.RelativePath,
		table.Video.Title,
		table.Video.FileName,
		table.Video.Height,
		table.Video.Width,
		table.Video.Runtime,
		table.Video.Size,
		table.Video.Checksum,
	).
		MODELS(videoModels).
		RETURNING(table.Video.ID)

	var newVideos []struct {
		model.Video
	}
	err = insertStatement.Query(db, &newVideos)
	er.CheckError(err)
}

func getOrCreateLibraryPathID(db *sql.DB, path string) uuid.UUID {
	libraryPathId, err := getExistingLibraryPathID(db)
	if err != nil {
		libraryPathId = createLibWithPath(db, path)
	}
	return libraryPathId
}

func getExistingLibraryPathID(db *sql.DB) (uuid.UUID, error) {
	selectQuery := table.LibraryPath.
		SELECT(table.LibraryPath.ID).
		FROM(table.LibraryPath)

	var libraryPath []struct {
		model.LibraryPath
	}

	err := selectQuery.Query(db, &libraryPath)
	er.CheckError(err)

	if len(libraryPath) == 0 {
		return uuid.Nil, errors.New("no library path was found, first creat a library")
	}

	return libraryPath[0].ID, nil
}

func createLibWithPath(db *sql.DB, path string) uuid.UUID {
	newLib := model.Library{
		Name: "New Lib",
	}

	insertStatement := table.Library.INSERT(table.Library.Name).
		MODEL(newLib).
		RETURNING(table.Library.ID)

	var library []struct {
		model.Library
	}

	err := insertStatement.Query(db, &library)
	er.CheckError(err)

	newLibPath := model.LibraryPath{
		LibraryID: library[0].ID,
		Path:      path,
	}

	insertStatement = table.LibraryPath.INSERT(
		table.LibraryPath.LibraryID,
		table.LibraryPath.Path,
	).
		MODEL(newLibPath).
		RETURNING(table.LibraryPath.ID)

	var libraryPath []struct {
		model.LibraryPath
	}

	err = insertStatement.Query(db, &libraryPath)
	er.CheckError(err)

	return libraryPath[0].ID
}

func setupDB() *sql.DB {
	host := os.Getenv("DATABASE_HOST")
	port := os.Getenv("DATABASE_PORT")
	user := os.Getenv("DATABASE_USER")
	password := os.Getenv("DATABASE_PASSWORD")
	dbname := os.Getenv("DATABASE_NAME")

	fmt.Printf("host=%s port=%s user=%s password=%s database=%s", host, port, user, password, dbname)
	psqlconn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	fmt.Println("Opening DB")
	db, err := sql.Open("postgres", psqlconn)
	er.CheckError(err)
	defer db.Close()

	return db
}
