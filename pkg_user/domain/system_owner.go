package domain

import (
	"context"

	"github.com/go-playground/validator"

	"github.com/kujilabo/cocotola-api/pkg_lib/log"
)

const SystemOwnerID = 2

type SystemOwner interface {
	AppUser

	GetOrganization(ctxc context.Context) (Organization, error)

	FindAppUserByID(ctx context.Context, id AppUserID) (AppUser, error)

	FindAppUserByLoginID(ctx context.Context, loginID string) (AppUser, error)

	AddAppUser(ctx context.Context, param *AppUserAddParameter) (AppUserID, error)
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

	v := validator.New()
	return m, v.Struct(m)
}

func (s *systemOwner) GetOrganization(ctx context.Context) (Organization, error) {
	return s.rf.NewOrganizationRepository().GetOrganization(ctx, s)
}

func (s *systemOwner) FindAppUserByID(ctx context.Context, id AppUserID) (AppUser, error) {
	return s.rf.NewAppUserRepository().FindAppUserByID(ctx, s, id)
}

func (s *systemOwner) FindAppUserByLoginID(ctx context.Context, loginID string) (AppUser, error) {
	return s.rf.NewAppUserRepository().FindAppUserByLoginID(ctx, s, loginID)
}

func (s *systemOwner) AddAppUser(ctx context.Context, param *AppUserAddParameter) (AppUserID, error) {
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
	if err := s.rf.NewGroupUserRepository().AddGroupUser(ctx, s, AppUserGroupID(publicGroup.GetID()), AppUserID(appUser.GetID())); err != nil {
		return 0, err
	}

	spaceID, err := s.rf.NewSpaceRepository().AddPersonalSpace(ctx, s, appUser)
	if err != nil {
		return 0, err
	}

	logger.Infof("Personal spaceID: %d", spaceID)

	spaceWriter := NewSpaceWriterRole(spaceID)
	spaceObject := NewSpaceObject(spaceID)
	userSubject := NewUserObject(appUserID)

	rbacRepo := s.rf.NewRBACRepository()
	rbacRepo.AddNamedPolicy(spaceWriter, spaceObject, "read")
	rbacRepo.AddNamedPolicy(spaceWriter, spaceObject, "write")
	rbacRepo.AddNamedGroupingPolicy(userSubject, spaceWriter)

	// defaultSpace, err := s.rf.NewSpaceRepository().FindDefaultSpace(ctx, s)
	// if err != nil {
	// 	return 0, err
	// }

	// if err := s.rf.NewUserSpaceRepository().Add(ctx, appUser, SpaceID(defaultSpace.GetID())); err != nil {
	// 	return 0, err
	// }

	return appUserID, nil
}
