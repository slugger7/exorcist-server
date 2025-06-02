package job

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/slugger7/exorcist/internal/db/exorcist/public/model"
	"github.com/slugger7/exorcist/internal/dto"
	"github.com/stretchr/testify/assert"
)

func Test_CreateGenerateChecksumJob(t *testing.T) {
	jobId, _ := uuid.NewRandom()
	id, _ := uuid.NewRandom()

	actual, err := CreateGenerateChecksumJob(id, jobId)
	assert.Nil(t, err)

	actualData := *actual.Data
	actual.Data = nil

	expectedData := fmt.Sprintf(`{"mediaId":"%v"}`, id)
	expected := model.Job{
		JobType:  model.JobTypeEnum_GenerateChecksum,
		Status:   model.JobStatusEnum_NotStarted,
		Data:     nil,
		Priority: dto.JobPriority_Low,
		Parent:   &jobId,
	}

	assert.Equal(t, expected, *actual)
	assert.Equal(t, expectedData, actualData)
}
