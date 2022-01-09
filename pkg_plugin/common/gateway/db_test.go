package gateway

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database"
	"gorm.io/gorm"
)

func dbList() map[string]*gorm.DB {
	dbList := make(map[string]*gorm.DB)
	m, err := openMySQLForTest()
	if err != nil {
		panic(err)
	}

	dbList["mysql"] = m

	// s, err := openSQLiteForTest()
	// if err != nil {
	// 	panic(err)
	// }
	// dbList["sqlite3"] = s

	return dbList
}

func setupDB(db *gorm.DB, driverName string, withInstance func(sqlDB *sql.DB) (database.Driver, error)) {
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal(err)
	}
	defer sqlDB.Close()

	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	pos := strings.Index(wd, "pkg_plugin")
	dir := wd[0:pos] + "sqls/" + driverName

	driver, err := withInstance(sqlDB)
	if err != nil {
		log.Fatal(fmt.Errorf("failed to WithInstance. err: %w", err))
	}
	m, err := migrate.NewWithDatabaseInstance("file://"+dir, driverName, driver)
	if err != nil {
		log.Fatal(fmt.Errorf("failed to NewWithDatabaseInstance. err: %w", err))
	}

	if err := m.Up(); err != nil {
		if !errors.Is(err, migrate.ErrNoChange) {
			log.Fatal(fmt.Errorf("failed to Up. driver:%s, err: %w", driverName, err))
		}
	}
}
