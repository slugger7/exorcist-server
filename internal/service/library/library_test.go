package libraryService

import (
	"errors"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/slugger7/exorcist/internal/db/exorcist/public/model"
	"github.com/slugger7/exorcist/internal/mocks/mrepository"
)

func beforeEach() (*LibraryService, mrepository.MockLibraryRepo) {
	mockRepo := mrepository.SetupMockRespository()
	ls := &LibraryService{repo: mockRepo}
	return ls, mockRepo.MockLibraryRepo
}

func Test_CreateLibrary_ProduceErrorWhileFetchingExistingLibraries(t *testing.T) {
	ls, mlr := beforeEach()
	expectedErr := errors.New("expected error")
	mlr.MockError[0] = expectedErr
	lib := model.Library{}

	expectedErrorMessage := fmt.Sprintf("github.com/slugger7/exorcist/internal/service/library.(*LibraryService).CreateLibrary: Could not fetch library by name \n%v", expectedErr.Error())
	if _, err := ls.CreateLibrary(lib); err.Error() != expectedErrorMessage {
		t.Errorf("Encountered an unexpected error creating library\nExpected: %v\nGot: %v", expectedErrorMessage, err.Error())
	}
}

func Test_CreateLibrary_WithExistingLibrary_ShouldThrowError(t *testing.T) {
	ls, mlr := beforeEach()
	expectedId, _ := uuid.NewRandom()
	mlr.MockModel[0] = nil
	mlr.MockModel[1] = &model.Library{ID: expectedId}

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
	mlr.MockError[0] = expectedError
	wrappedError := fmt.Sprintf("github.com/slugger7/exorcist/internal/service/library.(*LibraryService).GetLibraries: error getting libraries in repo\n%v", expectedError.Error())
	if _, err := ls.GetLibraries(); err.Error() != wrappedError {
		t.Errorf("Expected: %v\nGot: %v", wrappedError, err.Error())
	}
}

func Test_GetLibraries_ReturnsLibraries(t *testing.T) {
	ls, mlr := beforeEach()
	expectedName := "expected library name"
	mlr.MockModels[0] = []model.Library{{Name: expectedName}}
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
