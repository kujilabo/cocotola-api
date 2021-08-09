package domain

import (
	"context"

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

	CheckLimit(ctx context.Context, name string) error

	IncrementQuotaUsage(ctx context.Context, name string) error

	DecrementQuotaUsage(ctx context.Context, name string) error

	FindRecordbook(ctx context.Context, workbookID WorkbookID, studyType string) (Recordbook, error)
}

type student struct {
	user.AppUser
	rf       RepositoryFactory
	userRepo user.RepositoryFactory
}

func NewStudent(repositoryFactory RepositoryFactory, userRepo user.RepositoryFactory, appUser user.AppUser) (Student, error) {
	m := &student{
		AppUser:  appUser,
		rf:       repositoryFactory,
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
	workbookRepo, err := s.rf.NewWorkbookRepository(ctx)
	if err != nil {
		return nil, xerrors.Errorf("failed to NewWorkbookRepository. err: %w", err)
	}

	return workbookRepo.FindWorkbookByName(ctx, s, name)
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

func (s *student) CheckLimit(ctx context.Context, name string) error {
	// config, ok := s.quotaLimitConfigs[name]
	// if !ok {
	// 	return fmt.Errorf("NotFound %s", name)
	// }

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
	return nil
}

func (s *student) IncrementQuotaUsage(ctx context.Context, name string) error {
	// config, ok := s.quotaLimitConfigs[name]
	// if !ok {
	// 	return fmt.Errorf("NotFound %s", name)
	// }

	// switch name {
	// case EnglishWordProblemUpdateLimit, EnglishPhraseProblemUpdateLimit, TemplateProblemUpdateLimit:
	// 	exceeded, err := s.rf.NewQuotaLimitRepository().Increment(ctx, s.OrganizationID(), s.ID(), name, 1, config.Unit, config.Limit)
	// 	if err != nil {
	// 		return err
	// 	}
	// 	if exceeded {
	// 		return NewLimitExceededError(name)
	// 	}
	// 	return nil
	// default:
	// 	return errors.New("xxxxxxx")
	// }
	return nil
}

func (s *student) DecrementQuotaUsage(ctx context.Context, name string) error {
	return nil
}

func (s *student) FindRecordbook(ctx context.Context, workbookID WorkbookID, studyType string) (Recordbook, error) {
	repo, err := s.rf.NewStudyResultRepository(ctx)
	if err != nil {
		return nil, err
	}

	studyResults, err := repo.FindStudyResults(ctx, s, workbookID, studyType)
	if err != nil {
		return nil, err
	}

	workbook, err := s.FindWorkbookByID(ctx, workbookID)
	if err != nil {
		return nil, err
	}

	problemIDs, err := workbook.FindProblemIDs(ctx, s)
	if err != nil {
		return nil, err
	}

	results := make(map[ProblemID]int)
	for _, problemID := range problemIDs {
		if level, ok := studyResults[problemID]; ok {
			results[problemID] = level
		} else {
			results[problemID] = 0
		}
	}

	return NewRecordbook(workbookID, results)
}
