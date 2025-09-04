package server

import (
	"net/http"

	"github.com/gorilla/websocket"
)

func (s *server) wsUpgrader() websocket.Upgrader {
	return websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			origin := r.Header.Get("Origin")

			for _, o := range s.env.CorsOrigins {
				if origin == o {
					return true
				}
			}

			return false
		},
	}
}
