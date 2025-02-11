package server

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/slugger7/exorcist/internal/db/exorcist/public/model"
	"github.com/slugger7/exorcist/internal/mocks/mservice"
)

func Test_AuthRequiredMiddleware_Fails(t *testing.T) {
	r := setupEngine()
	s := &Server{}

	r.Use(s.AuthRequired)
	r.GET("/", func(ctx *gin.Context) {
		t.Errorf("Should not have run this function")
	})

	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	expectedStatusCode := http.StatusUnauthorized
	if rr.Code != expectedStatusCode {
		t.Errorf("wrong status code returned\nexpected %v but got %v", expectedStatusCode, rr.Code)
	}
}

func Test_AuthRequiredMiddleware_Succeeds(t *testing.T) {
	r := setupEngine()
	s := &Server{}

	r.Use(s.AuthRequired)
	expectedStatusCode := http.StatusOK
	r.GET("/", func(ctx *gin.Context) {
		ctx.JSON(expectedStatusCode, gin.H{"message": "success"})
	})

	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	setupCookies(req, r)
	r.ServeHTTP(rr, req)
	if rr.Code != expectedStatusCode {
		t.Errorf("wrong status code returned\nexpected %v but got %v", expectedStatusCode, rr.Code)
	}
	expectedBody := `{"message":"success"}`
	if body := rr.Body.String(); body != expectedBody {
		t.Errorf("incorrect body\nexpected %v but got %v", expectedBody, body)
	}
}

func Test_Login_InvalidBody(t *testing.T) {
	r := setupEngine()
	s := &Server{}

	r.POST("/", s.Login)

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

func Test_Login_InvalidParametersInBody(t *testing.T) {
	r := setupEngine()
	s := &Server{}

	r.POST("/", s.Login)

	req, err := http.NewRequest("POST", "/", body(`{"username": " ", "password": " "}`))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)
	expectedStatusCode := http.StatusBadRequest
	if rr.Code != expectedStatusCode {
		t.Errorf("wrong status code returned\nexpected %v but got %v", expectedStatusCode, rr.Code)
	}
	expectedBody := `{"error":"parameters can't be empty"}`
	if body := rr.Body.String(); body != expectedBody {
		t.Errorf("incorrect body\nexpected %v but got %v", expectedBody, body)
	}
}

func Test_Login_NoUserFromValidateUser(t *testing.T) {
	r := setupEngine()
	s := &Server{}
	svc, _ := mservice.SetupMockService()
	s.service = svc

	r.POST("/", s.Login)

	req, err := http.NewRequest("POST", "/", body(`{"username": "admin", "password": "admin"}`))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)
	expectedStatusCode := http.StatusUnauthorized
	if rr.Code != expectedStatusCode {
		t.Errorf("wrong status code returned\nexpected %v but got %v", expectedStatusCode, rr.Code)
	}
	expectedBody := `{"error":"could not authenticate with credentials"}`
	if body := rr.Body.String(); body != expectedBody {
		t.Errorf("incorrect body\nexpected %v but got %v", expectedBody, body)
	}
}

func Test_Login_Success(t *testing.T) {
	r := setupEngine()
	s := setupServer()
	svc, mSvc := mservice.SetupMockService()
	s.service = svc
	id, err := uuid.NewRandom()
	if err != nil {
		t.Fatalf("could not generate random uuid %v", err)
	}
	mSvc.UserService.MockModel[0] = &model.User{Username: "admin", ID: id}

	r.POST("/", s.Login)

	req, err := http.NewRequest("POST", "/", body(`{"username": "admin", "password": "admin"}`))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)
	expectedStatusCode := http.StatusCreated
	if rr.Code != expectedStatusCode {
		t.Errorf("wrong status code returned\nexpected %v but got %v", expectedStatusCode, rr.Code)
	}

	expectedBody := `{"message":"successfully authenticated user"}`
	if body := rr.Body.String(); body != expectedBody {
		t.Errorf("incorrect body\nexpected %v but got %v", expectedBody, body)
	}

	cookie := strings.Trim(rr.Header().Get("Set-Cookie"), " ")
	if cookie == "" {
		t.Errorf("No header is being set for exorcist")
	}
	if !strings.Contains(cookie, "exorcist") {
		t.Errorf("cookie was not set up correctly: %v", cookie)
	}
}

func Test_Logout_InvalidSessionToken(t *testing.T) {
	r := setupEngine()
	s := setupServer()
	svc, _ := mservice.SetupMockService()
	s.service = svc

	r.GET("/", s.Logout)

	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)
	expectedStatusCode := http.StatusBadRequest
	if rr.Code != expectedStatusCode {
		t.Errorf("wrong status code returned\nexpected %v but got %v", expectedStatusCode, rr.Code)
	}
	expectedBody := `{"error":"invalid session token"}`
	if body := rr.Body.String(); body != expectedBody {
		t.Errorf("incorrect body\nexpected %v but got %v", expectedBody, body)
	}
}

func Test_Logout_Success(t *testing.T) {
	r := setupEngine()
	svc, _ := mservice.SetupMockService()
	s := setupServer()
	s.service = svc

	r.GET("/", s.Logout)

	rr := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)

	setupCookies(req, r)

	r.ServeHTTP(rr, req)
	expectedStatusCode := http.StatusOK
	if rr.Code != expectedStatusCode {
		t.Errorf("wrong status code returned\nexpected %v but got %v", expectedStatusCode, rr.Code)
	}
	expectedBody := `{"message":"successfully logged out"}`
	if body := rr.Body.String(); body != expectedBody {
		t.Errorf("incorrect body\nexpected %v but got %v", expectedBody, body)
	}
}
