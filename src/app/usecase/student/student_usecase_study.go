package student

import (
	"context"

	"gorm.io/gorm"

	"github.com/kujilabo/cocotola-api/src/app/domain"
	"github.com/kujilabo/cocotola-api/src/app/service"
	"github.com/kujilabo/cocotola-api/src/app/usecase"
	liberrors "github.com/kujilabo/cocotola-api/src/lib/errors"
	userD "github.com/kujilabo/cocotola-api/src/user/domain"
	userS "github.com/kujilabo/cocotola-api/src/user/service"
)

type StudentUsecaseStudy interface {

	// study
	FindResults(ctx context.Context, organizationID userD.OrganizationID, operatorID userD.AppUserID, workbookID domain.WorkbookID, studyType string) ([]domain.StudyRecordWithProblemID, error)

	GetCompletionRate(ctx context.Context, organizationID userD.OrganizationID, operatorID userD.AppUserID, workbookID domain.WorkbookID) (map[string]int, error)

	// FindAllProblemsByWorkbookID(ctx context.Context, organizationID, operatorID, workbookID uint, studyTypeID domain.StudyTypeID) (domain.WorkbookWithProblems, error)
	SetResult(ctx context.Context, organizationID userD.OrganizationID, operatorID userD.AppUserID, workbookID domain.WorkbookID, studyType string, problemID domain.ProblemID, result, memorized bool) error
}

type studentUsecaseStudy struct {
	db         *gorm.DB
	pf         service.ProcessorFactory
	rfFunc     service.RepositoryFactoryFunc
	userRfFunc userS.RepositoryFactoryFunc
}

func NewStudentUsecaseStudy(db *gorm.DB, pf service.ProcessorFactory, rfFunc service.RepositoryFactoryFunc, userRfFunc userS.RepositoryFactoryFunc) StudentUsecaseStudy {
	return &studentUsecaseStudy{
		db:         db,
		pf:         pf,
		rfFunc:     rfFunc,
		userRfFunc: userRfFunc,
	}
}

func (s *studentUsecaseStudy) FindResults(ctx context.Context, organizationID userD.OrganizationID, operatorID userD.AppUserID, workbookID domain.WorkbookID, studyType string) ([]domain.StudyRecordWithProblemID, error) {
	var results []domain.StudyRecordWithProblemID
	if err := s.db.Transaction(func(tx *gorm.DB) error {
		student, err := s.findStudent(ctx, tx, organizationID, operatorID)
		if err != nil {
			return liberrors.Errorf("failed to findStudent. err: %w", err)
		}
		recordbook, err := student.FindRecordbook(ctx, workbookID, studyType)
		if err != nil {
			return liberrors.Errorf("failed to FindRecordbook. err: %w", err)
		}
		tmpResults, err := recordbook.GetResultsSortedLevel(ctx)
		if err != nil {
			return liberrors.Errorf("failed to GetResultsSortedLevel. err: %w", err)
		}
		results = tmpResults
		return nil
	}); err != nil {
		return nil, err
	}

	return results, nil
}

func (s *studentUsecaseStudy) GetCompletionRate(ctx context.Context, organizationID userD.OrganizationID, operatorID userD.AppUserID, workbookID domain.WorkbookID) (map[string]int, error) {
	var results map[string]int
	if err := s.db.Transaction(func(tx *gorm.DB) error {
		student, err := s.findStudent(ctx, tx, organizationID, operatorID)
		if err != nil {
			return liberrors.Errorf("failed to findStudent. err: %w", err)
		}
		recordbookSummary, err := student.FindRecordbookSummary(ctx, workbookID)
		if err != nil {
			return liberrors.Errorf("failed to FindRecordbook. err: %w", err)
		}
		tmpResults, err := recordbookSummary.GetCompletionRate(ctx)
		if err != nil {
			return liberrors.Errorf("failed to GetResultsSortedLevel. err: %w", err)
		}
		results = tmpResults
		return nil
	}); err != nil {
		return nil, err
	}

	return results, nil
}

func (s *studentUsecaseStudy) SetResult(ctx context.Context, organizationID userD.OrganizationID, operatorID userD.AppUserID, workbookID domain.WorkbookID, studyType string, problemID domain.ProblemID, result, memorized bool) error {
	if err := s.db.Transaction(func(tx *gorm.DB) error {
		student, err := s.findStudent(ctx, tx, organizationID, operatorID)
		if err != nil {
			return liberrors.Errorf("failed to findStudent. err: %w", err)
		}
		workbook, err := student.FindWorkbookByID(ctx, workbookID)
		if err != nil {
			return err
		}
		recordbook, err := student.FindRecordbook(ctx, workbookID, studyType)
		if err != nil {
			return liberrors.Errorf("failed to FindRecordbook. err: %w", err)
		}
		if err := recordbook.SetResult(ctx, workbook.GetProblemType(), problemID, result, memorized); err != nil {
			return liberrors.Errorf("failed to SetResult. err: %w", err)
		}
		return nil
	}); err != nil {
		return err
	}
	return nil
}
func (s *studentUsecaseStudy) findStudent(ctx context.Context, db *gorm.DB, organizationID userD.OrganizationID, operatorID userD.AppUserID) (service.Student, error) {
	rf, err := s.rfFunc(ctx, db)
	if err != nil {
		return nil, liberrors.Errorf("failed to rfFunc. err: %w", err)
	}
	userRepo, err := s.userRfFunc(ctx, db)
	if err != nil {
		return nil, liberrors.Errorf("failed to userRepo. err: %w", err)
	}
	student, err := usecase.FindStudent(ctx, s.pf, rf, userRepo, organizationID, operatorID)
	if err != nil {
		return nil, liberrors.Errorf("failed to findStudent. err: %w", err)
	}

	return student, nil
}
