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

type mockLibraryRepo struct {
	mockLibraryStatement
}

type mockLibraryStatement struct {
	mockError  error
	mockModels map[int][]struct{ model.Library }
}

func (mlr mockLibraryRepo) CreateLibraryStatement(name string) libraryRepository.ILibraryStatement {
	return mlr.mockLibraryStatement
}
func (mlr mockLibraryRepo) GetLibraryByName(name string) libraryRepository.ILibraryStatement {
	return mlr.mockLibraryStatement
}

var count = 0

func (mls mockLibraryStatement) Query(destination interface{}) error {
	if mls.mockModels[count] != nil {
		typedDest := *destination.(*[]struct{ model.Library })
		mockedReturn := mls.mockModels[count]
		copy(typedDest, mockedReturn)
	}
	count = count + 1
	return mls.mockError
}
func (mls mockLibraryStatement) Sql() string {
	panic("not implemented")
}

func Test_CreateLibrary_ProduceErrorWhileFetchingExistingLibraries(t *testing.T) {
	expectedErr := errors.New("expected error")
	ls := &LibraryService{repo: mockRepo{mockLibraryRepo{mockLibraryStatement{mockError: expectedErr}}}}
	lib := model.Library{}

	if _, err := ls.CreateLibrary(lib); err.Error() != expectedErr.Error() {
		t.Errorf("Encountered an unexpected error creating library %v", expectedErr.Error())
	}
}

func Test_CreateLibrary_WithExistingLibrary_ShouldThrowError(t *testing.T) {
	count = 0
	expectedId, _ := uuid.NewRandom()
	var mockModels = make(map[int][]struct{ model.Library })
	mockModels[0] = []struct{ model.Library }{}
	mockModels[1] = []struct{ model.Library }{{Library: model.Library{ID: expectedId}}}
	ls := &LibraryService{repo: mockRepo{mockLibraryRepo{mockLibraryStatement{mockModels: mockModels}}}}

	lib := model.Library{}
	library, err := ls.CreateLibrary(lib)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(library)
}
