package job

import (
	"github.com/google/uuid"
	"github.com/slugger7/exorcist/internal/db/exorcist/public/model"
)

type GenerateChecksumData struct {
	VideoId uuid.UUID `json:"videoId"`
}

func (jr *JobRunner) GenerateChecksum(job *model.Job) error {
	panic("not implemented")
}
