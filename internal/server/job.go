package server

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/slugger7/exorcist/internal/models"
)

// https://medium.com/@abhishekranjandev/building-a-production-grade-websocket-for-notifications-with-golang-and-gin-a-detailed-guide-5b676dcfbd5a

func (s *Server) withJobRoutes(r *gin.RouterGroup, route Route) *Server {
	r.GET(fmt.Sprintf("%v/start-runner", route), s.startJobRunner)
	return s
}

func (s *Server) withJobCreate(r *gin.RouterGroup, route Route) *Server {
	r.POST(route, s.CreateJob)
	return s
}

func (s *Server) withJobGetAll(r *gin.RouterGroup, route Route) *Server {
	r.GET(route, s.getAllJobs)
	return s
}

func (s *Server) startJobRunner(c *gin.Context) {
	s.jobCh <- true
	c.JSON(http.StatusOK, nil)
}

const (
	ErrJobCreate ApiError = "could not create job"
)

func (s *Server) CreateJob(c *gin.Context) {
	var cm models.CreateJobDTO
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

	jobDto := (&models.JobDTO{}).FromModel(*job)
	message := models.WSMessage[models.JobDTO]{
		Topic: models.WSTopic_JobCreate,
		Data:  *jobDto,
	}

	message.SendToAll(s.websockets)
	c.JSON(http.StatusOK, job)
}

const ErrGetAllJobs ApiError = "could not get all jobs"

func (s *Server) getAllJobs(c *gin.Context) {
	var jobSearch models.JobSearchDTO

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

	jobDtos := make([]models.JobDTO, len(jobsPage.Data))
	for i, j := range jobsPage.Data {
		jobDtos[i] = *(&models.JobDTO{}).FromModel(j)
	}

	c.JSON(http.StatusOK, models.DataToPage(jobDtos, *jobsPage))
}
