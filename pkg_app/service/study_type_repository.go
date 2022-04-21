//go:generate mockery --output mock --name StudyTypeRepository
package service

import (
	"context"

	"github.com/kujilabo/cocotola-api/pkg_app/domain"
)

type StudyTypeRepository interface {
	FindAllStudyTypes(ctx context.Context) ([]domain.StudyType, error)
}
