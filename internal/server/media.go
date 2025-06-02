package server

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/slugger7/exorcist/internal/dto"
	"github.com/slugger7/exorcist/internal/models"
)

func (s *Server) withMediaSearch(r *gin.RouterGroup, route Route) *Server {
	r.GET(route, s.getMedia)
	return s
}

func (s *Server) withMediaGet(r *gin.RouterGroup, route Route) *Server {
	r.GET(fmt.Sprintf("%v/:id", route), s.getMediaById)
	return s
}

func (s *Server) getMediaById(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		e := fmt.Sprintf((ErrIdParse), c.Param(("id")))
		s.logger.Error(e)
		c.JSON(http.StatusUnprocessableEntity, createError(e))
		return
	}

	m, err := s.repo.Media().GetById(id)
	if err != nil {
		s.logger.Errorf("could not get media by id: %v", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not get media by id"})
		return
	}

	if m == nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	c.JSON(http.StatusOK, (&dto.MediaDTO{}).FromModel(*m))
}

func (s *Server) getMedia(c *gin.Context) {
	var search dto.MediaSearchDTO

	if err := c.ShouldBindQuery(&search); err != nil {
		c.AbortWithError(http.StatusUnprocessableEntity, err)
		return
	}

	if search.Limit == 0 {
		search.Limit = 100
	}

	result, err := s.repo.Media().GetAll(search)
	if err != nil {
		s.logger.Errorf("could not get media from repo: %v", err.Error())
		c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("could not get media"))
		return
	}

	c.JSON(http.StatusOK, models.DataToPage(result.Data, *result))
}
