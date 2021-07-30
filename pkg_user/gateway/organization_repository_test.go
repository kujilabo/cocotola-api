package gateway

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"

	"github.com/kujilabo/cocotola-api/pkg_user/domain"
)

func TestGetOrganization(t *testing.T) {
	// logrus.SetLevel(logrus.DebugLevel)
	bg := context.Background()

	domain.InitSystemAdmin(nil)
	sysAd := domain.SystemAdminInstance()
	firstOwnerAddParam, err := domain.NewFirstOwnerAddParameter("LOGIN_ID", "USERNAME", "")
	assert.NoError(t, err)
	orgAddParam, err := domain.NewOrganizationAddParameter("ORG_NAME", firstOwnerAddParam)
	assert.NoError(t, err)
	for i, db := range dbList() {
		log.Printf("%d", i)
		sqlDB, err := db.DB()
		assert.NoError(t, err)
		defer sqlDB.Close()

		// delete all organizations
		db.Where("true").Delete(&organizationEntity{})

		orgRepo := NewOrganizationRepository(db)

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

	domain.InitSystemAdmin(nil)
	sysAd := domain.SystemAdminInstance()
	firstOwnerAddParam, err := domain.NewFirstOwnerAddParameter("LOGIN_ID", "USERNAME", "")
	assert.NoError(t, err)
	orgAddParam, err := domain.NewOrganizationAddParameter("ORG_NAME", firstOwnerAddParam)
	assert.NoError(t, err)
	for i, db := range dbList() {
		log.Printf("%d", i)
		sqlDB, err := db.DB()
		assert.NoError(t, err)
		defer sqlDB.Close()

		// delete all organizations
		db.Where("true").Delete(&organizationEntity{})

		orgRepo := NewOrganizationRepository(db)

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

	domain.InitSystemAdmin(nil)
	sysAd := domain.SystemAdminInstance()
	firstOwnerAddParam, err := domain.NewFirstOwnerAddParameter("LOGIN_ID", "USERNAME", "")
	assert.NoError(t, err)
	orgAddParam, err := domain.NewOrganizationAddParameter("ORG_NAME", firstOwnerAddParam)
	assert.NoError(t, err)
	for i, db := range dbList() {
		log.Printf("%d", i)
		sqlDB, err := db.DB()
		assert.NoError(t, err)
		defer sqlDB.Close()

		// delete all organizations
		db.Where("true").Delete(&organizationEntity{})

		orgRepo := NewOrganizationRepository(db)

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
