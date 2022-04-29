package domain

type OwnerModel interface {
	AppUserModel
}

type ownerModel struct {
	// rf RepositoryFactory
	AppUserModel
}

func NewOwner(
	// rf RepositoryFactory,
	appUser AppUserModel) OwnerModel {
	return &ownerModel{
		// rf:      rf,
		AppUserModel: appUser,
	}
}
