package server

import (
	"fmt"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/slugger7/exorcist/internal/environment"
)

const (
	root         string = "/"
	userRoute    string = "/users"
	libraryRoute string = "/libraries"
	videoRoute   string = "/videos"
)

func (s *Server) RegisterRoutes() http.Handler {
	if s.env.AppEnv == environment.AppEnvEnum.Prod {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}
	r := gin.Default()

	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:5173"}
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"}
	config.AllowHeaders = []string{"Accept", "Authorization", "Content-Type", "Origin"}
	config.ExposeHeaders = []string{"Content-Length"}
	config.AllowCredentials = true

	r.Use(cors.New(config))

	s.withCookieStore(r)

	// Register authentication routes
	s.withAuthLogin(&r.RouterGroup, fmt.Sprintf("%v/login", root)).
		withAuthLogout(&r.RouterGroup, fmt.Sprintf("%v/logout", root))

	authenticated := r.Group("/api")
	authenticated.Use(s.AuthRequired)
	// Register user controller routes
	s.withUserCreate(authenticated, userRoute).
		withUserUpdatePassword(authenticated, userRoute)

	// Register library controller routes
	s.withLibraryGet(authenticated, libraryRoute).
		withLibraryGetAction(authenticated, libraryRoute).
		withLibraryPost(authenticated, libraryRoute)

	// Register library path controller routes
	s.withLibraryPathCreate(authenticated, libraryPathRoute).
		withLibraryPathGetAll(authenticated, libraryPathRoute)

	s.WithJobRoutes(authenticated)

	r.GET("/health", s.HealthHandler)
	return r
}

func (s *Server) HealthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, s.repo.Health())
}
