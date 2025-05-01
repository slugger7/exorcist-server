package server

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/slugger7/exorcist/internal/models"
)

const jobRoute = "/jobs"

func (s *Server) withJobRoutes(r *gin.RouterGroup, route Route) *Server {
	r.GET(fmt.Sprintf("%v/start-runner", jobRoute), s.startJobRunner)
	return s
}

func (s *Server) withJobCreate(r *gin.RouterGroup, route Route) *Server {
	r.POST(route, s.CreateJob)
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
	var cm models.CreateJob
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

	c.JSON(http.StatusOK, job)
}
