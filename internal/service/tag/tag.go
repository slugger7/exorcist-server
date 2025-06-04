package tagService

import (
	"github.com/slugger7/exorcist/internal/db/exorcist/public/model"
	"github.com/slugger7/exorcist/internal/environment"
	errs "github.com/slugger7/exorcist/internal/errors"
	"github.com/slugger7/exorcist/internal/logger"
	"github.com/slugger7/exorcist/internal/repository"
)

type TagService interface {
	Upsert(name string) (*model.Tag, error)
}

type tagService struct {
	env    *environment.EnvironmentVariables
	repo   repository.IRepository
	logger logger.ILogger
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

func New(repo repository.IRepository, env *environment.EnvironmentVariables) TagService {
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
