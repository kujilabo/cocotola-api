//go:generate mockery --output mock --name RepositoryFactory
package service

import (
	"context"
)

type RepositoryFactory interface {
	NewWorkbookRepository(ctx context.Context) (WorkbookRepository, error)

	NewProblemRepository(ctx context.Context, problemType string) (ProblemRepository, error)

	NewProblemTypeRepository(ctx context.Context) ProblemTypeRepository

	NewStudyTypeRepository(ctx context.Context) StudyTypeRepository

	NewRecordbookRepository(ctx context.Context) RecordbookRepository

	NewUserQuotaRepository(ctx context.Context) UserQuotaRepository
}
