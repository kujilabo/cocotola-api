package gateway

import (
	"database/sql"
	"errors"
	"os"

	"github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/mattn/go-sqlite3"
	"golang.org/x/xerrors"
	"gorm.io/gorm"
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

func ConvertRelationError(err error, newErr error) error {
	var mysqlErr *mysql.MySQLError
	if ok := errors.As(err, &mysqlErr); ok && mysqlErr.Number == 1452 {
		return newErr
	}

	return err
}

func migrateDB(db *gorm.DB, driverName string, withInstance func(sqlDB *sql.DB) (database.Driver, error)) error {
	sqlDB, err := db.DB()
	if err != nil {
		return xerrors.Errorf("failed to DB. err: %w", err)
	}

	wd, err := os.Getwd()
	if err != nil {
		return xerrors.Errorf("failed to Getwd. err: %w", err)
	}

	dir := wd + "/sqls/" + driverName

	driver, err := withInstance(sqlDB)
	if err != nil {
		return xerrors.Errorf("failed to withInstance. err: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance("file://"+dir, driverName, driver)
	if err != nil {
		return xerrors.Errorf("failed to NewWithDatabaseInstance. err: %w", err)
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return xerrors.Errorf("failed to Up. err: %w", err)
	}

	return nil
}
