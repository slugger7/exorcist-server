package server

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/slugger7/exorcist/internal/db/exorcist/public/model"
	"github.com/slugger7/exorcist/internal/dto"
)

func (s *server) withLibraryPost(r *gin.RouterGroup, route Route) *server {
	r.POST(route, s.CreateLibrary)
	return s
}

func (s *server) withLibraryGet(r *gin.RouterGroup, route Route) *server {
	r.GET(route, s.GetLibraries)
	return s
}

func (s *server) withLibraryGetPaths(r *gin.RouterGroup, route Route) *server {
	r.GET(fmt.Sprintf("%v/:%v/libraryPaths", route, idKey), s.LibraryGetPaths)
	return s
}

func (s *server) withLibraryGetMedia(r *gin.RouterGroup, route Route) *server {
	r.GET(fmt.Sprintf("%v/:%v/media", route, idKey), s.getMediaByLibrary)
	return s
}

const (
	ErrLibraryPathsForLibrary ApiError = "could not get library paths for library %v"
	ErrIdParse                ApiError = "could not parse id: %v"
)

func (s *server) getMediaByLibrary(c *gin.Context) {
	id, err := uuid.Parse(c.Param(idKey))
	if err != nil {
		c.AbortWithStatus(http.StatusUnprocessableEntity)
		return
	}

	var search dto.MediaSearchDTO

	if err := c.ShouldBindQuery(&search); err != nil {
		c.AbortWithError(http.StatusUnprocessableEntity, err)
		return
	}

	if search.Limit == 0 {
		search.Limit = 100
	}

	media, err := s.service.Library().GetMedia(id, search)
	if err != nil {
		s.logger.Errorf("colud not fetch media for library (%v): %v", id.String(), err.Error())
		c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("something went wrong while fetching media for library"))
		return
	}

	dtos := make([]dto.MediaOverviewDTO, len(media.Data))
	for i, m := range media.Data {
		dtos[i] = *(&dto.MediaOverviewDTO{}).FromModel(m)
	}

	c.JSON(http.StatusOK, dto.DataToPage(dtos, *media))
}

func (s *server) LibraryGetPaths(c *gin.Context) {
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

	libPathModels := make([]dto.LibraryPathDTO, len(libraryPaths))
	for i, m := range libraryPaths {
		libPathModels[i] = *(&dto.LibraryPathDTO{}).FromModel(m)
	}

	c.JSON(http.StatusOK, libPathModels)
}

const ErrCreatingLibrary ApiError = "could not create new library"

func (s *server) CreateLibrary(c *gin.Context) {
	var cm dto.CreateLibraryDTO
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

	l := dto.LibraryDTO{}

	c.JSON(http.StatusCreated, l.FromModel(*lib))
}

const ErrGetLibraries ApiError = "could not fetch libraries"

func (s *server) GetLibraries(c *gin.Context) {
	libs, err := s.service.Library().GetAll()
	if err != nil {
		s.logger.Errorf("could not get libraries: %v", err)
		c.JSON(http.StatusInternalServerError, createError(ErrGetLibraries))
		return
	}

	ms := []dto.LibraryDTO{}
	for _, l := range libs {
		ms = append(ms, *(&dto.LibraryDTO{}).FromModel(l))
	}

	c.JSON(http.StatusOK, ms)
}
