package libraryService

import (
	"errors"
	"fmt"
	"testing"

	"github.com/google/uuid"
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
	svc             *LibraryService
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

	ls := &LibraryService{repo: mockRepo}
	return &testService{ls, mockRepo, mockLibraryRepo, mockLibraryPathRepo, mockJobRepo}
}

func Test_CreateLibrary_ProduceErrorWhileFetchingExistingLibraries(t *testing.T) {
	s := setup(t)

	name := "someLibName"
	s.libraryRepo.EXPECT().
		GetLibraryByName(name).
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
		GetLibraryByName(name).
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
		GetLibraries().
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
		GetLibraries().
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

func Test_Action_InvalidAction(t *testing.T) {
	s := setup(t)

	id, _ := uuid.NewRandom()
	action := "non-existent-action"

	err := s.svc.Action(id, action)
	if err == nil {
		t.Error("Expected and error but was nil")
	}
	expectedError := fmt.Errorf(ErrActionNotFound, action)
	if expectedError.Error() != err.Error() {
		t.Errorf("Expected error: %v\nGot error: %v", expectedError.Error(), err.Error())
	}
}

func Test_Action_RepoErrorWhenGettingLibrary(t *testing.T) {
	s := setup(t)

	id, _ := uuid.NewRandom()

	s.libraryRepo.EXPECT().
		GetLibraryById(id).
		DoAndReturn(func(uuid.UUID) (*model.Library, error) {
			return nil, fmt.Errorf("expected error")
		})

	err := s.svc.Action(id, ActionScan)
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
	s := setup(t)

	id, _ := uuid.NewRandom()
	library := &model.Library{ID: id}

	s.libraryPathRepo.EXPECT().
		GetByLibraryId(id).
		DoAndReturn(func(uuid.UUID) ([]model.LibraryPath, error) {
			return nil, fmt.Errorf("expected error")
		}).
		Times(1)

	err := s.svc.actionScan(library)
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
	s := setup(t)

	id, _ := uuid.NewRandom()
	lib := model.Library{ID: id}
	libPathID, _ := uuid.NewRandom()
	libPath := model.LibraryPath{
		ID:        libPathID,
		LibraryID: id,
		Path:      "some path",
	}

	s.libraryPathRepo.EXPECT().
		GetByLibraryId(id).
		DoAndReturn(func(uuid.UUID) ([]model.LibraryPath, error) {
			return []model.LibraryPath{libPath}, nil
		}).
		Times(1)

	s.jobRepo.EXPECT().
		CreateAll(gomock.Any()).
		DoAndReturn(func([]model.Job) ([]model.Job, error) {
			return nil, fmt.Errorf("expected error")
		})

	err := s.svc.actionScan(&lib)
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
	s := setup(t)

	id, _ := uuid.NewRandom()
	lib := model.Library{ID: id}
	libPathID, _ := uuid.NewRandom()
	libPath := model.LibraryPath{
		ID:        libPathID,
		LibraryID: id,
		Path:      "some path",
	}

	s.libraryPathRepo.EXPECT().
		GetByLibraryId(id).
		DoAndReturn(func(uuid.UUID) ([]model.LibraryPath, error) {
			return []model.LibraryPath{libPath}, nil
		}).
		Times(1)

	s.jobRepo.EXPECT().
		CreateAll(gomock.Any()).
		DoAndReturn(func([]model.Job) ([]model.Job, error) {
			return nil, nil
		})

	err := s.svc.actionScan(&lib)
	if err != nil {
		t.Errorf("Expected no error but was %v", err)
	}
}
