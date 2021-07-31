package application

import (
	"context"

	"golang.org/x/xerrors"
	"gorm.io/gorm"

	"github.com/kujilabo/cocotola-api/pkg_app/domain"
	"github.com/kujilabo/cocotola-api/pkg_lib/log"
	user "github.com/kujilabo/cocotola-api/pkg_user/domain"
)

type ProblemService interface {
	FindProblemsByWorkbookID(ctx context.Context, organizationID user.OrganizationID, operatorID user.AppUserID, workbookID domain.WorkbookID, param domain.ProblemSearchCondition) (*domain.ProblemSearchResult, error)

	FindProblemsByProblemIDs(ctx context.Context, organizationID user.OrganizationID, operatorID user.AppUserID, workbookID domain.WorkbookID, param domain.ProblemIDsCondition) (*domain.ProblemSearchResult, error)

	FindProblemByID(ctx context.Context, organizationID user.OrganizationID, operatorID user.AppUserID, workbookID domain.WorkbookID, problemID domain.ProblemID) (domain.Problem, error)

	FindProblemIDs(ctx context.Context, organizationID user.OrganizationID, operatorID user.AppUserID, workbookID domain.WorkbookID) ([]domain.ProblemID, error)

	AddProblem(ctx context.Context, organizationID user.OrganizationID, operatorID user.AppUserID, param domain.ProblemAddParameter) (domain.ProblemID, error)

	RemoveProblem(ctx context.Context, organizationID user.OrganizationID, operatorID user.AppUserID, workbookID domain.WorkbookID, problemID domain.ProblemID, version int) error
}

type problemService struct {
	db       *gorm.DB
	repo     func(db *gorm.DB) domain.RepositoryFactory
	userRepo func(db *gorm.DB) user.RepositoryFactory
}

func NewProblemService(db *gorm.DB, repo func(db *gorm.DB) domain.RepositoryFactory, userRepo func(db *gorm.DB) user.RepositoryFactory) ProblemService {
	return &problemService{
		db:       db,
		repo:     repo,
		userRepo: userRepo,
	}
}

func (s *problemService) FindProblemsByWorkbookID(ctx context.Context, organizationID user.OrganizationID, operatorID user.AppUserID, workbookID domain.WorkbookID, param domain.ProblemSearchCondition) (*domain.ProblemSearchResult, error) {
	var result *domain.ProblemSearchResult
	if err := s.db.Transaction(func(tx *gorm.DB) error {
		repo := s.repo(tx)
		userRepo := s.userRepo(tx)
		student, err := findStudent(ctx, repo, userRepo, organizationID, operatorID)
		if err != nil {
			return err
		}
		workbook, err := student.FindWorkbookByID(ctx, workbookID)
		if err != nil {
			return err
		}
		tmpResult, err := workbook.FindProblems(ctx, student, param)
		if err != nil {
			return err
		}
		result = tmpResult
		return nil
	}); err != nil {
		return nil, err
	}
	return result, nil
}

func (s *problemService) FindProblemsByProblemIDs(ctx context.Context, organizationID user.OrganizationID, operatorID user.AppUserID, workbookID domain.WorkbookID, param domain.ProblemIDsCondition) (*domain.ProblemSearchResult, error) {
	var result *domain.ProblemSearchResult
	if err := s.db.Transaction(func(tx *gorm.DB) error {
		repo := s.repo(tx)
		userRepo := s.userRepo(tx)
		student, err := findStudent(ctx, repo, userRepo, organizationID, operatorID)
		if err != nil {
			return err
		}
		workbook, err := student.FindWorkbookByID(ctx, workbookID)
		if err != nil {
			return err
		}
		tmpResult, err := workbook.FindProblemsByProblemIDs(ctx, student, param)
		if err != nil {
			return err
		}
		result = tmpResult
		return nil
	}); err != nil {
		return nil, err
	}
	return result, nil
}

