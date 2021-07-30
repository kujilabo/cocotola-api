package gateway

import (
	"errors"

	"github.com/go-sql-driver/mysql"
	"github.com/mattn/go-sqlite3"
)

func ConvertDuplicatedError(err error, newErr error) error {
	var mysqlErr *mysql.MySQLError
	if ok := errors.As(err, &mysqlErr); ok && mysqlErr.Number == 1062 {
		return newErr
	}

	var sqlite3Err sqlite3.Error
	if ok := errors.As(err, &sqlite3Err); ok && int(sqlite3Err.ExtendedCode) == 2067 {
		return newErr
	}

	return err
}
