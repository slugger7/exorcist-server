package mediaService

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/slugger7/exorcist/internal/db/exorcist/public/model"
	"github.com/slugger7/exorcist/internal/environment"
	errs "github.com/slugger7/exorcist/internal/errors"
	"github.com/slugger7/exorcist/internal/logger"
	"github.com/slugger7/exorcist/internal/repository"
	personService "github.com/slugger7/exorcist/internal/service/person"
	tagService "github.com/slugger7/exorcist/internal/service/tag"
)

type MediaService interface {
	AddTag(id uuid.UUID, tagId uuid.UUID) (*model.MediaTag, error)
	AddPerson(id uuid.UUID, personId uuid.UUID) (*model.MediaPerson, error)
	Delete(id uuid.UUID, physical bool) error
}

type mediaService struct {
	env           *environment.EnvironmentVariables
	repo          repository.IRepository
	logger        logger.ILogger
	personService personService.IPersonService
	tagService    tagService.TagService
}

// Delete implements MediaService.
func (m *mediaService) Delete(id uuid.UUID, physical bool) error {
	panic("unimplemented")
}

// AddPerson implements MediaService.
func (m *mediaService) AddPerson(id uuid.UUID, personId uuid.UUID) (*model.MediaPerson, error) {
	mediaModel, err := m.repo.Media().GetById(id)
	if err != nil {
		return nil, errs.BuildError(err, "could not get media by id from repo: %v", id.String())
	}

	if mediaModel == nil {
		return nil, fmt.Errorf("could not find media by id: %v", id)
	}

	personModel, err := m.repo.Person().GetById(personId)
	if err != nil {
		return nil, errs.BuildError(err, "could not get person by id from repo: %v", personId.String())
	}

	if personModel == nil {
		return nil, fmt.Errorf("colud not find person by id: %v", personId)
	}

	for _, p := range mediaModel.People {
		if p.ID == personId {
			return &model.MediaPerson{MediaID: id, PersonID: personId}, nil
		}
	}

	mediaPeopleModels := []model.MediaPerson{
		{
			PersonID: personId,
			MediaID:  id,
		},
	}
	createdMediaModelPeople, err := m.repo.Person().AddToMedia(mediaPeopleModels)
	if err != nil {
		return nil, errs.BuildError(err, "could not add person (%v) to media (%v)", personId, id)
	}

	return &createdMediaModelPeople[0], nil
}

// AddTag implements MediaService.
func (m *mediaService) AddTag(id uuid.UUID, tagId uuid.UUID) (*model.MediaTag, error) {
	mediaModel, err := m.repo.Media().GetById(id)
	if err != nil {
		return nil, errs.BuildError(err, "could not get media by id from repo: %v", id.String())
	}

	if mediaModel == nil {
		return nil, fmt.Errorf("could not find media by id: %v", id)
	}

	tagModel, err := m.repo.Tag().GetById(tagId)
	if err != nil {
		return nil, errs.BuildError(err, "could not get tag by id from repo: %v", id.String())
	}

	if tagModel == nil {
		return nil, fmt.Errorf("could not find tag by id: %v", tagId)
	}

	for _, t := range mediaModel.Tags {
		if t.ID == tagId {
			return &model.MediaTag{MediaID: id, TagID: tagId}, nil
		}
	}

	mediaTagModels := []model.MediaTag{
		{
			TagID:   tagId,
			MediaID: id,
		},
	}
	createdMediaModelTags, err := m.repo.Tag().AddToMedia(mediaTagModels)
	if err != nil {
		return nil, errs.BuildError(err, "could not add tag (%v) to media (%v)", tagId, id)
	}

	return &createdMediaModelTags[0], nil
}

func lowerStringComparator(a string) func(string) bool {
	return func(b string) bool {
		return strings.ToLower(a) == strings.ToLower(b)
	}
}

var mediaServiceInstance *mediaService

func New(env *environment.EnvironmentVariables, repo repository.IRepository, personService personService.IPersonService, tagService tagService.TagService) MediaService {
	if mediaServiceInstance == nil {
		mediaServiceInstance = &mediaService{
			env:           env,
			repo:          repo,
			logger:        logger.New(env),
			personService: personService,
			tagService:    tagService,
		}

		mediaServiceInstance.logger.Info("Created media service instance")
	}

	return mediaServiceInstance
}
