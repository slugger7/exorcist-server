package libraryPathService

import (
	"errors"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/slugger7/exorcist/internal/db/exorcist/public/model"
	"github.com/slugger7/exorcist/internal/mocks/mrepository"
)

func setup() (*LibraryPathService, *mrepository.MockRepo) {
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

	libPathModel := &model.LibraryPath{}

	expectedError := "expected error"
	repo.MockLibraryRepo.MockError[0] = errors.New(expectedError)
	lib, err := ls.Create(libPathModel)
	if err == nil {
		t.Error("expecting an error but was nil")
	}
	expectedErrorMessage := fmt.Sprintf("github.com/slugger7/exorcist/internal/service/library_path.(*LibraryPathService).Create: could not get library by id\n%v", expectedError)
	if err.Error() != expectedErrorMessage {
		t.Errorf("Expected error: %v\nGot error: %v", expectedErrorMessage, err.Error())
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
	expectedError := "expected error"

	repo.MockLibraryRepo.MockModel[0] = library
	repo.MockLibraryPathRepo.MockError[1] = errors.New(expectedError)

	lib, err := ls.Create(libPathModel)
	if err == nil {
		t.Error("Expecting an error but was nil")
	}
	expectedErrorMessage := fmt.Sprintf("github.com/slugger7/exorcist/internal/service/library_path.(*LibraryPathService).Create: could not create new library path\n%v", expectedError)
	if expectedErrorMessage != err.Error() {
		t.Errorf("Expected error: %v\nGot error: %v", expectedErrorMessage, err.Error())
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
