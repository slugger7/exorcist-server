package job

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/slugger7/exorcist/internal/assert"
	"github.com/slugger7/exorcist/internal/db/exorcist/public/model"
)

func Test_CreateGenerateThumbnailJob(t *testing.T) {
	id, _ := uuid.NewRandom()
	imagePath := "some path"
	timestamp, height, width := 1337, 69, 420

	actual, err := CreateGenerateThumbnailJob(id, imagePath, timestamp, height, width)
	assert.ErrorNil(t, err)

	actualData := *actual.Data
	actual.Data = nil

	expectedData := fmt.Sprintf(`{"videoId":"%v","path":"%v","timestamp":%v,"height":%v,"width":%v}`, id, imagePath, timestamp, height, width)
	expected := model.Job{
		JobType: model.JobTypeEnum_GenerateThumbnail,
		Status:  model.JobStatusEnum_NotStarted,
		Data:    nil,
	}

	assert.Eq(t, expected, *actual)
	assert.Eq(t, expectedData, actualData)
}
