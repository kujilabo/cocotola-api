package gateway

import (
	"database/sql"
	"errors"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"golang.org/x/xerrors"
	"gorm.io/gorm"
)

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

func MigrateSQLiteDB(db *gorm.DB) error {
	return migrateDB(db, "sqlite3", func(sqlDB *sql.DB) (database.Driver, error) {
		return sqlite3.WithInstance(sqlDB, &sqlite3.Config{})
	})
}
