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

func NewProblemTypeRepository(db *gorm.DB) service.ProblemTypeRepository {
	return &problemTypeRepository{db: db}
}

func (r *problemTypeRepository) FindAllProblemTypes(ctx context.Context) ([]domain.ProblemType, error) {
	_, span := tracer.Start(ctx, "problemTypeRepository.FindAllProblemTypes")
	defer span.End()

	entities := []problemTypeEntity{}
	if err := r.db.Find(&entities).Error; err != nil {
		return nil, err
	}

	models := make([]domain.ProblemType, len(entities))
	for i, e := range entities {
		model, err := e.toModel()
		if err != nil {
			return nil, err
		}

		models[i] = model
	}

	return models, nil
}
