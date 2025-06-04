package server

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/slugger7/exorcist/internal/dto"
)

// https://medium.com/@abhishekranjandev/building-a-production-grade-websocket-for-notifications-with-golang-and-gin-a-detailed-guide-5b676dcfbd5a

func (s *server) withJobRoutes(r *gin.RouterGroup, route Route) *server {
	r.GET(fmt.Sprintf("%v/start-runner", route), s.startJobRunner)
	return s
}

func (s *server) withJobCreate(r *gin.RouterGroup, route Route) *server {
	r.POST(route, s.CreateJob)
	return s
}

func (s *server) withJobGetAll(r *gin.RouterGroup, route Route) *server {
	r.GET(route, s.getAllJobs)
	return s
}

func (s *server) startJobRunner(c *gin.Context) {
	s.jobCh <- true
	c.JSON(http.StatusOK, nil)
}

const (
	ErrJobCreate ApiError = "could not create job"
)

func (s *server) CreateJob(c *gin.Context) {
	var cm dto.CreateJobDTO
	if err := c.ShouldBindBodyWithJSON(&cm); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err})
		return
	}

	job, err := s.service.Job().Create(cm)
	if err != nil {
		s.logger.Errorf("could not create job: %v", err.Error())
		c.JSON(http.StatusInternalServerError, createError(ErrJobCreate))
		return
	}

	jobDto := (&dto.JobDTO{}).FromModel(*job)
	message := dto.WSMessage[dto.JobDTO]{
		Topic: dto.WSTopic_JobCreate,
		Data:  *jobDto,
	}

	message.SendToAll(s.websockets)
	c.JSON(http.StatusOK, job)
}

const ErrGetAllJobs ApiError = "could not get all jobs"

func (s *server) getAllJobs(c *gin.Context) {
	var jobSearch dto.JobSearchDTO

	if err := c.ShouldBindQuery(&jobSearch); err != nil {
		s.logger.Errorf("could not bind query to entity %v", err.Error())
		c.AbortWithStatus(http.StatusUnprocessableEntity)
		return
	}

	jobsPage, err := s.repo.Job().GetAll(jobSearch)
	if err != nil {
		s.logger.Errorf("colud not get jobs: %v", err.Error())
		c.JSON(http.StatusInternalServerError, errBody(ErrGetAllJobs))
		return
	}

	jobDtos := make([]dto.JobDTO, len(jobsPage.Data))
	for i, j := range jobsPage.Data {
		jobDtos[i] = *(&dto.JobDTO{}).FromModel(j)
	}

	c.JSON(http.StatusOK, dto.DataToPage(jobDtos, *jobsPage))
}
