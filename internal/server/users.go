package server

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) RegisterUserRoutes(r *gin.RouterGroup) *gin.RouterGroup {
	r.POST("/users", s.CreateUser)

	return r
}

func (s *Server) CreateUser(c *gin.Context) {
	log.Println("Creating user")
	var newUser struct {
		Username string
		Password string
	}

	if err := c.BindJSON(&newUser); err != nil {
		s.logger.Info("Colud not read body")
		c.JSON(http.StatusBadRequest, gin.H{"error": "could not read body of request"})
		return
	}

	user, err := s.service.UserService().CreateUser(newUser.Username, newUser.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, user)
}
