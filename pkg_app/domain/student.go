package domain

import (
	"context"
	"fmt"

	"github.com/go-playground/validator"
	"golang.org/x/xerrors"

	user "github.com/kujilabo/cocotola-api/pkg_user/domain"
)

type Student interface {
	user.AppUser

	GetDefaultSpace(ctx context.Context) (user.Space, error)
	GetPersonalSpace(ctx context.Context) (user.Space, error)

	FindWorkbooksFromPersonalSpace(ctx context.Context, condition WorkbookSearchCondition) (*WorkbookSearchResult, error)

	FindWorkbookByID(ctx context.Context, id WorkbookID) (Workbook, error)

	FindWorkbookByName(ctx context.Context, name string) (Workbook, error)

	AddWorkbookToPersonalSpace(ctx context.Context, parameter WorkbookAddParameter) (WorkbookID, error)

	UpdateWorkbook(ctx context.Context, workbookID WorkbookID, version int, parameter WorkbookUpdateParameter) error

	RemoveWorkbook(ctx context.Context, id WorkbookID, version int) error

	CheckQuota(ctx context.Context, problemType string, name QuotaName) error

	IncrementQuotaUsage(ctx context.Context, problemType string, name QuotaName) error

	DecrementQuotaUsage(ctx context.Context, problemType string, name QuotaName) error

	FindRecordbook(ctx context.Context, workbookID WorkbookID, studyType string) (Recordbook, error)
}

type student struct {
	user.AppUser
	rf       RepositoryFactory
	pf       ProcessorFactory
	userRepo user.RepositoryFactory
}

func NewStudent(pf ProcessorFactory, rf RepositoryFactory, userRepo user.RepositoryFactory, appUser user.AppUser) (Student, error) {
	m := &student{
		AppUser:  appUser,
		pf:       pf,
		rf:       rf,
		userRepo: userRepo,
	}

	v := validator.New()
	return m, v.Struct(m)
}

func (s *student) GetDefaultSpace(ctx context.Context) (user.Space, error) {
	return s.userRepo.NewSpaceRepository().FindDefaultSpace(ctx, s)
}

func (s *student) GetPersonalSpace(ctx context.Context) (user.Space, error) {
	return s.userRepo.NewSpaceRepository().FindPersonalSpace(ctx, s)
}

func (s *student) FindWorkbooksFromPersonalSpace(ctx context.Context, condition WorkbookSearchCondition) (*WorkbookSearchResult, error) {
	space, err := s.GetPersonalSpace(ctx)
	if err != nil {
		return nil, xerrors.Errorf("failed to GetPersonalSpace. err: %w", err)
	}

	// specify space
	newCondition, err := NewWorkbookSearchCondition(condition.GetPageNo(), condition.GetPageSize(), []user.SpaceID{user.SpaceID(space.GetID())})
	if err != nil {
		return nil, xerrors.Errorf("failed to NewWorkbookSearchCondition. err: %w", err)
	}

	workbookRepo, err := s.rf.NewWorkbookRepository(ctx)
	if err != nil {
		return nil, xerrors.Errorf("failed to NewWorkbookRepository. err: %w", err)
	}

	return workbookRepo.FindPersonalWorkbooks(ctx, s, newCondition)
}

func (s *student) FindWorkbookByID(ctx context.Context, id WorkbookID) (Workbook, error) {
	workbookRepo, err := s.rf.NewWorkbookRepository(ctx)
	if err != nil {
		return nil, xerrors.Errorf("failed to NewWorkbookRepository. err: %w", err)
	}

	return workbookRepo.FindWorkbookByID(ctx, s, id)
}

func (s *student) FindWorkbookByName(ctx context.Context, name string) (Workbook, error) {
	space, err := s.GetPersonalSpace(ctx)
	if err != nil {
		return nil, xerrors.Errorf("failed to GetPersonalSpace. err: %w", err)
	}

	workbookRepo, err := s.rf.NewWorkbookRepository(ctx)
	if err != nil {
		return nil, xerrors.Errorf("failed to NewWorkbookRepository. err: %w", err)
	}

	return workbookRepo.FindWorkbookByName(ctx, s, user.SpaceID(space.GetID()), name)
}

func (s *student) AddWorkbookToPersonalSpace(ctx context.Context, parameter WorkbookAddParameter) (WorkbookID, error) {
	space, err := s.GetPersonalSpace(ctx)
	if err != nil {
		return 0, xerrors.Errorf("failed to GetPersonalSpace. err: %w", err)
	}

	workbookRepo, err := s.rf.NewWorkbookRepository(ctx)
	if err != nil {
		return 0, xerrors.Errorf("failed to NewWorkbookRepository. err: %w", err)
	}

	workbookID, err := workbookRepo.AddWorkbook(ctx, s, user.SpaceID(space.GetID()), parameter)
	if err != nil {
		return 0, xerrors.Errorf("failed to AddWorkbook. err: %w", err)
	}

	return workbookID, nil
}

