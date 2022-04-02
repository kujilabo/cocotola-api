package service

import (
	"context"
)

type RepositoryFactory interface {
	NewAzureTranslationRepository(ctx context.Context) (AzureTranslationRepository, error)

	NewCustomTranslationRepository(ctx context.Context) (CustomTranslationRepository, error)

	NewTatoebaLinkRepository(ctx context.Context) (TatoebaLinkRepository, error)

	NewTatoebaSentenceRepository(ctx context.Context) (TatoebaSentenceRepository, error)
}
