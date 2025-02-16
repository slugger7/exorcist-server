package server

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

const jobRoute = "/jobs"

func (s *Server) WithJobRoutes(r *gin.RouterGroup) *Server {
	r.GET(fmt.Sprintf("%v/start-runner", jobRoute), s.StartJobRunner)
	return s
}

func (s *Server) StartJobRunner(c *gin.Context) {
	s.jobCh <- true
	c.JSON(http.StatusOK, nil)
}
