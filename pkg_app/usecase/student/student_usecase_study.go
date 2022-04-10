package student

import (
	"context"

	"github.com/kujilabo/cocotola-api/pkg_app/domain"
	"github.com/kujilabo/cocotola-api/pkg_app/service"
	"github.com/kujilabo/cocotola-api/pkg_app/usecase"
	user "github.com/kujilabo/cocotola-api/pkg_user/domain"
	userS "github.com/kujilabo/cocotola-api/pkg_user/service"
	"golang.org/x/xerrors"
	"gorm.io/gorm"
)

type StudentUsecaseStudy interface {

	// study
	FindResults(ctx context.Context, organizationID user.OrganizationID, operatorID user.AppUserID, workbookID domain.WorkbookID, studyType string) ([]domain.ProblemWithLevel, error)

	// FindAllProblemsByWorkbookID(ctx context.Context, organizationID, operatorID, workbookID uint, studyTypeID domain.StudyTypeID) (domain.WorkbookWithProblems, error)
	SetResult(ctx context.Context, organizationID user.OrganizationID, operatorID user.AppUserID, workbookID domain.WorkbookID, studyType string, problemID domain.ProblemID, result, memorized bool) error
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

func (s *studentUsecaseStudy) FindResults(ctx context.Context, organizationID user.OrganizationID, operatorID user.AppUserID, workbookID domain.WorkbookID, studyType string) ([]domain.ProblemWithLevel, error) {
	var results []domain.ProblemWithLevel
	if err := s.db.Transaction(func(tx *gorm.DB) error {
		rfFunc, err := s.rfFunc(ctx, tx)
		if err != nil {
			return xerrors.Errorf("failed to rfFunc. err: %w", err)
		}
		userRepo, err := s.userRfFunc(ctx, tx)
		if err != nil {
			return xerrors.Errorf("failed to userRepo. err: %w", err)
		}
		student, err := usecase.FindStudent(ctx, s.pf, rfFunc, userRepo, organizationID, operatorID)
		if err != nil {
			return xerrors.Errorf("failed to findStudent. err: %w", err)
		}
		tmpResult, err := student.FindRecordbook(ctx, workbookID, studyType)
		if err != nil {
			return xerrors.Errorf("failed to FindRecordbook. err: %w", err)
		}
		tmpResults, err := tmpResult.GetResultsSortedLevel(ctx)
		if err != nil {
			return xerrors.Errorf("failed to GetResultsSortedLevel. err: %w", err)
		}
		results = tmpResults
		return nil
	}); err != nil {
		return nil, err
	}

	return results, nil
}

func (s *studentUsecaseStudy) SetResult(ctx context.Context, organizationID user.OrganizationID, operatorID user.AppUserID, workbookID domain.WorkbookID, studyType string, problemID domain.ProblemID, result, memorized bool) error {
	if err := s.db.Transaction(func(tx *gorm.DB) error {
		rf, err := s.rfFunc(ctx, tx)
		if err != nil {
			return xerrors.Errorf("failed to rfFunc. err: %w", err)
		}
		userRepo, err := s.userRfFunc(ctx, tx)
		if err != nil {
			return xerrors.Errorf("failed to userRepo. err: %w", err)
		}
		student, err := usecase.FindStudent(ctx, s.pf, rf, userRepo, organizationID, operatorID)
		if err != nil {
			return xerrors.Errorf("failed to findStudent. err: %w", err)
		}
		workbook, err := student.FindWorkbookByID(ctx, workbookID)
		if err != nil {
			return err
		}
		tmpResult, err := student.FindRecordbook(ctx, workbookID, studyType)
		if err != nil {
			return xerrors.Errorf("failed to FindRecordbook. err: %w", err)
		}
		if err := tmpResult.SetResult(ctx, workbook.GetProblemType(), problemID, result, memorized); err != nil {
			return xerrors.Errorf("failed to SetResult. err: %w", err)
		}
		return nil
	}); err != nil {
		return err
	}
	return nil
}
