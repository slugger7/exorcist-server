package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/slugger7/exorcist/internal/assert"
	"github.com/slugger7/exorcist/internal/db/exorcist/public/model"
	"github.com/slugger7/exorcist/internal/dto"
	"github.com/slugger7/exorcist/internal/models"
	"go.uber.org/mock/gomock"
)

func Test_Create_InvalidBody(t *testing.T) {
	s := setupServer(t)

	s.server.withUserCreate(&s.engine.RouterGroup, "/")
	rr := s.withPostRequest(body("{invalid body}")).
		exec()

	assert.StatusCode(t, http.StatusUnprocessableEntity, rr.Code)
	expectedBody := `{"error":"invalid character 'i' looking for beginning of object key string"}`
	assert.Body(t, expectedBody, rr.Body.String())
}

func Test_Create_ServiceReturnsError(t *testing.T) {
	s := setupServer(t).
		withUserService()

	u := models.CreateUserDTO{
		Username: "someUsername",
		Password: "somePassword",
	}

	s.mockUserService.EXPECT().
		Create(gomock.Eq(u.Username), gomock.Eq(u.Password)).
		DoAndReturn(func(string, string) (*model.User, error) {
			return nil, fmt.Errorf("some error")
		}).
		Times(1)

	s.server.withUserCreate(&s.engine.RouterGroup, "/")
	rr := s.withPostRequest(bodyM(u)).
		exec()

	assert.StatusCode(t, http.StatusBadRequest, rr.Code)
	assert.Body(t, errBody(ErrCreateUser), rr.Body.String())
}

func Test_Create_Success(t *testing.T) {
	s := setupServer(t).
		withUserService()

	nu := &models.CreateUserDTO{
		Username: "expectedUsername",
		Password: "somePassword",
	}

	m := &model.User{
		Username: nu.Username,
		Password: nu.Password,
	}

	s.mockUserService.EXPECT().
		Create(gomock.Eq(nu.Username), gomock.Eq(nu.Password)).
		DoAndReturn(func(string, string) (*model.User, error) {
			return m, nil
		}).
		Times(1)

	s.server.withUserCreate(&s.engine.RouterGroup, "/")
	rr := s.withPostRequest(bodyM(nu)).
		exec()

	body, _ := json.Marshal(m)
	assert.StatusCode(t, http.StatusCreated, rr.Code)
	assert.Body(t, string(body), rr.Body.String())
}

func Test_UpdatePassword_InvalidBody(t *testing.T) {
	s := setupServer(t).
		withAuth()

	s.server.withUserUpdatePassword(s.authGroup, "/")
	rr := s.withAuthPutRequest(body("{invalid json body}"), "").
		withCookie(TestCookie{}).
		exec()

	assert.StatusCode(t, http.StatusUnprocessableEntity, rr.Code)
}

func Test_UpdatePassword_ServiceReturnsError(t *testing.T) {
	s := setupServer(t).
		withUserService().
		withAuth()

	rpm := dto.ResetPasswordDTO{
		OldPassword:    "good old boy",
		NewPassword:    "sparkly new",
		RepeatPassword: "sparkly new",
	}
	id, _ := uuid.NewRandom()

	s.mockUserService.EXPECT().
		UpdatePassword(gomock.Eq(id), gomock.Eq(rpm)).
		DoAndReturn(func(uuid.UUID, dto.ResetPasswordDTO) error {
			return fmt.Errorf("some error")
		}).
		Times(1)

	s.server.withUserUpdatePassword(s.authGroup, "/")
	rr := s.withAuthPutRequest(bodyM(rpm), "").
		withCookie(TestCookie{Value: id}).
		exec()

	assert.StatusCode(t, http.StatusInternalServerError, rr.Code)
	assert.Body(t, errBody(ErrUpdatePassword), rr.Body.String())
}

func Test_UpdatePasswrod_ServiceSucceeds(t *testing.T) {
	s := setupServer(t).
		withUserService().
		withAuth()

	rpm := dto.ResetPasswordDTO{
		OldPassword:    "good old boy",
		NewPassword:    "sparkly new",
		RepeatPassword: "sparkly new",
	}

	id, _ := uuid.NewRandom()

	s.mockUserService.EXPECT().
		UpdatePassword(gomock.Eq(id), gomock.Eq(rpm)).
		DoAndReturn(func(uuid.UUID, dto.ResetPasswordDTO) error {
			return nil
		}).
		Times(1)

	s.server.withUserUpdatePassword(s.authGroup, "/")
	rr := s.withAuthPutRequest(bodyM(rpm), "").
		withCookie(TestCookie{Value: id}).
		exec()

	assert.StatusCode(t, http.StatusOK, rr.Code)
	assert.Body(t, fmt.Sprintf(`{"message":"%v"}`, OkPasswordUpdate), rr.Body.String())
}
