package server

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const videoRoute = "/videos"

func (s *Server) WithVideoRoutes(r *gin.RouterGroup) *Server {
	r.GET(videoRoute, s.GetVideos)
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

func (s Server) GetVideos(c *gin.Context) {
	// vids, err := s.service.VideoService().GetVideos()
	// if err != nil {
	// 	s.logger.Errorf("could not fetch videos", err)
	// }
}
