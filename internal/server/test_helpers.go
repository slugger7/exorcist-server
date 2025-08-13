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
	mock_libraryPathService "github.com/slugger7/exorcist/internal/mock/service/library_path"
	mock_userService "github.com/slugger7/exorcist/internal/mock/service/user"
	libraryService "github.com/slugger7/exorcist/internal/service/library"
	libraryPathService "github.com/slugger7/exorcist/internal/service/library_path"
	userService "github.com/slugger7/exorcist/internal/service/user"
	"go.uber.org/mock/gomock"
)

const SET_COOKIE_URL string = "/set"
const AUTH_ROUTE string = "/authenticated"
const OK string = "ok"

type TestCookie struct {
	Value uuid.UUID `json:"value"`
}

type TestServer struct {
	server                 *server
	mockService            *mock_service.MockService
	mockUserService        *mock_userService.MockUserService
	mockLibraryService     *mock_libraryService.MockLibraryService
	mockLibraryPathService *mock_libraryPathService.MockLibraryPathService
	ctrl                   *gomock.Controller
	engine                 *gin.Engine
	authGroup              *gin.RouterGroup
	request                *http.Request
}

func setupServer(t *testing.T) *TestServer {
	ctrl := gomock.NewController(t)
	svc := mock_service.NewMockService(ctrl)
	env := environment.EnvironmentVariables{LogLevel: "none"}
	server := &server{logger: logger.New(&env), service: svc}
	engine := setupEngine()
	return &TestServer{server: server, mockService: svc, engine: engine, ctrl: ctrl}
}

func (s *TestServer) withUserService() *TestServer {
	us := mock_userService.NewMockUserService(s.ctrl)

	s.mockService.EXPECT().
		User().
		DoAndReturn(func() userService.UserService {
			return us
		}).
		AnyTimes()

	s.mockUserService = us

	return s
}

func (s *TestServer) withLibraryService() *TestServer {
	ls := mock_libraryService.NewMockLibraryService(s.ctrl)

	s.mockService.EXPECT().
		Library().
		DoAndReturn(func() libraryService.LibraryService {
			return ls
		}).
		AnyTimes()

	s.mockLibraryService = ls

	return s
}

func (s *TestServer) withLibraryPathService() *TestServer {
	ls := mock_libraryPathService.NewMockLibraryPathService(s.ctrl)

	s.mockService.EXPECT().
		LibraryPath().
		DoAndReturn(func() libraryPathService.LibraryPathService {
			return ls
		}).
		AnyTimes()

	s.mockLibraryPathService = ls

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

func (s *TestServer) withAuthGetRequest(params string) *TestServer {
	req, _ := http.NewRequest("GET", fmt.Sprintf("%v/%v", AUTH_ROUTE, params), nil)
	s.request = req
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

func errBody(e ApiError, args ...any) string {
	message := e
	if args != nil {
		message = fmt.Sprintf(e, args...)
	}
	return fmt.Sprintf(`{"error":"%v"}`, message)
}

func msgBody(msg string, args ...any) string {
	message := msg
	if args != nil {
		message = fmt.Sprintf(msg, args...)
	}
	return fmt.Sprintf(`{"message":"%v"}`, message)
}
