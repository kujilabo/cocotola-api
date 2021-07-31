package application

import (
	"context"

	"gorm.io/gorm"

	"github.com/kujilabo/cocotola-api/pkg_app/domain"
	user "github.com/kujilabo/cocotola-api/pkg_user/domain"
)

type StudyService interface {
	FindRecordbook(ctx context.Context, organizationID user.OrganizationID, operatorID user.AppUserID, workbookID domain.WorkbookID, studyType string) (domain.Recordbook, error)

	// FindAllProblemsByWorkbookID(ctx context.Context, organizationID, operatorID, workbookID uint, studyTypeID domain.StudyTypeID) (domain.WorkbookWithProblems, error)
	// SetProblemResult(ctx context.Context, organizationID, operatorID, workbookID uint, problemID uint, studyTypeID domain.StudyTypeID, result bool) error
}

type studyService struct {
	db       *gorm.DB
	repo     func(db *gorm.DB) domain.RepositoryFactory
	userRepo func(db *gorm.DB) user.RepositoryFactory
}

func NewStudyService(db *gorm.DB, repo func(db *gorm.DB) domain.RepositoryFactory, userRepo func(db *gorm.DB) user.RepositoryFactory) StudyService {
	return &studyService{
		db:       db,
		repo:     repo,
		userRepo: userRepo,
	}
}

func (s *studyService) FindRecordbook(ctx context.Context, organizationID user.OrganizationID, operatorID user.AppUserID, workbookID domain.WorkbookID, studyType string) (domain.Recordbook, error) {
	var result domain.Recordbook
	if err := s.db.Transaction(func(tx *gorm.DB) error {
		repo := s.repo(tx)
		userRepo := s.userRepo(tx)
		student, err := findStudent(ctx, repo, userRepo, organizationID, operatorID)
		if err != nil {
			return err
		}
		tmpResult, err := student.FindRecordbook(ctx, workbookID, studyType)
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

// func (s *studyService) SetProblemResult(ctx context.Context, organizationID, operatorID, workbookID, problemID uint, studyTypeID domain.StudyTypeID, result bool) error {
// 	student, err := findStudent(ctx, s.repositoryFactory, organizationID, operatorID)
// 	if err != nil {
// 		return err
// 	}
// 	workbook, err := student.FindWorkbookByID(ctx, workbookID)
// 	if err != nil {
// 		return err
// 	}
// 	studyResult, err := student.FindStudyResultByWorkbookID(ctx, workbookID, studyTypeID)
// 	if err != nil {
// 		return err
// 	}
// 	if err := studyResult.SetResult(ctx, workbook.ProblemTypeID(), problemID, result); err != nil {
// 		return err
// 	}
// 	return nil
// }
