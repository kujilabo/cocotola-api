package service

import (
	"context"

	"gorm.io/gorm"

	"github.com/kujilabo/cocotola-api/src/app/domain"
	userD "github.com/kujilabo/cocotola-api/src/user/domain"
	userS "github.com/kujilabo/cocotola-api/src/user/service"
)

type RepositoryFactoryFunc func(ctx context.Context, db *gorm.DB) (RepositoryFactory, error)

var (
	appPropertiesSystemSpaceID     = userD.SpaceID(0)
	appPropertiesSystemStudentID   = userD.AppUserID(0)
	appPropertiesTatoebaWorkbookID = domain.WorkbookID(0)
	SystemStudentLoginID           = "system-student"
	TatoebaWorkbookName            = "tatoeba"
	OrganizationName               = "cocotola"
	UserRfFunc                     userS.RepositoryFactoryFunc
	RfFunc                         RepositoryFactoryFunc
)

func InitAppProperties(systemSpaceID userD.SpaceID, systemStudentID userD.AppUserID, tatoebaWorkbookID domain.WorkbookID) {
	appPropertiesSystemSpaceID = systemSpaceID
	appPropertiesSystemStudentID = systemStudentID
	appPropertiesTatoebaWorkbookID = tatoebaWorkbookID
}

func GetSystemSpaceID() userD.SpaceID {
	return appPropertiesSystemSpaceID
}
func SetSystemSpaceID(propertiesSystemSpaceID userD.SpaceID) {
	appPropertiesSystemSpaceID = propertiesSystemSpaceID
}

func GetSystemStudentID() userD.AppUserID {
	return appPropertiesSystemStudentID
}
func SetSystemStudentID(propertiesSystemStudentID userD.AppUserID) {
	appPropertiesSystemStudentID = propertiesSystemStudentID
}

func GetTatoebaWorkbookID() domain.WorkbookID {
	return appPropertiesTatoebaWorkbookID
}
func SetTatoebaWorkbookID(propertiesTatoebaWorkbookID domain.WorkbookID) {
	appPropertiesTatoebaWorkbookID = propertiesTatoebaWorkbookID
}
