package job

import (
	"encoding/json"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/slugger7/exorcist/internal/db/exorcist/public/model"
	errs "github.com/slugger7/exorcist/internal/errors"
	"github.com/slugger7/exorcist/internal/media"
	"github.com/slugger7/exorcist/internal/models"
)

type GenerateChecksumData struct {
	VideoId uuid.UUID `json:"videoId"`
}

func CreateGenerateChecksumJob(videoId, jobId uuid.UUID) (*model.Job, error) {
	d := GenerateChecksumData{
		VideoId: videoId,
	}
	js, err := json.Marshal(d)
	if err != nil {
		return nil, errs.BuildError(err, "could not marshal generate checksum data for: %v", videoId)
	}
	data := string(js)
	job := model.Job{
		JobType:  model.JobTypeEnum_GenerateChecksum,
		Status:   model.JobStatusEnum_NotStarted,
		Data:     &data,
		Parent:   &jobId,
		Priority: models.JobPriority_Low,
	}

	return &job, nil
}

func (jr *JobRunner) GenerateChecksum(job *model.Job) error {
	var jobData GenerateChecksumData
	if err := json.Unmarshal([]byte(*job.Data), &jobData); err != nil {
		return errs.BuildError(err, "error parsing job data: %v", job.Data)
	}

	video, err := jr.repo.Video().GetByIdWithLibraryPath(jobData.VideoId)
	if err != nil {
		return errs.BuildError(err, "error fetching video with library path by id: %v", jobData.VideoId)
	}

	absolutePath := filepath.Join(video.LibraryPath.Path, video.RelativePath)
	jr.logger.Infof("Calculating checksum for %v", absolutePath)

	checksum, err := media.CalculateMD5(absolutePath)
	if err != nil {
		return errs.BuildError(err, "error calculating md5sum for %v", absolutePath)
	}

	video.Video.Checksum = &checksum

	if err := jr.repo.Video().UpdateChecksum(&video.Video); err != nil {
		return errs.BuildError(err, "error updating video checksum")
	}

	return nil
}
