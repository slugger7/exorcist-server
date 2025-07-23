package server

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/slugger7/exorcist/internal/dto"
	"github.com/slugger7/exorcist/internal/models"
)

func (s *server) withUserCreate(r *gin.RouterGroup, route Route) *server {
	r.POST(route, s.CreateUser)
	return s
}

func (s *server) withUserUpdatePassword(r *gin.RouterGroup, route Route) *server {
	r.PUT(route, s.UpdatePassword)
	return s
}

const ErrCreateUser ApiError = "could not create new user"

func (s *server) CreateUser(c *gin.Context) {
	var newUser models.CreateUserDTO

	if err := c.ShouldBindBodyWithJSON(&newUser); err != nil {
		s.logger.Info("Colud not read body")
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	user, err := s.service.User().Create(newUser.Username, newUser.Password)
	if err != nil {
		s.logger.Errorf("could not create new user: %v", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": ErrCreateUser})
		return
	}

	c.JSON(http.StatusCreated, user)
}

const ErrUpdatePassword string = "could not update password"
const OkPasswordUpdate string = "password updated"

func (s *server) UpdatePassword(c *gin.Context) {
	var model dto.ResetPasswordDTO
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
