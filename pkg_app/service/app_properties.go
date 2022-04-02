package service

import (
	"context"

	"gorm.io/gorm"

	"github.com/kujilabo/cocotola-api/pkg_app/domain"
	user "github.com/kujilabo/cocotola-api/pkg_user/domain"
	userS "github.com/kujilabo/cocotola-api/pkg_user/service"
)

type RepositoryFactoryFunc func(ctx context.Context, db *gorm.DB) (RepositoryFactory, error)

var (
	appPropertiesSystemSpaceID     = user.SpaceID(0)
	appPropertiesSystemStudentID   = user.AppUserID(0)
	appPropertiesTatoebaWorkbookID = domain.WorkbookID(0)
	SystemStudentLoginID           = "system-student"
	TatoebaWorkbookName            = "tatoeba"
	OrganizationName               = "cocotola"
	UserRfFunc                     userS.RepositoryFactoryFunc
	RfFunc                         RepositoryFactoryFunc
)

func InitAppProperties(systemSpaceID user.SpaceID, systemStudentID user.AppUserID, tatoebaWorkbookID domain.WorkbookID) {
	appPropertiesSystemSpaceID = systemSpaceID
	appPropertiesSystemStudentID = systemStudentID
	appPropertiesTatoebaWorkbookID = tatoebaWorkbookID
}

func GetSystemSpaceID() user.SpaceID {
	return appPropertiesSystemSpaceID
}
func SetSystemSpaceID(propertiesSystemSpaceID user.SpaceID) {
	appPropertiesSystemSpaceID = propertiesSystemSpaceID
}

func GetSystemStudentID() user.AppUserID {
	return appPropertiesSystemStudentID
}
func SetSystemStudentID(propertiesSystemStudentID user.AppUserID) {
	appPropertiesSystemStudentID = propertiesSystemStudentID
}

func GetTatoebaWorkbookID() domain.WorkbookID {
	return appPropertiesTatoebaWorkbookID
}
func SetTatoebaWorkbookID(propertiesTatoebaWorkbookID domain.WorkbookID) {
	appPropertiesTatoebaWorkbookID = propertiesTatoebaWorkbookID
}
