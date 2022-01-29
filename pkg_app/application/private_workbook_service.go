package application

import (
	"context"
	"fmt"

	"gorm.io/gorm"

	"github.com/kujilabo/cocotola-api/pkg_app/domain"
	user "github.com/kujilabo/cocotola-api/pkg_user/domain"
)

const DefaultPageNo = 1
const DefaultPageSize = 10

type PrivateWorkbookService interface {
	FindWorkbooks(ctx context.Context, organizationID user.OrganizationID, operatorID user.AppUserID) (*domain.WorkbookSearchResult, error)

	FindWorkbookByID(ctx context.Context, organizationID user.OrganizationID, operatorID user.AppUserID, workBookID domain.WorkbookID) (domain.Workbook, error)

	AddWorkbook(ctx context.Context, organizationID user.OrganizationID, operatorID user.AppUserID, parameter domain.WorkbookAddParameter) (domain.WorkbookID, error)

	UpdateWorkbook(ctx context.Context, organizationID user.OrganizationID, operatorID user.AppUserID, workbookID domain.WorkbookID, version int, parameter domain.WorkbookUpdateParameter) error

	RemoveWorkbook(ctx context.Context, organizationID user.OrganizationID, operatorID user.AppUserID, workbookID domain.WorkbookID, version int) error
}

type privateWorkbookService struct {
	db       *gorm.DB
	pf       domain.ProcessorFactory
	repo     func(db *gorm.DB) (domain.RepositoryFactory, error)
	userRepo func(db *gorm.DB) (user.RepositoryFactory, error)
}

func NewPrivateWorkbookService(db *gorm.DB, pf domain.ProcessorFactory, repo func(db *gorm.DB) (domain.RepositoryFactory, error), userRepo func(db *gorm.DB) (user.RepositoryFactory, error)) PrivateWorkbookService {
	return &privateWorkbookService{
		db:       db,
		pf:       pf,
		repo:     repo,
		userRepo: userRepo,
	}
}

func (s *privateWorkbookService) FindWorkbooks(ctx context.Context, organizationID user.OrganizationID, operatorID user.AppUserID) (*domain.WorkbookSearchResult, error) {
	var result *domain.WorkbookSearchResult
	if err := s.db.Transaction(func(tx *gorm.DB) error {
		repo, err := s.repo(tx)
		if err != nil {
			return err
		}
		userRepo, err := s.userRepo(tx)
		if err != nil {
			return err
		}
		student, err := findStudent(ctx, s.pf, repo, userRepo, organizationID, operatorID)
		if err != nil {
			return fmt.Errorf("failed to findStudent. err: %w", err)
		}

		condition, err := domain.NewWorkbookSearchCondition(DefaultPageNo, DefaultPageSize, []user.SpaceID{})
		if err != nil {
			return fmt.Errorf("failed to NewWorkbookSearchCondition. err: %w", err)
		}

		tmpResult, err := student.FindWorkbooksFromPersonalSpace(ctx, condition)
		if err != nil {
			return fmt.Errorf("failed to FindWorkbooksFromPersonalSpace. err: %w", err)
		}

		result = tmpResult
		return nil
	}); err != nil {
		return nil, err
	}
	return result, nil
}

func (s *privateWorkbookService) FindWorkbookByID(ctx context.Context, organizationID user.OrganizationID, operatorID user.AppUserID, workBookID domain.WorkbookID) (domain.Workbook, error) {
	var result domain.Workbook
	if err := s.db.Transaction(func(tx *gorm.DB) error {
		repo, err := s.repo(tx)
		if err != nil {
			return err
		}
		userRepo, err := s.userRepo(tx)
		if err != nil {
			return err
		}
		student, err := findStudent(ctx, s.pf, repo, userRepo, organizationID, operatorID)
		if err != nil {
			return fmt.Errorf("failed to findStudent. err: %w", err)
		}

		tmpResult, err := student.FindWorkbookByID(ctx, workBookID)
		if err != nil {
			return fmt.Errorf("failed to FindWorkbookByID. err: %w", err)
		}

		result = tmpResult
		return nil
	}); err != nil {
		return nil, err
	}
	return result, nil
}

func (s *privateWorkbookService) AddWorkbook(ctx context.Context, organizationID user.OrganizationID, operatorID user.AppUserID, parameter domain.WorkbookAddParameter) (domain.WorkbookID, error) {
	var result domain.WorkbookID
	if err := s.db.Transaction(func(tx *gorm.DB) error {
		repo, err := s.repo(tx)
		if err != nil {
			return err
		}
		userRepo, err := s.userRepo(tx)
		if err != nil {
			return err
		}
		student, err := findStudent(ctx, s.pf, repo, userRepo, organizationID, operatorID)
		if err != nil {
			return fmt.Errorf("failed to findStudent. err: %w", err)
		}

		tmpResult, err := student.AddWorkbookToPersonalSpace(ctx, parameter)
		if err != nil {
			return fmt.Errorf("failed to AddWorkbookToPersonalSpace. err: %w", err)
		}

		result = tmpResult
		return nil
	}); err != nil {
		return 0, err
	}
	return result, nil
}

func (s *privateWorkbookService) UpdateWorkbook(ctx context.Context, organizationID user.OrganizationID, operatorID user.AppUserID, workbookID domain.WorkbookID, version int, parameter domain.WorkbookUpdateParameter) error {
	if err := s.db.Transaction(func(tx *gorm.DB) error {
		repo, err := s.repo(tx)
		if err != nil {
			return err
		}
		userRepo, err := s.userRepo(tx)
		if err != nil {
			return err
		}
		student, err := findStudent(ctx, s.pf, repo, userRepo, organizationID, operatorID)
		if err != nil {
			return fmt.Errorf("failed to findStudent. err: %w", err)
		}

		return student.UpdateWorkbook(ctx, workbookID, version, parameter)
	}); err != nil {
		return err
	}
	return nil
}

func (s *privateWorkbookService) RemoveWorkbook(ctx context.Context, organizationID user.OrganizationID, operatorID user.AppUserID, workbookID domain.WorkbookID, version int) error {
	if err := s.db.Transaction(func(tx *gorm.DB) error {
		repo, err := s.repo(tx)
		if err != nil {
			return err
		}
		userRepo, err := s.userRepo(tx)
		if err != nil {
			return err
		}
		student, err := findStudent(ctx, s.pf, repo, userRepo, organizationID, operatorID)
		if err != nil {
			return fmt.Errorf("failed to findStudent. err: %w", err)
		}

		return student.RemoveWorkbook(ctx, workbookID, version)
	}); err != nil {
		return err
	}
	return nil
}
