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
	userD "github.com/kujilabo/cocotola-api/pkg_user/domain"
	userG "github.com/kujilabo/cocotola-api/pkg_user/gateway"
	"github.com/stretchr/testify/assert"
	"golang.org/x/xerrors"
	"gorm.io/gorm"
)

func dbList() map[string]*gorm.DB {
	dbList := make(map[string]*gorm.DB)
	m, err := openMySQLForTest()
	if err != nil {
		panic(err)
	}

	dbList["mysql"] = m

	s, err := openSQLiteForTest()
	if err != nil {
		panic(err)
	}
	dbList["sqlite3"] = s

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
	pos := strings.Index(wd, "pkg_app")
	dir := wd[0:pos] + "sqls/" + driverName

	driver, err := withInstance(sqlDB)
	if err != nil {
		log.Fatal(xerrors.Errorf("failed to WithInstance. err: %w", err))
	}
	m, err := migrate.NewWithDatabaseInstance("file://"+dir, driverName, driver)
	if err != nil {
		log.Fatal(xerrors.Errorf("failed to NewWithDatabaseInstance. err: %w", err))
	}

	if err := m.Up(); err != nil {
		if !errors.Is(err, migrate.ErrNoChange) {
			log.Fatal(xerrors.Errorf("failed to Up. driver:%s, err: %w", driverName, err))
		}
	}
}

func testInitOrganization(t *testing.T, db *gorm.DB) (userD.OrganizationID, userD.SystemOwner, userD.Owner) {
	log.Println("testInitOrganization")
	bg := context.Background()
	sysAd, err := userD.NewSystemAdminFromDB(db)
	assert.NoError(t, err)

	// delete all organizations
	result := db.Debug().Session(&gorm.Session{AllowGlobalUpdate: true}).Exec("delete from space")
	assert.NoError(t, result.Error)
	result = db.Session(&gorm.Session{AllowGlobalUpdate: true}).Exec("delete from app_user")
	assert.NoError(t, result.Error)
	result = db.Session(&gorm.Session{AllowGlobalUpdate: true}).Exec("delete from organization")
	assert.NoError(t, result.Error)

	firstOwnerAddParam, err := userD.NewFirstOwnerAddParameter("OWNER_ID", "OWNER_NAME", "")
	assert.NoError(t, err)
	orgAddParam, err := userD.NewOrganizationAddParameter("ORG_NAME", firstOwnerAddParam)
	assert.NoError(t, err)
	orgRepo := userG.NewOrganizationRepository(db)

	// register new organization
	orgID, err := orgRepo.AddOrganization(bg, sysAd, orgAddParam)
	assert.NoError(t, err)
	assert.Greater(t, int(uint(orgID)), 0)
	log.Printf("OrgID: %d \n", uint(orgID))
	org, err := orgRepo.FindOrganizationByID(bg, sysAd, orgID)
	assert.NoError(t, err)
	log.Printf("OrgID: %d \n", org.GetID())

	appUserRepo := userG.NewAppUserRepository(nil, db)
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

	spaceRepo := userG.NewSpaceRepository(db)
	_, err = spaceRepo.AddDefaultSpace(bg, sysOwner)
	assert.NoError(t, err)
	_, err = spaceRepo.AddPersonalSpace(bg, sysOwner, firstOwner)
	assert.NoError(t, err)

	return orgID, sysOwner, firstOwner
}

func testNewAppUserAddParameter(t *testing.T, loginID, username string) userD.AppUserAddParameter {
	p, err := userD.NewAppUserAddParameter(loginID, username, []string{}, map[string]string{})
	assert.NoError(t, err)
	return p
}
