package gateway_test

import (
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func Test_xxx_FindByText(t *testing.T) {
	logrus.SetLevel(logrus.DebugLevel)

	for driverName, db := range dbList() {
		logrus.Println(driverName)
		sqlDB, err := db.DB()
		assert.NoError(t, err)
		defer sqlDB.Close()
	}
}
