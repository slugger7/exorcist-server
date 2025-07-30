package server

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/slugger7/exorcist/internal/dto"
)

func (s *server) withVideoGet(r *gin.RouterGroup, route Route) *server {
	r.GET(fmt.Sprintf("%v/:%v", route, idKey), s.getVideoStream)
	return s
}

func (s *server) withVideoPut(r *gin.RouterGroup, route Route) *server {
	r.PUT(fmt.Sprintf("%v/:%v", route, idKey), s.putVideoProgress)
	return s
}

func (s *server) putVideoProgress(c *gin.Context) {
	id, err := uuid.Parse(c.Param(idKey))
	if err != nil {
		c.AbortWithStatus(http.StatusUnprocessableEntity)
		return
	}

	var progress dto.ProgressUpdateDTO
	if err := c.ShouldBindQuery(&progress); err != nil {
		c.AbortWithError(http.StatusUnprocessableEntity, err)
		return
	}

	userId, err := s.getUserId(c)
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	prog, err := s.service.Media().LogProgress(id, *userId, progress)
	if err != nil {
		s.logger.Errorf("colud not log progress for %v: %v", id.String(), err.Error())
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	progDto := (&dto.ProgressDTO{}).FromModel(*prog)
	c.JSON(http.StatusOK, progDto)
}

func (s *server) getVideoStream(c *gin.Context) {
	idString := c.Param(idKey)

	id, err := uuid.Parse(idString)
	if err != nil {
		s.logger.Errorf("Incorrect id format: %v", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": ErrInvalidIdFormat})
		return
	}

	med, err := s.repo.Video().GetByMediaId(id)
	if err != nil {
		s.logger.Errorf("Error getting video absolute path by id: %v", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": ErrGetVideoService})
		return
	}

	if med == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": ErrVideoNotFound})
		return
	}

	c.File(med.Path)
}
