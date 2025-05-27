package server

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (s *Server) withVideoGet(r *gin.RouterGroup, route Route) *Server {
	r.GET(fmt.Sprintf("%v/:id", route), s.getVideoStream)
	return s
}

func (s *Server) getVideoStream(c *gin.Context) {
	idString := c.Param("id")

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
