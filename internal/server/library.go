package server

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/slugger7/exorcist/internal/db/exorcist/public/model"
	"github.com/slugger7/exorcist/internal/models"
)

func (s *Server) withLibraryGetAction(r *gin.RouterGroup, route Route) *Server {
	r.GET(fmt.Sprintf("%v/:id/*action", route), s.LibraryAction)
	return s
}

func (s *Server) withLibraryPost(r *gin.RouterGroup, route Route) *Server {
	r.POST(route, s.CreateLibrary)
	return s
}

func (s *Server) withLibraryGet(r *gin.RouterGroup, route Route) *Server {
	r.GET(route, s.GetLibraries)
	return s
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

const (
	ErrIdParse       ApiError = "Could not parse id in path: %v"
	ErrLibraryAction ApiError = "could not perform %v on %v"
)

func (s *Server) LibraryAction(c *gin.Context) {
	id := c.Param("id")
	action := c.Param("action")

	libraryId, err := uuid.Parse(id)
	if err != nil {
		e := fmt.Sprintf(ErrIdParse, id)
		s.logger.Error(e)
		c.JSON(http.StatusBadRequest, createError(e))
		return
	}

	err = s.service.Library().Action(libraryId, action)
	if err != nil {
		s.logger.Errorf("Could not perform action %v on %v: %v", action, libraryId, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf(ErrLibraryAction, action, libraryId)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "started"})
}
