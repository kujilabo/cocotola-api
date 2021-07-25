package domain

import (
	"context"
	"fmt"

	"github.com/kujilabo/cocotola-api/pkg_lib/log"
)

const SystemAdminID = 1

var systemAdminInstance SystemAdmin

func InitSystemAdmin(rf RepositoryFactory) {
	systemAdminInstance = &systemAdmin{
		rf: rf,
	}
}

type SystemAdmin interface {
	ID() uint

	FindSystemOwnerByOrganizationID(ctx context.Context, organizationID OrganizationID) (SystemOwner, error)

	FindSystemOwnerByOrganizationName(ctx context.Context, organizationName string) (SystemOwner, error)

	FindOrganizationByName(ctx context.Context, name string) (Organization, error)

	AddOrganization(ctx context.Context, parma *OrganizationAddParameter) (OrganizationID, error)
}

type systemAdmin struct {
	rf RepositoryFactory
}

func SystemAdminInstance() SystemAdmin {
	return systemAdminInstance
}

func (s *systemAdmin) ID() uint {
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

func (s *systemAdmin) AddOrganization(ctx context.Context, param *OrganizationAddParameter) (OrganizationID, error) {
	logger := log.FromContext(ctx)
	// add organization
	organizationID, err := s.rf.NewOrganizationRepository().AddOrganization(ctx, s, param)
	if err != nil {
		return 0, fmt.Errorf("failed to AddOrganization. error: %w", err)
	}

	// add system owner
	systemOwnerID, err := s.rf.NewAppUserRepository().AddSystemOwner(ctx, s, organizationID)
	if err != nil {
		return 0, fmt.Errorf("failed to AddSystemOwner. error: %w", err)
	}

	systemOwner, err := s.rf.NewAppUserRepository().FindSystemOwnerByOrganizationName(ctx, s, param.Name)
	if err != nil {
		return 0, fmt.Errorf("failed to FindSystemOwnerByOrganizationName. error: %w", err)
	}

	// // add system student
	// systemStudentID, err := s.rf.NewAppUserRepository().AddSystemStudent(ctx, systemOwner)
	// if err != nil {
	// 	return 0, fmt.Errorf("failed to AddSystemStudent. error: %w", err)
	// }

	// add owner
	ownerID, err := s.rf.NewAppUserRepository().AddFirstOwner(ctx, systemOwner, param.FirstOwner)
	if err != nil {
		return 0, fmt.Errorf("failed to AddFirstOwner. error: %w", err)
	}

	owner, err := s.rf.NewAppUserRepository().FindOwnerByLoginID(ctx, systemOwner, param.FirstOwner.LoginID)
	if err != nil {
		return 0, fmt.Errorf("failed to FindOwnerByLoginID. error: %w", err)
	}

	// add public group
	publicGroupID, err := s.rf.NewAppUserGroupRepository().AddPublicGroup(ctx, systemOwner)
	if err != nil {
		return 0, fmt.Errorf("failed to AddPublicGroup. error: %w", err)
	}

	// public-group <-> owner
	if err := s.rf.NewGroupUserRepository().AddGroupUser(ctx, systemOwner, publicGroupID, ownerID); err != nil {
		return 0, fmt.Errorf("failed to AddGroupUser. error: %w", err)
	}

	// add default space
	spaceID, err := s.rf.NewSpaceRepository().AddDefaultSpace(ctx, systemOwner)
	if err != nil {
		return 0, fmt.Errorf("failed to AddDefaultSpace. error: %w", err)
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