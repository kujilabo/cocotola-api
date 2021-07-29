package gateway

import (
	"github.com/go-sql-driver/mysql"
	"github.com/mattn/go-sqlite3"
)

func convertDuplicatedError(err error, newErr error) error {
	if dbErr, ok := err.(*mysql.MySQLError); ok && dbErr.Number == 1062 {
		return newErr
	}
	if dbErr, ok := err.(sqlite3.Error); ok && int(dbErr.ExtendedCode) == 2067 {
		return newErr
	}
	return err
}
