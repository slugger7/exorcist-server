package dto

import (
	"time"

	"github.com/go-jet/jet/v2/postgres"
	"github.com/google/uuid"
	"github.com/slugger7/exorcist/internal/db/exorcist/public/model"
	"github.com/slugger7/exorcist/internal/db/exorcist/public/table"
)

type CreateJobDTO struct {
	Type     model.JobTypeEnum      `json:"type" binding:"required" tstype:"model.JobTypeEnum"`
	Data     map[string]interface{} `json:"data" tstype:"ScanPathData | GenerateThumbnailData"`
	Priority *JobPriority           `json:"priority"`
}

type JobPriority = int16

const (
	JobPriority_Highest JobPriority = iota
	JobPriority_High
	JobPriority_MediumHigh
	JobPriority_Medium
	JobPriority_MediumLow
	JobPriority_Low
	JobPriority_Lowest
)

type JobSearchDTO struct {
	PageRequestDTO
	Statuses []model.JobStatusEnum `form:"status" json:"statuses" tstype:"model.JobStatusEnum"`
	Parent   *string               `form:"parent" binding:"omitempty,uuid" json:"parent"`
	OrderBy  JobOrdinal            `form:"orderBy" json:"orderBy"`
	JobTypes []model.JobTypeEnum   `form:"type" tstype:"model.JobTypeEnum" json:"jobTypes"`
}

type JobOrdinal string

const (
	JobOrdinal_Created  JobOrdinal = "created"
	JobOrdinal_Modified JobOrdinal = "modified"
	JobOrdinal_Priority JobOrdinal = "priority"
)

func (o JobOrdinal) ToColumn() postgres.Column {
	switch o {
	case JobOrdinal_Created:
		return table.Job.Created
	case JobOrdinal_Modified:
		return table.Job.Modified
	case JobOrdinal_Priority:
		return table.Job.Priority
	default:
		return table.Job.Created
	}
}

type JobDTO struct {
	Id       uuid.UUID           `json:"id"`
	Parent   *uuid.UUID          `json:"parent,omitempty"`
	Priority int16               `json:"priority,omitempty"`
	JobType  model.JobTypeEnum   `json:"jobType,omitempty" tstype:"model.JobTypeEnum"`
	Status   model.JobStatusEnum `json:"status,omitempty" tstype:"model.JobStatusEnum"`
	Data     *string             `json:"data,omitempty"`
	Outcome  *string             `json:"outcome,omitempty"`
	Created  time.Time           `json:"created,omitempty"`
	Modified time.Time           `json:"modified,omitempty"`
}

func (j *JobDTO) FromModel(m model.Job) *JobDTO {
	j.Id = m.ID
	j.Parent = m.Parent
	j.Priority = m.Priority
	j.JobType = m.JobType
	j.Status = m.Status
	j.Data = m.Data
	j.Outcome = m.Outcome
	j.Created = m.Created
	j.Modified = m.Modified

	return j
}
