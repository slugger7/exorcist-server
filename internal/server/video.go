package server

import (
	"fmt"
	"net/http"
	"strconv"

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

func (s *Server) defaultInt(strVal string, def int) int {
	if strVal != "" {
		val, err := strconv.Atoi(strVal)
		if err != nil {
			s.logger.Warningf("could not parse %v to int", strVal)
		}

		return val
	}

	return def
}

func (s *Server) defualtBool(strVal string, def bool) bool {
	val, err := strconv.ParseBool(strVal)
	if err != nil {
		return def
	}
	return val
}

func (s *Server) GetVideos(c *gin.Context) {
	limit := s.defaultInt(c.Query("limit"), 48)
	skip := s.defaultInt(c.Query("skip"), 0)
	orderBy := c.Query("orderBy")
	asc := s.defualtBool(c.Query("asc"), true)

	vids, err := s.service.Video().GetOverview(limit, skip, (*models.VideoOrdinal)(&orderBy), asc)
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
