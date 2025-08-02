package jobService

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/slugger7/exorcist/internal/db/exorcist/public/model"
	"github.com/slugger7/exorcist/internal/dto"
	"github.com/slugger7/exorcist/internal/environment"
	errs "github.com/slugger7/exorcist/internal/errors"
	"github.com/slugger7/exorcist/internal/logger"
	"github.com/slugger7/exorcist/internal/repository"
)

type IJobService interface {
	Create(dto.CreateJobDTO) (*model.Job, error)
}

type jobService struct {
	env    *environment.EnvironmentVariables
	repo   repository.IRepository
	logger logger.ILogger
	jobCh  chan bool
	ctx    context.Context
}

var jobServiceInstance *jobService

func New(repo repository.IRepository, env *environment.EnvironmentVariables, jobCh chan bool, ctx context.Context) IJobService {
	if jobServiceInstance == nil {
		jobServiceInstance = &jobService{
			env:    env,
			repo:   repo,
			logger: logger.New(env),
			jobCh:  jobCh,
			ctx:    ctx,
		}

		jobServiceInstance.logger.Info("UserService instance created")
	}
	return jobServiceInstance
}

func (s *jobService) Create(m dto.CreateJobDTO) (*model.Job, error) {
	defaultJobPriority := dto.JobPriority_Medium
	if m.Priority == nil {
		m.Priority = &(defaultJobPriority)
	}
	data, err := json.Marshal(m.Data)
	strData := string(data)
	if err != nil {
		return nil, errs.BuildError(err, "could not marhsal data field")
	}
	var j *model.Job
	var e error
	switch m.Type {
	case model.JobTypeEnum_ScanPath:
		j, e = s.scanPath(strData, *m.Priority)
	case model.JobTypeEnum_GenerateThumbnail:
		j, e = s.generateThumbnail(strData, *m.Priority)
	case model.JobTypeEnum_RefreshMetadata:
		j, e = s.refreshMetadata(strData, *m.Priority)
	case model.JobTypeEnum_RefreshLibraryMetadata:
		j, e = s.refreshLibraryMetadata(strData, *m.Priority)
	default:
		return nil, fmt.Errorf("job type not implemented: %v", m.Type)
	}
	if e != nil {
		return nil, errs.BuildError(err, "error encountered while creating job")
	}

	job := model.Job{
		JobType:  m.Type,
		Status:   model.JobStatusEnum_NotStarted,
		Data:     j.Data,
		Priority: j.Priority,
	}

	jobs, err := s.repo.Job().CreateAll([]model.Job{job})
	if err != nil {
		return nil, errs.BuildError(err, "creating job")
	}

	if len(jobs) == 0 {
		return nil, fmt.Errorf("no jobs were returned after creating a job")
	}

	go s.startJobRunner()

	return &jobs[0], nil
}

func (i *jobService) refreshLibraryMetadata(data string, priority int16) (*model.Job, error) {
	var jobData dto.RefreshLibraryMetadata
	if err := json.Unmarshal([]byte(data), &jobData); err != nil {
		return nil, errs.BuildError(err, "unmarshalling data for refresh library metadata: %v", data)
	}

	library, err := i.repo.Library().GetById(jobData.LibraryId)
	if err != nil {
		return nil, errs.BuildError(err, "getting library by id: %v", jobData.LibraryId.String())
	}

	if library == nil {
		return nil, fmt.Errorf("no library found with id: %v", jobData.LibraryId.String())
	}

	return &model.Job{
		Data:     &data,
		Priority: priority,
	}, nil
}

func (i *jobService) refreshMetadata(data string, priority int16) (*model.Job, error) {
	var jobData dto.RefreshMetadata
	if err := json.Unmarshal([]byte(data), &jobData); err != nil {
		return nil, errs.BuildError(err, "unmarshalling data for refresh meta data: %v", data)
	}

	mediaEntity, err := i.repo.Media().GetById(jobData.MediaId)
	if err != nil {
		return nil, errs.BuildError(err, "fetching media entity by id: %v", jobData.MediaId.String())
	}

	if mediaEntity == nil {
		return nil, fmt.Errorf("no media entity found to refresh the metada of: %v", jobData.MediaId.String())
	}

	return &model.Job{
		Data:     &data,
		Priority: priority,
	}, nil

}

const ErrActionGenerateThumbnailVideoNotFound = "could not find video for generate thumbnail job: %v"

func (i *jobService) generateThumbnail(data string, priority int16) (*model.Job, error) {
	var generateThumbnailData dto.GenerateThumbnailData

	if err := json.Unmarshal([]byte(data), &generateThumbnailData); err != nil {
		return nil, errs.BuildError(err, "could not unmarshal data for job %v", data)
	}

	if _, err := i.repo.Video().GetByIdWithMedia(generateThumbnailData.VideoId); err != nil {
		return nil, errs.BuildError(
			err,
			ErrActionGenerateThumbnailVideoNotFound,
			generateThumbnailData.VideoId)
	}

	return &model.Job{
		Data:     &data,
		Priority: priority,
	}, nil
}

const ErrActionScanGetLibraryPaths = "could not get library paths in scan action"
const ErrCreatingJobs = "error creating jobs"

func (i *jobService) scanPath(data string, priority int16) (*model.Job, error) {
	var scanPathData dto.ScanPathData

	if err := json.Unmarshal([]byte(data), &scanPathData); err != nil {
		return nil, errs.BuildError(err, "could not unmarshall data for job %v", data)
	}

	// check to see if the library path actually exists before creating a job
	_, err := i.repo.LibraryPath().GetById(scanPathData.LibraryPathId)
	if err != nil {
		return nil, errs.BuildError(err, ErrActionScanGetLibraryPaths)
	}

	return &model.Job{
		Data:     &data,
		Priority: priority,
	}, nil
}

// We do this at the moment to stack a signal to the job runner if it is already running
func (i *jobService) startJobRunner() {
	i.logger.Debug("Starting a job runner")
	select {
	case <-i.ctx.Done():
		i.logger.Debug("Shutdown signal recieved. Not starting job runner")
		return
	default:
		i.logger.Debug("Starting job runner")
		if i.env.JobRunner {
			i.jobCh <- true
			i.logger.Debug("Job runner signal sent")
		}
	}
}
