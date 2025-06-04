package server

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/slugger7/exorcist/internal/db/exorcist/public/model"
	"github.com/slugger7/exorcist/internal/dto"
)

func (s *server) withLibraryPathCreate(r *gin.RouterGroup, route Route) *server {
	r.POST(route, s.CreateLibraryPath)
	return s
}
func (s *server) withLibraryPathGetAll(r *gin.RouterGroup, route Route) *server {
	r.GET(route, s.GetAllLibraryPaths)
	return s
}

func (s *server) withLibraryPathGet(r *gin.RouterGroup, route Route) *server {
	r.GET(fmt.Sprintf("%v/:id", route), s.GetLibraryPath)
	return s
}

func (s *server) GetLibraryPath(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		e := fmt.Sprintf(ErrIdParse, c.Param("id"))
		s.logger.Error(e)
		c.JSON(http.StatusBadRequest, createError(e))
		return
	}

	libraryPath, err := s.repo.LibraryPath().GetById(id)
	if err != nil {
		s.logger.Errorf("error fetching library path by id: %v", id)
		c.JSON(http.StatusInternalServerError, createError("could not get libray path by id"))
		return
	}

	libPathDto := *(&dto.LibraryPathDTO{}).FromModel(*libraryPath)

	c.JSON(http.StatusOK, libPathDto)
}

const ErrCreatingLibraryPath string = "colud not create new library path"

func (s *server) CreateLibraryPath(c *gin.Context) {
	var body dto.CreateLibraryPathModelDTO

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

	result := (&dto.LibraryPathDTO{}).FromModel(*libPath)

	c.JSON(http.StatusCreated, result)
}

const ErrGetAllLibraryPathsService = "could not get all library paths"

func (s *server) GetAllLibraryPaths(c *gin.Context) {
	libraryPaths, err := s.service.LibraryPath().GetAll()
	if err != nil {
		s.logger.Errorf("Error getting all libraries\n%v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": ErrGetAllLibraryPathsService})
		return
	}

	libPaths := []dto.LibraryPathDTO{}
	for _, l := range libraryPaths {
		libPaths = append(libPaths, *(&dto.LibraryPathDTO{}).FromModel(l))
	}

	c.JSON(http.StatusOK, libPaths)
}
