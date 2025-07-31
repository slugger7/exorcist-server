package playlistRepository

import (
	"context"
	"database/sql"
	"time"

	"github.com/go-jet/jet/v2/postgres"
	"github.com/google/uuid"
	"github.com/slugger7/exorcist/internal/db/exorcist/public/model"
	"github.com/slugger7/exorcist/internal/db/exorcist/public/table"
	"github.com/slugger7/exorcist/internal/dto"
	"github.com/slugger7/exorcist/internal/environment"
	errs "github.com/slugger7/exorcist/internal/errors"
	"github.com/slugger7/exorcist/internal/models"
	"github.com/slugger7/exorcist/internal/repository/helpers"
	"github.com/slugger7/exorcist/internal/repository/util"
)

type PlaylistRepository interface {
	GetById(id uuid.UUID) (*model.Playlist, error)
	GetAll() ([]model.Playlist, error)
	GetMedia(id, userId uuid.UUID, search dto.MediaSearchDTO) (*dto.PageDTO[models.MediaOverviewModel], error)
	CreateAll(playlists []model.Playlist) ([]model.Playlist, error)
	AddMedia(playlistMedia []model.PlaylistMedia) ([]model.PlaylistMedia, error)
	Update(m model.Playlist) (*model.Playlist, error)
}

type playlistRepository struct {
	env *environment.EnvironmentVariables
	db  *sql.DB
	ctx context.Context
}

// Update implements PlaylistRepository.
func (p *playlistRepository) Update(m model.Playlist) (*model.Playlist, error) {
	m.Modified = time.Now()

	statement := table.Playlist.UPDATE(table.Playlist.Modified, table.Playlist.Name).
		MODEL(m).
		WHERE(table.Playlist.ID.EQ(postgres.UUID(m.ID))).
		RETURNING(table.Playlist.AllColumns)

	util.DebugCheck(p.env, statement)

	var updatedModel model.Playlist
	if err := statement.QueryContext(p.ctx, p.db, &updatedModel); err != nil {
		return nil, errs.BuildError(err, "could not update playlist")
	}

	return &updatedModel, nil
}

// AddMedia implements PlaylistRepository.
func (p *playlistRepository) AddMedia(playlistMedia []model.PlaylistMedia) ([]model.PlaylistMedia, error) {
	if len(playlistMedia) == 0 {
		return nil, nil
	}

	statement := table.PlaylistMedia.INSERT(table.PlaylistMedia.PlaylistID, table.PlaylistMedia.MediaID).
		MODELS(playlistMedia).
		RETURNING(table.PlaylistMedia.AllColumns)

	util.DebugCheck(p.env, statement)

	var playlistMediaEntities []model.PlaylistMedia
	if err := statement.QueryContext(p.ctx, p.db, &playlistMediaEntities); err != nil {
		return nil, errs.BuildError(err, "could not create playlist media")
	}

	return playlistMediaEntities, nil
}

// GetMedia implements PlaylistRepository.
func (p *playlistRepository) GetMedia(id uuid.UUID, userId uuid.UUID, search dto.MediaSearchDTO) (*dto.PageDTO[models.MediaOverviewModel], error) {
	relationFn := func(rel postgres.ReadableTable) postgres.ReadableTable {
		return rel.INNER_JOIN(
			table.PlaylistMedia,
			table.PlaylistMedia.MediaID.EQ(table.Media.ID),
		)
	}

	wherFn := func(whr postgres.BoolExpression) postgres.BoolExpression {
		return whr.AND(table.PlaylistMedia.PlaylistID.EQ(postgres.UUID(id)))
	}

	return helpers.QueryMediaOverview(userId, search, relationFn, wherFn, p.ctx, p.db, p.env)
}

// GetById implements PlaylistRepository.
func (p *playlistRepository) GetById(id uuid.UUID) (*model.Playlist, error) {
	statement := table.Playlist.SELECT(table.Playlist.AllColumns).
		WHERE(table.Playlist.ID.EQ(postgres.UUID(id)))

	util.DebugCheck(p.env, statement)

	var playlists []model.Playlist
	if err := statement.QueryContext(p.ctx, p.db, &playlists); err != nil {
		return nil, errs.BuildError(err, "could not query playlists by id: %v", id.String())
	}

	if len(playlists) == 0 {
		return nil, nil
	}

	return &playlists[0], nil
}

// CreateAll implements PlaylistRepository.
func (p *playlistRepository) CreateAll(playlists []model.Playlist) ([]model.Playlist, error) {
	statement := table.Playlist.INSERT(table.Playlist.UserID, table.Playlist.Name).
		MODELS(playlists).
		RETURNING(table.Playlist.AllColumns)

	util.DebugCheck(p.env, statement)

	var playlistEntities []model.Playlist
	if err := statement.QueryContext(p.ctx, p.db, &playlistEntities); err != nil {
		return nil, errs.BuildError(err, "could not create and return playlists")
	}

	return playlistEntities, nil
}

// GetAll implements PlaylistRepository.
func (p *playlistRepository) GetAll() ([]model.Playlist, error) {
	statement := table.Playlist.SELECT(table.Playlist.AllColumns)

	util.DebugCheck(p.env, statement)

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
