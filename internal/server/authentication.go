package server

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

const userKey string = "userId"

func (s *Server) withCookieStore(r *gin.Engine) *Server {
	r.Use(sessions.Sessions("exorcist", cookie.NewStore([]byte(s.env.Secret))))
	return s
}

func (s *Server) withAuthLogin(r *gin.RouterGroup, route string) *Server {
	r.POST(route, s.Login)
	return s
}

func (s *Server) withAuthLogout(r *gin.RouterGroup, route string) *Server {
	r.GET(route, s.Logout)
	return s
}

const ErrUnauthorized ApiError = "unauthorized"

func (s *Server) AuthRequired(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get(userKey)

	if user == nil || user == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": ErrUnauthorized})
		return
	}

	c.Next()
}

type LoginModel struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

const MsgAuthSuccess string = "successfully authenticated user"

func (s *Server) Login(c *gin.Context) {
	session := sessions.Default(c)
	var userBody LoginModel
	if err := c.ShouldBindBodyWithJSON(&userBody); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	user, err := s.service.User().Validate(userBody.Username, userBody.Password)
	if err != nil || user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": ErrUnauthorized})
		return
	}

	session.Set(userKey, user.ID.String())
	if err := session.Save(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save session"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"userId": user.ID})
}

const (
	ErrInvalidSessionToken ApiError = "invalid session token"
	MsgLoggedOut           string   = "successfully logged out"
)

func (s *Server) Logout(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get(userKey)
	if user == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": ErrInvalidSessionToken})
		return
	}

	session.Delete(userKey)
	if err := session.Save(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save session"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": MsgLoggedOut})
}
