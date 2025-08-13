package tagService

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

type TagService interface {
	Upsert(name string) (*model.Tag, error)
	GetMedia(id, userId uuid.UUID, search dto.MediaSearchDTO) (*dto.PageDTO[models.MediaOverviewModel], error)
}

type tagService struct {
	env    *environment.EnvironmentVariables
	repo   repository.Repository
	logger logger.Logger
}

// GetMedia implements TagService.
func (p *tagService) GetMedia(id, userId uuid.UUID, search dto.MediaSearchDTO) (*dto.PageDTO[models.MediaOverviewModel], error) {
	tag, err := p.repo.Tag().GetById(id)
	if err != nil {
		return nil, errs.BuildError(err, "could net get tag by id from repo: %v", id)
	}

	if tag == nil {
		return nil, fmt.Errorf("no tag found with id: %v", id)
	}

	media, err := p.repo.Tag().GetMedia(id, userId, search)
	if err != nil {
		return nil, errs.BuildError(err, "could not get media by tag id from repo: %v", id)
	}

	return media, nil
}

// Upsert implements TagService.
func (p *tagService) Upsert(name string) (*model.Tag, error) {
	tag, err := p.repo.Tag().GetByName(name)
	if err != nil {
		return nil, errs.BuildError(err, "could not get tag by name from repo")
	}

	if tag == nil {
		tags, err := p.repo.Tag().Create([]string{name})
		if err != nil {
			return nil, errs.BuildError(err, "could not create tag by name")
		}
		tag = &tags[0]
	}

	return tag, nil
}

var tagServiceInstance *tagService

func New(repo repository.Repository, env *environment.EnvironmentVariables) TagService {
	if tagServiceInstance == nil {
		tagServiceInstance = &tagService{
			env:    env,
			repo:   repo,
			logger: logger.New(env),
		}

		tagServiceInstance.logger.Info("TagService instance created")
	}

	return tagServiceInstance
}
