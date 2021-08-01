package domain

import "github.com/casbin/casbin/v2"

type RBACUser string
type RBACRole string
type RBACObject string
type RBACAction string

type RBACRepository interface {
	Init() error

	AddNamedPolicy(subject RBACRole, object RBACObject, action RBACAction) error

	AddNamedGroupingPolicy(subject RBACUser, object RBACRole) error

	NewEnforcerWithRolesAndUsers(roles []RBACRole, users []RBACUser) (*casbin.Enforcer, error)
}
