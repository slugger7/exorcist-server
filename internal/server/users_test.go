package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/slugger7/exorcist/internal/assert"
	"github.com/slugger7/exorcist/internal/db/exorcist/public/model"
	"github.com/slugger7/exorcist/internal/models"
	"go.uber.org/mock/gomock"
)

func Test_Create_InvalidBody(t *testing.T) {
	r := setupEngine()
	s := setupOldServer()

	r.POST("/", s.server.CreateUser)

	req, err := http.NewRequest("POST", "/", body(`{invalid json}`))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)
	expectedStatusCode := http.StatusBadRequest
	if rr.Code != expectedStatusCode {
		t.Errorf("wrong status code returned\nexpected %v but got %v", expectedStatusCode, rr.Code)
	}
	expectedBody := `{"error":"could not read body of request"}`
	if body := rr.Body.String(); body != expectedBody {
		t.Errorf("incorrect body\nexpected %v but got %v", expectedBody, body)
	}
}

func Test_Create_ServiceReturnsError(t *testing.T) {
	r := setupEngine()
	s := setupOldServer()

	expectedErrorMessage := "expected error"
	s.mockService.User.MockError[0] = errors.New(expectedErrorMessage)
	r.POST("/", s.server.CreateUser)

	req, err := http.NewRequest("POST", "/", body(`{"username":"someUsername","password":"somePassword"}`))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)
	expectedStatusCode := http.StatusBadRequest
	if rr.Code != expectedStatusCode {
		t.Errorf("wrong status code returned\nexpected %v but got %v", expectedStatusCode, rr.Code)
	}
	expectedBody := fmt.Sprintf(`{"error":"%s"}`, expectedErrorMessage)
	if body := rr.Body.String(); body != expectedBody {
		t.Errorf("incorrect body\nexpected %v but got %v", expectedBody, body)
	}
}

func Test_Create_Success(t *testing.T) {
	s := setupServer(t).
		withUserService()

	nu := &CreateUserModel{
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

	rr := s.withPostEndpoint(s.server.CreateUser).
		withPostRequest(bodyM(nu)).
		exec()
	body, _ := json.Marshal(m)
	assert.StatusCode(t, http.StatusCreated, rr.Code)
	assert.Body(t, string(body), rr.Body.String())
}

func Test_UpdatePassword_InvalidBody(t *testing.T) {
	s := setupServer(t).
		withAuth()

	rr := s.withAuthPutEndpoint(s.server.UpdatePassword, "").
		withAuthPutRequest(body("{invalid json body}"), "").
		withCookie(TestCookie{}).
		exec()

	expectedStatus := http.StatusUnprocessableEntity
	if rr.Code != expectedStatus {
		t.Errorf("Expected status: %v\nGot status: %v", expectedStatus, rr.Code)
	}
}

func Test_UpdatePassword_ServiceReturnsError(t *testing.T) {
	s := setupServer(t).
		withUserService().
		withAuth()

	rpm := models.ResetPasswordModel{
		OldPassword:    "good old boy",
		NewPassword:    "sparkly new",
		RepeatPassword: "sparkly new",
	}
	id, _ := uuid.NewRandom()

	s.mockUserService.EXPECT().
		UpdatePassword(gomock.Eq(id), gomock.Eq(rpm)).
		DoAndReturn(func(uuid.UUID, models.ResetPasswordModel) error {
			return fmt.Errorf("some error")
		}).
		Times(1)

	rr := s.withAuthPutEndpoint(s.server.UpdatePassword, "").
		withAuthPutRequest(bodyM(rpm), "").
		withCookie(TestCookie{Value: id}).
		exec()

	assert.StatusCode(t, http.StatusInternalServerError, rr.Code)
	assert.Body(t, fmt.Sprintf(`{"error":"%v"}`, ErrUpdatePassword), rr.Body.String())
}

func Test_UpdatePasswrod_ServiceSucceeds(t *testing.T) {
	s := setupServer(t).
		withUserService().
		withAuth()

	rpm := models.ResetPasswordModel{
		OldPassword:    "good old boy",
		NewPassword:    "sparkly new",
		RepeatPassword: "sparkly new",
	}

	id, _ := uuid.NewRandom()

	s.mockUserService.EXPECT().
		UpdatePassword(gomock.Eq(id), gomock.Eq(rpm)).
		DoAndReturn(func(uuid.UUID, models.ResetPasswordModel) error {
			return nil
		}).
		Times(1)

	rr := s.withAuthPutEndpoint(s.server.UpdatePassword, "").
		withAuthPutRequest(bodyM(rpm), "").
		withCookie(TestCookie{Value: id}).
		exec()

	assert.StatusCode(t, http.StatusOK, rr.Code)
	assert.Body(t, fmt.Sprintf(`{"message":"%v"}`, OkPasswordUpdate), rr.Body.String())
}
