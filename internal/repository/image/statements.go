package imageRepository

import (
	"github.com/go-jet/jet/v2/postgres"
	"github.com/slugger7/exorcist/internal/db/exorcist/public/model"
	"github.com/slugger7/exorcist/internal/db/exorcist/public/table"
	"github.com/slugger7/exorcist/internal/repository/util"
)

type ImageStatment struct {
	postgres.Statement
	*ImageRepository
}

func (is *ImageStatment) Query(destination interface{}) error {
	util.DebugCheck(is.ImageRepository.Env, is.Statement)
	return is.Statement.QueryContext(is.ctx, is.db, destination)
}

func (ir *ImageRepository) createStatement(m *model.Image) *ImageStatment {
	statement := table.Image.INSERT(table.Image.Name, table.Image.Path).
		MODEL(m).
		RETURNING(table.Image.AllColumns)

	return &ImageStatment{statement, ir}
}

func (ir *ImageRepository) relateVideoStatement(m *model.VideoImage) *ImageStatment {
	statement := table.VideoImage.INSERT(table.VideoImage.VideoID, table.VideoImage.ImageID, table.VideoImage.VideoImageType).
		MODEL(m).
		RETURNING(table.VideoImage.AllColumns)

	return &ImageStatment{statement, ir}
}
