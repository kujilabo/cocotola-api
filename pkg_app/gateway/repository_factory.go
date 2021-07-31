package gateway

import (
	"context"

	"golang.org/x/xerrors"
	"gorm.io/gorm"

	"github.com/kujilabo/cocotola-api/pkg_app/domain"
	user "github.com/kujilabo/cocotola-api/pkg_user/domain"
)

type repositoryFactory struct {
	db                  *gorm.DB
	userRepo            func(db *gorm.DB) user.RepositoryFactory
	pf                  domain.ProcessorFactory
	problemRepositories map[string]func(*gorm.DB) (domain.ProblemRepository, error)
}

func NewRepositoryFactory(db *gorm.DB, userRepo func(db *gorm.DB) user.RepositoryFactory, pf domain.ProcessorFactory, problemRepositories map[string]func(*gorm.DB) (domain.ProblemRepository, error)) domain.RepositoryFactory {
	return &repositoryFactory{
		db:                  db,
		userRepo:            userRepo,
		pf:                  pf,
		problemRepositories: problemRepositories,
	}
}

func (f *repositoryFactory) NewWorkbookRepository(ctx context.Context) (domain.WorkbookRepository, error) {
	return NewWorkbookRepository(ctx, f, f.userRepo(f.db), f.pf, f.db)
}

func (f *repositoryFactory) NewProblemRepository(ctx context.Context, problemType string) (domain.ProblemRepository, error) {
	problemRepository, ok := f.problemRepositories[problemType]
	if !ok {
		return nil, xerrors.Errorf("problem repository not found. problemType: %s", problemType)
	}
	return problemRepository(f.db)
}

func (f *repositoryFactory) NewProblemTypeRepository(ctx context.Context) (domain.ProblemTypeRepository, error) {
	return NewProblemTypeRepository(f.db)
}

func (f *repositoryFactory) NewStudyTypeRepository(ctx context.Context) (domain.StudyTypeRepository, error) {
	return NewStudyTypeRepository(f.db)
}

func (f *repositoryFactory) NewStudyResultRepository(ctx context.Context) (domain.StudyResultRepository, error) {
	return NewStudyResultRepository(ctx, f, f.db)
}

func (f *repositoryFactory) NewAudioRepository(ctx context.Context) (domain.AudioRepository, error) {
	return NewAudioRepository(f.db), nil
}
