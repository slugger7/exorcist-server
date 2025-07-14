package libraryPathService

import (
	"errors"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/slugger7/exorcist/internal/db/exorcist/public/model"
	errs "github.com/slugger7/exorcist/internal/errors"
	mock_repository "github.com/slugger7/exorcist/internal/mock/repository"
	mock_libraryRepository "github.com/slugger7/exorcist/internal/mock/repository/library"
	mock_libraryPathRepository "github.com/slugger7/exorcist/internal/mock/repository/library_path"

	libraryRepository "github.com/slugger7/exorcist/internal/repository/library"
	libraryPathRepository "github.com/slugger7/exorcist/internal/repository/library_path"
	"go.uber.org/mock/gomock"
)

type testService struct {
	svc         *libraryPathService
	repo        *mock_repository.MockIRepository
	libRepo     *mock_libraryRepository.MockLibraryRepository
	libPathRepo *mock_libraryPathRepository.MockILibraryPathRepository
}

func setup(t *testing.T) *testService {
	ctrl := gomock.NewController(t)

	mockRepo := mock_repository.NewMockIRepository(ctrl)
	mockLibraryPathRepo := mock_libraryPathRepository.NewMockILibraryPathRepository(ctrl)
	mockLibraryRepo := mock_libraryRepository.NewMockLibraryRepository(ctrl)

	mockRepo.EXPECT().
		LibraryPath().
		DoAndReturn(func() libraryPathRepository.ILibraryPathRepository {
			return mockLibraryPathRepo
		}).
		AnyTimes()

	mockRepo.EXPECT().
		Library().
		DoAndReturn(func() libraryRepository.LibraryRepository {
			return mockLibraryRepo
		}).
		AnyTimes()

	lps := &libraryPathService{repo: mockRepo}
	return &testService{
		lps,
		mockRepo,
		mockLibraryRepo,
		mockLibraryPathRepo,
	}
}

func Test_Create_ModelPassedToFunctionNil(t *testing.T) {
	s := setup(t)

	lib, err := s.svc.Create(nil)
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
	s := setup(t)

	id, _ := uuid.NewRandom()
	libPathModel := &model.LibraryPath{LibraryID: id}

	s.libRepo.EXPECT().
		GetById(libPathModel.LibraryID).
		DoAndReturn(func(id uuid.UUID) (*model.Library, error) {
			return nil, fmt.Errorf("error")
		}).
		Times(1)

	lib, err := s.svc.Create(libPathModel)

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
	s := setup(t)

	id, err := uuid.NewRandom()
	libPathModel := &model.LibraryPath{LibraryID: id}

	s.libRepo.EXPECT().
		GetById(libPathModel.LibraryID).
		DoAndReturn(func(id uuid.UUID) (*model.Library, error) {
			return nil, nil
		}).
		Times(1)

	lib, err := s.svc.Create(libPathModel)
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
	s := setup(t)

	id, _ := uuid.NewRandom()
	libPathModel := &model.LibraryPath{Path: "/some/expected/path", LibraryID: id}
	library := &model.Library{ID: id}

	s.libRepo.EXPECT().
		GetById(id).
		DoAndReturn(func(id uuid.UUID) (*model.Library, error) {
			return library, nil
		}).
		Times(1)

	s.libPathRepo.EXPECT().
		Create(libPathModel.Path, libPathModel.LibraryID).
		DoAndReturn(func(p string, id uuid.UUID) (*model.LibraryPath, error) {
			return nil, fmt.Errorf("error")
		})

	lib, err := s.svc.Create(libPathModel)
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

func Test_Create_Success(t *testing.T) {
	s := setup(t)

	id, _ := uuid.NewRandom()
	libPathModel := &model.LibraryPath{Path: "/some/expected/path", LibraryID: id}
	library := &model.Library{ID: id}

	s.libRepo.EXPECT().
		GetById(id).
		DoAndReturn(func(id uuid.UUID) (*model.Library, error) {
			return library, nil
		}).
		Times(1)

	s.libPathRepo.EXPECT().
		Create(libPathModel.Path, libPathModel.LibraryID).
		DoAndReturn(func(p string, id uuid.UUID) (*model.LibraryPath, error) {
			return libPathModel, nil
		})

	lib, err := s.svc.Create(libPathModel)
	if err != nil {
		t.Errorf("Expected no error but got: %v", err)
	}
	if lib != libPathModel {
		t.Errorf("Expected: %v\nGot: %v", libPathModel, lib)
	}
}

func Test_GetAll_RepoReturnsError(t *testing.T) {
	s := setup(t)

	s.libPathRepo.EXPECT().
		GetAll().
		DoAndReturn(func() ([]model.LibraryPath, error) {
			return nil, fmt.Errorf("error")
		}).
		Times(1)

	libPaths, err := s.svc.GetAll()
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
	s := setup(t)

	id, _ := uuid.NewRandom()
	libPath := model.LibraryPath{ID: id}
	libPaths := []model.LibraryPath{libPath}

	s.libPathRepo.EXPECT().
		GetAll().
		DoAndReturn(func() ([]model.LibraryPath, error) {
			return libPaths, nil
		}).
		Times(1)

	libPaths, err := s.svc.GetAll()
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
