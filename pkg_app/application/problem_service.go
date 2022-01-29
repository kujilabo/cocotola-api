package application

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"

	"github.com/kujilabo/cocotola-api/pkg_app/domain"
	"github.com/kujilabo/cocotola-api/pkg_lib/log"
	user "github.com/kujilabo/cocotola-api/pkg_user/domain"
)

type ProblemService interface {
	FindProblemsByWorkbookID(ctx context.Context, organizationID user.OrganizationID, operatorID user.AppUserID, workbookID domain.WorkbookID, param domain.ProblemSearchCondition) (*domain.ProblemSearchResult, error)

	FindAllProblemsByWorkbookID(ctx context.Context, organizationID user.OrganizationID, operatorID user.AppUserID, workbookID domain.WorkbookID) (*domain.ProblemSearchResult, error)

	FindProblemsByProblemIDs(ctx context.Context, organizationID user.OrganizationID, operatorID user.AppUserID, workbookID domain.WorkbookID, param domain.ProblemIDsCondition) (*domain.ProblemSearchResult, error)

	FindProblemByID(ctx context.Context, organizationID user.OrganizationID, operatorID user.AppUserID, workbookID domain.WorkbookID, problemID domain.ProblemID) (domain.Problem, error)

	FindProblemIDs(ctx context.Context, organizationID user.OrganizationID, operatorID user.AppUserID, workbookID domain.WorkbookID) ([]domain.ProblemID, error)

	AddProblem(ctx context.Context, organizationID user.OrganizationID, operatorID user.AppUserID, param domain.ProblemAddParameter) (domain.ProblemID, error)

	UpdateProblem(ctx context.Context, organizationID user.OrganizationID, operatorID user.AppUserID, param domain.ProblemUpdateParameter) error

	RemoveProblem(ctx context.Context, organizationID user.OrganizationID, operatorID user.AppUserID, workbookID domain.WorkbookID, problemID domain.ProblemID, version int) error

	ImportProblems(ctx context.Context, organizationID user.OrganizationID, operatorID user.AppUserID, workbookID domain.WorkbookID, newIterator func(workbookID domain.WorkbookID, problemType string) (domain.ProblemAddParameterIterator, error)) error
}

type problemService struct {
	db       *gorm.DB
	pf       domain.ProcessorFactory
	repo     func(db *gorm.DB) (domain.RepositoryFactory, error)
	userRepo func(db *gorm.DB) (user.RepositoryFactory, error)
}

func NewProblemService(db *gorm.DB, pf domain.ProcessorFactory, repo func(db *gorm.DB) (domain.RepositoryFactory, error), userRepo func(db *gorm.DB) (user.RepositoryFactory, error)) ProblemService {
	return &problemService{
		db:       db,
		pf:       pf,
		repo:     repo,
		userRepo: userRepo,
	}
}

