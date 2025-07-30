package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *server) withPlaylistsGetAll(r *gin.RouterGroup, route Route) *server {
	r.GET(route, s.getAllPlayists)
	return s
}

func (s *server) getAllPlayists(c *gin.Context) {
	playlists, err := s.repo.Playlist().GetAll()
	if err != nil {
		s.logger.Errorf("could not get all playlists from repo: %v", err.Error())
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, playlists)
}
