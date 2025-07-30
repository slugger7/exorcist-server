package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/slugger7/exorcist/internal/dto"
)

func (s *server) withPlaylistsGetAll(r *gin.RouterGroup, route Route) *server {
	r.GET(route, s.getAllPlayists)
	return s
}

func (s *server) withPlaylistsCreate(r *gin.RouterGroup, route Route) *server {
	r.POST(route, s.createPlaylists)
	return s
}

func (s *server) createPlaylists(c *gin.Context) {
	var playlists []dto.CreatePlaylistDTO
	if err := c.ShouldBindBodyWithJSON(&playlists); err != nil {
		c.AbortWithError(http.StatusUnprocessableEntity, err)
	}

	if playlists == nil || len(playlists) == 0 {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	userId, err := s.getUserId(c)
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	newPlaylists, err := s.service.Playlist().CreateAll(*userId, playlists)
	if err != nil {
		s.logger.Errorf("error occured while creating playlists: %v", err.Error())
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	newPlaylistDtos := make([]dto.PlaylistDTO, len(newPlaylists))
	for i, m := range newPlaylists {
		newPlaylistDtos[i] = *(&dto.PlaylistDTO{}).FromModel(m)
	}

	c.JSON(http.StatusCreated, newPlaylistDtos)
}

func (s *server) getAllPlayists(c *gin.Context) {
	playlists, err := s.repo.Playlist().GetAll()
	if err != nil {
		s.logger.Errorf("could not get all playlists from repo: %v", err.Error())
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	playlistDtos := make([]dto.PlaylistDTO, len(playlists))
	for i, m := range playlists {
		playlistDtos[i] = *(&dto.PlaylistDTO{}).FromModel(m)
	}

	c.JSON(http.StatusOK, playlistDtos)
}
