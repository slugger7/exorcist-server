package mediaService

import (
	"errors"
	"fmt"
	"slices"
	"strings"

	"github.com/google/uuid"
	"github.com/slugger7/exorcist/internal/db/exorcist/public/model"
	"github.com/slugger7/exorcist/internal/environment"
	errs "github.com/slugger7/exorcist/internal/errors"
	"github.com/slugger7/exorcist/internal/logger"
	"github.com/slugger7/exorcist/internal/models"
	"github.com/slugger7/exorcist/internal/repository"
	personService "github.com/slugger7/exorcist/internal/service/person"
	tagService "github.com/slugger7/exorcist/internal/service/tag"
)

type MediaService interface {
	SetPeople(id uuid.UUID, people []string) (*models.Media, error)
	SetTags(id uuid.UUID, tags []string) (*models.Media, error)
}

type mediaService struct {
	env           *environment.EnvironmentVariables
	repo          repository.IRepository
	logger        logger.ILogger
	personService personService.IPersonService
	tagService    tagService.TagService
}

// SetTags implements MediaService.
func (m *mediaService) SetTags(id uuid.UUID, tags []string) (*models.Media, error) {
	mediaModel, err := m.repo.Media().GetById(id)
	if err != nil {
		return nil, errs.BuildError(err, "could not get media by id: %v", id.String())
	}

	if mediaModel == nil {
		return nil, fmt.Errorf("could not find media by id")
	}

	if len(tags) == 0 {
		return mediaModel, nil
	}

	uniqueTags := []string{tags[0]}
	for _, p := range tags {
		if !slices.ContainsFunc(uniqueTags, lowerStringComparator(p)) {
			uniqueTags = append(uniqueTags, p)
		}
	}
	tags = uniqueTags

	leftovers := []model.Tag{}
	for _, t := range mediaModel.Tags {
		if !slices.ContainsFunc(tags, lowerStringComparator(t.Name)) {
			m.repo.Tag().RemoveFromMedia(model.MediaTag{MediaID: id, TagID: t.ID})
		} else {
			leftovers = append(leftovers, t)
		}
	}
	mediaModel.Tags = leftovers

	tagModels := []model.Tag{}
	errorList := []error{}
	for _, p := range tags {
		tagModel, err := m.tagService.Upsert(p)
		if err != nil {
			errorList = append(errorList, err)
			continue
		}
		tagModels = append(tagModels, *tagModel)
	}

	if len(errorList) != 0 {
		joinedErrors := errors.Join(errorList...)
		m.logger.Errorf("error while upserting tags: %v", errs.BuildError(joinedErrors, "could not upsert some tags").Error())
	}

	mediaTagModels := []model.MediaTag{}
	for _, t := range tagModels {
		if slices.ContainsFunc(mediaModel.Tags, func(tag model.Tag) bool {
			return tag.ID == t.ID
		}) {
			continue
		}

		mediaTagModels = append(mediaTagModels, model.MediaTag{
			MediaID: id,
			TagID:   t.ID,
		})
	}

	mediaTagModels, err = m.repo.Tag().AddToMedia(mediaTagModels)
	if err != nil {
		return nil, errs.BuildError(err, "error linking media with tag")
	}

	mediaModel, err = m.repo.Media().GetById(id)
	if err != nil {
		return nil, errs.BuildError(err, "could not get media by id after adding tags: %v", id)
	}

	return mediaModel, nil
}

func lowerStringComparator(a string) func(string) bool {
	return func(b string) bool {
		return strings.ToLower(a) == strings.ToLower(b)
	}
}

// SetPeople implements MediaService.
func (m *mediaService) SetPeople(id uuid.UUID, people []string) (*models.Media, error) {
	mediaModel, err := m.repo.Media().GetById(id)
	if err != nil {
		return nil, errs.BuildError(err, "could not get media by id: %v", id.String())
	}

	if mediaModel == nil {
		return nil, fmt.Errorf("could not find media by id")
	}

	if len(people) == 0 {
		return mediaModel, nil
	}

	uniquePeople := []string{people[0]}
	for _, p := range people {
		if !slices.ContainsFunc(uniquePeople, lowerStringComparator(p)) {
			uniquePeople = append(uniquePeople, p)
		}
	}
	people = uniquePeople

	leftovers := []model.Person{}
	for _, p := range mediaModel.People {
		if !slices.ContainsFunc(people, lowerStringComparator(p.Name)) {
			m.repo.Person().RemoveFromMedia(model.MediaPerson{MediaID: id, PersonID: p.ID})
		} else {
			leftovers = append(leftovers, p)
		}
	}
	mediaModel.People = leftovers

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

	mediaPersonModels := []model.MediaPerson{}
	for _, p := range peopleModels {
		if slices.ContainsFunc(mediaModel.People, func(person model.Person) bool {
			return person.ID == p.ID
		}) {
			continue
		}

		mediaPersonModels = append(mediaPersonModels, model.MediaPerson{
			MediaID:  id,
			PersonID: p.ID,
		})
	}

	mediaPersonModels, err = m.repo.Person().AddToMedia(mediaPersonModels)
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
