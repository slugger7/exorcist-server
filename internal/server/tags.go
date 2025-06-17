package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/slugger7/exorcist/internal/dto"
)

func (s *server) withTagsGetAll(r *gin.RouterGroup, route Route) *server {
	r.GET(route, s.getAllTags)
	return s
}

func (s *server) getAllTags(c *gin.Context) {
	tags, err := s.repo.Tag().GetAll()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	tagDtos := make([]dto.TagDTO, len(tags))
	for i, t := range tags {
		tagDtos[i] = *(&dto.TagDTO{}).FromModel(&t)
	}

	c.JSON(http.StatusOK, tagDtos)
}
