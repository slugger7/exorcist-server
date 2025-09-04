package server

import (
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/slugger7/exorcist/internal/assert"
	"github.com/slugger7/exorcist/internal/db/exorcist/public/model"
	"go.uber.org/mock/gomock"
)

func Test_AuthRequiredMiddleware_Fails(t *testing.T) {
	s := setupServer(t).
		withAuth()

	s.authGroup.GET("/", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	rr := s.withAuthGetRequest("").
		exec()

	assert.StatusCode(t, http.StatusUnauthorized, rr.Code)
	assert.Body(t, errBody(ErrUnauthorized), rr.Body.String())

}

func Test_AuthRequiredMiddleware_Success(t *testing.T) {
	s := setupServer(t).
		withAuth()

	expectedStatusCode := http.StatusOK
	id, _ := uuid.NewRandom()

	s.authGroup.GET("/", func(ctx *gin.Context) {
		ctx.JSON(expectedStatusCode, gin.H{"message": "success"})
	})

	rr := s.withAuthGetRequest("").
		withCookie(TestCookie{Value: id}).
		exec()

	assert.StatusCode(t, expectedStatusCode, rr.Code)
	assert.Body(t, `{"message":"success"}`, rr.Body.String())
}

func Test_Login_InvalidBody(t *testing.T) {
	s := setupServer(t)

	s.server.withAuthLogin(&s.engine.RouterGroup, "/")
	rr := s.withPostRequest(body(`{some invalid body}`)).
		exec()

	assert.StatusCode(t, http.StatusUnprocessableEntity, rr.Code)
}

func Test_Login_ErrFromValidateUser(t *testing.T) {
	s := setupServer(t).withUserService()

	m := LoginModel{
		Username: "someUsername",
		Password: "somePassword",
	}
	s.mockUserService.EXPECT().
		Validate(gomock.Eq(m.Username), gomock.Eq(m.Password)).
		DoAndReturn(func(string, string) (*model.User, error) {
			return nil, fmt.Errorf("some error")
		}).
		Times(1)

	s.server.withAuthLogin(&s.engine.RouterGroup, "/")
	rr := s.withPostRequest(bodyM(m)).exec()

	assert.StatusCode(t, http.StatusUnauthorized, rr.Code)
	assert.Body(t, errBody(ErrUnauthorized), rr.Body.String())
}

func Test_Login_Success(t *testing.T) {
	s := setupServer(t).withUserService()

	id, _ := uuid.NewRandom()
	u := &model.User{Username: "admin", ID: id}

	l := LoginModel{
		Username: u.Username,
		Password: "some password",
	}

	s.mockUserService.EXPECT().
		Validate(gomock.Eq(l.Username), gomock.Eq(l.Password)).
		DoAndReturn(func(string, string) (*model.User, error) {
			return u, nil
		})

	s.server.withAuthLogin(&s.engine.RouterGroup, "/")
	rr := s.withPostRequest(bodyM(l)).exec()

	assert.StatusCode(t, http.StatusCreated, rr.Code)
	assert.Body(t, fmt.Sprintf(`{"userId":"%v","username":"%v"}`, id, u.Username), rr.Body.String())

	cookie := strings.Trim(rr.Header().Get("Set-Cookie"), " ")

	if cookie == "" {
		t.Errorf("No header is being set for exorcist")
	}
	if !strings.Contains(cookie, "exorcist") {
		t.Errorf("cookie was not set up correctly: %v", cookie)
	}
}

func Test_Logout_InvalidSessionToken(t *testing.T) {
	s := setupServer(t)

	s.server.withAuthLogout(&s.engine.RouterGroup, "/")
	rr := s.withGetRequest("").exec()

	assert.StatusCode(t, http.StatusBadRequest, rr.Code)
	assert.Body(t, errBody(ErrInvalidSessionToken), rr.Body.String())
}

func Test_Logout_Success(t *testing.T) {
	s := setupServer(t).
		withAuth()

	id, _ := uuid.NewRandom()

	s.server.withAuthLogout(&s.engine.RouterGroup, AUTH_ROUTE+"/logout")
	rr := s.withAuthGetRequest("logout").
		withCookie(TestCookie{Value: id}).
		exec()

	assert.StatusCode(t, http.StatusOK, rr.Code)
	assert.Body(t, msgBody(MsgLoggedOut), rr.Body.String())
}
