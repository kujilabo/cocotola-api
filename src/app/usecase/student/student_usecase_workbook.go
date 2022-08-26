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

const DefaultPageNo = 1
const DefaultPageSize = 10

type StudentUsecaseWorkbook interface {
	FindWorkbooks(ctx context.Context, organizationID userD.OrganizationID, operatorID userD.AppUserID) (service.WorkbookSearchResult, error)

	FindWorkbookByID(ctx context.Context, organizationID userD.OrganizationID, operatorID userD.AppUserID, workBookID domain.WorkbookID) (domain.WorkbookModel, error)

	AddWorkbook(ctx context.Context, organizationID userD.OrganizationID, operatorID userD.AppUserID, parameter service.WorkbookAddParameter) (domain.WorkbookID, error)

	UpdateWorkbook(ctx context.Context, organizationID userD.OrganizationID, operatorID userD.AppUserID, workbookID domain.WorkbookID, version int, parameter service.WorkbookUpdateParameter) error

	RemoveWorkbook(ctx context.Context, organizationID userD.OrganizationID, operatorID userD.AppUserID, workbookID domain.WorkbookID, version int) error
}

type studentUsecaseWorkbook struct {
	db         *gorm.DB
	pf         service.ProcessorFactory
	rfFunc     service.RepositoryFactoryFunc
	userRfFunc userS.RepositoryFactoryFunc
}

func NewStudentUsecaseWorkbook(db *gorm.DB, pf service.ProcessorFactory, rfFunc service.RepositoryFactoryFunc, userRfFunc userS.RepositoryFactoryFunc) StudentUsecaseWorkbook {
	return &studentUsecaseWorkbook{
		db:         db,
		pf:         pf,
		rfFunc:     rfFunc,
		userRfFunc: userRfFunc,
	}
}

func (s *studentUsecaseWorkbook) FindWorkbooks(ctx context.Context, organizationID userD.OrganizationID, operatorID userD.AppUserID) (service.WorkbookSearchResult, error) {
	var result service.WorkbookSearchResult
	if err := s.db.Transaction(func(tx *gorm.DB) error {
		rf, err := s.rfFunc(ctx, tx)
		if err != nil {
			return err
		}
		userRf, err := s.userRfFunc(ctx, tx)
		if err != nil {
			return err
		}
		student, err := usecase.FindStudent(ctx, s.pf, rf, userRf, organizationID, operatorID)
		if err != nil {
			return liberrors.Errorf("failed to findStudent. err: %w", err)
		}

		condition, err := service.NewWorkbookSearchCondition(DefaultPageNo, DefaultPageSize, []userD.SpaceID{})
		if err != nil {
			return liberrors.Errorf("failed to NewWorkbookSearchCondition. err: %w", err)
		}

		tmpResult, err := student.FindWorkbooksFromPersonalSpace(ctx, condition)
		if err != nil {
			return liberrors.Errorf("failed to FindWorkbooksFromPersonalSpace. err: %w", err)
		}

		result = tmpResult
		return nil
	}); err != nil {
		return nil, err
	}
	return result, nil
}

func (s *studentUsecaseWorkbook) FindWorkbookByID(ctx context.Context, organizationID userD.OrganizationID, operatorID userD.AppUserID, workBookID domain.WorkbookID) (domain.WorkbookModel, error) {
	var result domain.WorkbookModel
	if err := s.db.Transaction(func(tx *gorm.DB) error {
		rf, err := s.rfFunc(ctx, tx)
		if err != nil {
			return err
		}
		userRf, err := s.userRfFunc(ctx, tx)
		if err != nil {
			return err
		}
		student, err := usecase.FindStudent(ctx, s.pf, rf, userRf, organizationID, operatorID)
		if err != nil {
			return liberrors.Errorf("failed to findStudent. err: %w", err)
		}

		tmpResult, err := student.FindWorkbookByID(ctx, workBookID)
		if err != nil {
			return liberrors.Errorf("failed to FindWorkbookByID. err: %w", err)
		}

		result = tmpResult
		return nil
	}); err != nil {
		return nil, err
	}
	return result, nil
}

func (s *studentUsecaseWorkbook) AddWorkbook(ctx context.Context, organizationID userD.OrganizationID, operatorID userD.AppUserID, parameter service.WorkbookAddParameter) (domain.WorkbookID, error) {
	var result domain.WorkbookID
	if err := s.db.Transaction(func(tx *gorm.DB) error {
		rf, err := s.rfFunc(ctx, tx)
		if err != nil {
			return err
		}
		userRf, err := s.userRfFunc(ctx, tx)
		if err != nil {
			return err
		}
		student, err := usecase.FindStudent(ctx, s.pf, rf, userRf, organizationID, operatorID)
		if err != nil {
			return liberrors.Errorf("failed to findStudent. err: %w", err)
		}

		tmpResult, err := student.AddWorkbookToPersonalSpace(ctx, parameter)
		if err != nil {
			return liberrors.Errorf("failed to AddWorkbookToPersonalSpace. err: %w", err)
		}

		result = tmpResult
		return nil
	}); err != nil {
		return 0, err
	}
	return result, nil
}

func (s *studentUsecaseWorkbook) UpdateWorkbook(ctx context.Context, organizationID userD.OrganizationID, operatorID userD.AppUserID, workbookID domain.WorkbookID, version int, parameter service.WorkbookUpdateParameter) error {
	if err := s.db.Transaction(func(tx *gorm.DB) error {
		rf, err := s.rfFunc(ctx, tx)
		if err != nil {
			return err
		}
		userRf, err := s.userRfFunc(ctx, tx)
		if err != nil {
			return err
		}
		student, err := usecase.FindStudent(ctx, s.pf, rf, userRf, organizationID, operatorID)
		if err != nil {
			return liberrors.Errorf("failed to findStudent. err: %w", err)
		}

		return student.UpdateWorkbook(ctx, workbookID, version, parameter)
	}); err != nil {
		return err
	}
	return nil
}

func (s *studentUsecaseWorkbook) RemoveWorkbook(ctx context.Context, organizationID userD.OrganizationID, operatorID userD.AppUserID, workbookID domain.WorkbookID, version int) error {
	if err := s.db.Transaction(func(tx *gorm.DB) error {
		rf, err := s.rfFunc(ctx, tx)
		if err != nil {
			return err
		}
		userRf, err := s.userRfFunc(ctx, tx)
		if err != nil {
			return err
		}
		student, err := usecase.FindStudent(ctx, s.pf, rf, userRf, organizationID, operatorID)
		if err != nil {
			return liberrors.Errorf("failed to findStudent. err: %w", err)
		}

		return student.RemoveWorkbook(ctx, workbookID, version)
	}); err != nil {
		return err
	}
	return nil
}