func (s *problemService) FindProblemsByWorkbookID(ctx context.Context, organizationID user.OrganizationID, operatorID user.AppUserID, workbookID domain.WorkbookID, param domain.ProblemSearchCondition) (*domain.ProblemSearchResult, error) {
	var result *domain.ProblemSearchResult
	if err := s.db.Transaction(func(tx *gorm.DB) error {
		student, workbook, err := s.findStudentAndWorkbook(ctx, tx, organizationID, operatorID, workbookID)
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

func (s *problemService) FindAllProblemsByWorkbookID(ctx context.Context, organizationID user.OrganizationID, operatorID user.AppUserID, workbookID domain.WorkbookID) (*domain.ProblemSearchResult, error) {
	var result *domain.ProblemSearchResult
	if err := s.db.Transaction(func(tx *gorm.DB) error {
		student, workbook, err := s.findStudentAndWorkbook(ctx, tx, organizationID, operatorID, workbookID)
		if err != nil {
			return err
		}
		tmpResult, err := workbook.FindAllProblems(ctx, student)
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
		student, workbook, err := s.findStudentAndWorkbook(ctx, tx, organizationID, operatorID, workbookID)
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
		student, workbook, err := s.findStudentAndWorkbook(ctx, tx, organizationID, operatorID, workbookID)
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
		student, workbook, err := s.findStudentAndWorkbook(ctx, tx, organizationID, operatorID, workbookID)
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
		student, workbook, err := s.findStudentAndWorkbook(ctx, tx, organizationID, operatorID, param.GetWorkbookID())
		if err != nil {
			return err
		}
		tmpResult, err := s.addProblem(ctx, student, workbook, param)
		if err != nil {
			return fmt.Errorf("failed to AddProblem. err: %w", err)
		}
		result = tmpResult
		return nil
	}); err != nil {
		return 0, err
	}
	logger.Debug("problem id: %d", result)
	return result, nil
}

func (s *problemService) UpdateProblem(ctx context.Context, organizationID user.OrganizationID, operatorID user.AppUserID, param domain.ProblemUpdateParameter) error {
	if err := s.db.Transaction(func(tx *gorm.DB) error {
		student, workbook, err := s.findStudentAndWorkbook(ctx, tx, organizationID, operatorID, param.GetWorkbookID())
		if err != nil {
			return err
		}
		if err := s.updateProblem(ctx, student, workbook, param); err != nil {
			return fmt.Errorf("failed to UpdateProblem. err: %w", err)
		}
		return nil
	}); err != nil {
		return err
	}
	return nil
}

func (s *problemService) RemoveProblem(ctx context.Context, organizationID user.OrganizationID, operatorID user.AppUserID, workbookID domain.WorkbookID, problemID domain.ProblemID, version int) error {
	logger := log.FromContext(ctx)
	logger.Debug("ProblemService.RemoveProblem")

	if err := s.db.Transaction(func(tx *gorm.DB) error {
		student, workbook, err := s.findStudentAndWorkbook(ctx, tx, organizationID, operatorID, workbookID)
		if err != nil {
			return err
		}
		if err := workbook.RemoveProblem(ctx, student, problemID, version); err != nil {
			return err
		}
		problemType := workbook.GetProblemType()
		if err := student.DecrementQuotaUsage(ctx, problemType, "Size", 1); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return err
	}
	return nil
}

func (s *problemService) ImportProblems(ctx context.Context, organizationID user.OrganizationID, operatorID user.AppUserID, workbookID domain.WorkbookID, newIterator func(workbookID domain.WorkbookID, problemType string) (domain.ProblemAddParameterIterator, error)) error {
	logger := log.FromContext(ctx)
	logger.Debug("ProblemService.ImportProblems")

	var problemType string
	{
		_, workbook, err := s.findStudentAndWorkbook(ctx, s.db, organizationID, operatorID, workbookID)
		if err != nil {
			return err
		}
		problemType = workbook.GetProblemType()
	}
	iterator, err := newIterator(workbookID, problemType)
	if err != nil {
		return err
	}

	for {
		param, err := iterator.Next()
		if err != nil {
			return err
		}
		if param == nil {
			continue
		}

		logger.Infof("param.properties: %+v", param.GetProperties())

		if err := s.db.Transaction(func(tx *gorm.DB) error {
			student, workbook, err := s.findStudentAndWorkbook(ctx, tx, organizationID, operatorID, workbookID)
			if err != nil {
				return err
			}

			id, err := s.addProblem(ctx, student, workbook, param)
			if errors.Is(err, domain.ErrProblemAlreadyExists) {
				logger.Infof("Problem already exists. param: %+v", param)
				return nil
			}

			if err != nil {
				return fmt.Errorf("failed to addProblem. err: %w", err)
			}
			logger.Infof("%d", id)

			return nil
		}); err != nil {
			return err
		}
	}
}

func (s *problemService) findStudentAndWorkbook(ctx context.Context, tx *gorm.DB, organizationID user.OrganizationID, operatorID user.AppUserID, workbookID domain.WorkbookID) (domain.Student, domain.Workbook, error) {
	repo, err := s.repo(tx)
	if err != nil {
		return nil, nil, err
	}
	userRepo, err := s.userRepo(tx)
	if err != nil {
		return nil, nil, err
	}
	student, err := findStudent(ctx, s.pf, repo, userRepo, organizationID, operatorID)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to findStudent. err: %w", err)
	}
	workbook, err := student.FindWorkbookByID(ctx, workbookID)
	if err != nil {
		return nil, nil, err
	}
	return student, workbook, nil
}

func (s *problemService) addProblem(ctx context.Context, student domain.Student, workbook domain.Workbook, param domain.ProblemAddParameter) (domain.ProblemID, error) {
	problemType := workbook.GetProblemType()
	if err := student.CheckQuota(ctx, problemType, "Size"); err != nil {
		return 0, err
	}
	if err := student.CheckQuota(ctx, problemType, "Update"); err != nil {
		return 0, err
	}
	added, id, err := workbook.AddProblem(ctx, student, param)
	if err != nil {
		return 0, err
	}
	if err := student.IncrementQuotaUsage(ctx, problemType, "Size", int(added)); err != nil {
		return 0, err
	}
	if err := student.IncrementQuotaUsage(ctx, problemType, "Update", int(added)); err != nil {
		return 0, err
	}
	return id, nil
}

func (s *problemService) updateProblem(ctx context.Context, student domain.Student, workbook domain.Workbook, param domain.ProblemUpdateParameter) error {
	problemType := workbook.GetProblemType()
	if err := student.CheckQuota(ctx, problemType, "Size"); err != nil {
		return err
	}
	if err := student.CheckQuota(ctx, problemType, "Update"); err != nil {
		return err
	}
	added, updated, err := workbook.UpdateProblem(ctx, student, param)
	if err != nil {
		return err
	}
	if added > 0 {
		if err := student.IncrementQuotaUsage(ctx, problemType, "Size", int(added)); err != nil {
			return err
		}
	} else if added < 0 {
		if err := student.DecrementQuotaUsage(ctx, problemType, "Size", -int(added)); err != nil {
			return err
		}
	}
	if updated > 0 {
		if err := student.IncrementQuotaUsage(ctx, problemType, "Update", int(updated)); err != nil {
			return err
		}
	}
	return nil
}
