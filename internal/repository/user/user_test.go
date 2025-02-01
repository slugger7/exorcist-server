package userRepository_test

import (
	"testing"

	"github.com/slugger7/exorcist/internal/environment"
	userRepository "github.com/slugger7/exorcist/internal/repository/user"
)

func Test_GetUserByUsernameAndPassword(t *testing.T) {
	s := &userRepository.UserRepository{
		Env: &environment.EnvironmentVariables{DebugSql: false},
	}

	actual, _ := s.GetUserByUsernameAndPassword("someUsername", "somePassword").Sql()

	exected := "\nSELECT users.id AS \"users.id\",\n     users.username AS \"users.username\"\nFROM public.users\nWHERE (users.username = $1::text) AND (users.password = $2::text);\n"
	if exected != actual {
		t.Errorf("Expected %v but got %v", exected, actual)
	}
}
