package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const videoRoute = "/videos"

func (s *Server) WithVideoRoutes(r *gin.RouterGroup) *Server {
	r.GET(videoRoute, s.GetVideos)
	r.GET(fmt.Sprintf("%s/:id", videoRoute), s.GetVideo)
	return s
}

type CreateVideoDTO struct {
	LibraryPathId uuid.UUID `binding:"required,uuid4"`
	RelativePath  string    `binding:"required,unix_addr"`
	Title         string    `binding:"required"`
	FileName      string    `binding:"required"`
	Height        int32     `binding:"required"`
	Width         int32     `binding:"required"`
	Runtime       int64     `binding:"required"`
	Size          int64     `binding:"required"`
	Checksum      *string
	Deleted       *bool
	Exists        *bool
	Created       time.Time
	Modified      time.Time
}

func (s *Server) GetVideos(c *gin.Context) {
	vids, err := s.service.Video().GetAll()
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
