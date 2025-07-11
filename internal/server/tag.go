package server

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/slugger7/exorcist/internal/dto"
)

func (s *server) withTagGetAll(r *gin.RouterGroup, route Route) *server {
	r.GET(route, s.getAllTags)
	return s
}

func (s *server) withTagCreate(r *gin.RouterGroup, route Route) *server {
	r.POST(route, s.createTags)
	return s
}

func (s *server) withTagGetMedia(r *gin.RouterGroup, route Route) *server {
	r.GET(fmt.Sprintf("%v/:%v/media", route, idKey), s.getMediaByTag)
	return s
}

func (s *server) getMediaByTag(c *gin.Context) {
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

	media, err := s.service.Tag().GetMedia(id, search)
	if err != nil {
		s.logger.Errorf("could not get media from tag service for %v: %v", id, err.Error())
		c.AbortWithError(http.StatusUnprocessableEntity, fmt.Errorf("could not get media for tag"))
		return
	}

	dtos := make([]dto.MediaOverviewDTO, len(media.Data))
	for i, m := range media.Data {
		dtos[i] = *(&dto.MediaOverviewDTO{}).FromModel(m)
	}

	c.JSON(http.StatusOK, dto.DataToPage(dtos, *media))
}

func (s *server) createTags(c *gin.Context) {
	var tags []string
	if err := c.ShouldBindBodyWithJSON(&tags); err != nil {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, err)
		return
	}

	var createdTags []dto.TagDTO
	var accErrs error
	for _, t := range tags {
		createdTag, err := s.service.Tag().Upsert(strings.Trim(t, " \n\t"))
		if err != nil {
			accErrs = errors.Join(accErrs, err)
			continue
		}

		createdTags = append(createdTags, *(&dto.TagDTO{}).FromModel(createdTag))
	}
	if accErrs != nil {
		s.logger.Errorf("some errors while creating tags: %v", accErrs.Error())
	}

	c.JSON(http.StatusCreated, createdTags)
}

func (s *server) getAllTags(c *gin.Context) {
	tags, err := s.repo.Tag().GetAll()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	tagDtos := make([]dto.TagDTO, len(tags))
	for i, t := range tags {
		tagDtos[i] = *(&dto.TagDTO{}).FromModel(&t)
	}

	c.JSON(http.StatusOK, tagDtos)
}
