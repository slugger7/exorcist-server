package server

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/slugger7/exorcist/internal/dto"
)

func (s *server) withMediaSearch(r *gin.RouterGroup, route Route) *server {
	r.GET(route, s.getMedia)
	return s
}

func (s *server) withMediaGet(r *gin.RouterGroup, route Route) *server {
	r.GET(fmt.Sprintf("%v/:%v", route, idKey), s.getMediaById)
	return s
}

func (s *server) withMediaPutPeople(r *gin.RouterGroup, route Route) *server {
	r.PUT(fmt.Sprintf("%v/:%v/people", route, idKey), s.putMediaPeople)
	return s
}

func (s *server) withMediaPutTags(r *gin.RouterGroup, route Route) *server {
	r.PUT(fmt.Sprintf("%v/:%v/tags", route, idKey), s.putMediaTags)
	return s
}

func (s *server) putMediaTags(c *gin.Context) {
	id, err := uuid.Parse(c.Param(idKey))
	if err != nil {
		c.AbortWithStatus(http.StatusUnprocessableEntity)
	}

	var tags []string
	if err := c.ShouldBindBodyWithJSON(&tags); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "could not process body"})
		return
	}

	m, err := s.service.Media().SetTags(id, tags)
	c.JSON(http.StatusOK, (&dto.MediaDTO{}).FromModel(*m))
}

func (s *server) putMediaPeople(c *gin.Context) {
	id, err := uuid.Parse(c.Param(idKey))
	if err != nil {
		c.AbortWithStatus(http.StatusUnprocessableEntity)
	}

	var people []string
	if err := c.ShouldBindBodyWithJSON(&people); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "could not process body"})
		return
	}

	m, err := s.service.Media().SetPeople(id, people)
	c.JSON(http.StatusOK, (&dto.MediaDTO{}).FromModel(*m))
}

func (s *server) getMediaById(c *gin.Context) {
	id, err := uuid.Parse(c.Param(idKey))
	if err != nil {
		e := fmt.Sprintf((ErrIdParse), c.Param((idKey)))
		s.logger.Error(e)
		c.JSON(http.StatusUnprocessableEntity, createError(e))
		return
	}

	m, err := s.repo.Media().GetById(id)
	if err != nil {
		s.logger.Errorf("could not get media by id: %v", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not get media by id"})
		return
	}

	if m == nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	c.JSON(http.StatusOK, (&dto.MediaDTO{}).FromModel(*m))
}

func (s *server) getMedia(c *gin.Context) {
	var search dto.MediaSearchDTO

	if err := c.ShouldBindQuery(&search); err != nil {
		c.AbortWithError(http.StatusUnprocessableEntity, err)
		return
	}

	if search.Limit == 0 {
		search.Limit = 100
	}

	result, err := s.repo.Media().GetAll(search)
	if err != nil {
		s.logger.Errorf("could not get media from repo: %v", err.Error())
		c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("could not get media"))
		return
	}

	dtos := make([]dto.MediaOverviewDTO, len(result.Data))
	for i, m := range result.Data {
		dtos[i] = *(&dto.MediaOverviewDTO{}).FromModel(m)
	}

	c.JSON(http.StatusOK, dto.DataToPage(dtos, *result))
}
