package libraryService

import (
	"errors"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/slugger7/exorcist/internal/db/exorcist/public/model"
	errs "github.com/slugger7/exorcist/internal/errors"
	"github.com/slugger7/exorcist/internal/mocks/mrepository"
)

func setup() (*LibraryService, *mrepository.MockLibraryRepo) {
	mockRepo := mrepository.SetupMockRespository()
	ls := &LibraryService{repo: mockRepo}
	return ls, mockRepo.MockLibraryRepo
}

func Test_CreateLibrary_ProduceErrorWhileFetchingExistingLibraries(t *testing.T) {
	ls, mlr := setup()

	mlr.MockError[0] = errors.New("error")
	lib := model.Library{Name: "expected library"}

	expectedErrMsg := fmt.Sprintf(ErrLibraryByName, lib.Name)
	newLib, err := ls.Create(lib)
	if err != nil {
		var e errs.IError
		if errors.As(err, &e) {
			if e.Message() != expectedErrMsg {
				t.Errorf("Expected: %v\nGot: %v", expectedErrMsg, e.Message())
			}
		} else {
			t.Errorf("Expected a different error: %v", err)
		}
	}

	if newLib != nil {
		t.Fatal("Error was supposed to be thrown but new lib had a value")
	}
}

func Test_CreateLibrary_WithExistingLibrary_ShouldThrowError(t *testing.T) {
	ls, mlr := setup()
	expectedId, _ := uuid.NewRandom()
	mlr.MockModel[0] = nil
	mlr.MockModel[1] = &model.Library{ID: expectedId}

	lib := model.Library{}
	library, err := ls.Create(lib)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(library)
}

func Test_GetLibraries_RepoReturnsError_ShouldReturnError(t *testing.T) {
	ls, mlr := setup()

	mlr.MockError[0] = errors.New("error")
	libs, err := ls.GetAll()
	if err != nil {
		var e errs.IError
		if errors.As(err, &e) {
			if e.Message() != ErrGetLibraries {
				t.Errorf("Expected error: %v\nGot error: %v", ErrGetLibraries, e.Message())
			}
		} else {
			t.Errorf("Expected a specific error but got: %v", err)
		}
	}

	if libs != nil {
		t.Fatal("Expected an error but libs was defined")
	}
}

func Test_GetLibraries_ReturnsLibraries(t *testing.T) {
	ls, mlr := setup()
	expectedName := "expected library name"
	mlr.MockModels[0] = []model.Library{{Name: expectedName}}
	actual, err := ls.repo.Library().GetLibraries()
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
