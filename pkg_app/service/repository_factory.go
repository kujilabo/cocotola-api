package service

import (
	"context"
)

type RepositoryFactory interface {
	NewWorkbookRepository(ctx context.Context) (WorkbookRepository, error)

	NewProblemRepository(ctx context.Context, problemType string) (ProblemRepository, error)

	NewProblemTypeRepository(ctx context.Context) (ProblemTypeRepository, error)

	NewStudyTypeRepository(ctx context.Context) (StudyTypeRepository, error)

	NewAudioRepository(ctx context.Context) (AudioRepository, error)

	NewRecordbookRepository(ctx context.Context) (RecordbookRepository, error)

	NewUserQuotaRepository(ctx context.Context) (UserQuotaRepository, error)
}

type AudioRepositoryFactory interface {
	NewAudioRepository(ctx context.Context) (AudioRepository, error)
}