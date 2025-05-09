package server

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/slugger7/exorcist/internal/db/exorcist/public/model"
	"github.com/slugger7/exorcist/internal/models"
)

func (s *Server) withLibraryPost(r *gin.RouterGroup, route Route) *Server {
	r.POST(route, s.CreateLibrary)
	return s
}

func (s *Server) withLibraryGet(r *gin.RouterGroup, route Route) *Server {
	r.GET(route, s.GetLibraries)
	return s
}

func (s *Server) withLibraryGetPaths(r *gin.RouterGroup, route Route) *Server {
	r.GET(fmt.Sprintf("%v/:id/libraryPaths", route), s.LibraryGetPaths)
	return s
}

const (
	ErrLibraryPathsForLibrary ApiError = "could not get library paths for library %v"
	ErrIdParse                ApiError = "could not parse id: %v"
)

func (s *Server) LibraryGetPaths(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		e := fmt.Sprintf(ErrIdParse, c.Param("id"))
		s.logger.Error(e)
		c.JSON(http.StatusBadRequest, createError(e))
		return
	}

	libraryPaths, err := s.repo.LibraryPath().GetByLibraryId(id)
	if err != nil {
		s.logger.Errorf(ErrLibraryPathsForLibrary, id)
		c.JSON(http.StatusInternalServerError, createError("could not get library paths for library"))
		return
	}

	libPathModels := make([]models.LibraryPathModel, len(libraryPaths))
	for i, m := range libraryPaths {
		libPathModels[i] = *(&models.LibraryPathModel{}).FromModel(m)
	}

	c.JSON(http.StatusOK, libPathModels)
}

const ErrCreatingLibrary ApiError = "could not create new library"

func (s *Server) CreateLibrary(c *gin.Context) {
	var cm models.CreateLibraryModel
	if err := c.ShouldBindBodyWithJSON(&cm); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err})
		return
	}

	newLibrary := model.Library{
		Name: cm.Name,
	}

	lib, err := s.service.Library().Create(&newLibrary)
	if err != nil {
		s.logger.Errorf("could not create library: %v", err)
		c.JSON(http.StatusBadRequest, createError(ErrCreatingLibrary))
		return
	}

	l := models.Library{}

	c.JSON(http.StatusCreated, l.FromModel(*lib))
}

const ErrGetLibraries ApiError = "could not fetch libraries"

func (s *Server) GetLibraries(c *gin.Context) {
	libs, err := s.service.Library().GetAll()
	if err != nil {
		s.logger.Errorf("could not get libraries: %v", err)
		c.JSON(http.StatusInternalServerError, createError(ErrGetLibraries))
		return
	}

	ms := []models.Library{}
	for _, l := range libs {
		ms = append(ms, *(&models.Library{}).FromModel(l))
	}

	c.JSON(http.StatusOK, ms)
}
