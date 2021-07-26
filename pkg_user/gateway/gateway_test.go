package gateway

import "gorm.io/gorm"

func init() {
	initMySQL()
	initSQLite()
}

func dbList() []*gorm.DB {
	dbList := make([]*gorm.DB, 0)
	m, err := openMySQLForTest()
	if err != nil {
		panic(err)
	}

	dbList = append(dbList, m)

	s, err := openSQLiteForTest()
	if err != nil {
		panic(err)
	}
	dbList = append(dbList, s)

	return dbList
}
