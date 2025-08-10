package job

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/slugger7/exorcist/internal/db/exorcist/public/model"
	"github.com/slugger7/exorcist/internal/dto"
	"github.com/stretchr/testify/assert"
)

func Test_CreateGenerateThumbnailJob(t *testing.T) {
	id, _ := uuid.NewRandom()
	jobId, _ := uuid.NewRandom()
	imagePath := "some path"
	timestamp := float64(1337)
	height, width := 69, 420

	actual, err := CreateGenerateThumbnailJob(model.Video{ID: id}, &jobId, imagePath, timestamp, height, width, nil, nil)
	assert.Equal(t, err, nil, "Error should be nil")

	actualData := *actual.Data
	actual.Data = nil

	expectedData := fmt.Sprintf(`{"videoId":"%v","path":"%v","timestamp":%v,"height":%v,"width":%v}`, id, imagePath, timestamp, height, width)
	expected := model.Job{
		JobType:  model.JobTypeEnum_GenerateThumbnail,
		Status:   model.JobStatusEnum_NotStarted,
		Data:     nil,
		Parent:   &jobId,
		Priority: dto.JobPriority_MediumHigh,
	}

	assert.Equal(t, expected, *actual, "Expected job should be equal to actual job")
	assert.Equal(t, expectedData, actualData, "Expected data should be equal to actual data")
}
