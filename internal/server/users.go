package server

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	errs "github.com/slugger7/exorcist/internal/errors"
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
		log.Println("Colud not read body")
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	user, err := s.service.UserService().CreateUser(newUser.Username, newUser.Password)
	errs.CheckError(err)

	c.JSON(http.StatusCreated, user)
}
