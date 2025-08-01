package libraryRepository

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

type LibraryRepository interface {
	Create(name string) (*model.Library, error)
	GetByName(name string) (*model.Library, error)
	GetAll() ([]model.Library, error)
	GetById(uuid.UUID) (*model.Library, error)
	GetMedia(id, userId uuid.UUID, search dto.MediaSearchDTO) (*dto.PageDTO[models.MediaOverviewModel], error)
	Update(m model.Library) (*model.Library, error)
}

type libraryRepository struct {
	db  *sql.DB
	env *environment.EnvironmentVariables
	ctx context.Context
}

// Update implements LibraryRepository.
func (ls *libraryRepository) Update(m model.Library) (*model.Library, error) {
	m.Modified = time.Now()
	statement := table.Library.UPDATE(table.Library.Modified, table.Library.Name).
		MODEL(m).
		WHERE(table.Library.ID.EQ(postgres.UUID(m.ID))).
		RETURNING(table.Library.AllColumns)

	util.DebugCheck(ls.env, statement)

	var updatedModel model.Library
	if err := statement.QueryContext(ls.ctx, ls.db, &updatedModel); err != nil {
		return nil, errs.BuildError(err, "could not update library name")
	}

	return &updatedModel, nil
}

// GetMedia implements LibraryRepository.
func (ls *libraryRepository) GetMedia(id, userId uuid.UUID, search dto.MediaSearchDTO) (*dto.PageDTO[models.MediaOverviewModel], error) {
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

	mediaPage, err := helpers.QueryMediaOverview(userId, search, relationFn, whereFn, ls.ctx, ls.db, ls.env)
	if err != nil {
		return nil, errs.BuildError(err, "colud not query media overview from media repo")
	}

	return mediaPage, nil
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
