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

func (s *server) withTagPut(r *gin.RouterGroup, route Route) *server {
	r.PUT(fmt.Sprintf("%v/:%v", route, idKey), s.putTag)
	return s
}

func (s *server) putTag(c *gin.Context) {
	id, err := uuid.Parse(c.Param(idKey))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"error": "could not parse tag id"})
		return
	}

	var updateDto dto.TagUpdateDTO
	if err := c.ShouldBindBodyWithJSON(&updateDto); err != nil {
		c.AbortWithError(http.StatusUnprocessableEntity, err)
		return
	}

	updateModel := model.Tag{
		ID:   id,
		Name: updateDto.Name,
	}

	updatedModel, err := s.repo.Tag().Update(updateModel)
	if err != nil {
		s.logger.Errorf("error updating tag %v: %v", id.String(), err.Error())
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	updatedDto := (&dto.TagDTO{}).FromModel(updatedModel)

	c.JSON(http.StatusOK, updatedDto)
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

	userId, err := s.getUserId(c)
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	media, err := s.service.Tag().GetMedia(id, *userId, search)
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
	var search dto.TagSearchDTO

	if err := c.ShouldBindQuery(&search); err != nil {
		c.AbortWithError(http.StatusUnprocessableEntity, err)
		return
	}

	tags, err := s.repo.Tag().GetAll(search)
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
