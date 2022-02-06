package gateway

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"

	libG "github.com/kujilabo/cocotola-api/pkg_lib/gateway"
	"github.com/kujilabo/cocotola-api/pkg_plugin/common/domain"
)

type tatoebaLinkRepository struct {
	db *gorm.DB
}

type tatoebaLinkEntity struct {
	From int
	To   int
}

func (e *tatoebaLinkEntity) TableName() string {
	return "tatoeba_link"
}

func NewTatoebaLinkRepository(db *gorm.DB) domain.TatoebaLinkRepository {
	return &tatoebaLinkRepository{
		db: db,
	}
}

func (r *tatoebaLinkRepository) Add(ctx context.Context, param domain.TatoebaLinkAddParameter) error {
	entity := tatoebaLinkEntity{
		From: param.GetFrom(),
		To:   param.GetTo(),
	}

	if result := r.db.Create(&entity); result.Error != nil {
		if err := libG.ConvertDuplicatedError(result.Error, domain.ErrTatoebaLinkAlreadyExists); errors.Is(err, domain.ErrTatoebaLinkAlreadyExists) {
			return fmt.Errorf("failed to Add tatoebaLink. err: %w", err)
		}

		if err := libG.ConvertRelationError(result.Error, domain.ErrTatoebaLinkSourceNotFound); errors.Is(err, domain.ErrTatoebaLinkSourceNotFound) {
			return nil
		}

		return fmt.Errorf("failed to Add tatoebaLink. err: %w", result.Error)
	}

	return nil
}
