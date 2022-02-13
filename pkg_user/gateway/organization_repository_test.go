package gateway_test

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"

	"github.com/kujilabo/cocotola-api/pkg_user/domain"
	"github.com/kujilabo/cocotola-api/pkg_user/gateway"
)

func TestGetOrganization(t *testing.T) {
	// logrus.SetLevel(logrus.DebugLevel)
	bg := context.Background()

	userRfFunc := func(db *gorm.DB) (domain.RepositoryFactory, error) {
		return gateway.NewRepositoryFactory(db)
	}

	domain.InitSystemAdmin(userRfFunc)
	firstOwnerAddParam, err := domain.NewFirstOwnerAddParameter("LOGIN_ID", "USERNAME", "")
	assert.NoError(t, err)
	orgAddParam, err := domain.NewOrganizationAddParameter("ORG_NAME", firstOwnerAddParam)
	assert.NoError(t, err)
	for i, db := range dbList() {
		sysAd, err := domain.NewSystemAdminFromDB(db)
		assert.NoError(t, err)

		log.Printf("%d", i)
		sqlDB, err := db.DB()
		assert.NoError(t, err)
		defer sqlDB.Close()

		// delete all organizations
		db.Exec("delete from organization")

		orgRepo := gateway.NewOrganizationRepository(db)

		// register new organization
		orgID, err := orgRepo.AddOrganization(bg, sysAd, orgAddParam)
		assert.NoError(t, err)
		assert.Greater(t, int(uint(orgID)), 0)

		// get organization registered
		model, err := domain.NewModel(1, 1, time.Now(), time.Now(), 1, 1)
		assert.NoError(t, err)
		user, err := domain.NewAppUser(nil, model, orgID, "login_id", "username", []string{}, map[string]string{})
		assert.NoError(t, err)
		{
			org, err := orgRepo.GetOrganization(bg, user)
			assert.NoError(t, err)
			assert.Equal(t, "ORG_NAME", org.GetName())
		}

		// get organization unregistered
		otherUser, err := domain.NewAppUser(nil, model, orgID+1, "login_id", "username", []string{}, map[string]string{})
		assert.NoError(t, err)
		{
			_, err := orgRepo.GetOrganization(bg, otherUser)
			assert.Equal(t, domain.ErrOrganizationNotFound, err)
		}
	}
}

func TestFindOrganizationByName(t *testing.T) {
	// logrus.SetLevel(logrus.DebugLevel)
	bg := context.Background()

	userRfFunc := func(db *gorm.DB) (domain.RepositoryFactory, error) {
		return gateway.NewRepositoryFactory(db)
	}

	domain.InitSystemAdmin(userRfFunc)

	firstOwnerAddParam, err := domain.NewFirstOwnerAddParameter("LOGIN_ID", "USERNAME", "")
	assert.NoError(t, err)
	orgAddParam, err := domain.NewOrganizationAddParameter("ORG_NAME", firstOwnerAddParam)
	assert.NoError(t, err)
	for i, db := range dbList() {
		sysAd, err := domain.NewSystemAdminFromDB(db)
		assert.NoError(t, err)

		log.Printf("%d", i)
		sqlDB, err := db.DB()
		assert.NoError(t, err)
		defer sqlDB.Close()

		// delete all organizations
		db.Exec("delete from organization")
		// db.Where("true").Delete(&organizationEntity{})

		orgRepo := gateway.NewOrganizationRepository(db)

		// register new organization
		orgID, err := orgRepo.AddOrganization(bg, sysAd, orgAddParam)
		assert.NoError(t, err)
		assert.Greater(t, int(uint(orgID)), 0)

		// find organization registered by name
		{
			org, err := orgRepo.FindOrganizationByName(bg, sysAd, "ORG_NAME")
			assert.NoError(t, err)
			assert.Equal(t, "ORG_NAME", org.GetName())
		}

		// find organization unregistered by name
		{
			_, err := orgRepo.FindOrganizationByName(bg, sysAd, "NOT_FOUND")
			assert.Equal(t, domain.ErrOrganizationNotFound, err)
		}
	}
}

func TestAddOrganization(t *testing.T) {
	logrus.SetLevel(logrus.DebugLevel)
	bg := context.Background()

	userRfFunc := func(db *gorm.DB) (domain.RepositoryFactory, error) {
		return gateway.NewRepositoryFactory(db)
	}

	domain.InitSystemAdmin(userRfFunc)

	firstOwnerAddParam, err := domain.NewFirstOwnerAddParameter("LOGIN_ID", "USERNAME", "")
	assert.NoError(t, err)
	orgAddParam, err := domain.NewOrganizationAddParameter("ORG_NAME", firstOwnerAddParam)
	assert.NoError(t, err)
	for i, db := range dbList() {
		sysAd, err := domain.NewSystemAdminFromDB(db)
		assert.NoError(t, err)

		log.Printf("%d", i)
		sqlDB, err := db.DB()
		assert.NoError(t, err)
		defer sqlDB.Close()

		// delete all organizations
		db.Exec("delete from organization")
		// db.Where("true").Delete(&organizationEntity{})

		orgRepo := gateway.NewOrganizationRepository(db)

		// register new organization
		{
			orgID, err := orgRepo.AddOrganization(bg, sysAd, orgAddParam)
			assert.NoError(t, err)
			assert.Greater(t, int(uint(orgID)), 0)
		}

		// register new organization
		{
			_, err := orgRepo.AddOrganization(bg, sysAd, orgAddParam)
			assert.Equal(t, domain.ErrOrganizationAlreadyExists, err)
		}

	}
}
