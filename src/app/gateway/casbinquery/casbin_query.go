package casbinquery

import (
	"errors"
	"fmt"

	"gorm.io/gorm"
)

const mysqlObjectSelectSQL = `
SELECT SUBSTRING_INDEX(tp.v1, '_', -1) AS %s 
FROM casbin_rule tg
INNER JOIN casbin_rule tp ON tg.v1 = tp.v0
WHERE tg.v0 = ?
AND tg.ptype = 'g'
AND tp.ptype = 'p'
AND tp.v2 = ?
AND tp.v1 LIKE ?

UNION

SELECT SUBSTRING_INDEX(tp.v1, '_', -1) AS %s 
FROM casbin_rule tp
WHERE tp.v0 = ?
AND tp.ptype = 'p'
AND tp.v2 = ?
AND tp.v1 LIKE ?
`

const sqlite3ObjectSelectSQL = `
SELECT SUBSTR(tp.v1, INSTR(tp.v1, '_') + 1) AS %s 
FROM casbin_rule tg
INNER JOIN casbin_rule tp ON tg.v1 = tp.v0
WHERE tg.v0 = ?
AND tg.ptype = 'g'
AND tp.ptype = 'p'
AND tp.v2 = ?
AND tp.v1 LIKE ?

UNION

SELECT SUBSTR(tp.v1, INSTR(tp.v1, '_') + 1) AS %s 
FROM casbin_rule tp
WHERE tp.v0 = ?
AND tp.ptype = 'p'
AND tp.v2 = ?
AND tp.v1 LIKE ?
`

func QueryObject(db *gorm.DB, driverName, objectPrefix, columnName, subject, action string) (*gorm.DB, error) {
	if db == nil {
		return nil, errors.New("invalid argument")
	}
	objectgKeyword := objectPrefix + "%"

	var objectSelectSQL string
	switch driverName {
	case "mysql":
		objectSelectSQL = mysqlObjectSelectSQL
	case "sqlite3":
		objectSelectSQL = sqlite3ObjectSelectSQL
	default:
		return nil, errors.New("invalid driver name")
	}

	sql := fmt.Sprintf(objectSelectSQL, columnName, columnName)

	return db.Raw(sql, subject, action, objectgKeyword, subject, action, objectgKeyword), nil
}

const mysqlObjectFindSQL = `
SELECT SUBSTRING_INDEX(tp.v1, '_', -1) AS %s 
FROM casbin_rule tg
INNER JOIN casbin_rule tp ON tg.v1 = tp.v0
WHERE tg.v0 = ?
AND tg.ptype = 'g'
AND tp.ptype = 'p'
AND tp.v2 = ?
AND tp.v1 = ?

UNION

SELECT SUBSTRING_INDEX(tp.v1, '_', -1) AS %s 
FROM casbin_rule tp
WHERE tp.v0 = ?
AND tp.ptype = 'p'
AND tp.v2 = ?
AND tp.v1 = ?
`

const sqlite3ObjectFindSQL = `
SELECT SUBSTR(tp.v1, INSTR(tp.v1, '_') + 1) AS %s 
FROM casbin_rule tg
INNER JOIN casbin_rule tp ON tg.v1 = tp.v0
WHERE tg.v0 = ?
AND tg.ptype = 'g'
AND tp.ptype = 'p'
AND tp.v2 = ?
AND tp.v1 = ?

UNION

SELECT SUBSTR(tp.v1, INSTR(tp.v1, '_') + 1) AS %s 
FROM casbin_rule tp
WHERE tp.v0 = ?
AND tp.ptype = 'p'
AND tp.v2 = ?
AND tp.v1 = ?
`

func FindObject(db *gorm.DB, driverName, object, columnName, subject, action string) (*gorm.DB, error) {
	if db == nil {
		return nil, errors.New("invalid argument")
	}

	var objectSelectSQL string
	switch driverName {
	case "mysql":
		objectSelectSQL = mysqlObjectFindSQL
	case "sqlite3":
		objectSelectSQL = sqlite3ObjectFindSQL
	default:
		return nil, fmt.Errorf("invalid driver name. driver: %s", driverName)
	}

	sql := fmt.Sprintf(objectSelectSQL, columnName, columnName)

	return db.Raw(sql, subject, action, object, subject, action, object), nil
}
