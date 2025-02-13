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

type CreateVideoDTO struct {
	LibraryPathId uuid.UUID `binding:"required"`
	RelativePath  string    `binding:"required"`
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

func (s Server) CreateVideo(c *gin.Context) {
	newVideo := CreateVideoDTO{}
	if err := c.ShouldBindBodyWithJSON(&newVideo); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, newVideo)
}
