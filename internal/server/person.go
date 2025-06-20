package server

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/slugger7/exorcist/internal/dto"
)

func (s *server) withPersonGetAll(r *gin.RouterGroup, route Route) *server {
	r.GET(route, s.getAllPeople)
	return s
}

func (s *server) withPersonCreate(r *gin.RouterGroup, route Route) *server {
	r.POST(route, s.createPeople)
	return s
}

func (s *server) createPeople(c *gin.Context) {
	var people []string
	if err := c.ShouldBindBodyWithJSON(&people); err != nil {
		c.AbortWithError(http.StatusUnprocessableEntity, err)
		return
	}

	var createdPeople []dto.PersonDTO
	var accErrs error
	for _, p := range people {
		createdPerson, err := s.service.Person().Upsert(strings.Trim(p, " \n\t"))
		if err != nil {
			accErrs = errors.Join(accErrs, err)
			continue
		}

		createdPeople = append(createdPeople, *(&dto.PersonDTO{}).FromModel(createdPerson))
	}
	if accErrs != nil {
		s.logger.Errorf("some errors while creating people: %v", accErrs.Error())
	}

	c.JSON(http.StatusCreated, createdPeople)
}

func (s *server) getAllPeople(c *gin.Context) {
	people, err := s.repo.Person().GetAll()
	if err != nil {
		s.logger.Errorf("could not fetch people from repo: %v", err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "could not fetch people"})
		return
	}

	peopleDtos := make([]dto.PersonDTO, len(people))
	for i, p := range people {
		peopleDtos[i] = *(&dto.PersonDTO{}).FromModel(&p)
	}

	c.JSON(http.StatusOK, peopleDtos)
}
