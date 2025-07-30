package server

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/slugger7/exorcist/internal/db/exorcist/public/model"
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

func (s *server) withMediaPutTag(r *gin.RouterGroup, route Route) *server {
	r.PUT(fmt.Sprintf("%v/:%v/tags/:%v", route, idKey, tagIdKey), s.putMediaTag)
	return s
}

func (s *server) withMediaDeleteTag(r *gin.RouterGroup, route Route) *server {
	r.DELETE(fmt.Sprintf("%v/:%v/tags/:%v", route, idKey, tagIdKey), s.deleteMediaTag)
	return s
}

func (s *server) withMediaPutPerson(r *gin.RouterGroup, route Route) *server {
	r.PUT(fmt.Sprintf("%v/:%v/people/:%v", route, idKey, personIdKey), s.putMediaPerson)
	return s
}

func (s *server) withMediaDeletePerson(r *gin.RouterGroup, route Route) *server {
	r.DELETE(fmt.Sprintf("%v/:%v/people/:%v", route, idKey, personIdKey), s.deleteMediaPerson)
	return s
}

func (s *server) withMediaDelete(r *gin.RouterGroup, route Route) *server {
	r.DELETE(fmt.Sprintf("%v/:%v", route, idKey), s.deleteMedia)
	return s
}

func (s *server) deleteMedia(c *gin.Context) {
	id, err := uuid.Parse(c.Param(idKey))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"error": "could not parse media id"})
		return
	}

	var query dto.DeleteMediaDTO
	if err = c.ShouldBindQuery(&query); err != nil {
		c.AbortWithError(http.StatusUnprocessableEntity, err)
		return
	}

	if query.Physical == nil {
		b := false
		query.Physical = &b
	}

	err = s.service.Media().Delete(id, *query.Physical)
	if err != nil {
		s.logger.Errorf("error deleting video (%v): %v", id, err.Error())
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	// TODO: notify websockets of media deletion

	c.Status(http.StatusOK)
}

func (s *server) deleteMediaPerson(c *gin.Context) {
	id, err := uuid.Parse(c.Param(idKey))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"error": "could not parse media id"})
		return
	}

	personId, err := uuid.Parse(c.Param(personIdKey))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"error": "could not parse person id"})
		return
	}

	err = s.repo.Person().RemoveFromMedia(model.MediaPerson{MediaID: id, PersonID: personId})
	if err != nil {
		s.logger.Errorf("error removing person (%v) from media (%v): %v", personId.String(), id.String(), err.Error())
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)
}

func (s *server) putMediaPerson(c *gin.Context) {
	id, err := uuid.Parse(c.Param(idKey))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"error": "could not parse media id"})
		return
	}

	personId, err := uuid.Parse(c.Param(personIdKey))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"error": "could not parse person id"})
		return
	}

	m, err := s.service.Media().AddPerson(id, personId)
	if err != nil {
		s.logger.Errorf("error while adding person (%v) to media (%v): %v", personId.String(), id.String(), err.Error())
		c.AbortWithStatus(http.StatusInternalServerError)
	}

	c.JSON(http.StatusOK, m)
}

func (s *server) deleteMediaTag(c *gin.Context) {
	id, err := uuid.Parse(c.Param(idKey))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"error": "could not parse media id"})
		return
	}

	tagId, err := uuid.Parse(c.Param(tagIdKey))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"error": "could not parse tag id"})
		return
	}

	err = s.repo.Tag().RemoveFromMedia(model.MediaTag{MediaID: id, TagID: tagId})
	if err != nil {
		s.logger.Errorf("error while removing tag (%v) from media (%v): %v", tagId.String(), id.String(), err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "error while removing tag from media"})
		return
	}

	c.Status(http.StatusOK)
}

func (s *server) putMediaTag(c *gin.Context) {
	id, err := uuid.Parse(c.Param(idKey))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"error": "could not parse media id"})
		return
	}

	tagId, err := uuid.Parse(c.Param(tagIdKey))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"error": "could not parse tag id"})
		return
	}

	m, err := s.service.Media().AddTag(id, tagId)
	if err != nil {
		s.logger.Errorf("could not add tag %v to media %v: %v", tagId.String(), id.String(), err.Error())
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, m)
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

	userId, err := s.getUserId(c)
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	result, err := s.repo.Media().GetAll(*userId, search)
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
