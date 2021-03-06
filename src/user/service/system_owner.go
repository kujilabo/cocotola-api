package service

import (
	"context"

	libD "github.com/kujilabo/cocotola-api/src/lib/domain"
	"github.com/kujilabo/cocotola-api/src/lib/log"
	"github.com/kujilabo/cocotola-api/src/user/domain"
)

const SystemOwnerID = 2

type SystemOwner interface {
	AppUser

	GetOrganization(ctxc context.Context) (Organization, error)

	FindAppUserByID(ctx context.Context, id domain.AppUserID) (AppUser, error)

	FindAppUserByLoginID(ctx context.Context, loginID string) (AppUser, error)

	FindSystemSpace(ctx context.Context) (Space, error)

	AddAppUser(ctx context.Context, param AppUserAddParameter) (domain.AppUserID, error)

	AddSystemSpace(ctx context.Context) (domain.SpaceID, error)
}

type systemOwner struct {
	AppUser
	rf RepositoryFactory
}

func NewSystemOwner(rf RepositoryFactory, appUser AppUser) (SystemOwner, error) {
	m := &systemOwner{
		rf:      rf,
		AppUser: appUser,
	}

	return m, libD.Validator.Struct(m)
}

func (s *systemOwner) GetOrganization(ctx context.Context) (Organization, error) {
	return s.rf.NewOrganizationRepository().GetOrganization(ctx, s)
}

func (s *systemOwner) FindAppUserByID(ctx context.Context, id domain.AppUserID) (AppUser, error) {
	return s.rf.NewAppUserRepository().FindAppUserByID(ctx, s, id)
}

func (s *systemOwner) FindAppUserByLoginID(ctx context.Context, loginID string) (AppUser, error) {
	return s.rf.NewAppUserRepository().FindAppUserByLoginID(ctx, s, loginID)
}

func (s *systemOwner) FindSystemSpace(ctx context.Context) (Space, error) {
	return s.rf.NewSpaceRepository().FindSystemSpace(ctx, s)
}

func (s *systemOwner) AddAppUser(ctx context.Context, param AppUserAddParameter) (domain.AppUserID, error) {
	logger := log.FromContext(ctx)
	logger.Infof("AddStudent")
	appUserID, err := s.rf.NewAppUserRepository().AddAppUser(ctx, s, param)
	if err != nil {
		return 0, err
	}
	appUser, err := s.rf.NewAppUserRepository().FindAppUserByID(ctx, s, appUserID)
	if err != nil {
		return 0, err
	}

	// personalGroupID, err := s.rf.NewAppUserGroupRepository().AddPersonalGroup(s, studentID)
	// if err != nil {
	// 	return 0, err
	// }

	publicGroup, err := s.rf.NewAppUserGroupRepository().FindPublicGroup(ctx, s)
	if err != nil {
		return 0, err
	}
	if err := s.rf.NewGroupUserRepository().AddGroupUser(ctx, s, domain.AppUserGroupID(publicGroup.GetID()), domain.AppUserID(appUser.GetID())); err != nil {
		return 0, err
	}

	spaceID, err := s.rf.NewSpaceRepository().AddPersonalSpace(ctx, s, appUser)
	if err != nil {
		return 0, err
	}

	logger.Infof("Personal spaceID: %d", spaceID)

	spaceWriter := domain.NewSpaceWriterRole(spaceID)
	spaceObject := domain.NewSpaceObject(spaceID)
	userSubject := domain.NewUserObject(appUserID)

	rbacRepo := s.rf.NewRBACRepository()
	if err := rbacRepo.AddNamedPolicy(spaceWriter, spaceObject, "read"); err != nil {
		return 0, err
	}

	if err := rbacRepo.AddNamedPolicy(spaceWriter, spaceObject, "write"); err != nil {
		return 0, err
	}

	if err := rbacRepo.AddNamedGroupingPolicy(userSubject, spaceWriter); err != nil {
		return 0, err
	}

	// defaultSpace, err := s.rf.NewSpaceRepository().FindDefaultSpace(ctx, s)
	// if err != nil {
	// 	return 0, err
	// }

	// if err := s.rf.NewUserSpaceRepository().Add(ctx, appUser, SpaceID(defaultSpace.GetID())); err != nil {
	// 	return 0, err
	// }

	return appUserID, nil
}

func (s *systemOwner) AddSystemSpace(ctx context.Context) (domain.SpaceID, error) {
	logger := log.FromContext(ctx)
	logger.Infof("AddSystemSpace")

	spaceID, err := s.rf.NewSpaceRepository().AddSystemSpace(ctx, s)
	if err != nil {
		return 0, err
	}
	return spaceID, nil
}
