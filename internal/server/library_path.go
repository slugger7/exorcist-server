package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/slugger7/exorcist/internal/db/exorcist/public/model"
)

const libraryPathRoute = "/libraryPaths"

func (s *Server) WithLibraryPathRoutes(r *gin.RouterGroup) *Server {
	r.POST(libraryPathRoute, s.CreateLibraryPath)
	r.GET(libraryPathRoute, s.GetAllLibraryPaths)
	return s
}

func (s *Server) CreateLibraryPath(c *gin.Context) {
	var body struct {
		LibraryId uuid.UUID
		Path      string
	}

	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "could not read body of request"})
		return
	}
	if body.Path == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no path in body"})
		return
	}

	libPath := &model.LibraryPath{LibraryID: body.LibraryId, Path: body.Path}
	libPath, err := s.service.LibraryPathService().Create(libPath)
	if err != nil {
		s.logger.Errorf("Erorr creating library path\n%v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "colud not create new library path"})
		return
	}

	c.JSON(http.StatusCreated, libPath)
}

const ErrGetAllLibraryPathsService = "could not get all library paths"

func (s *Server) GetAllLibraryPaths(c *gin.Context) {
	libraryPaths, err := s.service.LibraryPathService().GetAll()
	if err != nil {
		s.logger.Errorf("Error getting all libraries\n%v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": ErrGetAllLibraryPathsService})
		return
	}
	c.JSON(http.StatusOK, libraryPaths)
}
