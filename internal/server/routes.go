package server

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/slugger7/exorcist/internal/environment"
)

type Route = string

const (
	root         Route = "/"
	userRoute    Route = "/users"
	libraryRoute Route = "/libraries"
	videoRoute   Route = "/videos"
	mediaRoute   Route = "/media"
)

func (s *Server) RegisterRoutes() http.Handler {
	if s.env.AppEnv == environment.AppEnvEnum.Prod {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}
	r := gin.Default()

	s.withCors(r).
		withStaticFiles(r).
		withCookieStore(r)

	// Register authentication routes
	s.withAuthLogin(&r.RouterGroup, fmt.Sprintf("%v/api/login", root)).
		withAuthLogout(&r.RouterGroup, fmt.Sprintf("%v/api/logout", root))

	authenticated := r.Group("/api")
	authenticated.Use(s.AuthRequired)

	// Register user controller routes
	s.withUserCreate(authenticated, userRoute).
		withUserUpdatePassword(authenticated, userRoute)

	// Register library controller routes
	s.withLibraryGet(authenticated, libraryRoute).
		//withLibraryGetAction(authenticated, libraryRoute).
		withLibraryPost(authenticated, libraryRoute).
		withLibraryGetPaths(authenticated, libraryRoute)

	// Register library path controller routes
	s.withLibraryPathCreate(authenticated, libraryPathRoute).
		withLibraryPathGetAll(authenticated, libraryPathRoute)

	// Register video controller routes
	s.withVideoGet(authenticated, videoRoute).
		withVideoGetById(authenticated, videoRoute)

	// Register media controller routes
	s.withMediaVideo(authenticated, mediaRoute).
		withMediaImage(authenticated, mediaRoute)

	s.withJobRoutes(authenticated)

	r.GET("/health", s.HealthHandler)
	return r
}

func (s *Server) HealthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, s.repo.Health())
}
