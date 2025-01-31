package server

import (
	"github.com/gin-gonic/gin"
	"github.com/slugger7/exorcist/internal/job"
)

func (s *Server) RegisterLibraryRoutes(r *gin.RouterGroup) *gin.RouterGroup {
	r.GET("/library/scan", s.ScanLibrary)
	return r
}

func (s *Server) ScanLibrary(c *gin.Context) {
	go job.ScanPath(s.repo)
}
