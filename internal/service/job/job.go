package jobService

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

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
	case model.JobTypeEnum_GenerateChapters:
		j, e = s.generateChapters(strData, *m.Priority)
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

func (i *jobService) generateChapters(data string, priority int16) (*model.Job, error) {
	var jobData dto.GenerateChaptersData
	if err := json.Unmarshal([]byte(data), &jobData); err != nil {
		return nil, errs.BuildError(err, "unmarshalling data for generate chapters data: %v", data)
	}

	if jobData.Interval == 0 {
		jobData.Interval = float64(((time.Minute * 5).Seconds()))
	}

	media, err := i.repo.Media().GetById(jobData.MediaId)
	if err != nil {
		return nil, errs.BuildError(err, "getting media by id: %v", jobData.MediaId.String())
	}

	if media == nil {
		return nil, fmt.Errorf("no media with id: %v", jobData.MediaId.String())
	}

	if media.Video == nil {
		return nil, fmt.Errorf("media is not of type video: %v", jobData.MediaId.String())
	}

	bytes, err := json.Marshal(jobData)
	if err != nil {
		return nil, errs.BuildError(err, "could not remarshall generate chapters data")
	}

	data = string(bytes)

	return &model.Job{
		Data:     &data,
		Priority: priority,
	}, nil
}

func (i *jobService) refreshLibraryMetadata(data string, priority int16) (*model.Job, error) {
	var jobData dto.RefreshLibraryMetadata
	if err := json.Unmarshal([]byte(data), &jobData); err != nil {
		return nil, errs.BuildError(err, "unmarshalling data for refresh library metadata: %v", data)
	}

	if jobData.RefreshFields == nil {
		jobData.RefreshFields = &dto.RefreshFields{
			Size:     true,
			Checksum: false,
		}
	}

	library, err := i.repo.Library().GetById(jobData.LibraryId)
	if err != nil {
		return nil, errs.BuildError(err, "getting library by id: %v", jobData.LibraryId.String())
	}

	if library == nil {
		return nil, fmt.Errorf("no library found with id: %v", jobData.LibraryId.String())
	}

	bytes, err := json.Marshal(jobData)
	if err != nil {
		return nil, errs.BuildError(err, "remarshalling refresh library metadata")
	}

	data = string(bytes)

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

	if generateThumbnailData.RelationType == nil {
		v := model.MediaRelationTypeEnum_Thumbnail
		generateThumbnailData.RelationType = &v
	}

	if _, err := i.repo.Video().GetByIdWithMedia(generateThumbnailData.MediaId); err != nil {
		return nil, errs.BuildError(
			err,
			ErrActionGenerateThumbnailVideoNotFound,
			generateThumbnailData.MediaId)
	}

	bytes, err := json.Marshal(generateThumbnailData)
	if err != nil {
		return nil, errs.BuildError(err, "could not remarshal generate thumbnail data")
	}

	data = string(bytes)

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
