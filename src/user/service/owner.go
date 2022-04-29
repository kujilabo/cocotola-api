package service

type Owner interface {
	AppUser
}

type owner struct {
	rf RepositoryFactory
	AppUser
}

func NewOwner(rf RepositoryFactory, appUser AppUser) Owner {
	return &owner{
		rf:      rf,
		AppUser: appUser,
	}
}
