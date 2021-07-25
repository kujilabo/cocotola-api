package domain

type RepositoryFactory interface {
	NewOrganizationRepository() OrganizationRepository
	NewSpaceRepository() SpaceRepository
	NewAppUserRepository() AppUserRepository
	NewAppUserGroupRepository() AppUserGroupRepository

	NewGroupUserRepository() GroupUserRepository
	NewUserSpaceRepository() UserSpaceRepository
	NewRBACRepository() RBACRepository
}
