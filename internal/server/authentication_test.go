package server

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/slugger7/exorcist/internal/db/exorcist/public/model"
	libraryService "github.com/slugger7/exorcist/internal/service/library"
	userService "github.com/slugger7/exorcist/internal/service/user"
)

type mockService struct {
	userService userService.IUserService
}

type mockUserService struct {
	validateUserModel model.User
	validateUserError error
}

func (ms mockService) UserService() userService.IUserService {
	return ms.userService
}
func (ms mockService) LibraryService() libraryService.ILibraryService {
	return nil
}

func (mus mockUserService) CreateUser(username, password string) (*model.User, error) {
	return nil, errors.New("Testing error")
}

func (mus mockUserService) ValidateUser(username, password string) (*model.User, error) {
	return &mus.validateUserModel, mus.validateUserError
}

func Test_Login_Success(t *testing.T) {
	s := &Server{}
	r := gin.New()
	r.Use(sessions.Sessions("exorcist", cookie.NewStore([]byte("cookieSecret"))))
	r.POST("/", s.Login)

	id, err := uuid.NewRandom()
	if err != nil {
		t.Fatalf("could not generate random uuid %v", err)
	}

	s.service = mockService{userService: mockUserService{
		validateUserModel: model.User{Username: "admin", ID: id},
	}}
	body := []byte(`{"username": "admin", "password": "admin"}`)

	req, err := http.NewRequest("POST", "/", bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Got an error %v %v", status, err)
	}

	// check to see if cookie was set
	cookie := strings.Trim(rr.Header().Get("Set-Cookie"), " ")
	if cookie == "" {
		t.Errorf("No header is being set for exorcist")
	}
	if !strings.Contains(cookie, "exorcist") {
		t.Errorf("cookie was not set up correctly: %v", cookie)
	}
}
