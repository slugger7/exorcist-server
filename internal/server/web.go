package server

import (
	"fmt"
	"strings"

	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
)

func (s *Server) withStaticFiles(r *gin.Engine) {
	if s.env.Web == nil {
		return
	}

	r.Use(static.Serve("/", static.LocalFile(*s.env.Web, false)))

	r.NoRoute(func(c *gin.Context) {
		path := c.Request.RequestURI
		if !strings.HasPrefix(path, "/api") {
			indexHtml := fmt.Sprintf("%v/index.html", *s.env.Web)
			s.logger.Debugf("Rerouting to frontend %v", indexHtml)

			c.File(indexHtml)
		}
	})
}
