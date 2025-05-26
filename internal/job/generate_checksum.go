package job

import (
	"encoding/json"

	"github.com/google/uuid"
	"github.com/slugger7/exorcist/internal/db/exorcist/public/model"
	errs "github.com/slugger7/exorcist/internal/errors"
	"github.com/slugger7/exorcist/internal/media"
	"github.com/slugger7/exorcist/internal/models"
)

type GenerateChecksumData struct {
	MediaId uuid.UUID `json:"mediaId"`
}

func CreateGenerateChecksumJob(mediaId, jobId uuid.UUID) (*model.Job, error) {
	d := GenerateChecksumData{
		MediaId: mediaId,
	}
	js, err := json.Marshal(d)
	if err != nil {
		return nil, errs.BuildError(err, "could not marshal generate checksum data for: %v", mediaId)
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

	jobMedia, err := jr.repo.Media().GetById(jobData.MediaId)
	if err != nil {
		return errs.BuildError(err, "error fetching video with library path by id: %v", jobData.MediaId)
	}

	jr.logger.Infof("Calculating checksum for %v", jobMedia.Path)

	checksum, err := media.CalculateMD5(jobMedia.Path)
	if err != nil {
		return errs.BuildError(err, "error calculating md5sum for %v", jobMedia.Path)
	}

	jobMedia.Checksum = &checksum

	if err := jr.repo.Media().UpdateChecksum(*jobMedia); err != nil {
		return errs.BuildError(err, "error updating video checksum")
	}

	return nil
}
