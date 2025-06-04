package mediaService

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/slugger7/exorcist/internal/db/exorcist/public/model"
	"github.com/slugger7/exorcist/internal/environment"
	errs "github.com/slugger7/exorcist/internal/errors"
	"github.com/slugger7/exorcist/internal/logger"
	"github.com/slugger7/exorcist/internal/models"
	"github.com/slugger7/exorcist/internal/repository"
	personService "github.com/slugger7/exorcist/internal/service/person"
)

type MediaService interface {
	AddPeople(id uuid.UUID, people []string) (*models.Media, error)
}

type mediaService struct {
	env           *environment.EnvironmentVariables
	repo          repository.IRepository
	logger        logger.ILogger
	personService personService.IPersonService
}

// AddPeople implements MediaService.
func (m *mediaService) AddPeople(id uuid.UUID, people []string) (*models.Media, error) {
	mediaModel, err := m.repo.Media().GetById(id)
	if err != nil {
		return nil, errs.BuildError(err, "could not get media by id: %v", id.String())
	}

	if mediaModel == nil {
		return nil, fmt.Errorf("could not find media by id")
	}

	peopleModels := []model.Person{}
	errorList := []error{}
	for _, p := range people {
		personModel, err := m.personService.Upsert(p)
		if err != nil {
			errorList = append(errorList, err)
			continue
		}
		peopleModels = append(peopleModels, *personModel)
	}

	if len(errorList) != 0 {
		joinedErrors := errors.Join(errorList...)
		m.logger.Errorf("error while upserting people: %v", errs.BuildError(joinedErrors, "could not upsert some people").Error())
	}

	mediaPersonModels := make([]model.MediaPerson, len(peopleModels))
	for i, p := range peopleModels {
		mediaPersonModels[i] = model.MediaPerson{
			MediaID:  id,
			PersonID: p.ID,
		}
	}

	mediaPersonModels, err = m.repo.Person().LinkWithMedia(mediaPersonModels)
	if err != nil {
		return nil, errs.BuildError(err, "error linking media with people")
	}

	mediaModel, err = m.repo.Media().GetById(id)
	if err != nil {
		return nil, errs.BuildError(err, "could not get media by id after adding people: %v", id)
	}

	return mediaModel, nil
}

var mediaServiceInstance *mediaService

func New(env *environment.EnvironmentVariables, repo repository.IRepository, personService personService.IPersonService) MediaService {
	if mediaServiceInstance == nil {
		mediaServiceInstance = &mediaService{
			env:           env,
			repo:          repo,
			logger:        logger.New(env),
			personService: personService,
		}

		mediaServiceInstance.logger.Info("Created media service instance")
	}

	return mediaServiceInstance
}
