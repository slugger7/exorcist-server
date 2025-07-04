package server

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/slugger7/exorcist/internal/dto"
)

func (s *server) withTagsGetAll(r *gin.RouterGroup, route Route) *server {
	r.GET(route, s.getAllTags)
	return s
}

func (s *server) withTagsCreate(r *gin.RouterGroup, route Route) *server {
	r.POST(route, s.createTags)
	return s
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
