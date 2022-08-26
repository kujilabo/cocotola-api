package gateway_test

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"os"
	"strings"
	"testing"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"

	liberrors "github.com/kujilabo/cocotola-api/src/lib/errors"
	"github.com/kujilabo/cocotola-api/src/user/domain"
	"github.com/kujilabo/cocotola-api/src/user/gateway"
	"github.com/kujilabo/cocotola-api/src/user/service"
)

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

func setupDB(db *gorm.DB, driverName string, withInstance func(sqlDB *sql.DB) (database.Driver, error)) {
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal(err)
	}
	defer sqlDB.Close()

	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	pos := strings.Index(wd, "src/user")
	dir := wd[0:pos] + "sqls/" + driverName

	driver, err := withInstance(sqlDB)
	if err != nil {
		log.Fatal(liberrors.Errorf("failed to WithInstance. err: %w", err))
	}
	m, err := migrate.NewWithDatabaseInstance("file://"+dir, driverName, driver)
	if err != nil {
		log.Fatal(liberrors.Errorf("failed to NewWithDatabaseInstance. err: %w", err))
	}

	if err := m.Up(); err != nil {
		if !errors.Is(err, migrate.ErrNoChange) {
			log.Fatal(liberrors.Errorf("failed to Up. driver:%s, err: %w", driverName, err))
		}
	}
}

func testInitOrganization(t *testing.T, db *gorm.DB) (domain.OrganizationID, service.Owner) {
	bg := context.Background()
	sysAd, err := service.NewSystemAdminFromDB(bg, db)
	assert.NoError(t, err)

	firstOwnerAddParam, err := service.NewFirstOwnerAddParameter("OWNER_ID", "OWNER_NAME", "")
	assert.NoError(t, err)
	orgAddParam, err := service.NewOrganizationAddParameter("ORG_NAME", firstOwnerAddParam)
	assert.NoError(t, err)

	// delete all organizations
	db.Exec("delete from space")
	db.Exec("delete from app_user")
	db.Exec("delete from organization")
	// db.Where("true").Delete(&spaceEntity{})
	// db.Where("true").Delete(&appUserEntity{})
	// db.Where("true").Delete(&organizationEntity{})

	orgRepo := gateway.NewOrganizationRepository(db)

	// register new organization
	orgID, err := orgRepo.AddOrganization(bg, sysAd, orgAddParam)
	assert.NoError(t, err)
	assert.Greater(t, int(uint(orgID)), 0)

	appUserRepo := gateway.NewAppUserRepository(nil, db)
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
	assert.NoError(t, err)

	spaceRepo := gateway.NewSpaceRepository(db)
	_, err = spaceRepo.AddDefaultSpace(bg, sysOwner)
	assert.NoError(t, err)
	_, err = spaceRepo.AddPersonalSpace(bg, sysOwner, firstOwner)
	assert.NoError(t, err)

	return orgID, firstOwner
}

func testNewAppUserAddParameter(t *testing.T, loginID, username string) service.AppUserAddParameter {
	p, err := service.NewAppUserAddParameter(loginID, username, []string{}, map[string]string{})
	assert.NoError(t, err)
	return p
}
