package domain

import (
	user "github.com/kujilabo/cocotola-api/pkg_user/domain"
	"gorm.io/gorm"
)

var (
	appPropertiesSystemSpaceID     = user.SpaceID(0)
	appPropertiesSystemStudentID   = user.AppUserID(0)
	appPropertiesTatoebaWorkbookID = WorkbookID(0)
	SystemStudentLoginID           = "system-student"
	TatoebaWorkbookName            = "tatoeba"
	OrganizationName               = "cocotola"
	UserRfFunc                     func(db *gorm.DB) (user.RepositoryFactory, error)
	RfFunc                         func(db *gorm.DB) (RepositoryFactory, error)
)

func InitAppProperties(systemSpaceID user.SpaceID, systemStudentID user.AppUserID, tatoebaWorkbookID WorkbookID) {
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

func GetTatoebaWorkbookID() WorkbookID {
	return appPropertiesTatoebaWorkbookID
}
func SetTatoebaWorkbookID(propertiesTatoebaWorkbookID WorkbookID) {
	appPropertiesTatoebaWorkbookID = propertiesTatoebaWorkbookID
}
