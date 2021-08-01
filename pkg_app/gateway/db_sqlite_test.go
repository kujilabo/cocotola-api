package gateway

import (
	"database/sql"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4/database"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/mattn/go-sqlite3"
	gorm_logrus "github.com/onrik/gorm-logrus"
	gormSQLite "gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var testDBFile string

func openSQLiteForTest() (*gorm.DB, error) {
	return gorm.Open(gormSQLite.Open(testDBFile), &gorm.Config{
		Logger: gorm_logrus.New(),
	})
}

func initSQLite() {
	testDBFile = "./test.db"
	os.Remove(testDBFile)
	setupSQLite()
}

func setupSQLite() {
	db, err := openSQLiteForTest()
	if err != nil {
		log.Fatal(err)
	}
	setupDB(db, "sqlite3", func(sqlDB *sql.DB) (database.Driver, error) {
		return sqlite3.WithInstance(sqlDB, &sqlite3.Config{})
	})
}
