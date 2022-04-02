package gateway

import (
	"context"

	"gorm.io/gorm"

	"github.com/kujilabo/cocotola-api/pkg_app/domain"
	"github.com/kujilabo/cocotola-api/pkg_app/service"
)

type problemTypeEntity struct {
	ID   uint
	Name string
}

func (e *problemTypeEntity) TableName() string {
	return "problem_type"
}

func (e *problemTypeEntity) toModel() (domain.ProblemType, error) {
	return domain.NewProblemType(e.ID, e.Name)
}

type problemTypeRepository struct {
	db *gorm.DB
}

func NewProblemTypeRepository(db *gorm.DB) (service.ProblemTypeRepository, error) {
	return &problemTypeRepository{
		db: db,
	}, nil
}

func (r *problemTypeRepository) FindAllProblemTypes(ctx context.Context) ([]domain.ProblemType, error) {
	entities := []problemTypeEntity{}
	if err := r.db.Find(&entities).Error; err != nil {
		return nil, err
	}
	models := make([]domain.ProblemType, 0)
	for _, e := range entities {
		model, err := e.toModel()
		if err != nil {
			return nil, err
		}
		models = append(models, model)
	}
	return models, nil
}
