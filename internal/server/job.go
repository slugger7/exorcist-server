package server

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

const jobRoute = "/jobs"

func (s *Server) withJobRoutes(r *gin.RouterGroup) *Server {
	r.GET(fmt.Sprintf("%v/start-runner", jobRoute), s.startJobRunner)
	return s
}

func (s *Server) startJobRunner(c *gin.Context) {
	s.jobCh <- true
	c.JSON(http.StatusOK, nil)
}
