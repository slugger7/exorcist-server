package server

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/slugger7/exorcist/internal/environment"
)

func (s *Server) withCors(r *gin.Engine) *Server {
	if s.env.AppEnv == environment.AppEnvEnum.Prod {
		return s
	}

	config := cors.DefaultConfig()
	config.AllowOrigins = s.env.CorsOrigins
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"}
	config.AllowHeaders = []string{"Accept", "Authorization", "Content-Type", "Origin"}
	config.ExposeHeaders = []string{"Content-Length"}
	config.AllowCredentials = true

	r.Use(cors.New(config))

	return s
}
