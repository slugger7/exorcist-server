package job

import (
	"encoding/json"
	"errors"
	"fmt"
	"path/filepath"
	"time"

	"github.com/slugger7/exorcist/internal/db/exorcist/public/model"
	"github.com/slugger7/exorcist/internal/dto"
	errs "github.com/slugger7/exorcist/internal/errors"
	"github.com/slugger7/exorcist/internal/ffmpeg"
)

func (jr *JobRunner) generateChapters(job *model.Job) error {
	var jobData dto.GenerateChaptersData
	if err := json.Unmarshal([]byte(*job.Data), &jobData); err != nil {
		return errs.BuildError(err, "error parsing job data for generate chapters: %v", job.Data)
	}

	media, err := jr.repo.Media().GetById(jobData.MediaId)
	if err != nil {
		return errs.BuildError(err, "could not find media by id for generate chapters job %v", jobData.MediaId.String())
	}

	if media == nil {
		return fmt.Errorf("media was nil for generate chapters job: %v", jobData.MediaId.String())
	}

	if media.Video == nil {
		return fmt.Errorf("media was not of type video: %v", jobData.MediaId.String())
	}

	runtimeDuration := time.Duration(int64(media.Video.Runtime * float64(time.Second)))
	intervalDuration := time.Duration(int64(jobData.Interval * float64(time.Second)))

	relationType := model.MediaRelationTypeEnum_Chapter

	if jobData.Height == 0 {
		jobData.Height = int(media.Video.Height)
	}

	if jobData.Width == 0 {
		jobData.Width = int(media.Video.Width)
	}

	if jobData.MaxDimension != 0 {
		if jobData.Width > jobData.MaxDimension {
			jobData.Height = ffmpeg.ScaleHeightByWidth(jobData.Height, jobData.Width, jobData.MaxDimension)
			jobData.Width = jobData.MaxDimension
		}

		if jobData.Height > jobData.MaxDimension {
			jobData.Width = ffmpeg.ScaleWidthByHeight(jobData.Height, jobData.Width, jobData.MaxDimension)
			jobData.Height = jobData.MaxDimension
		}
	}

	generateThumbnailJobs := []model.Job{}
	var accErr error
	for i := intervalDuration; i < runtimeDuration; i += intervalDuration {
		metadata := dto.ThumbnailMetadataDTO{
			Timestamp: i.Seconds(),
		}

		assetPath := filepath.Join(
			jr.env.Assets,
			media.Media.ID.String(),
			fmt.Sprintf(
				"%v.%v.%vx%v.%v.webp",
				filepath.Base(media.Media.Path),
				relationType.String(),
				jobData.Height,
				jobData.Width,
				i,
			))
		job, err := CreateGenerateThumbnailJob(*media.Video, &job.ID, assetPath, i.Seconds(), jobData.Height, jobData.Width, &relationType, &metadata)
		if err != nil {
			accErr = errors.Join(accErr, err)
			continue
		}

		generateThumbnailJobs = append(generateThumbnailJobs, *job)
	}

	if accErr != nil {
		jr.logger.Errorf("encountered while creating generate thumbnail jobs: %v", accErr.Error())
	}

	if len(generateThumbnailJobs) != 0 {
		if _, err := jr.repo.Job().CreateAll(generateThumbnailJobs); err != nil {
			return errs.BuildError(err, "creating generate thumbnail jobs")
		}
	}

	return nil
}
