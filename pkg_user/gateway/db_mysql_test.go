package gateway

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	gorm_logrus "github.com/onrik/gorm-logrus"
	gormMySQL "gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var testDBHost string
var testDBPort string
var testDBURL string

func openMySQLForTest() (*gorm.DB, error) {
	return gorm.Open(gormMySQL.Open(testDBURL), &gorm.Config{
		Logger: gorm_logrus.New(),
	})
}

func initMySQL() {
	testDBHost = os.Getenv("TEST_DB_HOST")
	if testDBHost == "" {
		testDBHost = "127.0.0.1"
	}

	testDBPort = os.Getenv("TEST_DB_PORT")
	if testDBPort == "" {
		testDBPort = "3307"
	}

	testDBURL = fmt.Sprintf("user:password@tcp(%s:%s)/testdb?charset=utf8&parseTime=True&loc=Asia%%2FTokyo", testDBHost, testDBPort)

	setupMySQL()
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
	pos := strings.Index(wd, "pkg")
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
		if err != migrate.ErrNoChange {
			log.Fatal(fmt.Errorf("failed to Up. err: %w", err))
		}
	}
}

func setupMySQL() {
	db, err := openMySQLForTest()
	if err != nil {
		log.Fatal(err)
	}
	setupDB(db, "mysql", func(sqlDB *sql.DB) (database.Driver, error) {
		return mysql.WithInstance(sqlDB, &mysql.Config{})
	})
}
