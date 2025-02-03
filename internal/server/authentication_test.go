package server

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/slugger7/exorcist/internal/db/exorcist/public/model"
	libraryService "github.com/slugger7/exorcist/internal/service/library"
	userService "github.com/slugger7/exorcist/internal/service/user"
)

type mockService struct {
	userService userService.IUserService
}

type mockUserService struct{}

func (ms mockService) UserService() userService.IUserService {
	fmt.Println("Return mock user service")
	return ms.userService
}
func (ms mockService) LibraryService() libraryService.ILibraryService {
	return nil
}

func (mus mockUserService) CreateUser(username, password string) (*model.User, error) {
	return nil, errors.New("Testing error")
}

func (mus mockUserService) ValidateUser(username, password string) (*model.User, error) {
	return nil, errors.New("Testing error")
}

func Test_Login(t *testing.T) {
	s := &Server{}
	r := gin.New()
	r.Use(sessions.Sessions("exorcist", cookie.NewStore([]byte("cookieSecret"))))
	r.POST("/", s.Login)
	s.service = mockService{userService: mockUserService{}}
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
}
