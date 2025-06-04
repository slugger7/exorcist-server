package server

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (s *server) withVideoGet(r *gin.RouterGroup, route Route) *server {
	r.GET(fmt.Sprintf("%v/:%v", route, idKey), s.getVideoStream)
	return s
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
