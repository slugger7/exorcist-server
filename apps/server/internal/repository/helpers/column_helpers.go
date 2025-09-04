package helpers

import "github.com/go-jet/jet/v2/postgres"

func OrderByDirectionColumn(asc bool, column postgres.Column, stmnt postgres.SelectStatement) postgres.SelectStatement {
	if !asc {
		return stmnt.ORDER_BY(column.DESC())
	}

	return stmnt.ORDER_BY(column.ASC())
}
