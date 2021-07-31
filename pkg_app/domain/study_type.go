package domain

import "github.com/go-playground/validator"

type StudyType interface {
	GetID() uint
	GetName() string
}

type studyType struct {
	ID   uint   `validate:"required,gte=1"`
	Name string `validate:"required"`
}

func NewStudyType(id uint, name string) (StudyType, error) {
	m := &studyType{
		ID:   id,
		Name: name,
	}

	v := validator.New()
	return m, v.Struct(m)
}

func (m *studyType) GetID() uint {
	return m.ID
}

func (m *studyType) GetName() string {
	return m.Name
}
