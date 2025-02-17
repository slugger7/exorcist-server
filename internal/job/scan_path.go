package job

import (
	"encoding/json"
	"slices"
	"strconv"

	"github.com/google/uuid"
	"github.com/slugger7/exorcist/internal/db/exorcist/public/model"
	errs "github.com/slugger7/exorcist/internal/errors"
	"github.com/slugger7/exorcist/internal/ffmpeg"
	"github.com/slugger7/exorcist/internal/media"
)

var extensions = [...]string{".mp4", ".m4v", ".mkv", ".avi", ".wmv", ".flv", ".webm", ".f4v", ".mpg", ".m2ts", ".mov"}

const batchSize = 100

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
		jr.removeVideos(nonExistentVideos)
	}

	accErrs := []error{}
	videoModels := []model.Video{}
	for i, v := range videosOnDisk {
		relativePath := media.GetRelativePath(libPath.Path, v.Path)

		if videoExsists(existingVideos, relativePath) {
			continue
		}

		data, err := ffmpeg.UnmarshalledProbe(v.Path)
		if err != nil {
			accErrs = append(accErrs, errs.BuildError(err, "could not get unmarshalled probe data: %v", v.Path))
			continue
		}

		width, height, err := ffmpeg.GetDimensions(data.Streams)
		if err != nil {
			jr.logger.Warningf("could not extract dimensions for %v. Setting to 0. Reason: %v", v.Path, err)
		}

		runtime, err := strconv.ParseFloat(data.Format.Duration, 32)
		if err != nil {
			jr.logger.Warningf("could not convert duration from string (%v) to float for video %v. Setting runtime to 0. Reason: %v", data.Format.Duration, v.Path, err)
		}

		size, err := strconv.Atoi(data.Format.Size)
		if err != nil {
			jr.logger.Warningf("could not convert size from string (%v) to int for video %v. Setting to 0. Reason: %v", data.Format.Size, v.Path, err)
		}

		videoModels = append(videoModels, model.Video{
			LibraryPathID: libPath.ID,
			RelativePath:  relativePath,
			Title:         v.Name,
			FileName:      v.FileName,
			Height:        int32(height),
			Width:         int32(width),
			Runtime:       int64(runtime), // FIXME: this value is off by a factor and needs fixing
			Size:          int64(size),
		})

		if i%batchSize == 0 {
			if err := jr.writeNewVideoBatch(videoModels); err != nil {
				jr.logger.Errorf("Error wirting batch %v to database: %v", int(i/batchSize), err)
			}

			videoModels = []model.Video{}
		}
	}

	if err := jr.writeNewVideoBatch(videoModels); err != nil {
		jr.logger.Errorf("Error writing last batch of videos to database: %v", err)
	}

	// TODO: generate checksum jobs

	job.Status = model.JobStatusEnum_Completed
	if err := jr.repo.Job().UpdateJobStatus(job); err != nil {
		return errs.BuildError(err, "could not update job status to %v", job.Status)
	}

	return nil
}

func (jr *JobRunner) writeNewVideoBatch(models []model.Video) error {
	if len(models) == 0 {
		return nil
	}

	jr.logger.Debug("Writing batch")

	if err := jr.repo.Video().Insert(models); err != nil {
		return errs.BuildError(err, "error writing batch of models to db")
	}
	return nil
}

func (jr *JobRunner) removeVideos(nonExistentVideos []model.Video) {
	for _, v := range nonExistentVideos {
		v.Exists = false
		err := jr.repo.Video().UpdateVideoExists(v)
		if err != nil {
			jr.logger.Errorf("Error occured while updating the existance state of the video '%v': %v", v.ID, err)
		}
	}
}

func videoExsists(existingVideos []model.Video, relativePath string) bool {
	return slices.ContainsFunc(existingVideos, func(existingVideo model.Video) bool {
		return existingVideo.RelativePath == relativePath
	})
}
