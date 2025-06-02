package job

import (
	"encoding/json"
	"errors"
	"fmt"
	"path/filepath"
	"slices"
	"strconv"

	"github.com/slugger7/exorcist/internal/db/exorcist/public/model"
	"github.com/slugger7/exorcist/internal/dto"
	errs "github.com/slugger7/exorcist/internal/errors"
	"github.com/slugger7/exorcist/internal/ffmpeg"
	"github.com/slugger7/exorcist/internal/media"
	"github.com/slugger7/exorcist/internal/models"
)

var videoExtensions = [...]string{".mp4", ".m4v", ".mkv", ".avi", ".wmv", ".flv", ".webm", ".f4v", ".mpg", ".m2ts", ".mov"}
var imageExtensions = [...]string{".jpg", ".png", ".webp"}

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
		return errs.BuildError(err, "could not get library path by id: %v", data.LibraryPathId)
	}
	if libPath == nil {
		return fmt.Errorf("library path not found: %v", data.LibraryPathId)
	}

	lib, err := jr.repo.Library().GetById(libPath.LibraryID)
	if err != nil {
		return errs.BuildError(err, "could not get library by id: %v", libPath.LibraryID)
	}
	if lib == nil {
		return fmt.Errorf("library not found: %v", libPath.LibraryID)
	}

	videoChan := make(chan []media.File)
	jr.wg.Add(1)
	go jr.getFilesByExtension(libPath.Path, videoExtensions[:], videoChan)

	imageChan := make(chan []media.File)
	jr.wg.Add(1)
	go jr.getFilesByExtension(libPath.Path, imageExtensions[:], imageChan)

	existingMedia, err := jr.repo.Media().GetByLibraryPathId(libPath.ID)
	if err != nil {
		return errs.BuildError(err, "could not get existing videos for library path: %v", libPath.ID)
	}

	select {
	case <-jr.shutdownCtx.Done():
		const msg string = "shutdown signal received. stopping"
		jr.logger.Warning(msg)
		return errors.New(msg)
	case imagesOnDisk := <-imageChan:
		_ = imagesOnDisk
		// TODO: handle images on disk
		return nil
	case videosOnDisk := <-videoChan:
		nonExistentMedia := media.FindNonExistentMedia(existingMedia, videosOnDisk)
		if len(nonExistentMedia) > 0 {
			jr.removeMedia(nonExistentMedia)
		}

		accErrs := []error{}
		for _, v := range videosOnDisk {
			select {
			case <-jr.shutdownCtx.Done():
				return fmt.Errorf("partially done, ended due to shutdown")
			default:
				if mediaExists(existingMedia, v.Path) {
					continue
				}

				data, err := ffmpeg.UnmarshalledProbe(v.Path)
				if err != nil {
					accErrs = append(accErrs, errs.BuildError(err, "could not get unmarshalled probe data: %v", v.Path))
					continue
				}

				newMediaModel := model.Media{
					LibraryPathID: libPath.ID,
					Title:         v.Name,
					Size:          v.Size,
					Path:          v.Path,
					MediaType:     model.MediaTypeEnum_Primary,
				}

				createdMedia, err := jr.repo.Media().Create([]model.Media{newMediaModel})
				if err != nil {
					accErrs = append(accErrs, errs.BuildError(err, "could not create media"))
					continue
				}
				if len(createdMedia) != 1 {
					accErrs = append(accErrs, fmt.Errorf("expected a created media but there was none"))
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

				mediaId := createdMedia[0].ID

				newVideoModel := model.Video{
					MediaID: mediaId,
					Height:  int32(height),
					Width:   int32(width),
					Runtime: float64(runtime),
				}

				createdVideos, err := jr.repo.Video().Insert([]model.Video{newVideoModel})
				if err != nil {
					accErrs = append(accErrs, errs.BuildError(err, "could not create video"))
				}

				dto := (&dto.MediaOverviewDTO{}).FromModel(&createdMedia[0], nil)
				jr.wsVideoCreate(*dto)

				checksumJob, err := CreateGenerateChecksumJob(mediaId, job.ID)
				if err != nil {
					accErrs = append(accErrs, errs.BuildError(err, "could not create checksum job for media %v in job %v", mediaId, job.ID))
				}

				maxDimension := 400
				if width > maxDimension {
					height = ffmpeg.ScaleHeightByWidth(height, width, maxDimension)
					width = maxDimension
				}

				if height > maxDimension {
					width = ffmpeg.ScaleWidthByHeight(height, width, maxDimension)
					height = maxDimension
				}

				assetPath := filepath.Join(jr.env.Assets, mediaId.String(), fmt.Sprintf(`%v.webp`, v.FileName))
				thumbnailJob, err := CreateGenerateThumbnailJob(createdVideos[0], job.ID, assetPath, 0, height, width)
				if err != nil {
					return errs.BuildError(err, "could not create generate thumbnail job")
				}

				jobs := []model.Job{*checksumJob, *thumbnailJob}

				jobs, err = jr.repo.Job().CreateAll(jobs)
				if err != nil {
					accErrs = append(accErrs, errs.BuildError(err, "could not create checksum and thumbnail job for video: %v", createdVideos[0].ID))
				}
			}

		}

		if len(accErrs) > 0 {
			return errors.Join(accErrs...)
		}

		return nil
	}
}

func (jr *JobRunner) removeMedia(nonExistentMedia []model.Media) {
	for _, v := range nonExistentMedia {
		select {
		case <-jr.shutdownCtx.Done():
			return
		default:
			v.Exists = false
			err := jr.repo.Media().UpdateExists(v)
			if err != nil {
				jr.logger.Errorf("Error occured while updating the existance state of the media '%v': %v", v.ID, err)
			}

			jr.wsVideoDelete(dto.MediaOverviewDTO{Id: v.ID, Deleted: true})
		}
	}
}

func mediaExists(existingMedia []model.Media, path string) bool {
	return slices.ContainsFunc(existingMedia, func(existingMedia model.Media) bool {
		return existingMedia.Path == path
	})
}
