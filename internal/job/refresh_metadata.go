package job

import (
	"encoding/json"
	"fmt"

	"github.com/go-jet/jet/v2/postgres"
	"github.com/google/uuid"
	"github.com/slugger7/exorcist/internal/db/exorcist/public/model"
	"github.com/slugger7/exorcist/internal/db/exorcist/public/table"
	"github.com/slugger7/exorcist/internal/dto"
	errs "github.com/slugger7/exorcist/internal/errors"
	"github.com/slugger7/exorcist/internal/media"
)

func CreateRefreshMetadataJob(media model.Media, jobId *uuid.UUID, refreshFields *dto.RefreshFields) (*model.Job, error) {
	localRefreshFields := *refreshFields
	if refreshFields == nil {
		localRefreshFields = dto.RefreshFields{
			Size:     true,
			Checksum: false,
		}
	}

	d := dto.RefreshMetadata{
		MediaId:       media.ID,
		RefreshFields: &localRefreshFields,
	}

	js, err := json.Marshal(d)
	if err != nil {
		return nil, errs.BuildError(err, "could not marshal refresh metadata dto")
	}

	data := string(js)
	job := &model.Job{
		JobType:  model.JobTypeEnum_RefreshMetadata,
		Status:   model.JobStatusEnum_NotStarted,
		Data:     &data,
		Parent:   jobId,
		Priority: dto.JobPriority_Low,
	}

	return job, nil
}

func (jr *JobRunner) RefreshMetadata(job *model.Job) error {
	var jobData dto.RefreshMetadata
	if err := json.Unmarshal([]byte(*job.Data), &jobData); err != nil {
		return errs.BuildError(err, "error parsing job data for refresh metadata: %v", job.Data)
	}

	mediaEntity, err := jr.repo.Media().GetById(jobData.MediaId)
	if err != nil {
		return errs.BuildError(err, "could not get media by id in refresh metadata: %v", jobData.MediaId.String())
	}

	if mediaEntity == nil {
		return fmt.Errorf("media entity was nil for %v", jobData.MediaId.String())
	}

	updateColumns := postgres.ColumnList{}

	if jobData.RefreshFields.Size {
		fileSize, err := media.GetFileSize(mediaEntity.Path)
		if err != nil {
			return errs.BuildError(err, "calculating file size for %v", mediaEntity.Path)
		}
		if fileSize != mediaEntity.Size {
			mediaEntity.Size = fileSize

			updateColumns = append(updateColumns, table.Media.Size)
		}
	}

	if jobData.RefreshFields.Checksum {
		checksum, err := media.CalculateMD5(mediaEntity.Path)
		if err != nil {
			return errs.BuildError(err, "calculating md5sum for %v", mediaEntity.Path)
		}
		if mediaEntity.Media.Checksum == nil || checksum != *mediaEntity.Media.Checksum {
			mediaEntity.Media.Checksum = &checksum

			updateColumns = append(updateColumns, table.Media.Checksum)
		}
	}

	if len(updateColumns) == 0 {
		return nil
	}

	if _, err := jr.repo.Media().Update(mediaEntity.Media, updateColumns); err != nil {
		return errs.BuildError(err, "saving updated metadata to media %v", mediaEntity.Media.ID.String())
	}

	return nil
}
