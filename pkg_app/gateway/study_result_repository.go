package gateway

import (
	"context"
	"time"

	"golang.org/x/xerrors"
	"gorm.io/gorm"

	"github.com/kujilabo/cocotola-api/pkg_app/domain"
	"github.com/kujilabo/cocotola-api/pkg_lib/log"
)

type studyResultEntity struct {
	AppUserID      uint
	WorkbookID     uint
	ProblemTypeID  uint
	StudyTypeID    uint
	ProblemID      uint
	ResultPrev1    *bool
	ResultPrev2    *bool
	ResultPrev3    *bool
	Level          int
	LastAnsweredAt time.Time
}

// type ProblemEntity interface {
// 	ToProblem() domain.Problem
// }

func (e *studyResultEntity) TableName() string {
	return "study_result"
}

type studyResultRepository struct {
	rf         domain.RepositoryFactory
	db         *gorm.DB
	studyTypes []domain.StudyType
}

func NewStudyResultRepository(ctx context.Context, rf domain.RepositoryFactory, db *gorm.DB) (domain.StudyResultRepository, error) {
	studyTypeRepo, err := rf.NewStudyTypeRepository(ctx)
	if err != nil {
		return nil, err
	}
	studyTypes, err := studyTypeRepo.FindAllStudyTypes(ctx)
	if err != nil {
		return nil, err
	}
	logger := log.FromContext(ctx)
	logger.Infof("study types: %+v", studyTypes)
	return &studyResultRepository{
		rf:         rf,
		db:         db,
		studyTypes: studyTypes,
	}, nil
}

// func (r *studyResultRepository) toStudyType(studyTypeID uint) string {
// 	for _, m := range r.studyTypes {
// 		if m.GetID() == studyTypeID {
// 			return m.GetName()
// 		}
// 	}
// 	return ""
// }

func (r *studyResultRepository) toStudyTypeID(studyType string) uint {
	for _, m := range r.studyTypes {
		if m.GetName() == studyType {
			return m.GetID()
		}
	}
	return 0
}

func (r *studyResultRepository) FindStudyResults(ctx context.Context, operator domain.Student, workbookID domain.WorkbookID, studyType string) (map[domain.ProblemID]int, error) {
	studyTypeID := r.toStudyTypeID(studyType)
	if studyTypeID == 0 {
		return nil, xerrors.Errorf("unsupported studyType. studyType: %s", studyType)
	}

	var entities []studyResultEntity
	if result := r.db.Where("workbook_id = ? and study_type_id = ?", uint(workbookID), studyTypeID).Find(&entities); result.Error != nil {
		return nil, result.Error
	}

	results := make(map[domain.ProblemID]int)
	for _, e := range entities {
		results[domain.ProblemID(e.ProblemID)] = e.Level
	}

	return results, nil
}
