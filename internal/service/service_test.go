package service

import (
	"testing"

	"github.com/slugger7/exorcist/internal/db/exorcist/public/model"
)

type mockUserService struct{}
type mockLibraryService struct{}

func (mus mockUserService) CreateUser(username, password string) (*model.User, error) {
	panic("not implemented")
}

func (mus mockUserService) ValidateUser(username, password string) (*model.User, error) {
	panic("not implemented")
}

func (mls mockLibraryService) CreateLibrary(_ *model.Library) error {
	panic("not implemented")
}

var service = &Service{env: nil, userService: mockUserService{}, libraryService: mockLibraryService{}}

func Test_UserService(t *testing.T) {
	actualUserService := service.UserService()
	if actualUserService == nil {
		t.Error("user service was nil")
	}
}

func Test_LibraryService(t *testing.T) {
	actualUserService := service.LibraryService()
	if actualUserService == nil {
		t.Error("library service was nil")
	}
}
