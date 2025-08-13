package personService

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/slugger7/exorcist/internal/db/exorcist/public/model"
	"github.com/slugger7/exorcist/internal/dto"
	"github.com/slugger7/exorcist/internal/environment"
	errs "github.com/slugger7/exorcist/internal/errors"
	"github.com/slugger7/exorcist/internal/logger"
	"github.com/slugger7/exorcist/internal/models"
	"github.com/slugger7/exorcist/internal/repository"
)

type PersonService interface {
	Upsert(name string) (*model.Person, error)
	GetMedia(id, userId uuid.UUID, search dto.MediaSearchDTO) (*dto.PageDTO[models.MediaOverviewModel], error)
}

type personService struct {
	env    *environment.EnvironmentVariables
	repo   repository.Repository
	logger logger.Logger
}

// GetMedia implements IPersonService.
func (p *personService) GetMedia(id, userId uuid.UUID, search dto.MediaSearchDTO) (*dto.PageDTO[models.MediaOverviewModel], error) {
	person, err := p.repo.Person().GetById(id)
	if err != nil {
		return nil, errs.BuildError(err, "could not get person by id from repo: %v", id)
	}

	if person == nil {
		return nil, fmt.Errorf("no person found with id: %v", id)
	}

	media, err := p.repo.Person().GetMedia(id, userId, search)
	if err != nil {
		return nil, errs.BuildError(err, "colud not get media by person id from repo: %v", id)
	}

	return media, nil
}

// Upsert implements IPersonService.
func (p *personService) Upsert(name string) (*model.Person, error) {
	person, err := p.repo.Person().GetByName(name)
	if err != nil {
		return nil, errs.BuildError(err, "could not get person by name from repo")
	}

	if person == nil {
		people, err := p.repo.Person().Create([]string{name})
		if err != nil {
			return nil, errs.BuildError(err, "could not create person by name")
		}
		person = &people[0]
	}

	return person, nil
}

var personServiceInstance *personService

func New(repo repository.Repository, env *environment.EnvironmentVariables) PersonService {
	if personServiceInstance == nil {
		personServiceInstance = &personService{
			env:    env,
			repo:   repo,
			logger: logger.New(env),
		}

		personServiceInstance.logger.Info("PersonService instance created")
	}

	return personServiceInstance
}
