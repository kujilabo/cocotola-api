package gateway

// import (
// 	"database/sql"

// 	"github.com/golang-migrate/migrate/v4/database"
// 	"github.com/golang-migrate/migrate/v4/database/mysql"
// 	_ "github.com/golang-migrate/migrate/v4/source/file"
// 	"gorm.io/gorm"
// )

// func MigrateMySQLDB(db *gorm.DB) error {
// 	return migrateDB(db, "mysql", func(sqlDB *sql.DB) (database.Driver, error) {
// 		return mysql.WithInstance(sqlDB, &mysql.Config{})
// 	})
// }
