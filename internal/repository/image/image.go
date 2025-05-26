package imageRepository

import (
	"context"
	"database/sql"

	"github.com/go-jet/jet/v2/postgres"
	"github.com/google/uuid"
	"github.com/slugger7/exorcist/internal/db/exorcist/public/model"
	"github.com/slugger7/exorcist/internal/db/exorcist/public/table"
	"github.com/slugger7/exorcist/internal/environment"
	errs "github.com/slugger7/exorcist/internal/errors"
	"github.com/slugger7/exorcist/internal/repository/util"
)

type MediaImage struct {
	model.Image
	model.Media
}

type IImageRepository interface {
	Create(m *model.Image) (*model.Image, error)
	GetById(uuid.UUID) (*MediaImage, error)
	GetByMediaId(uuid.UUID) (*MediaImage, error)
}

type ImageRepository struct {
	db  *sql.DB
	Env *environment.EnvironmentVariables
	ctx context.Context
}

// GetByMediaId implements IImageRepository.
func (i *ImageRepository) GetByMediaId(id uuid.UUID) (*MediaImage, error) {
	media := table.Media
	image := table.Image

	statement := image.SELECT(image.AllColumns, media.AllColumns).
		FROM(image.INNER_JOIN(
			media,
			image.MediaID.EQ(media.ID),
		)).
		WHERE(media.ID.EQ(postgres.UUID(id))).
		LIMIT(1)

	util.DebugCheck(i.Env, statement)

	var result MediaImage
	if err := statement.QueryContext(i.ctx, i.db, &result); err != nil {
		return nil, errs.BuildError(err, "could not get image by media id: %v", id)
	}

	return &result, nil
}

var imageRepoInstance *ImageRepository

func New(db *sql.DB, env *environment.EnvironmentVariables, context context.Context) IImageRepository {
	if imageRepoInstance != nil {
		return imageRepoInstance
	}

	imageRepoInstance = &ImageRepository{
		db:  db,
		Env: env,
		ctx: context,
	}

	return imageRepoInstance
}

func (i *ImageRepository) Create(m *model.Image) (*model.Image, error) {
	var img model.Image
	statement := table.Image.INSERT(table.Image.AllColumns).
		MODEL(m).
		RETURNING(table.Image.AllColumns)

	util.DebugCheck(i.Env, statement)

	if err := statement.QueryContext(i.ctx, i.db, &img); err != nil {
		return nil, errs.BuildError(err, "error creating image")
	}

	return &img, nil
}

func (i *ImageRepository) GetById(id uuid.UUID) (*MediaImage, error) {
	var img MediaImage
	statement := table.Image.SELECT(table.Image.AllColumns, table.Media.AllColumns).
		FROM(table.Image.INNER_JOIN(
			table.Image,
			table.Image.MediaID.EQ(table.Media.ID),
		)).
		WHERE(table.Image.ID.EQ(postgres.UUID(id))).
		LIMIT(1)

	if err := statement.Query(i.db, &img); err != nil {
		return nil, errs.BuildError(err, "error getting image by id")
	}

	return &img, nil
}
