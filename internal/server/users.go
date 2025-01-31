package server

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func (s *Server) RegisterUserRoutes(r *gin.RouterGroup) *gin.RouterGroup {
	r.GET("/users", s.GetUsers)

	return r
}

func (s *Server) GetUsers(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get(userKey)
	s.repo.JobRepo().FetchNextJob()
	c.JSON(http.StatusOK, gin.H{"user": user})
}
