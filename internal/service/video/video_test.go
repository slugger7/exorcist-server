package videoService

import (
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/slugger7/exorcist/internal/db/exorcist/public/model"
	errs "github.com/slugger7/exorcist/internal/errors"
	"github.com/slugger7/exorcist/internal/mocks/mrepository"
)

func setup() (*VideoService, *mrepository.MockRepository) {
	mockRepo := mrepository.SetupMockRespository()
	vs := &VideoService{repo: mockRepo}
	return vs, mockRepo
}

func Test_GetAll_ErrorFromRepo(t *testing.T) {
	vs, mr := setup()

	mr.MockVideoRepo.MockError[0] = errors.New("error")

	vids, err := vs.GetAll()
	if err == nil {
		t.Error("Expected error but got nil")
	}
	var e errs.IError
	if errors.As(err, &e) {
		if e.Message() != ErrGetAllVideos {
			t.Errorf("Expected error: %v\nGot error: %v", ErrGetAllVideos, e.Message())
		}
	} else {
		t.Errorf("Expected specific error but got: %v", err)
	}

	if vids != nil {
		t.Fatalf("Vids was supposed to be nil but was: %v", vids)
	}
}

func Test_GetAll_Success(t *testing.T) {
	vs, mr := setup()

	id, _ := uuid.NewRandom()
	videos := []model.Video{{ID: id}}
	mr.MockVideoRepo.MockModels[0] = videos

	vids, err := vs.GetAll()
	if err != nil {
		t.Fatal(err)
	}

	if len(vids) != 1 {
		t.Error("Incorrect videos length")
	}
	actualId := vids[len(vids)-1].ID
	if actualId != id {
		t.Errorf("Expected video with id: %v\nGot video with id: %v", id, actualId)
	}
}
