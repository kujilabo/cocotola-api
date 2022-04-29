//go:generate mockery --output mock --name SystemStudent
package service

import (
	"context"
	"errors"

	"golang.org/x/xerrors"

	"github.com/kujilabo/cocotola-api/src/app/domain"
	libD "github.com/kujilabo/cocotola-api/src/lib/domain"
	userD "github.com/kujilabo/cocotola-api/src/user/domain"
)

type SystemStudent interface {
	userD.AppUserModel

	FindWorkbookFromSystemSpace(ctx context.Context, name string) (Workbook, error)

	AddWorkbookToSystemSpace(ctx context.Context, parameter WorkbookAddParameter) (domain.WorkbookID, error)
}

type systemStudent struct {
	userD.AppUserModel
	rf RepositoryFactory
}

func NewSystemStudent(rf RepositoryFactory, appUser userD.AppUserModel) (SystemStudent, error) {
	m := &systemStudent{
		AppUserModel: appUser,
		rf:           rf,
	}

	return m, libD.Validator.Struct(m)
}

func (s *systemStudent) FindWorkbookFromSystemSpace(ctx context.Context, name string) (Workbook, error) {
	systemSpaceID := GetSystemSpaceID()

	workbookRepo, err := s.rf.NewWorkbookRepository(ctx)
	if err != nil {
		return nil, xerrors.Errorf("failed to NewWorkbookRepository. err: %w", err)
	}

	workbook, err := workbookRepo.FindWorkbookByName(ctx, s, systemSpaceID, name)
	if err != nil {
		return nil, xerrors.Errorf("failed to FindWorkbookByName. err: %w", err)
	}

	return workbook, nil
}

func (s *systemStudent) AddWorkbookToSystemSpace(ctx context.Context, parameter WorkbookAddParameter) (domain.WorkbookID, error) {
	systemSpaceID := GetSystemSpaceID()
	if uint(systemSpaceID) == 0 {
		return 0, errors.New("invalid system space ID")
	}

	workbookRepo, err := s.rf.NewWorkbookRepository(ctx)
	if err != nil {
		return 0, xerrors.Errorf("failed to NewWorkbookRepository. err: %w", err)
	}

	workbookID, err := workbookRepo.AddWorkbook(ctx, s, systemSpaceID, parameter)
	if err != nil {
		return 0, xerrors.Errorf("failed to AddWorkbook. err: %w", err)
	}

	return workbookID, nil
}
