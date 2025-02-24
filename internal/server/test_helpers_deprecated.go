package server

import (
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
}

// Deprecated: this is using the old mock service and should be using the newer mockgen mocks
func setupOldServer() *OldTestServer {
	svc, mSvc := mservice.SetupMockService()
	server := &Server{logger: logger.New(&environment.EnvironmentVariables{LogLevel: "none"}), service: svc}
	engine := setupEngine()
	return &OldTestServer{server: server, mockService: mSvc, engine: engine}
}
