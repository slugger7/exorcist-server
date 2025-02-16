package job

import (
	"encoding/json"

	"github.com/google/uuid"
	"github.com/slugger7/exorcist/internal/db/exorcist/public/model"
	errs "github.com/slugger7/exorcist/internal/errors"
	"github.com/slugger7/exorcist/internal/media"
)

var extensions = [...]string{".mp4", ".m4v", ".mkv", ".avi", ".wmv", ".flv", ".webm", ".f4v", ".mpg", ".m2ts", ".mov"}

func (jr *JobRunner) getFilesByExtension(path string, extensions []string, ch chan []media.File) {
	values, err := media.GetFilesByExtensions(path, extensions)
	if err != nil {
		jr.logger.Errorf("could not get files by extension: %v", err)
		ch <- nil
	}
	ch <- values
}

type ScanPathData struct {
	LibraryPathId uuid.UUID `json:"libraryPathId"`
}

func (jr *JobRunner) ScanPath(job *model.Job) error {
	var data ScanPathData
	if err := json.Unmarshal([]byte(*job.Data), &data); err != nil {
		return errs.BuildError(err, "could not unmarshal scan path job data: %v", err)
	}

	libPath, err := jr.repo.LibraryPath().GetById(data.LibraryPathId)
	if err != nil {
		return errs.BuildError(err, "could not get library by id: %v", data.LibraryPathId)
	}

	mediaChan := make(chan []media.File)

	go jr.getFilesByExtension(libPath.Path, extensions[:], mediaChan)

	existingVideos, err := jr.repo.Video().GetByLibraryPathId(libPath.ID)
	if err != nil {
		return errs.BuildError(err, "could not get existing videos for library path: %v", libPath.ID)
	}

	videosOnDisk := <-mediaChan

	nonExistentVideos := media.FindNonExistentVideos(existingVideos, videosOnDisk)
	if len(nonExistentVideos) > 0 {

	}

	return nil
}
