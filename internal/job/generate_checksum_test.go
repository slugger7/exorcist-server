package job

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/slugger7/exorcist/internal/assert"
	"github.com/slugger7/exorcist/internal/db/exorcist/public/model"
)

func Test_CreateGenerateChecksumJob(t *testing.T) {
	id, _ := uuid.NewRandom()

	actual, err := CreateGenerateChecksumJob(id)
	assert.ErrorNil(t, err)

	actualData := *actual.Data
	actual.Data = nil

	expectedData := fmt.Sprintf(`{"videoId":"%v"}`, id)
	expected := model.Job{
		JobType: model.JobTypeEnum_GenerateChecksum,
		Status:  model.JobStatusEnum_NotStarted,
		Data:    nil,
	}

	assert.Eq(t, expected, *actual)
	assert.Eq(t, expectedData, actualData)
}
