package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/slugger7/exorcist/internal/db/exorcist/public/model"
)

func (s *Server) RegisterLibraryRoutes(r *gin.RouterGroup) *gin.RouterGroup {
	r.POST("/libraries", s.CreateLibrary)
	r.GET("/libraries", s.GetLibraries)
	return r
}

func (s *Server) CreateLibrary(c *gin.Context) {
	var body struct {
		Name string
	}

	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "could not read body of request"})
		return
	}

	if body.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no value for name"})
		return
	}

	newLibrary := model.Library{
		Name: body.Name,
	}

	lib, err := s.service.LibraryService().CreateLibrary(newLibrary)
	if err != nil {
		s.logger.Errorf("could not create library: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "could not create new library"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": lib.ID})
}

func (s *Server) GetLibraries(c *gin.Context) {
	libs, err := s.service.LibraryService().GetLibraries()
	if err != nil {
		s.logger.Errorf("could not get libraries: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not fetch libraries"})
		return
	}

	c.JSON(http.StatusOK, libs)
}
