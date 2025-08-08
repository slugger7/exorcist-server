package job

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/slugger7/exorcist/internal/db/exorcist/public/model"
	"github.com/slugger7/exorcist/internal/dto"
	errs "github.com/slugger7/exorcist/internal/errors"
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
	intervalDuration := time.Duration(int64(jobData.Interval) * int64(time.Millisecond))

	for i := 0; i < int(runtimeDuration); i += int(intervalDuration) {
		// TODO: create generate thumbnail jobs
	}

	return nil
}
