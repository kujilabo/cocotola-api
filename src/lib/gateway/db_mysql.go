package gateway

import (
	"database/sql"
	"fmt"

	"github.com/golang-migrate/migrate/v4/database"
	migrate_mysql "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	gorm_logrus "github.com/onrik/gorm-logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func OpenMySQL(username, password, host string, port int, database string) (*gorm.DB, error) {
	// dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&multiStatements=true&&collation=utf8mb4_bin",
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=True&multiStatements=true&&collation=utf8mb4_bin", username, password, host, port, database)
	return gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: gorm_logrus.New(),
	})
}

func MigrateMySQLDB(db *gorm.DB) error {
	return migrateDB(db, "mysql", func(sqlDB *sql.DB) (database.Driver, error) {
		return migrate_mysql.WithInstance(sqlDB, &migrate_mysql.Config{})
	})
}
