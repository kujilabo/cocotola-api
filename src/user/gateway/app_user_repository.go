package gateway

import (
	"context"
	"errors"
	"time"

	"golang.org/x/xerrors"
	"gorm.io/gorm"

	libG "github.com/kujilabo/cocotola-api/src/lib/gateway"
	"github.com/kujilabo/cocotola-api/src/lib/passwordhelper"
	"github.com/kujilabo/cocotola-api/src/user/domain"
	"github.com/kujilabo/cocotola-api/src/user/service"
)

var (
	AppUserTableName = "app_user"

	SystemOwnerLoginID   = "system-owner"
	SystemStudentLoginID = "system-student"
	GuestLoginID         = "guest"

	AdministratorRole = "Administrator"
	OwnerRole         = "Owner"
	ManagerRole       = "Manager"
	UserRole          = "User"
	GuestRole         = "Guest"
	UnknownRole       = "Unknown"
)

type appUserRepository struct {
	rf service.RepositoryFactory
	db *gorm.DB
}

type appUserEntity struct {
	ID                   uint
	Version              int
	CreatedAt            time.Time
	UpdatedAt            time.Time
	CreatedBy            uint
	UpdatedBy            uint
	OrganizationID       uint
	LoginID              string
	Username             string
	HashedPassword       string
	Role                 string
	Provider             string
	ProviderID           string
	ProviderAccessToken  string
	ProviderRefreshToken string
	Removed              bool
}

func (e *appUserEntity) TableName() string {
	return AppUserTableName
}

// func toRole(role string) domain.Role {
// 	if role == "administrator" {
// 		return domain.AdministratorRole
// 	} else if role == OwnerRole {
// 		return domain.OwnerRole
// 	} else if role == "Manager" {
// 		return domain.ManagerRole
// 	} else if role == "User" {
// 		return domain.UserRole
// 	} else if role == "Guest" {
// 		return domain.GuestRole
// 	}
// 	return domain.UnknownRole
// }

// func fromRoleToString(role domain.Role) string {
// 	switch role {
// 	case domain.AdministratorRole:
// 		return AdministratorRole
// 	case domain.OwnerRole:
// 		return OwnerRole
// 	case domain.ManagerRole:
// 		return ManagerRole
// 	case domain.UserRole:
// 		return UserRole
// 	case domain.GuestRole:
// 		return GuestRole
// 	default:
// 		return UnknownRole
// 	}
// }

func (e *appUserEntity) toSystemOwner(rf service.RepositoryFactory) (service.SystemOwner, error) {
	if e.LoginID != SystemOwnerLoginID {
		return nil, xerrors.Errorf("invalid system owner. loginID: %s", e.LoginID)
	}

	model, err := domain.NewModel(e.ID, e.Version, e.CreatedAt, e.UpdatedAt, e.CreatedBy, e.UpdatedBy)
	if err != nil {
		return nil, err
	}

	appUserModel, err := domain.NewAppUserModel(model, domain.OrganizationID(e.OrganizationID), e.LoginID, e.Username, []string{"SystemOwner"}, map[string]string{})
	if err != nil {
		return nil, err
	}

	appUser, err := service.NewAppUser(rf, appUserModel)
	if err != nil {
		return nil, err
	}

	return service.NewSystemOwner(rf, appUser)
}

func (e *appUserEntity) toAppUser(rf service.RepositoryFactory, roles []string, properties map[string]string) (service.AppUser, error) {
	model, err := domain.NewModel(e.ID, e.Version, e.CreatedAt, e.UpdatedAt, e.CreatedBy, e.UpdatedBy)
	if err != nil {
		return nil, err
	}

	appUserModel, err := domain.NewAppUserModel(model, domain.OrganizationID(e.OrganizationID), e.LoginID, e.Username, roles, properties)
	if err != nil {
		return nil, err
	}

	appUser, err := service.NewAppUser(rf, appUserModel)
	if err != nil {
		return nil, err
	}

	return service.NewSystemOwner(rf, appUser)
}

func (e *appUserEntity) toOwner(rf service.RepositoryFactory, roles []string, properties map[string]string) (service.Owner, error) {
	model, err := domain.NewModel(e.ID, e.Version, e.CreatedAt, e.UpdatedAt, e.CreatedBy, e.UpdatedBy)
	if err != nil {
		return nil, err
	}

	appUserModel, err := domain.NewAppUserModel(model, domain.OrganizationID(e.OrganizationID), e.LoginID, e.Username, roles, properties)
	if err != nil {
		return nil, err
	}

	appUser, err := service.NewAppUser(rf, appUserModel)
	if err != nil {
		return nil, err
	}

	return service.NewOwner(rf, appUser), nil
}

func NewAppUserRepository(rf service.RepositoryFactory, db *gorm.DB) service.AppUserRepository {
	return &appUserRepository{
		rf: rf,
		db: db,
	}
}

func (r *appUserRepository) FindSystemOwnerByOrganizationID(ctx context.Context, operator domain.SystemAdminModel, organizationID domain.OrganizationID) (service.SystemOwner, error) {
	_, span := tracer.Start(ctx, "appUserRepository.FindSystemOwnerByOrganizationID")
	defer span.End()

	appUser := appUserEntity{}
	if result := r.db.Where("organization_id = ? and removed = 0", organizationID).
		Where("login_id = ?", SystemOwnerLoginID).
		First(&appUser); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, xerrors.Errorf("system owner not found. organization ID: %d, err: %w", organizationID, service.ErrSystemOwnerNotFound)
		}
		return nil, result.Error
	}
	return appUser.toSystemOwner(r.rf)
}

