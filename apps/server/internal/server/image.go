package server

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (s *server) withImageGet(r *gin.RouterGroup, route Route) *server {
	r.GET(fmt.Sprintf("%v/:id", route), s.getImage)
	return s
}

func (s *server) getImage(c *gin.Context) {
	idString := c.Param("id")

	id, err := uuid.Parse(idString)
	if err != nil {
		s.logger.Errorf("Incorrect id format: %v", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": ErrInvalidIdFormat})
		return
	}

	img, err := s.repo.Image().GetByMediaId(id)
	if err != nil {
		s.logger.Errorf("Error getting image by id: %v", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": ErrGetImageService})
		return
	}

	if img == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": ErrImageNotFound})
		return
	}

	c.File(img.Path)
}
