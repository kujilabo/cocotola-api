package gateway

import (
	"context"
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
		err := libG.ConvertDuplicatedError(result.Error, domain.ErrTatoebaLinkAlreadyExists)
		return fmt.Errorf("failed to Add tatoebaLink. err: %w", err)
	}

	return nil
}
