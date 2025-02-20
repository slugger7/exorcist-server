package server

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/slugger7/exorcist/internal/models"
)

const userRoute = "/user"

func (s *Server) WithUserRoutes(r *gin.RouterGroup) *Server {
	r.POST(userRoute, s.CreateUser)
	r.PUT(userRoute, s.UpdatePassword)
	return s
}

func (s *Server) CreateUser(c *gin.Context) {
	var newUser struct {
		Username string
		Password string
	}

	if err := c.BindJSON(&newUser); err != nil {
		s.logger.Info("Colud not read body")
		c.JSON(http.StatusBadRequest, gin.H{"error": "could not read body of request"})
		return
	}

	user, err := s.service.User().Create(newUser.Username, newUser.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, user)
}

const ErrUpdatePassword string = "could not update password"
const OkPasswordUpdate string = "password updated"

func (s *Server) UpdatePassword(c *gin.Context) {
	var model models.ResetPasswordModel
	if err := c.ShouldBindBodyWithJSON(&model); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	session := sessions.Default(c)
	user := session.Get(userKey).(string)
	id, err := uuid.Parse(user)
	if err != nil {
		s.logger.Errorf("could not parse user id (%v): %v", user, err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid session"})
		return
	}

	if err := s.service.User().UpdatePassword(id, model); err != nil {
		s.logger.Errorf("error updating password for %v: %v", id.String(), err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": ErrUpdatePassword})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": OkPasswordUpdate})
}
