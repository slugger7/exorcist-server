package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/slugger7/exorcist/internal/environment"
	"github.com/slugger7/exorcist/internal/logger"
	mock_service "github.com/slugger7/exorcist/internal/mock/service"
	mock_libraryService "github.com/slugger7/exorcist/internal/mock/service/library"
	mock_userService "github.com/slugger7/exorcist/internal/mock/service/user"
	"github.com/slugger7/exorcist/internal/mocks/mservice"
	libraryService "github.com/slugger7/exorcist/internal/service/library"
	userService "github.com/slugger7/exorcist/internal/service/user"
	"go.uber.org/mock/gomock"
)

const SET_COOKIE_URL string = "/set"
const AUTH_ROUTE string = "/authenticated"
const OK string = "ok"

type TestCookie struct {
	Value uuid.UUID `json:"value"`
}

// Deprecated: use test server and methods instead
type OldTestServer struct {
	server      *Server
	mockService *mservice.MockServices
	engine      *gin.Engine
	request     *http.Request
}

type TestServer struct {
	server             *Server
	mockService        *mock_service.MockIService
	mockUserService    *mock_userService.MockIUserService
	mockLibraryService *mock_libraryService.MockILibraryService
	ctrl               *gomock.Controller
	engine             *gin.Engine
	authGroup          *gin.RouterGroup
	request            *http.Request
}

func setupServer(t *testing.T) *TestServer {
	ctrl := gomock.NewController(t)
	svc := mock_service.NewMockIService(ctrl)
	env := environment.EnvironmentVariables{LogLevel: "none"}
	server := &Server{logger: logger.New(&env), service: svc}
	engine := setupEngine()
	return &TestServer{server: server, mockService: svc, engine: engine, ctrl: ctrl}
}

func (s *TestServer) withUserService() *TestServer {
	us := mock_userService.NewMockIUserService(s.ctrl)

	s.mockService.EXPECT().
		User().
		DoAndReturn(func() userService.IUserService {
			return us
		}).
		AnyTimes()

	s.mockUserService = us

	return s
}

func (s *TestServer) withLibraryService() *TestServer {
	ls := mock_libraryService.NewMockILibraryService(s.ctrl)

	s.mockService.EXPECT().
		Library().
		DoAndReturn(func() libraryService.ILibraryService {
			return ls
		}).
		AnyTimes()

	s.mockLibraryService = ls

	return s
}

func (s *TestServer) withCookie(cookie TestCookie) *TestServer {
	rr := httptest.NewRecorder()
	cookieReq, _ := http.NewRequest("GET", SET_COOKIE_URL, bodyM(cookie))
	s.engine.ServeHTTP(rr, cookieReq)

	setCookie := rr.Header().Values("Set-Cookie")

	s.request.Header.Set("Cookie", strings.Join(setCookie, "; "))

	return s
}

func (s *TestServer) withAuth() *TestServer {
	s.authGroup = s.engine.Group(AUTH_ROUTE)
	s.authGroup.Use(s.server.AuthRequired)

	return s
}

func setupEngine() *gin.Engine {
	r := gin.New()
	r.Use(sessions.Sessions("exorcist", cookie.NewStore([]byte("cookieSecret"))))

	r.GET(SET_COOKIE_URL, func(c *gin.Context) {
		session := sessions.Default(c)

		var cookieBody TestCookie

		if err := c.BindJSON(&cookieBody); err != nil {
			panic(err)
		}

		session.Set(userKey, cookieBody.Value.String())
		_ = session.Save()
		c.String(http.StatusOK, OK)
	})
	return r
}

func (s *TestServer) withGetEndpoint(f gin.HandlerFunc, extraPathParams string) *TestServer {
	s.engine.GET(fmt.Sprintf("/%v", extraPathParams), f)
	return s
}

func (s *TestServer) withPostEndpoint(f gin.HandlerFunc) *TestServer {
	s.engine.POST("/", f)
	return s
}

func (s *TestServer) withGetRequest(params string) *TestServer {
	req, _ := http.NewRequest("GET", fmt.Sprintf("/%v", params), nil)
	s.request = req
	return s
}

func (s *TestServer) withPostRequest(body io.Reader) *TestServer {
	req, _ := http.NewRequest("POST", "/", body)
	s.request = req
	return s
}

func (s *TestServer) withAuthGetEndpoint(f gin.HandlerFunc, extraPathParams string) *TestServer {
	s.authGroup.GET(fmt.Sprintf("/%v", extraPathParams), f)
	return s
}

func (s *TestServer) withAuthGetRequest(params string) *TestServer {
	req, _ := http.NewRequest("GET", fmt.Sprintf("%v/%v", AUTH_ROUTE, params), nil)
	s.request = req
	return s
}

func (s *TestServer) withAuthPutEndpoint(f gin.HandlerFunc, extraPathParams string) *TestServer {
	s.authGroup.PUT(fmt.Sprintf("/%v", extraPathParams), f)
	return s
}

func (s *TestServer) withAuthPutRequest(body io.Reader, params string) *TestServer {
	route := fmt.Sprintf("%v/%v", AUTH_ROUTE, params)
	req, _ := http.NewRequest("PUT", route, body)
	s.request = req
	return s
}

func (s *TestServer) exec() *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	s.engine.ServeHTTP(rr, s.request)
	return rr
}

func body(body string, args ...any) *bytes.Reader {
	message := body
	if args != nil {
		message = fmt.Sprintf(body, args...)
	}
	return bytes.NewReader([]byte(message))
}

// Marshals the model to json and creates the reader
func bodyM(model any) *bytes.Reader {
	b, _ := json.Marshal(model)

	return bytes.NewReader(b)
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
