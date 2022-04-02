package gateway

import (
	"context"
	"errors"

	"golang.org/x/xerrors"
	"gorm.io/gorm"

	libG "github.com/kujilabo/cocotola-api/pkg_lib/gateway"
	"github.com/kujilabo/cocotola-api/pkg_plugin/common/service"
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

func NewTatoebaLinkRepository(db *gorm.DB) service.TatoebaLinkRepository {
	return &tatoebaLinkRepository{
		db: db,
	}
}

func (r *tatoebaLinkRepository) Add(ctx context.Context, param service.TatoebaLinkAddParameter) error {
	entity := tatoebaLinkEntity{
		From: param.GetFrom(),
		To:   param.GetTo(),
	}

	if result := r.db.Create(&entity); result.Error != nil {
		if err := libG.ConvertDuplicatedError(result.Error, service.ErrTatoebaLinkAlreadyExists); errors.Is(err, service.ErrTatoebaLinkAlreadyExists) {
			return xerrors.Errorf("failed to Add tatoebaLink. err: %w", err)
		}

		if err := libG.ConvertRelationError(result.Error, service.ErrTatoebaLinkSourceNotFound); errors.Is(err, service.ErrTatoebaLinkSourceNotFound) {
			// nothing
			return nil
		}

		return xerrors.Errorf("failed to Add tatoebaLink. err: %w", result.Error)
	}

	return nil
}
