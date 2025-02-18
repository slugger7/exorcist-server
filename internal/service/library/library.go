package libraryService

import (
	"encoding/json"
	"errors"
	"fmt"
	"slices"

	"github.com/google/uuid"
	"github.com/slugger7/exorcist/internal/db/exorcist/public/model"
	"github.com/slugger7/exorcist/internal/environment"
	errs "github.com/slugger7/exorcist/internal/errors"
	"github.com/slugger7/exorcist/internal/logger"
	"github.com/slugger7/exorcist/internal/models"
	"github.com/slugger7/exorcist/internal/repository"
)

type ILibraryService interface {
	Create(newLibrary model.Library) (*model.Library, error)
	GetAll() ([]model.Library, error)
	Action(id uuid.UUID, action string) error
}

type LibraryService struct {
	Env    *environment.EnvironmentVariables
	repo   repository.IRepository
	logger logger.ILogger
	jobCh  chan bool
}

var libraryServiceInstance *LibraryService

func New(repo repository.IRepository, env *environment.EnvironmentVariables, jobCh chan bool) ILibraryService {
	if libraryServiceInstance == nil {
		libraryServiceInstance = &LibraryService{
			Env:    env,
			repo:   repo,
			logger: logger.New(env),
			jobCh:  jobCh,
		}

		libraryServiceInstance.logger.Info("LibraryService instance created")
	}
	return libraryServiceInstance
}

const ErrLibraryByName = "Could not fetch library by name %v"

func (i *LibraryService) Create(newLibrary model.Library) (*model.Library, error) {
	library, err := i.repo.Library().
		GetLibraryByName(newLibrary.Name)
	if err != nil {
		return nil, errs.BuildError(err, ErrLibraryByName, newLibrary.Name)
	}
	if library != nil {
		return nil, fmt.Errorf("library named %v already exists", newLibrary.Name)
	}

	library, err = i.repo.Library().
		CreateLibrary(newLibrary.Name)
	if err != nil {
		return nil, errs.BuildError(err, "could not create library with name %v", newLibrary.Name)
	}

	return library, nil
}

const ErrGetLibraries = "could not getting libraries in repo"

func (i *LibraryService) GetAll() ([]model.Library, error) {
	libraries, err := i.repo.Library().GetLibraries()
	if err != nil {
		return nil, errs.BuildError(err, ErrGetLibraries)
	}

	return libraries, nil
}

const (
	ActionScan = "/scan"
)

var Actions = []string{ActionScan}

const ErrActionNotFound = "action was not found: %v"
const ErrFindInRepo = "error finding library in repo with id %v"

func (i *LibraryService) Action(id uuid.UUID, action string) error {
	if !slices.Contains(Actions, action) {
		return fmt.Errorf(ErrActionNotFound, action)
	}

	lib, err := i.repo.Library().GetLibraryById(id)
	if err != nil {
		return errs.BuildError(err, ErrFindInRepo, id)
	}

	switch action {
	case ActionScan:
		err := i.actionScan(lib)
		if err != nil {
			return errs.BuildError(err, "error setting up action scan")
		}
		return nil
	default:
		panic("Action was not found after being found")
	}
}

const ErrActionScanGetLibraryPaths = "could not get library paths in scan action"
const ErrCreatingJobs = "error creating jobs"

func (i *LibraryService) actionScan(library *model.Library) error {
	libraryPaths, err := i.repo.LibraryPath().GetByLibraryId(library.ID)
	if err != nil {
		return errs.BuildError(err, ErrActionScanGetLibraryPaths)
	}

	jobs := []model.Job{}
	erro := []error{}
	for _, l := range libraryPaths {
		data, err := json.Marshal(models.ScanPathData{
			LibraryPathId: l.ID,
		})
		if err != nil {
			erro = append(erro, errs.BuildError(err, "could not marshal data for scan job"))
			continue
		}
		strData := string(data)
		job := model.Job{
			JobType: model.JobTypeEnum_ScanPath,
			Status:  model.JobStatusEnum_NotStarted,
			Data:    &strData,
		}
		jobs = append(jobs, job)
	}

	if len(erro) != 0 {
		i.logger.Warningf("encountered some errors while creating scan path jobs: %v", errors.Join(erro...).Error())
	}

	jobs, err = i.repo.Job().CreateAll(jobs)
	if err != nil {
		return errs.BuildError(err, ErrCreatingJobs)
	}

	go i.startScanPathJob()

	return nil
}

// We do this at the moment to stack a signal to the job runner if it is already running
func (i *LibraryService) startScanPathJob() {
	i.jobCh <- true
}
