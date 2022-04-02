package gateway

import (
	"context"

	"github.com/kujilabo/cocotola-api/pkg_plugin/common/service"
	"gorm.io/gorm"
)

type repositoryFactory struct {
	db         *gorm.DB
	driverName string
}

func NewRepositoryFactory(ctx context.Context, db *gorm.DB, driverName string) (service.RepositoryFactory, error) {
	return &repositoryFactory{
		db:         db,
		driverName: driverName,
	}, nil
}

func (f *repositoryFactory) NewAzureTranslationRepository(ctx context.Context) (service.AzureTranslationRepository, error) {
	return NewAzureTranslationRepository(f.db), nil
}

func (f *repositoryFactory) NewCustomTranslationRepository(ctx context.Context) (service.CustomTranslationRepository, error) {
	return NewCustomTranslationRepository(f.db), nil
}

func (f *repositoryFactory) NewTatoebaSentenceRepository(ctx context.Context) (service.TatoebaSentenceRepository, error) {
	return NewTatoebaSentenceRepository(f.db)
}

func (f *repositoryFactory) NewTatoebaLinkRepository(ctx context.Context) (service.TatoebaLinkRepository, error) {
	return NewTatoebaLinkRepository(f.db), nil
}
