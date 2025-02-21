package server

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
	"github.com/slugger7/exorcist/internal/environment"
	"github.com/slugger7/exorcist/internal/logger"
	"github.com/slugger7/exorcist/internal/mocks/mservice"
)

// Deprecated: use test server and methods instead
type OldTestServer struct {
	server      *Server
	mockService *mservice.MockServices
	engine      *gin.Engine
	request     *http.Request
}

// Deprecated
func (s *OldTestServer) withGetEndpoint(f gin.HandlerFunc, extraPathParams string) *OldTestServer {
	s.engine.GET(fmt.Sprintf("/%v", extraPathParams), f)
	return s
}

// Deprecated
func (s *OldTestServer) withPostEndpoint(f gin.HandlerFunc) *OldTestServer {
	s.engine.POST("/", f)
	return s
}

// Deprecated
func (s *OldTestServer) withGetRequest(body io.Reader, params string) *OldTestServer {
	req, _ := http.NewRequest("GET", fmt.Sprintf("/%v", params), body)
	s.request = req
	return s
}

// Deprecated
func (s *OldTestServer) withPostRequest(body io.Reader) *OldTestServer {
	req, _ := http.NewRequest("POST", "/", body)
	s.request = req
	return s
}

// Deprecated
func (s *OldTestServer) exec() *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	s.engine.ServeHTTP(rr, s.request)
	return rr
}

// Deprecated: this is using the old mock service and should be using the newer mockgen mocks
func setupOldServer() *OldTestServer {
	svc, mSvc := mservice.SetupMockService()
	server := &Server{logger: logger.New(&environment.EnvironmentVariables{LogLevel: "none"}), service: svc}
	engine := setupEngine()
	return &OldTestServer{server: server, mockService: mSvc, engine: engine}
}

// Deprecated: create a method to do this on the controller (library controller as example)
func (s *TestServer) withAuthGetEndpoint(f gin.HandlerFunc, extraPathParams string) *TestServer {
	s.authGroup.GET(fmt.Sprintf("/%v", extraPathParams), f)
	return s
}
