package server

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/slugger7/exorcist/internal/dto"
)

type key = string

const personName key = "name"

func (s *Server) withPersonUpsert(r *gin.RouterGroup, route Route) *Server {
	r.PUT(fmt.Sprintf("%v/:%v", route, personName), s.putPerson)
	return s
}

func (s *Server) putPerson(c *gin.Context) {
	name := c.Param(personName)
	if name == "" {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	person, err := s.service.Person().Upsert(name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not upsert person by name"})
		return
	}

	personDto := (&dto.PersonDTO{}).FromModel(person)

	c.JSON(http.StatusOK, personDto)
}
