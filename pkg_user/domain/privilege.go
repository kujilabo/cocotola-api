package domain

type Privileges interface {
	HasPrivilege(privilege RBACAction) bool
}

type privileges struct {
	values map[RBACAction]bool
}

func NewPrivileges(privs []RBACAction) Privileges {
	values := make(map[RBACAction]bool, 0)
	for _, p := range privs {
		values[p] = true
	}
	return &privileges{
		values: values,
	}
}

func (p *privileges) HasPrivilege(privilege RBACAction) bool {
	_, ok := p.values[privilege]
	return ok
}
