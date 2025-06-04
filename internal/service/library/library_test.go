package libraryService

import (
	"errors"
	"fmt"
	"testing"

	"github.com/slugger7/exorcist/internal/db/exorcist/public/model"
	errs "github.com/slugger7/exorcist/internal/errors"
	mock_repository "github.com/slugger7/exorcist/internal/mock/repository"
	mock_jobRepository "github.com/slugger7/exorcist/internal/mock/repository/job"
	mock_libraryRepository "github.com/slugger7/exorcist/internal/mock/repository/library"
	mock_libraryPathRepository "github.com/slugger7/exorcist/internal/mock/repository/library_path"
	jobRepository "github.com/slugger7/exorcist/internal/repository/job"
	libraryRepository "github.com/slugger7/exorcist/internal/repository/library"
	libraryPathRepository "github.com/slugger7/exorcist/internal/repository/library_path"
	"go.uber.org/mock/gomock"
)

type testService struct {
	svc             *libraryService
	repo            *mock_repository.MockIRepository
	libraryRepo     *mock_libraryRepository.MockILibraryRepository
	libraryPathRepo *mock_libraryPathRepository.MockILibraryPathRepository
	jobRepo         *mock_jobRepository.MockIJobRepository
}

func setup(t *testing.T) *testService {
	ctrl := gomock.NewController(t)

	mockRepo := mock_repository.NewMockIRepository(ctrl)
	mockLibraryRepo := mock_libraryRepository.NewMockILibraryRepository(ctrl)
	mockLibraryPathRepo := mock_libraryPathRepository.NewMockILibraryPathRepository(ctrl)
	mockJobRepo := mock_jobRepository.NewMockIJobRepository(ctrl)

	mockRepo.EXPECT().
		Library().
		DoAndReturn(func() libraryRepository.ILibraryRepository {
			return mockLibraryRepo
		}).
		AnyTimes()

	mockRepo.EXPECT().
		LibraryPath().
		DoAndReturn(func() libraryPathRepository.ILibraryPathRepository {
			return mockLibraryPathRepo
		}).
		AnyTimes()

	mockRepo.EXPECT().
		Job().
		DoAndReturn(func() jobRepository.IJobRepository {
			return mockJobRepo
		}).
		AnyTimes()

	ls := &libraryService{repo: mockRepo}
	return &testService{ls, mockRepo, mockLibraryRepo, mockLibraryPathRepo, mockJobRepo}
}

func Test_CreateLibrary_ProduceErrorWhileFetchingExistingLibraries(t *testing.T) {
	s := setup(t)

	name := "someLibName"
	s.libraryRepo.EXPECT().
		GetByName(name).
		DoAndReturn(func(string) (*model.Library, error) {
			return nil, errors.New("error")
		}).
		Times(1)

	lib := model.Library{Name: name}
	expectedErrMsg := fmt.Sprintf(ErrLibraryByName, name)
	newLib, err := s.svc.Create(&lib)
	if err != nil {
		var e errs.IError
		if errors.As(err, &e) {
			if e.Message() != expectedErrMsg {
				t.Errorf("Expected: %v\nGot: %v", expectedErrMsg, e.Message())
			}
		} else {
			t.Errorf("Expected a different error: %v", err.Error())
		}
	}

	if newLib != nil {
		t.Fatal("Error was supposed to be thrown but new lib had a value")
	}
}

func Test_CreateLibrary_WithExistingLibrary_ShouldThrowError(t *testing.T) {
	s := setup(t)

	name := "someName"

	s.libraryRepo.EXPECT().
		GetByName(name).
		DoAndReturn(func(string) (*model.Library, error) {
			return &model.Library{Name: name}, nil
		})

	expectedErr := fmt.Sprintf(ErrLibraryExists, name)

	lib := model.Library{Name: name}
	library, err := s.svc.Create(&lib)
	if err != nil {
		if err.Error() != expectedErr {
			t.Errorf("Expected: %v\nGot: %v", expectedErr, err.Error())
		}
	} else {
		t.Fatal("Expected an error but none was given")
	}

	if library != nil {
		t.Fatal("Expected library to be nil but it was defined")
	}
}

func Test_GetAll_RepoReturnsError_ShouldReturnError(t *testing.T) {
	s := setup(t)

	s.libraryRepo.EXPECT().
		GetAll().
		DoAndReturn(func() ([]model.Library, error) {
			return nil, fmt.Errorf("some error")
		})

	libs, err := s.svc.GetAll()
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
	s := setup(t)

	expectedName := "expected library name"

	s.libraryRepo.EXPECT().
		GetAll().
		DoAndReturn(func() ([]model.Library, error) {
			return []model.Library{{Name: expectedName}}, nil
		})

	actual, err := s.svc.GetAll()
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
