package server

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/slugger7/exorcist/internal/db/exorcist/public/model"
	"github.com/slugger7/exorcist/internal/job"
)

func (s *Server) RegisterLibraryRoutes(r *gin.RouterGroup) *gin.RouterGroup {
	r.GET("/libraries/scan", s.ScanLibrary)
	r.POST("/libraries", s.CreateLibrary)
	return r
}

func (s *Server) ScanLibrary(c *gin.Context) {
	go job.ScanPath(s.repo)
}

func (s *Server) CreateLibrary(c *gin.Context) {
	var body struct {
		Name string
	}

	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "could not read body of request"})
		return
	}

	newLibrary := &model.Library{
		Name: body.Name,
	}

	if err := s.service.LibraryService().CreateLibrary(newLibrary); err != nil {
		log.Printf("Something went wrong creating a library: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "could not create new library"})
		return
	}

	c.JSON(http.StatusCreated, newLibrary)
}
