package gateway_test

import (
	"database/sql"
	"fmt"
	"log"
	"os"

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

func setupMySQL() {
	db, err := openMySQLForTest()
	if err != nil {
		log.Fatal(err)
	}
	setupDB(db, "mysql", func(sqlDB *sql.DB) (database.Driver, error) {
		return mysql.WithInstance(sqlDB, &mysql.Config{})
	})
}
