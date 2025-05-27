package server

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/slugger7/exorcist/internal/models"
)

func (s *Server) withMediaSearch(r *gin.RouterGroup, route Route) *Server {
	r.GET(route, s.getMedia)
	return s
}

func (s *Server) getMedia(c *gin.Context) {
	var search models.MediaSearchDTO

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

	dtos := make([]models.MediaOverviewDTO, len(result.Data))
	for i, m := range result.Data {
		dtos[i] = *m.ToDTO()
	}

	c.JSON(http.StatusOK, models.DataToPage(dtos, *result))
}
