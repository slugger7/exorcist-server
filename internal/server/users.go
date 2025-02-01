package server

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/slugger7/exorcist/internal/db/exorcist/public/model"
	errs "github.com/slugger7/exorcist/internal/errors"
)

func (s *Server) RegisterUserRoutes(r *gin.RouterGroup) *gin.RouterGroup {
	r.GET("/users", s.GetUsers)
	r.POST("/users", s.CreateUser)

	return r
}

func (s *Server) GetUsers(c *gin.Context) {
	//session := sessions.Default(c)
	//user := session.Get(userKey)
	log.Println("Getting users")
	var users []struct {
		model.User
	}
	err := s.repo.UserRepo().
		GetUserByUsernameAndPassword("admin", "admin").
		Query(&users)
	errs.CheckError(err)

	c.JSON(http.StatusOK, gin.H{"user": users[len(users)-1].Username})
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
	}

	user, err := s.service.UserService().CreateUser(newUser.Username, newUser.Password)
	errs.CheckError(err)

	c.JSON(http.StatusCreated, user)
}
