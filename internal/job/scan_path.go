package job

import (
	"log"
	"slices"
	"strconv"
	"sync"

	"github.com/google/uuid"
	"github.com/slugger7/exorcist/internal/db/exorcist/public/model"
	"github.com/slugger7/exorcist/internal/environment"
	errs "github.com/slugger7/exorcist/internal/errors"
	ff "github.com/slugger7/exorcist/internal/ffmpeg"
	"github.com/slugger7/exorcist/internal/media"
	"github.com/slugger7/exorcist/internal/repository"
)

func getFilesByExtensionAsync(path string, extensions []string, ch chan []media.File, wg *sync.WaitGroup) {
	log.Println("Fetching files async")
	defer wg.Done()

	values, err := media.GetFilesByExtensions(path, extensions)
	errs.CheckError(err)
	ch <- values
}

func ScanPath(repo repository.IRepository) {
	env := environment.GetEnvironmentVariables()
	mediaFilesChannel := make(chan []media.File)
	var wg sync.WaitGroup
	wg.Add(1)
	go getFilesByExtensionAsync(
		env.MediaPath,
		[]string{".mp4", ".m4v", ".mkv", ".avi", ".wmv", ".flv", ".webm", ".f4v", ".mpg", ".m2ts", ".mov"},
		mediaFilesChannel,
		&wg,
	)

	libraryPath := getOrCreateLibraryPath(repo, env.MediaPath)
	log.Printf("Library path id %v\n", libraryPath.ID)

	existingVideos := getVideosInLibraryPath(repo, libraryPath.ID)

	log.Printf("Existing video count %v\n", len(existingVideos))

	values := <-mediaFilesChannel
	wg.Wait()

	nonExsistentVideos := media.FindNonExistentVideos(existingVideos, values)
	if len(nonExsistentVideos) > 0 {
		removeVideos(repo, nonExsistentVideos)
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
			Runtime:       int64(runtime), // FIXME: this value is off by a factor and needs fixing
			Size:          int64(size),
			Checksum:      nil,
		})

		if i%5 == 0 {
			writeModelsTodbBatch(repo, videoModels)

			videoModels = []model.Video{}
		}
	}

	writeModelsTodbBatch(repo, videoModels)

	GenerateChecksums(repo)
}

func removeVideos(repo repository.IRepository, nonExistentVideos []model.Video) {
	for _, v := range nonExistentVideos {
		v.Exists = false
		_, err := repo.VideoRepo().UpdateVideoExistsStatement(v).Exec()
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

func getVideosInLibraryPath(repo repository.IRepository, libraryPathId uuid.UUID) []struct{ model.Video } {
	var videos []struct {
		model.Video
	}
	err := repo.VideoRepo().
		GetVideosInLibraryPath(libraryPathId).
		Query(&videos)
	errs.CheckError(err)

	return videos
}

func writeModelsTodbBatch(repo repository.IRepository, models []model.Video) {
	if len(models) == 0 {
		return
	}
	log.Println("Writing batch")

	_, err := repo.VideoRepo().InsertVideosStatement(models).Exec()
	if err != nil {
		log.Printf("Error inserting new videos: %v", err)
	}
}

func printPercentage(index, total int) {
	log.Printf("Index: %v Total: %v Progress: %v%%\n", index, total, int(float64(index)/float64(total)*100.0))
}

func getOrCreateLibraryPath(repo repository.IRepository, path string) (libraryPath model.LibraryPath) {
	var libraryPaths []struct {
		model.LibraryPath
	}
	err := repo.LibraryPathRepo().GetLibraryPathsSelect().
		Query(&libraryPaths)
	if err == nil && libraryPaths != nil && len(libraryPaths) > 0 {
		libraryPath = libraryPaths[len(libraryPaths)-1].LibraryPath
	} else {
		if err != nil {
			log.Printf("An error occurred while looking for a library path %v", err.Error())
		}
		log.Println("Could not find a library path. Creating one")
		libraryPath = createLibWithPath(repo, path)
	}

	return libraryPath
}

func createLibWithPath(repo repository.IRepository, path string) model.LibraryPath {

	library, err := repo.LibraryRepo().
		CreateLibrary("New Lib")
	errs.CheckError(err)

	var libraryPaths []struct {
		model.LibraryPath
	}
	err = repo.LibraryPathRepo().
		CreateLibraryPath(library.ID, path).
		Query(&libraryPaths)
	errs.CheckError(err)

	return libraryPaths[len(libraryPaths)-1].LibraryPath
}
