package playlistRepository

import (
	"context"
	"database/sql"

	"github.com/slugger7/exorcist/internal/db/exorcist/public/model"
	"github.com/slugger7/exorcist/internal/db/exorcist/public/table"
	"github.com/slugger7/exorcist/internal/environment"
	errs "github.com/slugger7/exorcist/internal/errors"
)

type PlaylistRepository interface {
	GetAll() ([]model.Playlist, error)
	CreateAll(playlists []model.Playlist) ([]model.Playlist, error)
}

type playlistRepository struct {
	env *environment.EnvironmentVariables
	db  *sql.DB
	ctx context.Context
}

// CreateAll implements PlaylistRepository.
func (p *playlistRepository) CreateAll(playlists []model.Playlist) ([]model.Playlist, error) {
	statement := table.Playlist.INSERT(table.Playlist.UserID, table.Playlist.Name).
		MODELS(playlists).
		RETURNING(table.Playlist.AllColumns)

	var playlistEntities []model.Playlist
	if err := statement.QueryContext(p.ctx, p.db, &playlistEntities); err != nil {
		return nil, errs.BuildError(err, "could not create and return playlists")
	}

	return playlistEntities, nil
}

// GetAll implements PlaylistRepository.
func (p *playlistRepository) GetAll() ([]model.Playlist, error) {
	statement := table.Playlist.SELECT(table.Playlist.AllColumns)

	var res []model.Playlist
	if err := statement.QueryContext(p.ctx, p.db, &res); err != nil {
		return nil, errs.BuildError(err, "could not query playlists")
	}

	return res, nil
}

var playlistRepositoryInstance *playlistRepository

func New(env *environment.EnvironmentVariables, db *sql.DB, context context.Context) PlaylistRepository {
	if playlistRepositoryInstance != nil {
		return playlistRepositoryInstance
	}

	playlistRepositoryInstance = &playlistRepository{
		env: env,
		db:  db,
		ctx: context,
	}

	return playlistRepositoryInstance
}
