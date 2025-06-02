package jobService

import (
	"encoding/json"
	"fmt"

	"github.com/slugger7/exorcist/internal/db/exorcist/public/model"
	"github.com/slugger7/exorcist/internal/dto"
	"github.com/slugger7/exorcist/internal/environment"
	errs "github.com/slugger7/exorcist/internal/errors"
	"github.com/slugger7/exorcist/internal/logger"
	"github.com/slugger7/exorcist/internal/models"
	"github.com/slugger7/exorcist/internal/repository"
)

type IJobService interface {
	Create(dto.CreateJobDTO) (*model.Job, error)
}

type JobService struct {
	Env    *environment.EnvironmentVariables
	repo   repository.IRepository
	logger logger.ILogger
	jobCh  chan bool
}

var jobServiceInstance *JobService

func New(repo repository.IRepository, env *environment.EnvironmentVariables, jobCh chan bool) IJobService {
	if jobServiceInstance == nil {
		jobServiceInstance = &JobService{
			Env:    env,
			repo:   repo,
			logger: logger.New(env),
			jobCh:  jobCh,
		}

		jobServiceInstance.logger.Info("UserService instance created")
	}
	return jobServiceInstance
}

func (s *JobService) Create(m dto.CreateJobDTO) (*model.Job, error) {
	defaultJobPriority := dto.JobPriority_Medium
	if m.Priority == nil {
		m.Priority = &(defaultJobPriority)
	}
	data, err := json.Marshal(m.Data)
	strData := string(data)
	if err != nil {
		return nil, errs.BuildError(err, "could not marhsal data field")
	}
	switch m.Type {
	case model.JobTypeEnum_ScanPath:
		return s.scanPath(strData, *m.Priority)
	case model.JobTypeEnum_GenerateThumbnail:
		return s.generateThumbnail(strData, *m.Priority)
	default:
		return nil, fmt.Errorf("job type not implemented: %v", m.Type)
	}
}

const ErrActionGenerateThumbnailVideoNotFound = "could not find video for generate thumbnail job: %v"

func (i *JobService) generateThumbnail(data string, priority int16) (*model.Job, error) {
	var generateThumbnailData models.GenerateThumbnailData

	if err := json.Unmarshal([]byte(data), &generateThumbnailData); err != nil {
		return nil, errs.BuildError(err, "could not unmarshal data for job %v", data)
	}

	if _, err := i.repo.Video().GetByIdWithMedia(generateThumbnailData.VideoId); err != nil {
		return nil, errs.BuildError(
			err,
			ErrActionGenerateThumbnailVideoNotFound,
			generateThumbnailData.VideoId)
	}

	job := model.Job{
		JobType:  model.JobTypeEnum_GenerateThumbnail,
		Status:   model.JobStatusEnum_NotStarted,
		Data:     &data,
		Priority: priority,
	}

	jobs, err := i.repo.Job().CreateAll([]model.Job{job})
	if err != nil {
		return nil, errs.BuildError(err, ErrCreatingJobs)
	}

	go i.startJobRunner()

	return &jobs[0], nil
}

const ErrActionScanGetLibraryPaths = "could not get library paths in scan action"
const ErrCreatingJobs = "error creating jobs"

func (i *JobService) scanPath(data string, priority int16) (*model.Job, error) {
	var scanPathData models.ScanPathData

	if err := json.Unmarshal([]byte(data), &scanPathData); err != nil {
		return nil, errs.BuildError(err, "could not unmarshall data for job %v", data)
	}

	// check to see if the library path actually exists before creating a job
	_, err := i.repo.LibraryPath().GetById(scanPathData.LibraryPathId)
	if err != nil {
		return nil, errs.BuildError(err, ErrActionScanGetLibraryPaths)
	}

	job := model.Job{
		JobType:  model.JobTypeEnum_ScanPath,
		Status:   model.JobStatusEnum_NotStarted,
		Data:     &data,
		Priority: priority,
	}

	jobs, err := i.repo.Job().CreateAll([]model.Job{job})
	if err != nil {
		return nil, errs.BuildError(err, ErrCreatingJobs)
	}

	go i.startJobRunner()

	return &jobs[0], nil
}

// We do this at the moment to stack a signal to the job runner if it is already running
func (i *JobService) startJobRunner() {
	if i.Env.JobRunner {
		i.jobCh <- true
	}
}
