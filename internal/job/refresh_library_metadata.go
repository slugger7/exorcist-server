package job

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/slugger7/exorcist/internal/db/exorcist/public/model"
	"github.com/slugger7/exorcist/internal/dto"
	errs "github.com/slugger7/exorcist/internal/errors"
)

func (jr *JobRunner) refreshLibraryMetadata(job *model.Job) error {
	var jobData dto.RefreshLibraryMetadata
	if err := json.Unmarshal([]byte(*job.Data), &jobData); err != nil {
		return errs.BuildError(err, "error parsing job data for refresh library metadata: %v", job.Data)
	}

	skip := 0
	for {
		batchNr := 1
		var pageRequest *dto.PageRequestDTO
		if jobData.BatchSize != 0 {
			batchNr = skip/jobData.BatchSize + 1
			pageRequest = &dto.PageRequestDTO{
				Skip:  skip,
				Limit: jobData.BatchSize,
			}
		}

		jr.logger.Infof("Batch: %v", batchNr)

		mediaPage, err := jr.repo.Media().GetByLibraryId(jobData.LibraryId, pageRequest, nil)
		if err != nil {
			return errs.BuildError(err, "fetching batch of media entities from repo")
		}

		if len(mediaPage.Data) == 0 {
			break
		}

		var accErr error
		refreshJobs := []model.Job{}
		for _, o := range mediaPage.Data {
			j, err := CreateRefreshMetadataJob(o, &job.ID, jobData.RefreshFields)
			if err != nil {
				accErr = errors.Join(accErr, err)
				continue
			}
			refreshJobs = append(refreshJobs, *j)
		}

		if accErr != nil {
			jr.logger.Errorf("encountered errors while processing batch %v: %v", batchNr, accErr.Error())
		}

		jobs, err := jr.repo.Job().CreateAll(refreshJobs)
		if err != nil {
			return errs.BuildError(err, "creating refresh metadata jobs for %v", jobData.LibraryId)
		}

		if len(jobs) != len(refreshJobs) {
			return fmt.Errorf("jobs created (%v) and jobs saved to database (%v) differed in batch %v", len(refreshJobs), len(jobs), batchNr)
		}

		skip = skip + jobData.BatchSize

		if jobData.BatchSize == 0 {
			break
		}
	}

	return nil
}
