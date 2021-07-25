package domain

type Role interface {
	Name() string
}

type role struct {
	name string
}

func (r *role) Name() string {
	return r.name
}

var AdministratorRole = &role{name: "Administrator"}
var OwnerRole = &role{name: "Owner"}
var ManagerRole = &role{name: "Manager"}
var UserRole = &role{name: "User"}
var GuestRole = &role{name: "Guest"}
var UnknwonRole = &role{name: "Unknwon"}
