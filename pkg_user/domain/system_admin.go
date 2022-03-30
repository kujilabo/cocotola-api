package domain

import (
	"context"
	"errors"

	"golang.org/x/xerrors"
	"gorm.io/gorm"

	"github.com/kujilabo/cocotola-api/pkg_lib/log"
)

const SystemAdminID = 1

// var systemAdminInstance SystemAdmin
var rfFunc func(ctx context.Context, db *gorm.DB) (RepositoryFactory, error)

func InitSystemAdmin(rfFuncArg func(ctx context.Context, db *gorm.DB) (RepositoryFactory, error)) {
	if rfFuncArg == nil {
		panic(errors.New("invalid argment"))
	}
	rfFunc = rfFuncArg
}

type SystemAdmin interface {
	GetID() uint

	FindSystemOwnerByOrganizationID(ctx context.Context, organizationID OrganizationID) (SystemOwner, error)

	FindSystemOwnerByOrganizationName(ctx context.Context, organizationName string) (SystemOwner, error)

	FindOrganizationByName(ctx context.Context, name string) (Organization, error)

	AddOrganization(ctx context.Context, parma OrganizationAddParameter) (OrganizationID, error)
}

type systemAdmin struct {
	rf RepositoryFactory
}

// func SystemAdminInstance() SystemAdmin {
// 	return systemAdminInstance
// }

func NewSystemAdmin(rf RepositoryFactory) SystemAdmin {
	return &systemAdmin{
		rf: rf,
	}
}
func NewSystemAdminFromDB(ctx context.Context, db *gorm.DB) (SystemAdmin, error) {
	rf, err := rfFunc(ctx, db)
	if err != nil {
		return nil, err
	}
	return &systemAdmin{
		rf: rf,
	}, nil
}

func (s *systemAdmin) GetID() uint {
	return SystemAdminID
}

func (s *systemAdmin) FindSystemOwnerByOrganizationID(ctx context.Context, organizationID OrganizationID) (SystemOwner, error) {
	return s.rf.NewAppUserRepository().FindSystemOwnerByOrganizationID(ctx, s, organizationID)
}

func (s *systemAdmin) FindSystemOwnerByOrganizationName(ctx context.Context, organizationName string) (SystemOwner, error) {
	return s.rf.NewAppUserRepository().FindSystemOwnerByOrganizationName(ctx, s, organizationName)
}

func (s *systemAdmin) FindOrganizationByName(ctx context.Context, name string) (Organization, error) {
	return s.rf.NewOrganizationRepository().FindOrganizationByName(ctx, s, name)
}

func (s *systemAdmin) AddOrganization(ctx context.Context, param OrganizationAddParameter) (OrganizationID, error) {
	logger := log.FromContext(ctx)
	// add organization
	organizationID, err := s.rf.NewOrganizationRepository().AddOrganization(ctx, s, param)
	if err != nil {
		return 0, xerrors.Errorf("failed to AddOrganization. error: %w", err)
	}

	// add system owner
	systemOwnerID, err := s.rf.NewAppUserRepository().AddSystemOwner(ctx, s, organizationID)
	if err != nil {
		return 0, xerrors.Errorf("failed to AddSystemOwner. error: %w", err)
	}

	systemOwner, err := s.rf.NewAppUserRepository().FindSystemOwnerByOrganizationName(ctx, s, param.GetName())
	if err != nil {
		return 0, xerrors.Errorf("failed to FindSystemOwnerByOrganizationName. error: %w", err)
	}

	// // add system student
	// systemStudentID, err := s.rf.NewAppUserRepository().AddSystemStudent(ctx, systemOwner)
	// if err != nil {
	// 	return 0, fmt.Errorf("failed to AddSystemStudent. error: %w", err)
	// }

	// add owner
	ownerID, err := s.rf.NewAppUserRepository().AddFirstOwner(ctx, systemOwner, param.GetFirstOwner())
	if err != nil {
		return 0, xerrors.Errorf("failed to AddFirstOwner. error: %w", err)
	}

	owner, err := s.rf.NewAppUserRepository().FindOwnerByLoginID(ctx, systemOwner, param.GetFirstOwner().GetLoginID())
	if err != nil {
		return 0, xerrors.Errorf("failed to FindOwnerByLoginID. error: %w", err)
	}

	// add public group
	publicGroupID, err := s.rf.NewAppUserGroupRepository().AddPublicGroup(ctx, systemOwner)
	if err != nil {
		return 0, xerrors.Errorf("failed to AddPublicGroup. error: %w", err)
	}

	// public-group <-> owner
	if err := s.rf.NewGroupUserRepository().AddGroupUser(ctx, systemOwner, publicGroupID, ownerID); err != nil {
		return 0, xerrors.Errorf("failed to AddGroupUser. error: %w", err)
	}

	// add default space
	spaceID, err := s.rf.NewSpaceRepository().AddDefaultSpace(ctx, systemOwner)
	if err != nil {
		return 0, xerrors.Errorf("failed to AddDefaultSpace. error: %w", err)
	}

	logger.Infof("SystemOwnerID:%d, owner: %+v, spaceID: %d", systemOwnerID, owner, spaceID)
	// logger.Infof("SystemOwnerID:%d, SystemStudentID:%d, owner: %+v, spaceID: %d", systemOwnerID, systemStudentID, owner, spaceID)

	// // add personal group
	// personalGroupID, err := s.appUserGroupRepositor.AddPublicGroup(owner)
	// if err != nil {
	// 	return 0, fmt.Errorf("failed to AddPersonalGroup. error: %w", err)
	// }

	// // personal-group <-> owner
	// if err := s.groupUserRepository.AddGroupUser(systemOwner, personalGroupID, ownerID); err != nil {
	// 	return 0, fmt.Errorf("failed to AddGroupUser. error: %w", err)
	// }

	return organizationID, nil
}
