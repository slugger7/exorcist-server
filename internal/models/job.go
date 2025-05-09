package models

import (
	"time"

	"github.com/go-jet/jet/v2/postgres"
	"github.com/google/uuid"
	"github.com/slugger7/exorcist/internal/db/exorcist/public/model"
	"github.com/slugger7/exorcist/internal/db/exorcist/public/table"
)

type CreateJobDTO struct {
	Type     model.JobTypeEnum      `json:"type" binding:"required"`
	Data     map[string]interface{} `json:"data"`
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
	Skip     int                   `form:"skip"`
	Limit    int                   `form:"limit"`
	Statuses []model.JobStatusEnum `form:"status"`
	Parent   *string               `form:"parent" binding:"omitempty,uuid"`
	OrderBy  JobOrdinal            `form:"orderBy"`
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
	Parent   *uuid.UUID          `json:"parent"`
	Priority int16               `json:"priority"`
	JobType  model.JobTypeEnum   `json:"jobType"`
	Status   model.JobStatusEnum `json:"status"`
	Data     *string             `json:"data"`
	Outcome  *string             `json:"outcome"`
	Created  time.Time           `json:"created"`
	Modified time.Time           `json:"modified"`
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
