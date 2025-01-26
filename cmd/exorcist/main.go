package main

import (
	"database/sql"
	"log"
	"slices"
	"strconv"

	"github.com/go-jet/jet/v2/postgres"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/slugger7/exorcist/internal/constants/environment"
	"github.com/slugger7/exorcist/internal/db"
	"github.com/slugger7/exorcist/internal/db/exorcist/public/model"
	"github.com/slugger7/exorcist/internal/db/exorcist/public/table"
	er "github.com/slugger7/exorcist/internal/errors"
	ff "github.com/slugger7/exorcist/internal/ffmpeg"
	"github.com/slugger7/exorcist/internal/media"

	libRepo "github.com/slugger7/exorcist/internal/repository/library"
	libPathRepo "github.com/slugger7/exorcist/internal/repository/library_path"
)

func main() {
	err := godotenv.Load()
	er.CheckError(err)
	env := environment.GetEnvironmentVariables()

	db := db.NewDatabase(env)
	defer db.Close()

	libraryPath := getOrCreateLibraryPath(db, env.MediaPath)
	log.Printf("Library path id %v\n", libraryPath.ID)

	existingVideos := getVideosInLibraryPath(db, libraryPath.ID)

	log.Printf("Existing video count %v\n", len(existingVideos))

	values, err := media.GetFilesByExtensions(env.MediaPath, []string{".mp4", ".m4v", ".mkv", ".avi", ".wmv", ".flv", ".webm", ".f4v", ".mpg", ".m2ts", ".mov"})
	er.CheckError(err)

	nonExsistentVideos := media.FindNonExistentVideos(existingVideos, values)
	if len(nonExsistentVideos) > 0 {
		log.Println("Found some videos that do not exist any more on disk. Marking them as deleted.")
		removeVideos(db, nonExsistentVideos)
	}

	log.Println("Printing out results")
	videoModels := []model.Video{}
	for i, v := range values {
		printPercentage(i, len(values))
		relativePath := media.GetRelativePath(libraryPath.Path, v.Path)

		if videoExsists(existingVideos, relativePath) {
			continue
		}

		data, err := ff.UnmarshalledProbe(v.Path)
		if err != nil {
			log.Printf("Unmarshaling failed for %v\nThe error was %v", v.Path, err.Error())
			continue
		}

		width, height, err := ff.GetDimensions(data.Streams)
		if err != nil {
			log.Printf("Colud not extract dimensions. Setting to 0 %v\n", err.Error())
		}

		runtime, err := strconv.ParseFloat(data.Format.Duration, 32)
		if err != nil {
			log.Printf("Could not convert duration from string (%v) to float for video %v. Setting runtime to 0\n", data.Format.Duration, v)
			runtime = 0
		}
		size, err := strconv.Atoi(data.Format.Size)
		if err != nil {
			log.Printf("Could not convert size from string (%v) to int for video %v. Setting size to 0\n", data.Format.Size, v)
			size = 0
		}

		videoModels = append(videoModels, model.Video{
			LibraryPathID: libraryPath.ID,
			RelativePath:  relativePath,
			Title:         v.Name,
			FileName:      v.FileName,
			Height:        int32(height),
			Width:         int32(width),
			Runtime:       int64(runtime),
			Size:          int64(size),
			Checksum:      nil,
		})

		if i%5 == 0 {
			writeModelsToDatabaseBatch(db, videoModels)

			videoModels = []model.Video{}
		}
	}

	writeModelsToDatabaseBatch(db, videoModels)
}

func removeVideos(db *sql.DB, nonExistentVideos []model.Video) {
	for _, v := range nonExistentVideos {
		updateStmnt := table.Video.UPDATE().
			SET(table.Video.Deleted.SET(postgres.Bool(true))).
			MODEL(v).
			WHERE(table.Video.ID.EQ(postgres.UUID(v.ID)))
		dbgSql := updateStmnt.DebugSql()
		log.Println(dbgSql)
		_, err := updateStmnt.Exec(db)
		if err != nil {
			log.Printf("Could not update video %v to be deleted: %v", v.ID, err.Error())
			continue
		}
	}
}

func videoExsists(existingVideos []struct{ model.Video }, relativePath string) bool {
	return slices.ContainsFunc(existingVideos, func(existingVideo struct{ model.Video }) bool {
		return existingVideo.RelativePath == relativePath
	})
}

func getVideosInLibraryPath(db *sql.DB, libraryPathId uuid.UUID) []struct{ model.Video } {
	findStatement := table.Video.SELECT(table.Video.RelativePath, table.Video.ID).
		FROM(table.Video.Table).
		WHERE(table.Video.LibraryPathID.EQ(postgres.UUID(libraryPathId)))

	var videos []struct {
		model.Video
	}
	err := findStatement.Query(db, &videos)
	er.CheckError(err)

	return videos
}

func writeModelsToDatabaseBatch(db *sql.DB, models []model.Video) {
	if len(models) == 0 {
		return
	}
	log.Println("Writing batch")

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
		MODELS(models).
		RETURNING(table.Video.ID)

	var newVideos []struct {
		model.Video
	}
	err := insertStatement.Query(db, &newVideos)
	er.CheckError(err)
}

func printPercentage(index, total int) {
	log.Printf("Index: %v Total: %v Progress: %v\n", index, total, int(float64(index)/float64(total)*100.0))
}

func getOrCreateLibraryPath(db *sql.DB, path string) (libraryPath model.LibraryPath) {
	selectLibraryPath := libPathRepo.GetLibraryPathsSelect()
	libraryPaths, err := libPathRepo.ExecuteSelect(db, selectLibraryPath)
	if err == nil && libraryPaths != nil && len(libraryPaths) > 0 {
		libraryPath = libraryPaths[len(libraryPaths)-1].LibraryPath
	} else {
		if err != nil {
			log.Printf("An error occurred while looking for a library path %v", err.Error())
		}
		log.Println("Could not find a library path. Creating one")
		libraryPath = createLibWithPath(db, path)
	}

	return libraryPath
}

func createLibWithPath(db *sql.DB, path string) model.LibraryPath {
	libraryInsertStament := libRepo.CreateLibraryStatement("New Lib")
	libraries, err := libRepo.ExecuteInsert(db, libraryInsertStament)
	er.CheckError(err)

	libraryPathInsertStatement := libPathRepo.CreateLibraryPath(libraries[len(libraries)-1].ID, path)
	libraryPaths, err := libPathRepo.ExecuteInsert(db, libraryPathInsertStatement)
	er.CheckError(err)

	return libraryPaths[len(libraryPaths)-1].LibraryPath
}
