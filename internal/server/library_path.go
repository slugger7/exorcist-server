package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/slugger7/exorcist/internal/db/exorcist/public/model"
	"github.com/slugger7/exorcist/internal/models"
)

const libraryPathRoute string = "/libraryPaths"

func (s *Server) withLibraryPathCreate(r *gin.RouterGroup, route string) *Server {
	r.POST(route, s.CreateLibraryPath)
	return s
}
func (s *Server) withLibraryPathGetAll(r *gin.RouterGroup, route string) *Server {
	r.GET(route, s.GetAllLibraryPaths)
	return s
}

const ErrCreatingLibraryPath string = "colud not create new library path"

func (s *Server) CreateLibraryPath(c *gin.Context) {
	var body models.CreateLibraryPathModel

	if err := c.ShouldBindBodyWithJSON(&body); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	libPath := &model.LibraryPath{LibraryID: body.LibraryId, Path: body.Path}
	libPath, err := s.service.LibraryPath().Create(libPath)
	if err != nil {
		s.logger.Errorf("Erorr creating library path\n%v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": ErrCreatingLibraryPath})
		return
	}

	result := (&models.LibraryPathModel{}).FromModel(*libPath)

	c.JSON(http.StatusCreated, result)
}

const ErrGetAllLibraryPathsService = "could not get all library paths"

func (s *Server) GetAllLibraryPaths(c *gin.Context) {
	libraryPaths, err := s.service.LibraryPath().GetAll()
	if err != nil {
		s.logger.Errorf("Error getting all libraries\n%v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": ErrGetAllLibraryPathsService})
		return
	}

	libPaths := []models.LibraryPathModel{}
	for _, l := range libraryPaths {
		libPaths = append(libPaths, *(&models.LibraryPathModel{}).FromModel(l))
	}

	c.JSON(http.StatusOK, libPaths)
}