func (r *appUserRepository) FindSystemOwnerByOrganizationName(ctx context.Context, operator domain.SystemAdminModel, organizationName string) (service.SystemOwner, error) {
	_, span := tracer.Start(ctx, "appUserRepository.FindSystemOwnerByOrganizationName")
	defer span.End()

	appUser := appUserEntity{}
	if result := r.db.Table("organization").Select("app_user.*").
		Where("organization.name = ? and app_user.removed = 0", organizationName).
		Where("login_id = ?", SystemOwnerLoginID).
		Joins("inner join app_user on organization.id = app_user.organization_id").
		First(&appUser); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, xerrors.Errorf("system owner not found. organization name: %s, err: %w", organizationName, service.ErrSystemOwnerNotFound)
		}

		return nil, result.Error
	}
	return appUser.toSystemOwner(r.rf)
}

func (r *appUserRepository) FindAppUserByID(ctx context.Context, operator domain.AppUserModel, id domain.AppUserID) (service.AppUser, error) {
	_, span := tracer.Start(ctx, "appUserRepository.FindAppUserByID")
	defer span.End()

	appUser := appUserEntity{}
	if result := r.db.Where(&appUserEntity{
		OrganizationID: uint(operator.GetOrganizationID()),
		ID:             uint(id),
	}).First(&appUser); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, service.ErrAppUserNotFound
		}

		return nil, result.Error
	}

	roles := []string{appUser.Role}
	properties := map[string]string{}

	return appUser.toAppUser(r.rf, roles, properties)
}

func (r *appUserRepository) FindAppUserByLoginID(ctx context.Context, operator domain.AppUserModel, loginID string) (service.AppUser, error) {
	_, span := tracer.Start(ctx, "appUserRepository.FindAppUserByLoginID")
	defer span.End()

	if loginID == "" {
		return nil, errors.New("invalid parameter")
	}
	appUser := appUserEntity{}
	if result := r.db.Where(&appUserEntity{
		OrganizationID: uint(operator.GetOrganizationID()),
		LoginID:        loginID,
	}).First(&appUser); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, service.ErrAppUserNotFound
		}

		return nil, result.Error
	}

	roles := []string{appUser.Role}
	properties := map[string]string{}

	return appUser.toAppUser(r.rf, roles, properties)
}

func (r *appUserRepository) FindOwnerByLoginID(ctx context.Context, operator domain.SystemOwnerModel, loginID string) (service.Owner, error) {
	_, span := tracer.Start(ctx, "appUserRepository.FindOwnerByLoginID")
	defer span.End()

	appUser := appUserEntity{}
	if result := r.db.Where(&appUserEntity{
		OrganizationID: uint(operator.GetOrganizationID()),
		LoginID:        loginID,
	}).First(&appUser); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, service.ErrAppUserNotFound
		}

		return nil, result.Error
	}

	roles := []string{appUser.Role}
	properties := map[string]string{}

	return appUser.toOwner(r.rf, roles, properties)
}

func (r *appUserRepository) addAppUser(ctx context.Context, appUserEntity *appUserEntity) (domain.AppUserID, error) {
	if result := r.db.Create(appUserEntity); result.Error != nil {
		return 0, libG.ConvertDuplicatedError(result.Error, service.ErrAppUserAlreadyExists)
	}
	return domain.AppUserID(appUserEntity.ID), nil
}

func (r *appUserRepository) AddAppUser(ctx context.Context, operator domain.OwnerModel, param service.AppUserAddParameter) (domain.AppUserID, error) {
	_, span := tracer.Start(ctx, "appUserRepository.AddAppUser")
	defer span.End()

	hashedPassword := ""
	password, ok := param.GetProperties()["password"]
	if ok {
		hashed, err := passwordhelper.HashPassword(password)
		if err != nil {
			return 0, err
		}

		hashedPassword = hashed
	}

	appUserEntity := appUserEntity{
		Version:        1,
		CreatedBy:      operator.GetID(),
		UpdatedBy:      operator.GetID(),
		OrganizationID: uint(operator.GetOrganizationID()),
		LoginID:        param.GetLoginID(),
		Username:       param.GetUsername(),
		HashedPassword: hashedPassword,
		Role:           UserRole,
	}
	return r.addAppUser(ctx, &appUserEntity)
}

func (r *appUserRepository) AddSystemOwner(ctx context.Context, operator domain.SystemAdminModel, organizationID domain.OrganizationID) (domain.AppUserID, error) {
	_, span := tracer.Start(ctx, "appUserRepository.AddSystemOwner")
	defer span.End()

	appUserEntity := appUserEntity{
		Version:        1,
		CreatedBy:      operator.GetID(),
		UpdatedBy:      operator.GetID(),
		OrganizationID: uint(organizationID),
		LoginID:        SystemOwnerLoginID,
		Username:       "SystemOwner",
		Role:           OwnerRole,
	}
	return r.addAppUser(ctx, &appUserEntity)
}

func (r *appUserRepository) AddFirstOwner(ctx context.Context, operator domain.SystemOwnerModel, param service.FirstOwnerAddParameter) (domain.AppUserID, error) {
	_, span := tracer.Start(ctx, "appUserRepository.AddFirstOwner")
	defer span.End()

	hashedPassword, err := passwordhelper.HashPassword(param.GetPassword())
	if err != nil {
		return 0, err
	}

	appUserEntity := appUserEntity{
		Version:        1,
		CreatedBy:      operator.GetID(),
		UpdatedBy:      operator.GetID(),
		OrganizationID: uint(operator.GetOrganizationID()),
		LoginID:        param.GetLoginID(),
		Username:       param.GetUsername(),
		HashedPassword: hashedPassword,
		Role:           OwnerRole,
	}
	return r.addAppUser(ctx, &appUserEntity)
}
