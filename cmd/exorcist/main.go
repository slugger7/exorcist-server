package main

import (
	"database/sql"
	"log"
	"slices"
	"strconv"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/slugger7/exorcist/internal/db"
	"github.com/slugger7/exorcist/internal/db/exorcist/public/model"
	"github.com/slugger7/exorcist/internal/environment"
	errs "github.com/slugger7/exorcist/internal/errors"
	ff "github.com/slugger7/exorcist/internal/ffmpeg"
	"github.com/slugger7/exorcist/internal/job"
	"github.com/slugger7/exorcist/internal/media"

	libraryRepository "github.com/slugger7/exorcist/internal/repository/library"
	libraryPathRepository "github.com/slugger7/exorcist/internal/repository/library_path"
	videoRepository "github.com/slugger7/exorcist/internal/repository/video"
)

func main() {
	err := godotenv.Load()
	errs.CheckError(err)
	env := environment.GetEnvironmentVariables()

	database := db.NewDatabase(env)
	defer database.Close()

	err = db.RunMigrations(database, env)
	if err != nil {
		log.Printf("Error occured when running migrations: %v", err.Error())
	}

	libraryPath := getOrCreateLibraryPath(database, env.MediaPath)
	log.Printf("Library path id %v\n", libraryPath.ID)

	existingVideos := getVideosInLibraryPath(database, libraryPath.ID)

	log.Printf("Existing video count %v\n", len(existingVideos))

	values, err := media.GetFilesByExtensions(env.MediaPath, []string{".mp4", ".m4v", ".mkv", ".avi", ".wmv", ".flv", ".webm", ".f4v", ".mpg", ".m2ts", ".mov"})
	errs.CheckError(err)

	nonExsistentVideos := media.FindNonExistentVideos(existingVideos, values)
	if len(nonExsistentVideos) > 0 {
		removeVideos(database, nonExsistentVideos)
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
			writeModelsToDatabaseBatch(database, videoModels)

			videoModels = []model.Video{}
		}
	}

	writeModelsToDatabaseBatch(database, videoModels)

	job.GenerateChecksums(database) // this is a canditate to move to a goroutine
}

func removeVideos(db *sql.DB, nonExistentVideos []model.Video) {
	for _, v := range nonExistentVideos {
		v.Exists = false
		err := videoRepository.ExecuteUpdate(db, videoRepository.UpdateVideoExistsStatement(v))
		if err != nil {
			log.Printf("Error occured while updating the existance state of the video '%v': %v", v.ID, err)
		}
	}
}

func videoExsists(existingVideos []struct{ model.Video }, relativePath string) bool {
	return slices.ContainsFunc(existingVideos, func(existingVideo struct{ model.Video }) bool {
		return existingVideo.RelativePath == relativePath
	})
}

func getVideosInLibraryPath(db *sql.DB, libraryPathId uuid.UUID) []struct{ model.Video } {
	videos, err := videoRepository.QuerySelect(db, videoRepository.GetVideosInLibraryPath(libraryPathId))
	errs.CheckError(err)

	return videos
}

func writeModelsToDatabaseBatch(db *sql.DB, models []model.Video) {
	if len(models) == 0 {
		return
	}
	log.Println("Writing batch")

	err := videoRepository.ExecuteInsert(db, videoRepository.InsertVideosStatement(models))
	if err != nil {
		log.Printf("Error inserting new videos: %v", err)
	}
}

func printPercentage(index, total int) {
	log.Printf("Index: %v Total: %v Progress: %v%%\n", index, total, int(float64(index)/float64(total)*100.0))
}

func getOrCreateLibraryPath(db *sql.DB, path string) (libraryPath model.LibraryPath) {
	selectLibraryPath := libraryPathRepository.GetLibraryPathsSelect()
	libraryPaths, err := libraryPathRepository.QuerySelect(db, selectLibraryPath)
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
	libraryInsertStament := libraryRepository.CreateLibraryStatement("New Lib")
	libraries, err := libraryRepository.QueryInsert(db, libraryInsertStament)
	errs.CheckError(err)

	libraryPathInsertStatement := libraryPathRepository.CreateLibraryPath(libraries[len(libraries)-1].ID, path)
	libraryPaths, err := libraryPathRepository.QueryInsert(db, libraryPathInsertStatement)
	errs.CheckError(err)

	return libraryPaths[len(libraryPaths)-1].LibraryPath
}
