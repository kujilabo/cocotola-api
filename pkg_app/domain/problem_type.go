package domain

import "github.com/go-playground/validator"

type ProblemType interface {
	GetID() uint
	GetName() string
}

type problemType struct {
	ID   uint   `validate:"required,gte=1"`
	Name string `validate:"required"`
}

func NewProblemType(id uint, name string) (ProblemType, error) {
	m := &problemType{
		ID:   id,
		Name: name,
	}

	v := validator.New()
	return m, v.Struct(m)
}

func (m *problemType) GetID() uint {
	return m.ID
}

func (m *problemType) GetName() string {
	return m.Name
}
