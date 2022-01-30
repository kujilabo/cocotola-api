package domain

import (
	"context"
)

type RepositoryFactory interface {
	NewAzureTranslationRepository(ctx context.Context) (AzureTranslationRepository, error)

	NewCustomTranslationRepository(ctx context.Context) (CustomTranslationRepository, error)
}
