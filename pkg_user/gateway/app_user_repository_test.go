package gateway

// import (
// 	"context"
// 	"testing"
// 	"time"

// 	"github.com/stretchr/testify/assert"

// 	user "github.com/kujilabo/cocotola-api/pkg_user/domain"
// )

// func TestAddUser(t *testing.T) {
// 	for _, db := range dbList() {
// 		bg := context.Background()
// 		sqlDB, err := db.DB()
// 		assert.NoError(t, err)
// 		defer sqlDB.Close()

// 		model := user.NewModel(1, 1, time.Now(), time.Now(), 1, 1)
// 		appUser, err := user.NewAppUser(nil, model, user.OrganizationID(1), "loginid", "username", []string{}, map[string]string{})
// 		assert.NoError(t, err)

// 		db.Debug().Where("id <> ?", 1).Delete(&appUserEntity{})
// 		repo := NewAppUserRepository(nil, db)
// 		_, err = repo.FindAppUserByID(bg, appUser, user.AppUserID(1))
// 		assert.NoError(t, err)
// 		// db.Delete(&organizationEntity{})

// 		// organizationID, err := initialize(db)
// 		// assert.NoError(t, err)
// 		// assert.Greater(t, organizationID, uint(0))
// 	}
// }
