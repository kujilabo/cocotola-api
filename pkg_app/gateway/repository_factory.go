package gateway

import (
	"context"

	"golang.org/x/xerrors"
	"gorm.io/gorm"

	"github.com/kujilabo/cocotola-api/pkg_app/domain"
	"github.com/kujilabo/cocotola-api/pkg_lib/log"
	user "github.com/kujilabo/cocotola-api/pkg_user/domain"
)

type repositoryFactory struct {
	db                  *gorm.DB
	driverName          string
	userRepo            func(db *gorm.DB) user.RepositoryFactory
	pf                  domain.ProcessorFactory
	problemRepositories map[string]func(*gorm.DB) (domain.ProblemRepository, error)
}

func NewRepositoryFactory(db *gorm.DB, driverName string, userRepo func(db *gorm.DB) user.RepositoryFactory, pf domain.ProcessorFactory, problemRepositories map[string]func(*gorm.DB) (domain.ProblemRepository, error)) domain.RepositoryFactory {
	return &repositoryFactory{
		db:                  db,
		driverName:          driverName,
		userRepo:            userRepo,
		pf:                  pf,
		problemRepositories: problemRepositories,
	}
}

func (f *repositoryFactory) NewWorkbookRepository(ctx context.Context) (domain.WorkbookRepository, error) {
	problemTypeRepo, err := f.NewProblemTypeRepository(ctx)
	if err != nil {
		return nil, err
	}
	problemTypes, err := problemTypeRepo.FindAllProblemTypes(ctx)
	if err != nil {
		return nil, err
	}
	logger := log.FromContext(ctx)
	logger.Infof("problem types: %+v", problemTypes)
	return NewWorkbookRepository(ctx, f.driverName, f, f.userRepo(f.db), f.pf, f.db, problemTypes), nil
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
