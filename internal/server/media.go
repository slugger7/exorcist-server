package server

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (s *Server) withMediaVideo(r *gin.RouterGroup, route Route) *Server {
	r.GET(fmt.Sprintf("%v/video/:id", route), s.getVideoStream)
	return s
}

func (s *Server) withMediaImage(r *gin.RouterGroup, route Route) *Server {
	r.GET(fmt.Sprintf("%v/image/:id", route), s.getMediaImage)
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

	vid, err := s.service.Video().GetByIdWithLibraryPath(id)
	if err != nil {
		s.logger.Errorf("Error getting video absolute path by id: %v", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": ErrGetVideoService})
		return
	}

	if vid == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": ErrVideoNotFound})
		return
	}

	absolutePath := fmt.Sprintf("%v%v", vid.LibraryPath.Path, vid.Video.RelativePath)
	c.File(absolutePath)
}

const (
	ErrGetImageService string = "error getting image by id from service"
	ErrImageNotFound   string = "image not found"
)

func (s *Server) getMediaImage(c *gin.Context) {
	idString := c.Param("id")

	id, err := uuid.Parse(idString)
	if err != nil {
		s.logger.Errorf("Incorrect id format: %v", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": ErrInvalidIdFormat})
		return
	}

	img, err := s.service.Image().GetById(id)
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
