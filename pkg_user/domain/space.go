package domain

import (
	"time"
)

type SpaceID uint

type Space interface {
	Model
	Name() string
}

type space struct {
	model
	organizationID uint
	spaceType      int
	key            string
	name           string
	description    string
}

func NewSpace(id uint, version int, createdAt, updatedAt time.Time, createdBy, updatedBy uint, organizationID uint, spaceType int, key, name, description string) Space {
	return &space{
		model: model{
			ID:        id,
			Version:   version,
			CreatedAt: createdAt,
			UpdatedAt: updatedAt,
			CreatedBy: createdBy,
			UpdatedBy: updatedBy,
		},
		organizationID: organizationID,
		spaceType:      spaceType,
		key:            key,
		name:           name,
		description:    description,
	}
}

func (m *space) Key() string {
	return m.key
}

func (m *space) Name() string {
	return m.name
}

func (m *space) Description() string {
	return m.description
}

func (m *space) String() string {
	return m.name
}
