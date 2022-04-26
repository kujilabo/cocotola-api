//go:generate mockery --output mock --name Model
package domain

import (
	"time"

	libD "github.com/kujilabo/cocotola-api/pkg_lib/domain"
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
	ID        uint `validate:"gte=0"`
	Version   int  `validate:"required,gte=1"`
	CreatedAt time.Time
	UpdatedAt time.Time
	CreatedBy uint `validate:"gte=0"`
	UpdatedBy uint `validate:"gte=0"`
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
	return m, libD.Validator.Struct(m)
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
