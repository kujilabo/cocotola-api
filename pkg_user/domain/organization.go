package domain

type OrganizationID uint

type Organization interface {
	Model
	Name() string
}

type organization struct {
	Model
	name string
}

func NewOrganization(model Model, name string) Organization {
	return &organization{
		Model: model,
		name:  name,
	}
}

func (m *organization) Name() string {
	return m.name
}

func (m *organization) String() string {
	return m.name
}

// type OrganizationAddParameter struct {
// 	Name       string
// 	FirstOwner *FirstOwnerAddParameter
// }

// func NewOrganizationAddParameter(name string, firstOwner *FirstOwnerAddParameter) *OrganizationAddParameter {
// 	return &OrganizationAddParameter{
// 		Name:       name,
// 		FirstOwner: firstOwner,
// 	}
// }

// type OrganizationNotFoundError struct {
// 	id   uint
// 	text string
// }

// func NewOrganizationNotFoundError(id uint) *OrganizationNotFoundError {
// 	return &OrganizationNotFoundError{
// 		id:   id,
// 		text: fmt.Sprintf("Organization not found. Organization ID: %d", id),
// 	}
// }

// func (e *OrganizationNotFoundError) Error() string {
// 	return e.text
// }
