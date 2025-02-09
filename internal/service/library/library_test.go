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
	mockModels    map[int]*model.Library
	mockErrors    map[int]error
	mockLibraries map[int][]model.Library
}

func (mlr mockLibraryRepo) CreateLibrary(name string) (*model.Library, error) {
	count = count + 1
	return mlr.mockModels[count-1], mlr.mockErrors[count-1]
}
func (mlr mockLibraryRepo) GetLibraryByName(name string) (*model.Library, error) {
	count = count + 1
	return mlr.mockModels[count-1], mlr.mockErrors[count-1]
}
func (mlr mockLibraryRepo) GetLibraries() ([]model.Library, error) {
	count = count + 1
	return mlr.mockLibraries[count-1], mlr.mockErrors[count-1]
}

func beforeEach() (*LibraryService, mockLibraryRepo) {
	count = 0
	mockModels := make(map[int]*model.Library)
	mockErrors := make(map[int]error)
	mockLibraries := make(map[int][]model.Library)
	mlr := mockLibraryRepo{mockModels, mockErrors, mockLibraries}
	ls := &LibraryService{repo: mockRepo{mlr}}
	return ls, mlr
}

func Test_CreateLibrary_ProduceErrorWhileFetchingExistingLibraries(t *testing.T) {
	ls, mlr := beforeEach()
	expectedErr := errors.New("expected error")
	mlr.mockErrors[0] = expectedErr
	lib := model.Library{}

	expectedErrorMessage := fmt.Sprintf("github.com/slugger7/exorcist/internal/service/library.LibraryService.CreateLibrary: Could not fetch library by name \n%v", expectedErr.Error())
	if _, err := ls.CreateLibrary(lib); err.Error() != expectedErrorMessage {
		t.Errorf("Encountered an unexpected error creating library\nExpected: %v\nGot: %v", expectedErrorMessage, err.Error())
	}
}

func Test_CreateLibrary_WithExistingLibrary_ShouldThrowError(t *testing.T) {
	ls, mlr := beforeEach()
	expectedId, _ := uuid.NewRandom()
	mlr.mockModels[0] = nil
	mlr.mockModels[1] = &model.Library{ID: expectedId}

	lib := model.Library{}
	library, err := ls.CreateLibrary(lib)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(library)
}

func Test_GetLibraries_RepoReturnsErro_ShouldReturnError(t *testing.T) {
	ls, mlr := beforeEach()
	expectedError := errors.New("expected error")
	mlr.mockErrors[0] = expectedError
	wrappedError := fmt.Sprintf("github.com/slugger7/exorcist/internal/service/library.LibraryService.GetLibraries: error getting libraries in repo\n%v", expectedError.Error())
	if _, err := ls.GetLibraries(); err.Error() != wrappedError {
		t.Errorf("Expected: %v\nGot: %v", wrappedError, err.Error())
	}
}

func Test_GetLibraries_ReturnsLibraries(t *testing.T) {
	ls, mlr := beforeEach()
	expectedName := "expected library name"
	mlr.mockLibraries[0] = []model.Library{{Name: expectedName}}
	actual, err := ls.repo.LibraryRepo().GetLibraries()
	if err != nil {
		t.Fatal(err)
	}
	if len(actual) != 1 {
		t.Error("Length of libraries did not match expected length of 1")
	}
	if actual[0].Name != expectedName {
		t.Errorf("Expected name: %v\nGot: %v", expectedName, actual[0].Name)
	}
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
