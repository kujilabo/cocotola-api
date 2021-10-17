package application

import (
	"context"

	"golang.org/x/xerrors"
	"gorm.io/gorm"

	"github.com/kujilabo/cocotola-api/pkg_app/domain"
	user "github.com/kujilabo/cocotola-api/pkg_user/domain"
)

type StudyService interface {
	FindResults(ctx context.Context, organizationID user.OrganizationID, operatorID user.AppUserID, workbookID domain.WorkbookID, studyType string) ([]domain.ProblemWithLevel, error)

	// FindAllProblemsByWorkbookID(ctx context.Context, organizationID, operatorID, workbookID uint, studyTypeID domain.StudyTypeID) (domain.WorkbookWithProblems, error)
	SetResult(ctx context.Context, organizationID user.OrganizationID, operatorID user.AppUserID, workbookID domain.WorkbookID, studyType string, problemID domain.ProblemID, result bool) error
}

type studyService struct {
	db       *gorm.DB
	repo     func(db *gorm.DB) (domain.RepositoryFactory, error)
	userRepo func(db *gorm.DB) (user.RepositoryFactory, error)
}

func NewStudyService(db *gorm.DB, repo func(db *gorm.DB) (domain.RepositoryFactory, error), userRepo func(db *gorm.DB) (user.RepositoryFactory, error)) StudyService {
	return &studyService{
		db:       db,
		repo:     repo,
		userRepo: userRepo,
	}
}

func (s *studyService) FindResults(ctx context.Context, organizationID user.OrganizationID, operatorID user.AppUserID, workbookID domain.WorkbookID, studyType string) ([]domain.ProblemWithLevel, error) {
	var results []domain.ProblemWithLevel
	if err := s.db.Transaction(func(tx *gorm.DB) error {
		repo, err := s.repo(tx)
		if err != nil {
			return xerrors.Errorf("failed to repo. err: %w", err)
		}
		userRepo, err := s.userRepo(tx)
		if err != nil {
			return xerrors.Errorf("failed to userRepo. err: %w", err)
		}
		student, err := findStudent(ctx, repo, userRepo, organizationID, operatorID)
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

func (s *studyService) SetResult(ctx context.Context, organizationID user.OrganizationID, operatorID user.AppUserID, workbookID domain.WorkbookID, studyType string, problemID domain.ProblemID, result bool) error {
	if err := s.db.Transaction(func(tx *gorm.DB) error {
		repo, err := s.repo(tx)
		if err != nil {
			return xerrors.Errorf("failed to repo. err: %w", err)
		}
		userRepo, err := s.userRepo(tx)
		if err != nil {
			return xerrors.Errorf("failed to userRepo. err: %w", err)
		}
		student, err := findStudent(ctx, repo, userRepo, organizationID, operatorID)
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
		if err := tmpResult.SetResult(ctx, workbook.GetProblemType(), problemID, result); err != nil {
			return xerrors.Errorf("failed to SetResult. err: %w", err)
		}
		return nil
	}); err != nil {
		return err
	}
	return nil
}

// func (s *studyService) FindAllProblemsByWorkbookID(ctx context.Context, organizationID, operatorID, workbookID uint, studyTypeID domain.StudyTypeID) (domain.WorkbookWithProblems, error) {
// 	student, err := findStudent(ctx, s.repositoryFactory, organizationID, operatorID)
// 	if err != nil {
// 		return nil, err
// 	}
// 	workbook, err := student.FindWorkbookByID(ctx, workbookID)
// 	if err != nil {
// 		return nil, err
// 	}
// 	problems, err := workbook.FindAllProblems(ctx)
// 	if err != nil {
// 		return nil, err
// 	}
// 	problemWithResultsList := make([]domain.ProblemWithResults, 0)
// 	for _, p := range problems {
// 		problemWithResultsList = append(problemWithResultsList, domain.NewProblemWithResults(p, []bool{}, 0))
// 	}
// 	return domain.NewWorkbookWithProblems(workbook, problemWithResultsList), nil
// }
