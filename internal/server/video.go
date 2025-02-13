package server

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const videoRoute = "/videos"

func (s *Server) RegisterVideoRoutes(r *gin.RouterGroup) *gin.RouterGroup {
	r.POST(videoRoute, s.CreateVideo)
	return r
}

type VideoDTO struct {
	ID            *uuid.UUID
	LibraryPathId uuid.UUID
	RelativePath  string
	Title         string
	FileName      string
	Height        int32
	Width         int32
	Runtime       int64
	Size          int64
	Checksum      *string
	Deleted       *bool
	Exists        *bool
	Created       time.Time
	Modified      time.Time
}

func (s Server) CreateVideo(c *gin.Context) {
	newVideo := VideoDTO{}
	if err := c.BindJSON(&newVideo); err != nil {
		s.logger.Info("Could not read body")
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "colud not read body of request"})
		return
	}

	c.JSON(http.StatusCreated, newVideo)
}
