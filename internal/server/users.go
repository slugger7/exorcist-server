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

	return r
}

func (s *Server) GetUsers(c *gin.Context) {
	//session := sessions.Default(c)
	//user := session.Get(userKey)
	log.Println("Getting users")
	var users []struct {
		model.Users
	}
	err := s.repo.UserRepo().
		GetUserByUsernameAndPassword("admin", "admin").
		Query(&users)
	errs.CheckError(err)

	c.JSON(http.StatusOK, gin.H{"user": users[len(users)-1].Username})
}
