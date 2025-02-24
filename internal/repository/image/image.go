package imageRepository

import (
	"database/sql"

	"github.com/slugger7/exorcist/internal/db/exorcist/public/model"
	"github.com/slugger7/exorcist/internal/environment"
)

type IImageRepository interface {
	Create(m *model.Image) (*model.Image, error)
}

type ImageRepository struct {
	db  *sql.DB
	Env *environment.EnvironmentVariables
}

// Create implements IImageRepository.
func (i *ImageRepository) Create(model *model.Image) (*model.Image, error) {
	panic("unimplemented")
}

var imageRepoInstance *ImageRepository

func New(db *sql.DB, env *environment.EnvironmentVariables) IImageRepository {
	if imageRepoInstance == nil {
		imageRepoInstance = &ImageRepository{
			db:  db,
			Env: env,
		}
	}
	return imageRepoInstance
}
