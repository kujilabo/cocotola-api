//go:generate mockery --output mock --name RepositoryFactory
//go:generate mockery --output mock --name AudioRepositoryFactory
package service

import (
	"context"
)

type RepositoryFactory interface {
	NewWorkbookRepository(ctx context.Context) (WorkbookRepository, error)

	NewProblemRepository(ctx context.Context, problemType string) (ProblemRepository, error)

	NewProblemTypeRepository(ctx context.Context) ProblemTypeRepository

	NewStudyTypeRepository(ctx context.Context) StudyTypeRepository

	NewAudioRepository(ctx context.Context) AudioRepository

	NewRecordbookRepository(ctx context.Context) RecordbookRepository

	NewUserQuotaRepository(ctx context.Context) UserQuotaRepository
}

type AudioRepositoryFactory interface {
	NewAudioRepository(ctx context.Context) AudioRepository
}
