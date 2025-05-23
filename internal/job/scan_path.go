package job

import (
	"encoding/json"
	"errors"
	"fmt"
	"path/filepath"
	"slices"
	"strconv"

	"github.com/google/uuid"
	"github.com/slugger7/exorcist/internal/db/exorcist/public/model"
	errs "github.com/slugger7/exorcist/internal/errors"
	"github.com/slugger7/exorcist/internal/ffmpeg"
	"github.com/slugger7/exorcist/internal/media"
	"github.com/slugger7/exorcist/internal/models"
)

var extensions = [...]string{".mp4", ".m4v", ".mkv", ".avi", ".wmv", ".flv", ".webm", ".f4v", ".mpg", ".m2ts", ".mov"}

const batchSize = 100

func (jr *JobRunner) getFilesByExtension(path string, extensions []string, ch chan []media.File) {
	defer jr.wg.Done()

	values, err := media.GetFilesByExtensions(path, extensions)
	if err != nil {
		jr.logger.Errorf("could not get files by extension: %v", err)
		ch <- nil
	}
	ch <- values
}

func (jr *JobRunner) ScanPath(job *model.Job) error {
	var data models.ScanPathData
	if err := json.Unmarshal([]byte(*job.Data), &data); err != nil {
		return errs.BuildError(err, "could not unmarshal scan path job data: %v", err)
	}

	libPath, err := jr.repo.LibraryPath().GetById(data.LibraryPathId)
	if err != nil {
		return errs.BuildError(err, "could not get library by id: %v", data.LibraryPathId)
	}
	if libPath == nil {
		return fmt.Errorf("library path not found: %v", data.LibraryPathId)
	}

	mediaChan := make(chan []media.File)
	jr.wg.Add(1)
	go jr.getFilesByExtension(libPath.Path, extensions[:], mediaChan)

	existingVideos, err := jr.repo.Video().GetByLibraryPathId(libPath.ID)
	if err != nil {
		return errs.BuildError(err, "could not get existing videos for library path: %v", libPath.ID)
	}

	select {
	case <-jr.shutdownCtx.Done():
		const msg string = "shutdown signal received. stopping"
		jr.logger.Warning(msg)
		return errors.New(msg)
	case videosOnDisk := <-mediaChan:

		nonExistentVideos := media.FindNonExistentVideos(existingVideos, videosOnDisk)
		if len(nonExistentVideos) > 0 {
			jr.removeVideos(nonExistentVideos)
		}

		accErrs := []error{}
		videoModels := []model.Video{}
		for i, v := range videosOnDisk {
			select {
			case <-jr.shutdownCtx.Done():
				return fmt.Errorf("partially done, ended due to shutdown")
			default:
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
					Runtime:       float64(runtime),
					Size:          int64(size),
				})

				if i%batchSize == 0 {
					batch := int(i / batchSize)
					batches := int(len(videosOnDisk) / batchSize)
					jr.logger.Infof("Writing batch %v/%v", batch, batches)
					if err := jr.writeNewVideoBatch(videoModels, job.ID); err != nil {
						jr.logger.Errorf("Error writing batch %v to database: %v", batch, err)
					}

					videoModels = []model.Video{}
				}
			}

		}

		if err := jr.writeNewVideoBatch(videoModels, job.ID); err != nil {
			jr.logger.Errorf("Error writing last batch of videos to database: %v", err)
		}

		if len(accErrs) > 0 {
			return errors.Join(accErrs...)
		}

		return nil
	}
}

func (jr *JobRunner) writeNewVideoBatch(videoModels []model.Video, jobId uuid.UUID) error {
	if len(videoModels) == 0 {
		return nil
	}

	vids, err := jr.repo.Video().Insert(videoModels)
	if err != nil {
		return errs.BuildError(err, "error writing batch of models to db")
	}

	jobs := []model.Job{}
	for _, v := range vids {
		dto := (&models.VideoOverviewDTO{}).FromModel(&v, nil)
		jr.wsVideoCreate(*dto)

		select {
		case <-jr.shutdownCtx.Done():
			return fmt.Errorf("shutdown signal received")
		default:
			checksumJob, err := CreateGenerateChecksumJob(v.ID, jobId)
			if err != nil {
				return errs.BuildError(err, "could not create checksum job")
			}
			jobs = append(jobs, *checksumJob)

			maxDimension := 400
			width, height := int(v.Width), int(v.Height)
			if v.Width > int32(maxDimension) {
				height = ffmpeg.ScaleHeightByWidth(height, width, maxDimension)
				width = maxDimension
			}

			if v.Height > int32(maxDimension) {
				width = ffmpeg.ScaleWidthByHeight(height, width, maxDimension)
				height = maxDimension
			}

			assetPath := filepath.Join(jr.env.Assets, v.ID.String(), fmt.Sprintf(`%v.webp`, v.FileName))
			thumbnailJob, err := CreateGenerateThumbnailJob(v, jobId, assetPath, 0, height, width)
			if err != nil {
				return errs.BuildError(err, "could not create generate thumbnail job")
			}

			jobs = append(jobs, *thumbnailJob)
		}
	}

	if _, err = jr.repo.Job().CreateAll(jobs); err != nil {
		return errs.BuildError(err, "error creating checksum jobs")
	}
	return nil
}

func (jr *JobRunner) removeVideos(nonExistentVideos []model.Video) {
	for _, v := range nonExistentVideos {
		select {
		case <-jr.shutdownCtx.Done():
			return
		default:
			v.Exists = false
			err := jr.repo.Video().UpdateExists(&v)
			if err != nil {
				jr.logger.Errorf("Error occured while updating the existance state of the video '%v': %v", v.ID, err)
			}

			jr.wsVideoDelete(models.VideoOverviewDTO{Id: v.ID, Deleted: true})
		}
	}
}

func videoExsists(existingVideos []model.Video, relativePath string) bool {
	return slices.ContainsFunc(existingVideos, func(existingVideo model.Video) bool {
		return existingVideo.RelativePath == relativePath
	})
}
