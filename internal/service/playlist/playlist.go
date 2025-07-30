package playlistService

import (
	"github.com/google/uuid"
	"github.com/slugger7/exorcist/internal/db/exorcist/public/model"
	"github.com/slugger7/exorcist/internal/dto"
	"github.com/slugger7/exorcist/internal/environment"
	"github.com/slugger7/exorcist/internal/logger"
	"github.com/slugger7/exorcist/internal/repository"
)

type PlaylistService interface {
	CreateAll(userId uuid.UUID, playlists []dto.CreatePlaylistDTO) ([]model.Playlist, error)
}

type playlistService struct {
	env    *environment.EnvironmentVariables
	repo   repository.IRepository
	logger logger.ILogger
}

// CreateAll implements PlaylistService.
func (p *playlistService) CreateAll(userId uuid.UUID, playlists []dto.CreatePlaylistDTO) ([]model.Playlist, error) {
	if playlists == nil || len(playlists) == 0 {
		return nil, nil
	}

	playlistModels := make([]model.Playlist, len(playlists))
	for i, p := range playlists {
		playlistModels[i] = model.Playlist{
			UserID: userId,
			Name:   p.Name,
		}
	}

	return p.repo.Playlist().CreateAll(playlistModels)
}

var playlistServiceInstance *playlistService

func New(env *environment.EnvironmentVariables, repo repository.IRepository) PlaylistService {
	if playlistServiceInstance != nil {
		return playlistServiceInstance
	}

	playlistServiceInstance = &playlistService{
		env:    env,
		repo:   repo,
		logger: logger.New(env),
	}

	playlistServiceInstance.logger.Info("Createde PlaylistService")

	return playlistServiceInstance
}
