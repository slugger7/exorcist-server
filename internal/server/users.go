package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
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

type ResetPasswordModel struct {
	OldPassword    string `json:"oldPassword" binding:"required"`
	NewPassword    string `json:"newPassword" binding:"required"`
	RepeatPassword string `json:"repeatPassword" binding:"required,eqfield=NewPassword"`
}

func (s *Server) UpdatePassword(c *gin.Context) {
	var model ResetPasswordModel
	if err := c.ShouldBindBodyWithJSON(&model); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, model)
}
