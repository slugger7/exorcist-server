package server

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/slugger7/exorcist/internal/db/exorcist/public/model"
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

func (s *server) withPersonGetMedia(r *gin.RouterGroup, route Route) *server {
	r.GET(fmt.Sprintf(`%v/:%v/media`, route, personIdKey), s.getMediaByPerson)
	return s
}

func (s *server) withPersonPut(r *gin.RouterGroup, route Route) *server {
	r.PUT(fmt.Sprintf("%v/:%v", route, idKey), s.putPerson)
	return s
}

func (s *server) putPerson(c *gin.Context) {
	id, err := uuid.Parse(c.Param(idKey))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"error": "could not parse person id"})
		return
	}

	var updateDto dto.PersonUpdateDTO
	if err := c.ShouldBindBodyWithJSON(&updateDto); err != nil {
		c.AbortWithError(http.StatusUnprocessableEntity, err)
		return
	}

	updateModel := model.Person{
		ID:   id,
		Name: updateDto.Name,
	}

	updatedModel, err := s.repo.Person().Update(updateModel)
	if err != nil {
		s.logger.Errorf("error while updating person by id %v: %v", id.String(), err.Error())
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	updatedDto := (&dto.PersonDTO{}).FromModel(updatedModel)

	c.JSON(http.StatusOK, updatedDto)
}

func (s *server) getMediaByPerson(c *gin.Context) {
	id, err := uuid.Parse(c.Param(personIdKey))
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

	userId, err := s.getUserId(c)
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	media, err := s.service.Person().GetMedia(id, *userId, search)
	if err != nil {
		s.logger.Errorf("could not get media from person service for %v: %v", id, err.Error())
		c.AbortWithError(http.StatusUnprocessableEntity, fmt.Errorf("could not get media for person"))
		return
	}

	dtos := make([]dto.MediaOverviewDTO, len(media.Data))
	for i, m := range media.Data {
		dtos[i] = *(&dto.MediaOverviewDTO{}).FromModel(m)
	}

	c.JSON(http.StatusOK, dto.DataToPage(dtos, *media))
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
	var search dto.PersonSearchDTO

	if err := c.ShouldBindQuery(&search); err != nil {
		c.AbortWithError(http.StatusUnprocessableEntity, err)
		return
	}

	people, err := s.repo.Person().GetAll(search)
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
