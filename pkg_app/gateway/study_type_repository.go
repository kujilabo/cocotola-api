package gateway

import (
	"context"

	"gorm.io/gorm"

	"github.com/kujilabo/cocotola-api/pkg_app/domain"
	"github.com/kujilabo/cocotola-api/pkg_app/service"
)

type studyTypeEntity struct {
	ID   uint
	Name string
}

func (e *studyTypeEntity) TableName() string {
	return "study_type"
}

func (e *studyTypeEntity) toModel() (domain.StudyType, error) {
	return domain.NewStudyType(e.ID, e.Name)
}

type studyTypeRepository struct {
	db *gorm.DB
}

func NewStudyTypeRepository(db *gorm.DB) service.StudyTypeRepository {
	return &studyTypeRepository{db: db}
}

func (r *studyTypeRepository) FindAllStudyTypes(ctx context.Context) ([]domain.StudyType, error) {
	_, span := tracer.Start(ctx, "studyTypeRepository.FindAllStudyTypes")
	defer span.End()

	entities := []studyTypeEntity{}
	if err := r.db.Find(&entities).Error; err != nil {
		return nil, err
	}

	models := make([]domain.StudyType, len(entities))
	for i, e := range entities {
		model, err := e.toModel()
		if err != nil {
			return nil, err
		}

		models[i] = model
	}

	return models, nil
}
