package libraryRepository

import (
	"context"
	"database/sql"

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

type LibraryRepository interface {
	Create(name string) (*model.Library, error)
	GetByName(name string) (*model.Library, error)
	GetAll() ([]model.Library, error)
	GetById(uuid.UUID) (*model.Library, error)
	GetMedia(id uuid.UUID, search dto.MediaSearchDTO) (*dto.PageDTO[models.MediaOverviewModel], error)
}

type libraryRepository struct {
	db  *sql.DB
	env *environment.EnvironmentVariables
	ctx context.Context
}

// GetMedia implements LibraryRepository.
func (ls *libraryRepository) GetMedia(id uuid.UUID, search dto.MediaSearchDTO) (*dto.PageDTO[models.MediaOverviewModel], error) {
	relationFn := func(relationTable postgres.ReadableTable) postgres.ReadableTable {
		return relationTable.INNER_JOIN(
			table.LibraryPath,
			table.Media.LibraryPathID.EQ(table.LibraryPath.ID),
		)
	}

	whereFn := func(whr postgres.BoolExpression) postgres.BoolExpression {
		return whr.
			AND(table.LibraryPath.LibraryID.EQ(postgres.UUID(id)))
	}

	selectStatement, countStatement := helpers.MediaOverviewStatement(search, relationFn, whereFn)

	util.DebugCheck(ls.env, selectStatement)
	util.DebugCheck(ls.env, countStatement)

	var total struct {
		Total int
	}
	if err := countStatement.QueryContext(ls.ctx, ls.db, &total); err != nil {
		return nil, errs.BuildError(err, "could not query media total by library id: %v", id.String())
	}

	var mediaResult []models.MediaOverviewModel
	if err := selectStatement.QueryContext(ls.ctx, ls.db, &mediaResult); err != nil {
		return nil, errs.BuildError(err, "could not query media by library id: %v", id.String())
	}

	return &dto.PageDTO[models.MediaOverviewModel]{
		Data:  mediaResult,
		Limit: search.Limit,
		Skip:  search.Skip,
		Total: total.Total,
	}, nil
}

var libraryRepoInstance *libraryRepository

func New(db *sql.DB, env *environment.EnvironmentVariables, context context.Context) LibraryRepository {
	if libraryRepoInstance != nil {
		return libraryRepoInstance
	}
	libraryRepoInstance = &libraryRepository{
		db:  db,
		env: env,
		ctx: context,
	}
	return libraryRepoInstance
}

func (ls *LibraryStatement) Query(destination interface{}) error {
	return ls.Statement.QueryContext(ls.ctx, ls.db, destination)
}

func (ls *LibraryStatement) Sql() string {
	sql, _ := ls.Statement.Sql()
	return sql
}

func (ls *libraryRepository) Create(name string) (*model.Library, error) {
	var library struct{ model.Library }
	if err := ls.createLibraryStatement(name).Query(&library); err != nil {
		return nil, errs.BuildError(err, "error while creating library")
	}
	return &library.Library, nil
}

func (ls *libraryRepository) GetByName(name string) (*model.Library, error) {
	var libraries []struct{ model.Library }
	if err := ls.getLibraryByNameStatement(name).Query(&libraries); err != nil {
		return nil, errs.BuildError(err, "could not get library by name '%v'", name)
	}
	var library *model.Library
	if len(libraries) > 0 {
		library = &libraries[len(libraries)-1].Library
	}
	return library, nil
}

func (ls *libraryRepository) GetAll() ([]model.Library, error) {
	var libraries []struct{ model.Library }
	if err := ls.getLibrariesStatement().Query(&libraries); err != nil {
		return nil, errs.BuildError(err, "could not get libraries")
	}
	var libs []model.Library
	for _, v := range libraries {
		libs = append(libs, v.Library)
	}

	return libs, nil
}

func (ls *libraryRepository) GetById(id uuid.UUID) (*model.Library, error) {
	var library struct{ model.Library }
	if err := ls.getById(id).Query(&library); err != nil {
		return nil, errs.BuildError(err, "could not get library by id: %v", id)
	}
	return &library.Library, nil
}
