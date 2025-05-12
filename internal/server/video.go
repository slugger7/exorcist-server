package server

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/slugger7/exorcist/internal/models"
)

func (s *Server) withVideoGet(r *gin.RouterGroup, route Route) *Server {
	r.GET(route, s.GetVideos)
	return s
}

func (s *Server) withVideoGetById(r *gin.RouterGroup, route Route) *Server {
	r.GET(fmt.Sprintf("%s/:id", route), s.GetVideo)
	return s
}

func (s *Server) GetVideos(c *gin.Context) {
	var search models.VideoSearchDTO

	if err := c.ShouldBindQuery(&search); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	if search.Limit == 0 {
		search.Limit = 48
	}

	vids, err := s.service.Video().GetOverview(search)
	if err != nil {
		s.logger.Errorf("could not fetch videos", err)
	}
	c.JSON(http.StatusOK, vids)
}

const ErrInvalidIdFormat = "invalid id format"
const ErrGetVideoService = "could not get video"
const ErrVideoNotFound = "video not found"

func (s *Server) GetVideo(c *gin.Context) {
	idString := c.Param("id")
	s.logger.Debugf("Getting video by id: %v", idString)

	id, err := uuid.Parse(idString)
	if err != nil {
		s.logger.Errorf("Incorrect id format: %v", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": ErrInvalidIdFormat})
		return
	}

	video, err := s.service.Video().GetById(id)
	if err != nil {
		s.logger.Errorf("Error getting video by id: %v", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": ErrGetVideoService})
		return
	}

	if video == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": ErrVideoNotFound})
		return
	}

	c.JSON(http.StatusOK, video)
}
