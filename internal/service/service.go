package service

import (
	"context"

	"github.com/slugger7/exorcist/internal/environment"
	"github.com/slugger7/exorcist/internal/logger"
	"github.com/slugger7/exorcist/internal/repository"
	jobService "github.com/slugger7/exorcist/internal/service/job"
	libraryService "github.com/slugger7/exorcist/internal/service/library"
	libraryPathService "github.com/slugger7/exorcist/internal/service/library_path"
	mediaService "github.com/slugger7/exorcist/internal/service/media"
	personService "github.com/slugger7/exorcist/internal/service/person"
	playlistService "github.com/slugger7/exorcist/internal/service/playlist"
	tagService "github.com/slugger7/exorcist/internal/service/tag"
	userService "github.com/slugger7/exorcist/internal/service/user"
)

type Service interface {
	User() userService.UserService
	Library() libraryService.LibraryService
	LibraryPath() libraryPathService.LibraryPathService
	Job() jobService.JobService
	Person() personService.PersonService
	Tag() tagService.TagService
	Media() mediaService.MediaService
	Playlist() playlistService.PlaylistService
}

type service struct {
	env         *environment.EnvironmentVariables
	logger      logger.Logger
	user        userService.UserService
	library     libraryService.LibraryService
	libraryPath libraryPathService.LibraryPathService
	job         jobService.JobService
	person      personService.PersonService
	tag         tagService.TagService
	media       mediaService.MediaService
	playlist    playlistService.PlaylistService
	ctx         context.Context
}

var serviceInstance *service

func New(repo repository.Repository, env *environment.EnvironmentVariables, jobCh chan bool, ctx context.Context) Service {
	if serviceInstance == nil {
		personService := personService.New(repo, env)
		tagService := tagService.New(repo, env)
		serviceInstance = &service{
			env:         env,
			logger:      logger.New(env),
			user:        userService.New(repo, env),
			library:     libraryService.New(repo, env),
			libraryPath: libraryPathService.New(repo, env),
			job:         jobService.New(repo, env, jobCh, ctx),
			person:      personService,
			tag:         tagService,
			media:       mediaService.New(env, repo, personService, tagService),
			playlist:    playlistService.New(env, repo),
			ctx:         ctx,
		}

		serviceInstance.logger.Info("Service instance created")
	}
	return serviceInstance
}

func (s *service) User() userService.UserService {
	s.logger.Debug("Getting UserService")
	return s.user
}

func (s *service) Library() libraryService.LibraryService {
	s.logger.Debug("Getting LibraryService")
	return s.library
}

func (s *service) LibraryPath() libraryPathService.LibraryPathService {
	s.logger.Debug("Getting LibraryPathService")
	return s.libraryPath
}

func (s *service) Job() jobService.JobService {
	s.logger.Debug("Getting jobService")
	return s.job
}

func (s *service) Person() personService.PersonService {
	s.logger.Debug("Getting personService")
	return s.person
}

func (s *service) Media() mediaService.MediaService {
	s.logger.Debug("Getting mediaService")
	return s.media
}

func (s *service) Tag() tagService.TagService {
	s.logger.Debug("Getting tagService")
	return s.tag
}

func (s *service) Playlist() playlistService.PlaylistService {
	s.logger.Debug("Getting playlistService")
	return s.playlist
}
