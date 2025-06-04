package server

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/slugger7/exorcist/internal/dto"
)

func (s *server) withPersonUpsert(r *gin.RouterGroup, route Route) *server {
	r.PUT(fmt.Sprintf("%v/:%v", route, nameKey), s.putPerson)
	return s
}

func (s *server) putPerson(c *gin.Context) {
	name := c.Param(nameKey)
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
