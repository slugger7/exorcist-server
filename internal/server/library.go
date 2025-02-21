package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/slugger7/exorcist/internal/db/exorcist/public/model"
)

const libraryRoute = "/libraries"

func (s *Server) WithLibraryRoutes(r *gin.RouterGroup) *Server {
	r.POST(libraryRoute, s.CreateLibrary)
	r.GET(libraryRoute, s.GetLibraries)
	r.GET(fmt.Sprintf("%v/:id/*action", libraryRoute), s.LibraryAction)
	return s
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

	lib, err := s.service.Library().Create(newLibrary)
	if err != nil {
		s.logger.Errorf("could not create library: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "could not create new library"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": lib.ID})
}

type LibraryModel struct {
	Id       uuid.UUID `json:"id,omitempty"`
	Name     string    `json:"name,omitempty"`
	Created  time.Time `json:"created,omitempty"`
	Modified time.Time `json:"modified,omitempty"`
}

func (lm *LibraryModel) From(m model.Library) *LibraryModel {
	lm.Id = m.ID
	lm.Name = m.Name
	lm.Created = m.Created
	lm.Modified = m.Modified
	return lm
}

func (s *Server) GetLibraries(c *gin.Context) {
	libs, err := s.service.Library().GetAll()
	if err != nil {
		s.logger.Errorf("could not get libraries: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not fetch libraries"})
		return
	}

	ms := []LibraryModel{}
	for _, l := range libs {
		ms = append(ms, *(&LibraryModel{}).From(l))
	}

	c.JSON(http.StatusOK, ms)
}

const ErrIdParse = "Could not parse id in path: %v"
const ErrLibraryAction = "could not perform %v on %v"

func (s *Server) LibraryAction(c *gin.Context) {
	id := c.Param("id")
	action := c.Param("action")

	libraryId, err := uuid.Parse(id)
	if err != nil {
		e := fmt.Sprintf(ErrIdParse, id)
		s.logger.Error(e)
		c.JSON(http.StatusBadRequest, gin.H{"error": e})
		return
	}

	err = s.service.Library().Action(libraryId, action)
	if err != nil {
		s.logger.Errorf("Could not perform action %v on %v: %v", action, libraryId, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf(ErrLibraryAction, action, libraryId)})
	}
}
