package libraryRepository

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/slugger7/exorcist/internal/db/exorcist/public/model"
	"github.com/slugger7/exorcist/internal/environment"
)

type ILibraryRepository interface {
	CreateLibrary(name string) (*model.Library, error)
	GetLibraryByName(name string) (*model.Library, error)
	GetLibraries() ([]model.Library, error)
}

type LibraryRepository struct {
	db  *sql.DB
	Env *environment.EnvironmentVariables
}

var libraryRepoInstance *LibraryRepository

func New(db *sql.DB, env *environment.EnvironmentVariables) ILibraryRepository {
	if libraryRepoInstance != nil {
		return libraryRepoInstance
	}
	libraryRepoInstance = &LibraryRepository{
		db:  db,
		Env: env,
	}
	return libraryRepoInstance
}

func (ls *LibraryStatement) Query(destination interface{}) error {
	return ls.Statement.Query(ls.db, destination)
}

func (ls *LibraryStatement) Sql() string {
	sql, _ := ls.Statement.Sql()
	return sql
}

func (ls *LibraryRepository) CreateLibrary(name string) (*model.Library, error) {
	var library struct{ model.Library }
	if err := ls.createLibraryStatement(name).Query(&library); err != nil {
		log.Println("something went wrong creating the library")
		return nil, err
	}
	return &library.Library, nil
}

func (ls *LibraryRepository) GetLibraryByName(name string) (*model.Library, error) {
	var libraries []struct{ model.Library }
	if err := ls.getLibraryByNameStatement(name).Query(&libraries); err != nil {
		return nil, errors.Join(errors.New(fmt.Sprintf("something went wrong getting the library by name: %v", err)), err)
	}
	var library *model.Library
	if len(libraries) > 0 {
		library = &libraries[len(libraries)-1].Library
	}
	return library, nil
}

func (ls *LibraryRepository) GetLibraries() ([]model.Library, error) {
	var libraries []struct{ model.Library }
	if err := ls.getLibrariesStatement().Query(&libraries); err != nil {
		return nil, errors.Join(errors.New("error getting libraries from database"), err)
	}
	var libs []model.Library
	for _, v := range libraries {
		libs = append(libs, v.Library)
	}

	return libs, nil
}
