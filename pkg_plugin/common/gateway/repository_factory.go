package gateway

import (
	"context"

	"github.com/kujilabo/cocotola-api/pkg_plugin/common/domain"
	"gorm.io/gorm"
)

type repositoryFactory struct {
	db         *gorm.DB
	driverName string
}

func NewRepositoryFactory(ctx context.Context, db *gorm.DB, driverName string) (domain.RepositoryFactory, error) {
	return &repositoryFactory{
		db:         db,
		driverName: driverName,
	}, nil
}

func (f *repositoryFactory) NewAzureTranslationRepository(ctx context.Context) (domain.AzureTranslationRepository, error) {
	return NewAzureTranslationRepository(f.db), nil
}

func (f *repositoryFactory) NewCustomTranslationRepository(ctx context.Context) (domain.CustomTranslationRepository, error) {
	return NewCustomTranslationRepository(f.db), nil
}

func (f *repositoryFactory) NewTatoebaSentenceRepository(ctx context.Context) (domain.TatoebaSentenceRepository, error) {
	return NewTatoebaSentenceRepository(f.db)
}

func (f *repositoryFactory) NewTatoebaLinkRepository(ctx context.Context) (domain.TatoebaLinkRepository, error) {
	return NewTatoebaLinkRepository(f.db), nil
}
