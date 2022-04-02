package domain

import (
	lib "github.com/kujilabo/cocotola-api/pkg_lib/domain"
	user "github.com/kujilabo/cocotola-api/pkg_user/domain"
)

type SystemStudentModel interface {
	user.AppUserModel

	// FindWorkbookFromSystemSpace(ctx context.Context, name string) (Workbook, error)

	// AddWorkbookToSystemSpace(ctx context.Context, parameter WorkbookAddParameter) (WorkbookID, error)
}

type systemStudentModel struct {
	user.AppUserModel
	// rf RepositoryFactory
}

func NewSystemStudentModel(
	// rf RepositoryFactory,
	appUser user.AppUserModel) (SystemStudentModel, error) {
	m := &systemStudentModel{
		AppUserModel: appUser,
		// rf:      rf,
	}

	return m, lib.Validator.Struct(m)
}

// func (s *systemStudent) FindWorkbookFromSystemSpace(ctx context.Context, name string) (Workbook, error) {
// 	systemSpaceID := GetSystemSpaceID()

// 	workbookRepo, err := s.rf.NewWorkbookRepository(ctx)
// 	if err != nil {
// 		return nil, xerrors.Errorf("failed to NewWorkbookRepository. err: %w", err)
// 	}

// 	workbook, err := workbookRepo.FindWorkbookByName(ctx, s, systemSpaceID, name)
// 	if err != nil {
// 		return nil, xerrors.Errorf("failed to FindWorkbookByName. err: %w", err)
// 	}

// 	return workbook, nil
// }

// func (s *systemStudent) AddWorkbookToSystemSpace(ctx context.Context, parameter WorkbookAddParameter) (WorkbookID, error) {
// 	systemSpaceID := GetSystemSpaceID()
// 	if uint(systemSpaceID) == 0 {
// 		return 0, errors.New("invalid system space ID")
// 	}

// 	workbookRepo, err := s.rf.NewWorkbookRepository(ctx)
// 	if err != nil {
// 		return 0, xerrors.Errorf("failed to NewWorkbookRepository. err: %w", err)
// 	}

// 	workbookID, err := workbookRepo.AddWorkbook(ctx, s, systemSpaceID, parameter)
// 	if err != nil {
// 		return 0, xerrors.Errorf("failed to AddWorkbook. err: %w", err)
// 	}

// 	return workbookID, nil
// }
