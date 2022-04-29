package domain

type Role interface {
	GetName() string
}

type role struct {
	Name string `validate:"required"`
}

func (r *role) GetName() string {
	return r.Name
}

var AdministratorRole = &role{Name: "Administrator"}
var OwnerRole = &role{Name: "Owner"}
var ManagerRole = &role{Name: "Manager"}
var UserRole = &role{Name: "User"}
var GuestRole = &role{Name: "Guest"}
var UnknownRole = &role{Name: "Unknown"}
