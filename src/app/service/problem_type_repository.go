//go:generate mockery --output mock --name ProblemTypeRepository
package service

import (
	"context"

	"github.com/kujilabo/cocotola-api/src/app/domain"
)

type ProblemTypeRepository interface {
	FindAllProblemTypes(ctx context.Context) ([]domain.ProblemType, error)
}
