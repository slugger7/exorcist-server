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

type ImageRepository interface {
	Create(m *model.Image) (*model.Image, error)
	GetById(uuid.UUID) (*MediaImage, error)
	GetByMediaId(uuid.UUID) (*MediaImage, error)
}

type imageRepository struct {
	db  *sql.DB
	env *environment.EnvironmentVariables
	ctx context.Context
}

// GetByMediaId implements IImageRepository.
func (i *imageRepository) GetByMediaId(id uuid.UUID) (*MediaImage, error) {
	media := table.Media
	image := table.Image

	statement := image.SELECT(image.AllColumns, media.AllColumns).
		FROM(image.INNER_JOIN(
			media,
			image.MediaID.EQ(media.ID),
		)).
		WHERE(media.ID.EQ(postgres.UUID(id))).
		LIMIT(1)

	util.DebugCheck(i.env, statement)

	var result MediaImage
	if err := statement.QueryContext(i.ctx, i.db, &result); err != nil {
		return nil, errs.BuildError(err, "could not get image by media id: %v", id)
	}

	return &result, nil
}

var imageRepoInstance *imageRepository

func New(db *sql.DB, env *environment.EnvironmentVariables, context context.Context) ImageRepository {
	if imageRepoInstance != nil {
		return imageRepoInstance
	}

	imageRepoInstance = &imageRepository{
		db:  db,
		env: env,
		ctx: context,
	}

	return imageRepoInstance
}

func (i *imageRepository) Create(m *model.Image) (*model.Image, error) {
	image := table.Image
	var img model.Image
	statement := image.INSERT(
		image.MediaID,
		image.Height,
		image.Width,
	).
		MODEL(m).
		RETURNING(image.AllColumns)

	util.DebugCheck(i.env, statement)

	if err := statement.QueryContext(i.ctx, i.db, &img); err != nil {
		return nil, errs.BuildError(err, "error creating image")
	}

	return &img, nil
}

func (i *imageRepository) GetById(id uuid.UUID) (*MediaImage, error) {
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