func (s *student) UpdateWorkbook(ctx context.Context, workbookID WorkbookID, version int, parameter WorkbookUpdateParameter) error {
	workbook, err := s.FindWorkbookByID(ctx, workbookID)
	if err != nil {
		return err
	}

	return workbook.UpdateWorkbook(ctx, s, version, parameter)
}

func (s *student) RemoveWorkbook(ctx context.Context, workbookID WorkbookID, version int) error {
	workbook, err := s.FindWorkbookByID(ctx, workbookID)
	if err != nil {
		return err
	}

	return workbook.RemoveWorkbook(ctx, s, version)
}

func (s *student) CheckQuota(ctx context.Context, problemType string, name QuotaName) error {
	processor, err := s.pf.NewProblemQuotaProcessor(problemType)
	if err != nil {
		return err
	}

	userQuotaRepo, err := s.rf.NewUserQuotaRepository(ctx)
	if err != nil {
		return xerrors.Errorf("failed to NewProblemRepository. err: %w", err)
	}

	switch name {
	case QuotaNameSize:
		unit := processor.GetUnitForSizeQuota()
		limit := processor.GetLimitForSizeQuota()
		isExceeded, err := userQuotaRepo.IsExceeded(ctx, s, problemType+"_size", unit, limit)
		if err != nil {
			return err
		}

		if isExceeded {
			return ErrQuotaExceeded
		}

		return nil
	case QuotaNameUpdate:
		unit := processor.GetUnitForUpdateQuota()
		limit := processor.GetLimitForUpdateQuota()
		isExceeded, err := userQuotaRepo.IsExceeded(ctx, s, problemType+"_update", unit, limit)
		if err != nil {
			return err
		}

		if isExceeded {
			return ErrQuotaExceeded
		}

		return nil
	default:
		return fmt.Errorf("invalid name. name: %s", name)
	}

	// isExceeded, err := processor.IsExceeded(ctx, s.rf, s, name)
	// if err != nil {
	// 	return err
	// }

	// if isExceeded {
	// 	return ErrQuotaExceeded
	// }

	// return nil

	// switch name {
	// case EnglishWordProblemSizeLimit:
	// 	count, err := s.rf.NewEnglishWordProblemRepository().CountByAppUser(ctx, s)
	// 	if err != nil {
	// 		return err
	// 	}
	// 	if count > config.Limit {
	// 		return NewLimitExceededError(name)
	// 	}
	// 	return nil
	// case EnglishPhraseProblemSizeLimit:
	// 	count, err := s.rf.NewEnglishPhraseProblemRepository().CountByAppUser(ctx, s)
	// 	if err != nil {
	// 		return err
	// 	}
	// 	if count > config.Limit {
	// 		return NewLimitExceededError(name)
	// 	}
	// 	return nil
	// case TemplateProblemSizeLimit:
	// 	count, err := s.rf.NewTemplateProblemRepository().CountByAppUser(ctx, s)
	// 	if err != nil {
	// 		return err
	// 	}
	// 	if count > config.Limit {
	// 		return NewLimitExceededError(name)
	// 	}
	// 	return nil
	// case EnglishWordProblemUpdateLimit, EnglishPhraseProblemUpdateLimit, TemplateProblemUpdateLimit:
	// 	exceeded, err := s.rf.NewQuotaLimitRepository().IsExceeded(ctx, s.OrganizationID(), s.ID(), name, config.Unit, config.Limit)
	// 	if err != nil {
	// 		return err
	// 	}
	// 	if exceeded {
	// 		return NewLimitExceededError(name)
	// 	}
	// 	return nil
	// default:
	// 	return fmt.Errorf("Invalid name. name2: %s", name)
	// }
	// return nil
}

func (s *student) IncrementQuotaUsage(ctx context.Context, problemType string, name QuotaName) error {
	processor, err := s.pf.NewProblemQuotaProcessor(problemType)
	if err != nil {
		return err
	}

	userQuotaRepo, err := s.rf.NewUserQuotaRepository(ctx)
	if err != nil {
		return xerrors.Errorf("failed to NewProblemRepository. err: %w", err)
	}

	switch name {
	case QuotaNameSize:
		unit := processor.GetUnitForSizeQuota()
		limit := processor.GetLimitForSizeQuota()
		isExceeded, err := userQuotaRepo.Increment(ctx, s, problemType+"_size", unit, limit, 1)
		if err != nil {
			return err
		}

		if isExceeded {
			return ErrQuotaExceeded
		}

		return nil
	case QuotaNameUpdate:
		unit := processor.GetUnitForUpdateQuota()
		limit := processor.GetLimitForUpdateQuota()
		isExceeded, err := userQuotaRepo.Increment(ctx, s, problemType+"_update", unit, limit, 1)
		if err != nil {
			return err
		}

		if isExceeded {
			return ErrQuotaExceeded
		}

		return nil
	default:
		return fmt.Errorf("invalid name. name: %s", name)
	}
}

func (s *student) DecrementQuotaUsage(ctx context.Context, problemType string, name QuotaName) error {
	return nil
}

func (s *student) FindRecordbook(ctx context.Context, workbookID WorkbookID, studyType string) (Recordbook, error) {
	return NewRecordbook(s.rf, s, workbookID, studyType)
}
