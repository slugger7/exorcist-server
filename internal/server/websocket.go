package server

import (
	"fmt"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (s *Server) withWS(r *gin.RouterGroup, route Route) *Server {
	r.GET(fmt.Sprintf("%v/ws", route), s.ws)
	return s
}

func (s *Server) ws(c *gin.Context) {
	session := sessions.Default(c)
	if userRaw, ok := session.Get(userKey).(string); ok {
		userId, err := uuid.Parse(userRaw)
		if err != nil {
			s.logger.Warningf("could not parse user id from string: %v", userRaw)
			c.AbortWithStatus(http.StatusUnprocessableEntity)
		}

		upgrader := s.wsUpgrader()
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			s.logger.Infof("Origin: %v", c.Request.Host)
			s.logger.Errorf("failed to upgrade connection to web socket: %v", err.Error())
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		s.websocketMutex.Lock()
		defer s.websocketMutex.Unlock()
		s.websockets[userId] = conn
	} else {
		c.AbortWithStatus(http.StatusUnauthorized)
	}

}
