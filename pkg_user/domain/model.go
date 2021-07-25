package domain

import "time"

type Model interface {
	GetID() uint
	GetVersion() int
	GetCreatedAt() time.Time
	GetUpdatedAt() time.Time
	GetCreatedBy() uint
	GetUpdatedBy() uint
}

type model struct {
	ID        uint
	Version   int
	CreatedAt time.Time
	UpdatedAt time.Time
	CreatedBy uint
	UpdatedBy uint
}

func NewModel(id uint, version int, createdAt, updatedAt time.Time, createdBy, updatedBy uint) Model {
	return &model{
		ID:        id,
		Version:   version,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
		CreatedBy: createdBy,
		UpdatedBy: updatedBy,
	}
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
