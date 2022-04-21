package gateway

import (
	"context"
	"errors"
	"time"

	"golang.org/x/xerrors"
	"gorm.io/gorm"

	"github.com/kujilabo/cocotola-api/pkg_app/domain"
	"github.com/kujilabo/cocotola-api/pkg_app/service"
	"github.com/kujilabo/cocotola-api/pkg_lib/log"
)

type recordbookEntity struct {
	AppUserID      uint
	WorkbookID     uint
	ProblemTypeID  uint
	StudyTypeID    uint
	ProblemID      uint
	ResultPrev1    *bool
	ResultPrev2    *bool
	ResultPrev3    *bool
	Level          int
	Memorized      bool
	LastAnsweredAt time.Time
}

// type ProblemEntity interface {
// 	ToProblem() domain.Problem
// }

func (e *recordbookEntity) TableName() string {
	return "recordbook"
}

type recordbookRepository struct {
	rf           service.RepositoryFactory
	db           *gorm.DB
	problemTypes []domain.ProblemType
	studyTypes   []domain.StudyType
}

func NewRecordbookRepository(ctx context.Context, rf service.RepositoryFactory, db *gorm.DB, problemTypes []domain.ProblemType) (service.RecordbookRepository, error) {
	studyTypeRepo := rf.NewStudyTypeRepository(ctx)
	studyTypes, err := studyTypeRepo.FindAllStudyTypes(ctx)
	if err != nil {
		return nil, err
	}
	logger := log.FromContext(ctx)
	logger.Infof("study types: %+v", studyTypes)
	return &recordbookRepository{
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

func (r *recordbookRepository) toProblemTypeID(problemType string) uint {
	for _, m := range r.problemTypes {
		if m.GetName() == problemType {
			return m.GetID()
		}
	}
	return 0
}

func (r *recordbookRepository) toStudyTypeID(studyType string) uint {
	for _, m := range r.studyTypes {
		if m.GetName() == studyType {
			return m.GetID()
		}
	}
	return 0
}

func (r *recordbookRepository) FindStudyResults(ctx context.Context, operator domain.StudentModel, workbookID domain.WorkbookID, studyType string) (map[domain.ProblemID]domain.StudyStatus, error) {
	studyTypeID := r.toStudyTypeID(studyType)
	if studyTypeID == 0 {
		return nil, xerrors.Errorf("unsupported studyType. studyType: %s", studyType)
	}

	var entities []recordbookEntity
	if result := r.db.Where("workbook_id = ?", uint(workbookID)).
		Where("study_type_id = ?", studyTypeID).
		Where("app_user_id = ?", operator.GetID()).
		Find(&entities); result.Error != nil {
		return nil, result.Error
	}

	results := make(map[domain.ProblemID]domain.StudyStatus)
	for _, e := range entities {
		results[domain.ProblemID(e.ProblemID)] = domain.StudyStatus{
			Level:     e.Level,
			Memorized: e.Memorized,
		}
	}

	return results, nil
}

func (r *recordbookRepository) SetResult(ctx context.Context, operator domain.StudentModel, workbookID domain.WorkbookID, studyType string, problemType string, problemID domain.ProblemID, studyResult, memorized bool) error {

	studyTypeID := r.toStudyTypeID(studyType)
	if studyTypeID == 0 {
		return xerrors.Errorf("unsupported studyType. studyType: %s", studyType)
	}

	problemTypeID := r.toProblemTypeID(problemType)
	if studyTypeID == 0 {
		return xerrors.Errorf("unsupported problemType. problemType: %s", problemType)
	}

	if memorized {
		return r.setMemorized(ctx, operator, workbookID, studyTypeID, problemTypeID, problemID)
	}

	return r.setResult(ctx, operator, workbookID, studyTypeID, problemTypeID, problemID, studyResult)
}

func (r *recordbookRepository) setResult(ctx context.Context, operator domain.StudentModel, workbookID domain.WorkbookID, studyTypeID uint, problemTypeID uint, problemID domain.ProblemID, studyResult bool) error {
	logger := log.FromContext(ctx)
	var entity recordbookEntity
	if result := r.db.Where("workbook_id = ?", uint(workbookID)).
		Where("study_type_id = ?", studyTypeID).
		Where("problem_id = ?", uint(problemID)).
		Where("app_user_id = ?", operator.GetID()).
		First(&entity); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			logger.Debugf("workbook_id = %d and study_type_id = %d and problem_id = %d", uint(workbookID), studyTypeID, uint(problemID))

			prev := false
			level := 0
			if studyResult {
				prev = true
				level = 1
			}
			entity = recordbookEntity{
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

	if result := r.db.Where("workbook_id = ?", uint(workbookID)).
		Where("study_type_id = ?", studyTypeID).
		Where("problem_id = ?", uint(problemID)).
		Where("app_user_id = ?", operator.GetID()).
		Updates(&entity); result.Error != nil {
		return result.Error
	}

	return nil
}

func (r *recordbookRepository) setMemorized(ctx context.Context, operator domain.StudentModel, workbookID domain.WorkbookID, studyTypeID uint, problemTypeID uint, problemID domain.ProblemID) error {
	logger := log.FromContext(ctx)

	var entity recordbookEntity
	if result := r.db.Where("workbook_id = ?", uint(workbookID)).
		Where("study_type_id = ?", studyTypeID).
		Where("problem_id = ?", uint(problemID)).
		Where("app_user_id = ?", operator.GetID()).
		First(&entity); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			logger.Debugf("workbook_id = %d and study_type_id = %d and problem_id = %d", uint(workbookID), studyTypeID, uint(problemID))
			entity = recordbookEntity{
				AppUserID:      operator.GetID(),
				WorkbookID:     uint(workbookID),
				ProblemTypeID:  problemTypeID,
				StudyTypeID:    studyTypeID,
				ProblemID:      uint(problemID),
				ResultPrev1:    nil,
				ResultPrev2:    nil,
				ResultPrev3:    nil,
				Level:          0,
				Memorized:      true,
				LastAnsweredAt: time.Now(),
			}
			if result := r.db.Create(&entity); result.Error != nil {
				return result.Error
			}
			return nil
		}
		return result.Error
	}

	entity.Memorized = true
	entity.LastAnsweredAt = time.Now()

	if result := r.db.Where("workbook_id = ?", uint(workbookID)).
		Where("study_type_id = ?", studyTypeID).
		Where("problem_id = ?", uint(problemID)).
		Where("app_user_id = ?", operator.GetID()).
		Updates(&entity); result.Error != nil {
		return result.Error
	}

	return nil
}
