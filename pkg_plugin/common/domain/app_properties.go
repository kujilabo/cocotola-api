package domain

import (
	"context"

	"gorm.io/gorm"
)

type RepositoryFactoryFunc func(ctx context.Context, db *gorm.DB) (RepositoryFactory, error)
