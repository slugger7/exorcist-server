package libraryPathService

import (
	"errors"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/slugger7/exorcist/internal/db/exorcist/public/model"
	errs "github.com/slugger7/exorcist/internal/errors"
	"github.com/slugger7/exorcist/internal/mocks/mrepository"
)

func setup() (*LibraryPathService, *mrepository.MockRepository) {
	mockRepo := mrepository.SetupMockRespository()
	ls := &LibraryPathService{repo: mockRepo}
	return ls, mockRepo
}

func Test_Create_ModelPassedToFunctionNil(t *testing.T) {
	ls, _ := setup()

	lib, err := ls.Create(nil)
	if err == nil {
		t.Error("Expected an error but got nothing")
	}
	if err.Error() != LibraryPathWasNilErr {
		t.Errorf("Expected error: %v\nGot error: %v", LibraryPathWasNilErr, err.Error())
	}

	if lib != nil {
		t.Fatal("Library should have been nil")
	}
}

func Test_Create_ErrorWhileGettingLibraryByIdFromRepo(t *testing.T) {
	ls, repo := setup()

	id, _ := uuid.NewRandom()
	libPathModel := &model.LibraryPath{LibraryID: id}

	repo.MockLibraryRepo.MockError[0] = errors.New("error")

	lib, err := ls.Create(libPathModel)

	if err == nil {
		t.Error("expecting an error but was nil")
	}
	var e errs.IError
	if errors.As(err, &e) {
		expectedErr := fmt.Sprintf(ErrGetLibraryById, id)
		if e.Message() != expectedErr {
			t.Errorf("Expected error: %v\nGot error: %v", expectedErr, e.Message())
		}
	} else {
		t.Errorf("Expected a specific error but got: %v", err)
	}

	if lib != nil {
		t.Fatal("error was thrown lib should be nil")
	}
}

func Test_Create_LibraryNilFromRepo(t *testing.T) {
	ls, repo := setup()

	id, err := uuid.NewRandom()
	libPathModel := &model.LibraryPath{LibraryID: id}

	repo.MockLibraryRepo.MockModel[0] = nil
	lib, err := ls.Create(libPathModel)
	if err == nil {
		t.Error("expecting an error but was nil")
	}
	expectedErrorMessage := fmt.Sprintf(LibraryNilErr, id)
	if err.Error() != expectedErrorMessage {
		t.Errorf("Expected error: %v\nGot error: %v", expectedErrorMessage, err.Error())
	}

	if lib != nil {
		t.Fatal("error was thrown lib should be nil")
	}
}

func Test_Create_LibraryExists_CreatingLibraryPathReturnsError(t *testing.T) {
	ls, repo := setup()

	id, _ := uuid.NewRandom()
	libPathModel := &model.LibraryPath{Path: "/some/expected/path", LibraryID: id}
	library := &model.Library{ID: id}

	repo.MockLibraryRepo.MockModel[0] = library
	repo.MockLibraryPathRepo.MockError[1] = errors.New("error")

	lib, err := ls.Create(libPathModel)
	if err == nil {
		t.Error("Expecting an error but was nil")
	}
	var e errs.IError
	if errors.As(err, &e) {
		expectedError := fmt.Sprintf(ErrCreateLibraryPath)
		if e.Message() != expectedError {
			t.Errorf("Expected error: %v\nGot error: %v", expectedError, e.Message())
		}
	} else {
		t.Errorf("Expected specific error but got: %v", err)
	}

	if lib != nil {
		t.Fatal("lib and error was not nil")
	}
}

func Test_Create_Succcess(t *testing.T) {
	ls, repo := setup()

	id, _ := uuid.NewRandom()
	libPathModel := &model.LibraryPath{Path: "/some/expected/path", LibraryID: id}
	library := &model.Library{ID: id}

	repo.MockLibraryRepo.MockModel[0] = library
	repo.MockLibraryPathRepo.MockModel[1] = libPathModel

	lib, err := ls.Create(libPathModel)
	if err != nil {
		t.Errorf("Expected no error but got: %v", err)
	}
	if lib != libPathModel {
		t.Errorf("Expected: %v\nGot: %v", libPathModel, lib)
	}
}

func Test_GetAll_RepoReturnsError(t *testing.T) {
	ls, repo := setup()

	repo.MockLibraryPathRepo.MockError[0] = errors.New("error")

	libPaths, err := ls.GetAll()
	if err == nil {
		t.Error("expected error but was nil")
	}
	var e errs.IError
	if errors.As(err, &e) {
		expectedErr := fmt.Sprintf(ErrGetAllLibraryPaths)
		if e.Message() != expectedErr {
			t.Errorf("Expected error: %v\nGot error: %v", expectedErr, e.Message())
		}
	} else {
		t.Errorf("Expected specific error but got: %v", err)
	}

	if libPaths != nil {
		t.Fatal("error received but lib paths was not nil")
	}
}

func Test_GetAll_Success(t *testing.T) {
	ls, repo := setup()

	id, _ := uuid.NewRandom()
	libPath := model.LibraryPath{ID: id}
	libPaths := []model.LibraryPath{libPath}
	repo.MockLibraryPathRepo.MockModels[0] = libPaths

	libPaths, err := ls.GetAll()
	if err != nil {
		t.Errorf("Expected no error but got: %v", err)
	}
	if libPaths == nil {
		t.Error("Expected library paths but was nil")
	}
	if len(libPaths) != 1 {
		t.Errorf("Expected result to be of length 1 but was: %v", len(libPaths))
	}
	libPath = libPaths[len(libPaths)-1]
	if libPath.ID != id {
		t.Errorf("Expected lib path to have id %v but was %v", id, libPath.ID)
	}
}
