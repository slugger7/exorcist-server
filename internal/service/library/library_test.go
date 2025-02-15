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

func setup() (*LibraryService, *mrepository.MockRepository) {
	mockRepo := mrepository.SetupMockRespository()
	ls := &LibraryService{repo: mockRepo}
	return ls, mockRepo
}

func Test_CreateLibrary_ProduceErrorWhileFetchingExistingLibraries(t *testing.T) {
	ls, mr := setup()

	mr.MockLibraryRepo.MockError[0] = errors.New("error")
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
	mlr.MockLibraryRepo.MockModel[0] = nil
	mlr.MockLibraryRepo.MockModel[1] = &model.Library{ID: expectedId}

	lib := model.Library{}
	library, err := ls.Create(lib)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(library)
}

func Test_GetLibraries_RepoReturnsError_ShouldReturnError(t *testing.T) {
	ls, mr := setup()

	mr.MockLibraryRepo.MockError[0] = errors.New("error")
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
	mlr.MockLibraryRepo.MockModels[0] = []model.Library{{Name: expectedName}}
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

func Test_Action_InvalidAction(t *testing.T) {
	ls, _ := setup()

	id, _ := uuid.NewRandom()
	action := "non-existent-action"
	err := ls.Action(id, action)
	if err == nil {
		t.Error("Expected and error but was nil")
	}
	expectedError := fmt.Errorf(ErrActionNotFound, action)
	if expectedError.Error() != err.Error() {
		t.Errorf("Expected error: %v\nGot error: %v", expectedError.Error(), err.Error())
	}
}

func Test_Action_RepoErrorWhenGettingLibrary(t *testing.T) {
	ls, mlr := setup()

	id, _ := uuid.NewRandom()
	mlr.MockLibraryPathRepo.MockError[0] = fmt.Errorf("expected error")

	err := ls.Action(id, ActionScan)
	if err == nil {
		t.Error("Expected error but was nil")
	}
	var e errs.IError
	if errors.As(err, &e) {
		expectedErr := fmt.Errorf(ErrFindInRepo, id)
		if e.Message() != expectedErr.Error() {
			t.Errorf("Expected error: %v\nGot error: %v", expectedErr.Error(), e.Message())
		}
	} else {
		t.Errorf("Expected specific error but got: %v", err)
	}
}

func Test_ActionScan_GettingLibPathsReturnsErr(t *testing.T) {
	ls, mlr := setup()

	id, _ := uuid.NewRandom()
	library := &model.Library{ID: id}

	mlr.MockLibraryPathRepo.MockError[0] = fmt.Errorf("error")

	err := ls.actionScan(library)
	if err == nil {
		t.Error("Expected err but was nil")
	}
	var e errs.IError
	if errors.As(err, &e) {
		if e.Message() != ErrActionScanGetLibraryPaths {
			t.Errorf("Expected error: %v\nGot error: %v", ErrActionScanGetLibraryPaths, e.Message())
		}
	} else {
		t.Errorf("Expected specific error but got %v", err)
	}
}

func Test_ActionScan_WithLibraryPath_CreateThrowsError(t *testing.T) {
	ls, mlr := setup()

	id, _ := uuid.NewRandom()
	lib := model.Library{ID: id}
	libPathID, _ := uuid.NewRandom()
	libPath := model.LibraryPath{
		ID:        libPathID,
		LibraryID: id,
		Path:      "some path",
	}

	mlr.MockLibraryPathRepo.MockModels[0] = []model.LibraryPath{libPath}
	mlr.MockJobRepo.MockError[1] = fmt.Errorf("error")

	err := ls.actionScan(&lib)
	if err == nil {
		t.Error("Expected error but was nil")
		return
	}
	var e errs.IError
	if errors.As(err, &e) {
		if e.Message() != ErrCreatingJobs {
			t.Errorf("Expected error: %v\nGot error: %v", ErrCreatingJobs, e.Message())
		}
	} else {
		t.Errorf("Expected specific error but was %v", err)
	}
}

func Test_ActionScann_Success(t *testing.T) {
	ls, mlr := setup()

	id, _ := uuid.NewRandom()
	lib := model.Library{ID: id}
	libPathID, _ := uuid.NewRandom()
	libPath := model.LibraryPath{
		ID:        libPathID,
		LibraryID: id,
		Path:      "some path",
	}

	mlr.MockLibraryPathRepo.MockModels[0] = []model.LibraryPath{libPath}

	err := ls.actionScan(&lib)
	if err != nil {
		t.Errorf("Expected no error but was %v", err)
	}
}
