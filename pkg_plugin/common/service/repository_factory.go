package service

import (
	"context"
)

type RepositoryFactory interface {
	NewTatoebaLinkRepository(ctx context.Context) (TatoebaLinkRepository, error)

	NewTatoebaSentenceRepository(ctx context.Context) (TatoebaSentenceRepository, error)
}
