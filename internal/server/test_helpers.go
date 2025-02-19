package server

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/slugger7/exorcist/internal/environment"
	"github.com/slugger7/exorcist/internal/logger"
	mock_service "github.com/slugger7/exorcist/internal/mock/service"
	mock_userService "github.com/slugger7/exorcist/internal/mock/service/user"
	"github.com/slugger7/exorcist/internal/mocks/mservice"
	"go.uber.org/mock/gomock"
)

const SET_COOKIE_URL = "/set"
const OK = "ok"

type OldTestServer struct {
	server      *Server
	mockService *mservice.MockServices
	engine      *gin.Engine
	request     *http.Request
}

type TestServer struct {
	server          *Server
	mockService     *mock_service.MockIService
	mockUserService *mock_userService.MockIUserService
	ctrl            *gomock.Controller
	engine          *gin.Engine
	request         *http.Request
}

func setupEngine() *gin.Engine {
	r := gin.New()
	r.Use(sessions.Sessions("exorcist", cookie.NewStore([]byte("cookieSecret"))))

	r.GET(SET_COOKIE_URL, func(c *gin.Context) {
		session := sessions.Default(c)

		var cookieBody struct {
			value string
		}

		_ = c.BindJSON(&cookieBody)

		session.Set(userKey, cookieBody.value)
		_ = session.Save()
		c.String(http.StatusOK, OK)
	})
	return r
}

func (s *OldTestServer) withGetEndpoint(f gin.HandlerFunc, extraPathParams string) *OldTestServer {
	s.engine.GET(fmt.Sprintf("/%v", extraPathParams), f)
	return s
}

func (s *OldTestServer) withPostEndpoint(f gin.HandlerFunc) *OldTestServer {
	s.engine.POST("/", f)
	return s
}

func (s *OldTestServer) withGetRequest(body io.Reader, params string) *OldTestServer {
	req, _ := http.NewRequest("GET", fmt.Sprintf("/%v", params), body)
	s.request = req
	return s
}

func (s *OldTestServer) withPostRequest(body io.Reader) *OldTestServer {
	req, _ := http.NewRequest("POST", "/", body)
	s.request = req
	return s
}

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

func setupServer(t *testing.T) *TestServer {
	ctrl := gomock.NewController(t)
	svc := mock_service.NewMockIService(ctrl)
	env := environment.EnvironmentVariables{LogLevel: "none"}
	server := &Server{logger: logger.New(&env), service: svc}
	engine := setupEngine()
	return &TestServer{server: server, mockService: svc, engine: engine, ctrl: ctrl}
}

func setupCookies(req *http.Request, r *gin.Engine) {
	res := httptest.NewRecorder()
	cookieReq, _ := http.NewRequest("GET", SET_COOKIE_URL, body(`{"value": "val"}`))
	r.ServeHTTP(res, cookieReq)

	req.Header.Set("Cookie", strings.Join(res.Header().Values("Set-Cookie"), "; "))
}

func body(body string, args ...any) *bytes.Reader {
	message := body
	if args != nil {
		message = fmt.Sprintf(body, args...)
	}
	return bytes.NewReader([]byte(message))
}
