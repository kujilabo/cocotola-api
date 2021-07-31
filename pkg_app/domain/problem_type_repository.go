package domain

import "context"

type ProblemTypeRepository interface {
	FindAllProblemTypes(ctx context.Context) ([]ProblemType, error)
}
