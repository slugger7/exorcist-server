package server

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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

func (s *server) withPlaylistsMedia(r *gin.RouterGroup, route Route) *server {
	r.GET(fmt.Sprintf("%v/:%v/media", route, idKey), s.getMediaByPlaylist)
	return s
}

func (s *server) withPlaylistMediaAdd(r *gin.RouterGroup, route Route) *server {
	r.PUT(fmt.Sprintf("%v/:%v/media", route, idKey), s.putPlaylistMedia)
	return s
}

func (s *server) putPlaylistMedia(c *gin.Context) {
	playlistId, err := uuid.Parse(c.Param(idKey))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"error": "could not parse playlist id"})
		return
	}

	var playlistMediaDtos []dto.CreatePlaylistMediaDTO
	if err := c.ShouldBindBodyWithJSON(&playlistMediaDtos); err != nil {
		c.AbortWithError(http.StatusUnprocessableEntity, err)
		return
	}

	playlistMedia, err := s.service.Playlist().AddMedia(playlistId, playlistMediaDtos)
	if err != nil {
		s.logger.Errorf("error adding media to playlist %v: %v", playlistId.String(), err.Error())
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, playlistMedia)
}

func (s *server) getMediaByPlaylist(c *gin.Context) {
	playlistId, err := uuid.Parse(c.Param(idKey))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"error": "could not parse playlist id"})
		return
	}

	var search dto.MediaSearchDTO
	if err := c.ShouldBindQuery(&search); err != nil {
		c.AbortWithError(http.StatusUnprocessableEntity, err)
		return
	}

	if search.Limit == 0 {
		search.Limit = 50
	}

	userId, err := s.getUserId(c)
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	mediaOverviewModels, err := s.service.Playlist().GetMedia(playlistId, *userId, search)
	if err != nil {
		s.logger.Errorf("error fetching media for playlist %v: %v", playlistId.String(), err.Error())
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	dtos := make([]dto.MediaOverviewDTO, len(mediaOverviewModels.Data))
	for i, m := range mediaOverviewModels.Data {
		dtos[i] = *(&dto.MediaOverviewDTO{}).FromModel(m)
	}

	c.JSON(http.StatusOK, dto.DataToPage(dtos, *mediaOverviewModels))
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
