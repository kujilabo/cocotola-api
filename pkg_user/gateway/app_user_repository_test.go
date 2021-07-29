package gateway

import (
	"context"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/xerrors"
	"gorm.io/gorm"

	"github.com/kujilabo/cocotola-api/pkg_user/domain"
)

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

// func Test_appUserRepository_addAppUser(t *testing.T) {
// 	type fields struct {
// 		rf domain.RepositoryFactory
// 		db *gorm.DB
// 	}
// 	type args struct {
// 		ctx           context.Context
// 		appUserEntity *appUserEntity
// 	}
// 	tests := []struct {
// 		name    string
// 		fields  fields
// 		args    args
// 		want    domain.AppUserID
// 		wantErr bool
// 	}{
// 		// TODO: Add test cases.
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			r := &appUserRepository{
// 				rf: tt.fields.rf,
// 				db: tt.fields.db,
// 			}
// 			got, err := r.addAppUser(tt.args.ctx, tt.args.appUserEntity)
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("appUserRepository.addAppUser() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}
// 			if !reflect.DeepEqual(got, tt.want) {
// 				t.Errorf("appUserRepository.addAppUser() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }

func testInitOrganization(t *testing.T, db *gorm.DB) (domain.OrganizationID, domain.Owner) {
	bg := context.Background()
	sysAd := domain.SystemAdminInstance()

	firstOwnerAddParam, err := domain.NewFirstOwnerAddParameter("OWNER_ID", "OWNER_NAME", "")
	assert.NoError(t, err)
	orgAddParam, err := domain.NewOrganizationAddParameter("ORG_NAME", firstOwnerAddParam)
	assert.NoError(t, err)

	// delete all organizations
	db.Where("true").Delete(&appUserEntity{})
	db.Where("true").Delete(&organizationEntity{})

	orgRepo := NewOrganizationRepository(db)

	// register new organization
	orgID, err := orgRepo.AddOrganization(bg, sysAd, orgAddParam)
	assert.NoError(t, err)
	assert.Greater(t, int(uint(orgID)), 0)

	appUserRepo := NewAppUserRepository(nil, db)
	sysOwnerID, err := appUserRepo.AddSystemOwner(bg, sysAd, orgID)
	assert.NoError(t, err)
	assert.Greater(t, int(uint(sysOwnerID)), 0)

	sysOwner, err := appUserRepo.FindSystemOwnerByOrganizationName(bg, sysAd, "ORG_NAME")
	assert.NoError(t, err)
	assert.Greater(t, int(uint(sysOwnerID)), 0)

	firstOwnerID, err := appUserRepo.AddFirstOwner(bg, sysOwner, firstOwnerAddParam)
	assert.NoError(t, err)
	assert.Greater(t, int(uint(firstOwnerID)), 0)

	firstOwner, err := appUserRepo.FindOwnerByLoginID(bg, sysOwner, "OWNER_ID")

	return orgID, firstOwner
}

func testNewAppUserAddParameter(t *testing.T, loginID, username string) domain.AppUserAddParameter {
	p, err := domain.NewAppUserAddParameter(loginID, username, []string{}, map[string]string{})
	assert.NoError(t, err)
	return p
}

func Test_appUserRepository_AddAppUser(t *testing.T) {
	// logrus.SetLevel(logrus.DebugLevel)
	bg := context.Background()

	domain.InitSystemAdmin(nil)
	for i, db := range dbList() {
		log.Printf("%d", i)
		sqlDB, err := db.DB()
		assert.NoError(t, err)
		defer sqlDB.Close()

		_, owner := testInitOrganization(t, db)

		type args struct {
			operator domain.Owner
			param    domain.AppUserAddParameter
		}
		tests := []struct {
			name string
			args args
			err  error
		}{
			{
				name: "success",
				args: args{
					operator: owner,
					param:    testNewAppUserAddParameter(t, "LOGIN_ID", "USERNAME"),
				},
				err: nil,
			},
			{
				name: "duplicated",
				args: args{
					operator: owner,
					param:    testNewAppUserAddParameter(t, "LOGIN_ID", "USERNAME"),
				},
				err: domain.ErrAppUserAlreadyExists,
			},
		}
		appUserRepo := NewAppUserRepository(nil, db)
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				got, err := appUserRepo.AddAppUser(bg, tt.args.operator, tt.args.param)
				if err != nil && !xerrors.Is(err, tt.err) {
					t.Errorf("AddAppUser() error = %v, wantErr %v", err, tt.err)
					return
				}
				if err == nil {
					assert.Greater(t, uint(got), uint(0))
				}
			})
		}
	}
}