func (s *problemService) FindProblemByID(ctx context.Context, organizationID user.OrganizationID, operatorID user.AppUserID, workbookID domain.WorkbookID, problemID domain.ProblemID) (domain.Problem, error) {
	var result domain.Problem
	if err := s.db.Transaction(func(tx *gorm.DB) error {
		repo := s.repo(tx)
		userRepo := s.userRepo(tx)
		student, err := findStudent(ctx, repo, userRepo, organizationID, operatorID)
		if err != nil {
			return xerrors.Errorf("failed to findStudent. err: %w", err)
		}
		workbook, err := student.FindWorkbookByID(ctx, workbookID)
		if err != nil {
			return err
		}
		tmpResult, err := workbook.FindProblemByID(ctx, student, problemID)
		if err != nil {
			return err
		}
		result = tmpResult
		return nil
	}); err != nil {
		return nil, err
	}
	return result, nil
}

func (s *problemService) FindProblemIDs(ctx context.Context, organizationID user.OrganizationID, operatorID user.AppUserID, workbookID domain.WorkbookID) ([]domain.ProblemID, error) {
	var result []domain.ProblemID
	if err := s.db.Transaction(func(tx *gorm.DB) error {
		repo := s.repo(tx)
		userRepo := s.userRepo(tx)
		student, err := findStudent(ctx, repo, userRepo, organizationID, operatorID)
		if err != nil {
			return xerrors.Errorf("failed to findStudent. err: %w", err)
		}
		workbook, err := student.FindWorkbookByID(ctx, workbookID)
		if err != nil {
			return err
		}
		tmpResult, err := workbook.FindProblemIDs(ctx, student)
		if err != nil {
			return err
		}
		result = tmpResult
		return nil
	}); err != nil {
		return nil, err
	}
	return result, nil
}

func (s *problemService) AddProblem(ctx context.Context, organizationID user.OrganizationID, operatorID user.AppUserID, param domain.ProblemAddParameter) (domain.ProblemID, error) {
	logger := log.FromContext(ctx)
	var result domain.ProblemID
	if err := s.db.Transaction(func(tx *gorm.DB) error {
		repo := s.repo(tx)
		userRepo := s.userRepo(tx)
		student, err := findStudent(ctx, repo, userRepo, organizationID, operatorID)
		if err != nil {
			return xerrors.Errorf("failed to findStudent. err: %w", err)
		}
		workbook, err := student.FindWorkbookByID(ctx, param.GetWorkbookID())
		if err != nil {
			return err
		}
		problemType := workbook.GetProblemType()
		sizeLimitName := problemType + "Size"
		updateLimitName := problemType + "Update"
		if err := student.CheckLimit(ctx, sizeLimitName); err != nil {
			return err
		}
		if err := student.CheckLimit(ctx, updateLimitName); err != nil {
			return err
		}
		tmpResult, err := workbook.AddProblem(ctx, student, param)
		if err != nil {
			return err
		}
		if err := student.IncrementQuotaUsage(ctx, sizeLimitName); err != nil {
			return err
		}
		result = tmpResult
		return nil
	}); err != nil {
		return 0, err
	}
	logger.Debug("problem id: %d", result)
	return result, nil
}

func (s *problemService) RemoveProblem(ctx context.Context, organizationID user.OrganizationID, operatorID user.AppUserID, workbookID domain.WorkbookID, problemID domain.ProblemID, version int) error {
	logger := log.FromContext(ctx)
	logger.Debug("ProblemService.RemoveProblem")

	if err := s.db.Transaction(func(tx *gorm.DB) error {
		repo := s.repo(tx)
		userRepo := s.userRepo(tx)
		student, err := findStudent(ctx, repo, userRepo, organizationID, operatorID)
		if err != nil {
			return xerrors.Errorf("failed to findStudent. err: %w", err)
		}
		workbook, err := student.FindWorkbookByID(ctx, workbookID)
		if err != nil {
			return err
		}
		if err := workbook.RemoveProblem(ctx, student, problemID, version); err != nil {
			return err
		}
		problemType := workbook.GetProblemType()
		sizeLimitName := problemType + "Size"
		if err := student.DecrementQuotaUsage(ctx, sizeLimitName); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return err
	}
	return nil
}
