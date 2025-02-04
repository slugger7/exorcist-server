package libraryService

import (
	"errors"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/slugger7/exorcist/internal/db/exorcist/public/model"
	jobRepository "github.com/slugger7/exorcist/internal/repository/job"
	libraryRepository "github.com/slugger7/exorcist/internal/repository/library"
	libraryPathRepository "github.com/slugger7/exorcist/internal/repository/library_path"
	userRepository "github.com/slugger7/exorcist/internal/repository/user"
	videoRepository "github.com/slugger7/exorcist/internal/repository/video"
)

type mockRepo struct {
	mockLibraryRepo
}

func (mr mockRepo) LibraryRepo() libraryRepository.ILibraryRepository {
	return mr.mockLibraryRepo
}

var count = 0

type mockLibraryRepo struct {
	mockModels map[int]*model.Library
	mockError  error
}

func (mlr mockLibraryRepo) CreateLibrary(name string) (*model.Library, error) {
	if len(mlr.mockModels) > count {
		return mlr.mockModels[count], mlr.mockError
	}
	return nil, mlr.mockError
}
func (mlr mockLibraryRepo) GetLibraryByName(name string) (*model.Library, error) {
	if len(mlr.mockModels) > count {
		return mlr.mockModels[count], mlr.mockError
	}
	return nil, mlr.mockError
}

func Test_CreateLibrary_ProduceErrorWhileFetchingExistingLibraries(t *testing.T) {
	expectedErr := errors.New("expected error")
	ls := &LibraryService{repo: mockRepo{mockLibraryRepo{mockError: expectedErr}}}
	lib := model.Library{}

	if _, err := ls.CreateLibrary(lib); err.Error() != expectedErr.Error() {
		t.Errorf("Encountered an unexpected error creating library\nExpected: %v\nGot: %v", expectedErr.Error(), err.Error())
	}
}

func Test_CreateLibrary_WithExistingLibrary_ShouldThrowError(t *testing.T) {
	expectedId, _ := uuid.NewRandom()
	var mockModels = make(map[int]*model.Library)
	mockModels[0] = nil
	mockModels[1] = &model.Library{ID: expectedId}
	ls := &LibraryService{repo: mockRepo{mockLibraryRepo{mockModels: mockModels}}}

	lib := model.Library{}
	library, err := ls.CreateLibrary(lib)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(library)
}

func (mr mockRepo) Health() map[string]string {
	panic("not implemented")
}
func (mr mockRepo) Close() error {
	panic("not implemented")
}
func (mr mockRepo) JobRepo() jobRepository.IJobRepository {
	panic("not implemented")
}
func (mr mockRepo) LibraryPathRepo() libraryPathRepository.ILibraryPathRepository {
	panic("not implemented")
}
func (mr mockRepo) VideoRepo() videoRepository.IVideoRepository {
	panic("not implemented")
}
func (mr mockRepo) UserRepo() userRepository.IUserRepository {
	panic("not implemented")
}
