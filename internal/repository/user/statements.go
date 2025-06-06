package userRepository

import (
	"strings"
	"time"

	"github.com/go-jet/jet/v2/postgres"
	"github.com/google/uuid"
	"github.com/slugger7/exorcist/internal/db/exorcist/public/model"
	"github.com/slugger7/exorcist/internal/db/exorcist/public/table"
	"github.com/slugger7/exorcist/internal/repository/util"
)

func (ur *userRepository) getUserByUsernameAndPasswordStatement(username, password string) *UserStatement {
	statement := table.User.SELECT(table.User.ID, table.User.Username).
		FROM(table.User).
		WHERE(table.User.Username.EQ(postgres.String(username)).
			AND(table.User.Password.EQ(postgres.String(password))).
			AND(table.User.Active.IS_TRUE()))

	util.DebugCheck(ur.env, statement)
	return &UserStatement{statement, ur.db, ur.ctx}
}

func (ur *userRepository) getUserByUsernameStatement(username string, columns ...postgres.Projection) *UserStatement {
	if len(columns) == 0 {
		columns = []postgres.Projection{table.User.Username}
	}
	statement := table.User.SELECT(columns[0], columns[1:]...).
		FROM(table.User).
		WHERE(postgres.LOWER(table.User.Username).EQ(postgres.String(strings.ToLower(username))).
			AND(table.User.Active.IS_TRUE()))

	util.DebugCheck(ur.env, statement)
	return &UserStatement{statement, ur.db, ur.ctx}
}

func (ur *userRepository) createStatement(user model.User) *UserStatement {
	statement := table.User.INSERT(table.User.Username, table.User.Password).
		MODEL(user).
		RETURNING(table.User.ID, table.User.Username, table.User.Active, table.User.Created, table.User.Modified)

	util.DebugCheck(ur.env, statement)
	return &UserStatement{statement, ur.db, ur.ctx}
}

func (ur *userRepository) getByIdStatement(id uuid.UUID) *UserStatement {
	statement := table.User.SELECT(table.User.AllColumns).
		FROM(table.User).
		WHERE(table.User.ID.EQ(postgres.UUID(id))).
		LIMIT(1)

	util.DebugCheck(ur.env, statement)

	return &UserStatement{statement, ur.db, ur.ctx}
}

func (ur *userRepository) updatePasswordStatement(user *model.User) *UserStatement {
	user.Modified = time.Now()
	statement := table.User.UPDATE(table.User.Password, table.User.Modified).
		MODEL(user).
		WHERE(table.User.ID.EQ(postgres.UUID(user.ID)))

	util.DebugCheck(ur.env, statement)

	return &UserStatement{statement, ur.db, ur.ctx}
}
