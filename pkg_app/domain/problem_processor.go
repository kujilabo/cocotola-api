package domain

import (
	"context"
)

type ProblemAddProcessor interface {
	AddProblem(ctx context.Context, repo RepositoryFactory, operator Student, param *ProblemAddParameter) (ProblemID, error)
}

type ProblemRemoveProcessor interface {
	RemoveProblem(ctx context.Context, repo RepositoryFactory, operator Student, problemID ProblemID, version int) error
}
