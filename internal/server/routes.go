package server

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/slugger7/exorcist/internal/environment"
)

type Route = string

const (
	root        Route = "/"
	users       Route = "/users"
	libraries   Route = "/libraries"
	videos      Route = "/videos"
	images      Route = "/images"
	media       Route = "/media"
	jobs        Route = "/jobs"
	libraryPath Route = "/libraryPaths"
	people      Route = "/people"
	tags        Route = "/tags"
	playlists   Route = "/playlists"
)

type key = string

const (
	nameKey     key = "name"
	idKey       key = "id"
	tagIdKey    key = "tagIdKey"
	personIdKey key = "personIdKey"
)

func (s *server) RegisterRoutes() http.Handler {
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
	s.withUserCreate(authenticated, users).
		withUserUpdatePassword(authenticated, users)

	// Register library controller routes
	s.withLibraryGet(authenticated, libraries).
		withLibraryPost(authenticated, libraries).
		withLibraryGetPaths(authenticated, libraries).
		withLibraryGetMedia(authenticated, libraries)

	// Register library path controller routes
	s.withLibraryPathCreate(authenticated, libraryPath).
		withLibraryPathGetAll(authenticated, libraryPath).
		withLibraryPathGet(authenticated, libraryPath).
		withLibraryPut(authenticated, libraries)

	// Register media controller routes
	s.withMediaSearch(authenticated, media).
		withMediaGet(authenticated, media).
		withMediaPutTag(authenticated, media).
		withMediaDeleteTag(authenticated, media).
		withMediaPutPerson(authenticated, media).
		withMediaDeletePerson(authenticated, media).
		withMediaDelete(authenticated, media).
		withMediaPut(authenticated, media)

	s.withImageGet(authenticated, images).
		withVideoGet(authenticated, videos).
		withVideoPut(authenticated, videos)

	// Register job controller routes
	s.withJobRoutes(authenticated, jobs).
		withJobCreate(authenticated, jobs).
		withJobGetAll(authenticated, jobs)

	// Register person controller routes
	s.withPersonGetAll(authenticated, people).
		withPersonCreate(authenticated, people).
		withPersonGetMedia(authenticated, people).
		withPersonPut(authenticated, people)

	// Regsiter tags controller routes
	s.withTagGetAll(authenticated, tags).
		withTagCreate(authenticated, tags).
		withTagGetMedia(authenticated, tags).
		withTagPut(authenticated, tags)

	// Register playlist controller routes
	s.withPlaylistsGetAll(authenticated, playlists).
		withPlaylistsCreate(authenticated, playlists).
		withPlaylistsMedia(authenticated, playlists).
		withPlaylistMediaAdd(authenticated, playlists).
		withPlaylistPut(authenticated, playlists)

	s.withWS(authenticated, root)

	r.GET("/health", s.HealthHandler)
	return r
}

func (s *server) HealthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, s.repo.Health())
}
