package gateway

import (
	"context"
	"errors"
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
	rf           domain.RepositoryFactory
	db           *gorm.DB
	problemTypes []domain.ProblemType
	studyTypes   []domain.StudyType
}

func NewStudyResultRepository(ctx context.Context, rf domain.RepositoryFactory, db *gorm.DB, problemTypes []domain.ProblemType) (domain.StudyResultRepository, error) {
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
		rf:           rf,
		db:           db,
		problemTypes: problemTypes,
		studyTypes:   studyTypes,
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

func (r *studyResultRepository) toProblemTypeID(problemType string) uint {
	for _, m := range r.problemTypes {
		if m.GetName() == problemType {
			return m.GetID()
		}
	}
	return 0
}
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

func (r *studyResultRepository) SetResult(ctx context.Context, operator domain.Student, workbookID domain.WorkbookID, studyType string, problemType string, problemID domain.ProblemID, studyResult bool) error {
	logger := log.FromContext(ctx)

	studyTypeID := r.toStudyTypeID(studyType)
	if studyTypeID == 0 {
		return xerrors.Errorf("unsupported studyType. studyType: %s", studyType)
	}

	problemTypeID := r.toProblemTypeID(problemType)
	if studyTypeID == 0 {
		return xerrors.Errorf("unsupported problemType. problemType: %s", problemType)
	}

	var entity studyResultEntity
	if result := r.db.
		Where("workbook_id = ? and study_type_id = ? and problem_id = ?",
			uint(workbookID), studyTypeID, uint(problemID)).
		First(&entity); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			logger.Debugf("workbook_id = %d and study_type_id = %d and problem_id = %d", uint(workbookID), studyTypeID, uint(problemID))

			prev := false
			level := 0
			if studyResult {
				prev = true
				level = 1
			}
			entity = studyResultEntity{
				AppUserID:      operator.GetID(),
				WorkbookID:     uint(workbookID),
				ProblemTypeID:  problemTypeID,
				StudyTypeID:    studyTypeID,
				ProblemID:      uint(problemID),
				ResultPrev1:    &prev,
				ResultPrev2:    nil,
				ResultPrev3:    nil,
				Level:          level,
				LastAnsweredAt: time.Now(),
			}
			if result := r.db.Create(&entity); result.Error != nil {
				return result.Error
			}
			return nil
		}
		return result.Error
	}

	if studyResult {
		if entity.Level < domain.StudyMaxLevel {
			entity.Level++
		}
	} else {
		if entity.Level > domain.StudyMinLevel {
			entity.Level--
		}
	}

	if entity.ResultPrev2 != nil {
		b := *entity.ResultPrev2
		entity.ResultPrev3 = &b
	}
	if entity.ResultPrev1 != nil {
		b := *entity.ResultPrev1
		entity.ResultPrev2 = &b
	}
	*entity.ResultPrev1 = studyResult
	entity.LastAnsweredAt = time.Now()

	if result := r.db.
		Where("workbook_id = ? and study_type_id = ? and problem_id = ?",
			uint(workbookID), studyTypeID, uint(problemID)).
		Updates(&entity); result.Error != nil {
		return result.Error
	}

	return nil
}
