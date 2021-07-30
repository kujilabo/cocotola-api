package domain

import (
	"time"

	"github.com/go-playground/validator"
)

type Model interface {
	GetID() uint
	GetVersion() int
	GetCreatedAt() time.Time
	GetUpdatedAt() time.Time
	GetCreatedBy() uint
	GetUpdatedBy() uint
}

type model struct {
	ID        uint `validate:"required,gte=1"`
	Version   int  `validate:"required,gte=1"`
	CreatedAt time.Time
	UpdatedAt time.Time
	CreatedBy uint `validate:"required,gte=1"`
	UpdatedBy uint `validate:"required,gte=1"`
}

func NewModel(id uint, version int, createdAt, updatedAt time.Time, createdBy, updatedBy uint) (Model, error) {
	m := &model{
		ID:        id,
		Version:   version,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
		CreatedBy: createdBy,
		UpdatedBy: updatedBy,
	}

	v := validator.New()
	return m, v.Struct(m)
}

func (m *model) GetID() uint {
	return m.ID
}

func (m *model) GetVersion() int {
	return m.Version
}

func (m *model) GetCreatedAt() time.Time {
	return m.CreatedAt
}

func (m *model) GetUpdatedAt() time.Time {
	return m.UpdatedAt
}

func (m *model) GetCreatedBy() uint {
	return m.CreatedBy
}

func (m *model) GetUpdatedBy() uint {
	return m.UpdatedBy
}
