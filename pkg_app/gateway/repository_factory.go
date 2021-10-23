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
	driverName          string
	userRepo            func(db *gorm.DB) (user.RepositoryFactory, error)
	pf                  domain.ProcessorFactory
	problemRepositories map[string]func(*gorm.DB) (domain.ProblemRepository, error)
	problemTypes        []domain.ProblemType
}

func NewRepositoryFactory(ctx context.Context, db *gorm.DB, driverName string, userRepo func(db *gorm.DB) (user.RepositoryFactory, error), pf domain.ProcessorFactory, problemRepositories map[string]func(*gorm.DB) (domain.ProblemRepository, error)) (domain.RepositoryFactory, error) {
	problemTypeRepo, err := NewProblemTypeRepository(db)
	if err != nil {
		return nil, err
	}
	problemTypes, err := problemTypeRepo.FindAllProblemTypes(ctx)
	if err != nil {
		return nil, err
	}
	return &repositoryFactory{
		db:                  db,
		driverName:          driverName,
		userRepo:            userRepo,
		pf:                  pf,
		problemRepositories: problemRepositories,
		problemTypes:        problemTypes,
	}, nil
}

func (f *repositoryFactory) NewWorkbookRepository(ctx context.Context) (domain.WorkbookRepository, error) {
	// problemTypeRepo, err := f.NewProblemTypeRepository(ctx)
	// if err != nil {
	// 	return nil, err
	// }
	// problemTypes, err := problemTypeRepo.FindAllProblemTypes(ctx)
	// if err != nil {
	// 	return nil, err
	// }
	// logger := log.FromContext(ctx)
	// logger.Infof("problem types: %+v", problemTypes)
	userRepo, err := f.userRepo(f.db)
	if err != nil {
		return nil, err
	}
	return NewWorkbookRepository(ctx, f.driverName, f, userRepo, f.pf, f.db, f.problemTypes), nil
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
	return NewStudyResultRepository(ctx, f, f.db, f.problemTypes)
}

func (f *repositoryFactory) NewAudioRepository(ctx context.Context) (domain.AudioRepository, error) {
	return NewAudioRepository(f.db), nil
}
