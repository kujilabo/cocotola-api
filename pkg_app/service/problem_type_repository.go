package service

import (
	"context"

	"github.com/kujilabo/cocotola-api/pkg_app/domain"
)

type ProblemTypeRepository interface {
	FindAllProblemTypes(ctx context.Context) ([]domain.ProblemType, error)
}
