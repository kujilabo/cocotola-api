package gateway

// import (
// 	"fmt"

// 	"github.com/onrik/gorm-logrus"
// 	"gorm.io/driver/mysql"
// 	"gorm.io/gorm"
// )

// func OpenMySQL(username, password, host string, port int, database string) (*gorm.DB, error) {
// 	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Asia%%2FTokyo", username, password, host, port, database)
// 	return gorm.Open(mysql.Open(dsn), &gorm.Config{
// 		Logger: gorm_logrus.New(),
// 	})
// }
