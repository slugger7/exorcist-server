package server

import (
	"log"
	"net/http"
	"strings"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

var secret = []byte("secret")

const userKey = "userId"

func (s *Server) RegisterAuthenticationRoutes(r *gin.Engine) *gin.Engine {
	r.Use(sessions.Sessions("mysession", cookie.NewStore([]byte(s.env.Secret))))

	r.POST("/login", s.Login)
	r.GET("/logout", s.Logout)

	return r
}

func (s *Server) AuthRequired(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get(userKey)

	if user == nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	c.Next()
}

func (s *Server) Login(c *gin.Context) {
	session := sessions.Default(c)
	var userBody struct {
		Username string
		Password string
	}

	if err := c.BindJSON(&userBody); err != nil {
		log.Println("Could not read body on login")
		c.JSON(http.StatusBadRequest, gin.H{"error": "could not read body of request"})
		return
	}

	if strings.Trim(userBody.Username, " ") == "" || strings.Trim(userBody.Password, " ") == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Parameters can't be empty"})
		return
	}

	user, err := s.service.UserService().ValidateUser(userBody.Username, userBody.Password)
	if err != nil || user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "could not authenticate with credentials"})
		return
	}

	session.Set(userKey, user.ID.String())
	if err := session.Save(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save session"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Successfully authenticated user"})
}

func (s *Server) Logout(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get(userKey)
	if user == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid session token"})
		return
	}

	session.Delete(userKey)
	if err := session.Save(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save session"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Successfully logged out"})
}
