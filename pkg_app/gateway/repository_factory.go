package gateway

import (
	"context"

	"golang.org/x/xerrors"
	"gorm.io/gorm"

	"github.com/kujilabo/cocotola-api/pkg_app/domain"
	"github.com/kujilabo/cocotola-api/pkg_app/service"
	libD "github.com/kujilabo/cocotola-api/pkg_lib/domain"
	"github.com/kujilabo/cocotola-api/pkg_lib/log"
	userS "github.com/kujilabo/cocotola-api/pkg_user/service"
)

type repositoryFactory struct {
	db                  *gorm.DB
	driverName          string
	userRfFunc          userS.RepositoryFactoryFunc
	pf                  service.ProcessorFactory
	problemRepositories map[string]func(context.Context, *gorm.DB) (service.ProblemRepository, error)
	problemTypes        []domain.ProblemType
}

func NewRepositoryFactory(ctx context.Context, db *gorm.DB, driverName string, userRfFunc userS.RepositoryFactoryFunc, pf service.ProcessorFactory, problemRepositories map[string]func(context.Context, *gorm.DB) (service.ProblemRepository, error)) (service.RepositoryFactory, error) {
	if db == nil {
		return nil, libD.ErrInvalidArgument
	}

	problemTypeRepo := NewProblemTypeRepository(db)
	problemTypes, err := problemTypeRepo.FindAllProblemTypes(ctx)
	if err != nil {
		return nil, err
	}

	return &repositoryFactory{
		db:                  db,
		driverName:          driverName,
		userRfFunc:          userRfFunc,
		pf:                  pf,
		problemRepositories: problemRepositories,
		problemTypes:        problemTypes,
	}, nil
}

func (f *repositoryFactory) NewWorkbookRepository(ctx context.Context) (service.WorkbookRepository, error) {
	userRf, err := f.userRfFunc(ctx, f.db)
	if err != nil {
		return nil, err
	}

	return NewWorkbookRepository(ctx, f.driverName, f, userRf, f.pf, f.db, f.problemTypes), nil
}

func (f *repositoryFactory) NewProblemRepository(ctx context.Context, problemType string) (service.ProblemRepository, error) {
	logger := log.FromContext(ctx)
	logger.Infof("problemType: %s", problemType)
	problemRepository, ok := f.problemRepositories[problemType]
	if !ok {
		return nil, xerrors.Errorf("problem repository not found. problemType: %s", problemType)
	}
	return problemRepository(ctx, f.db)
}

func (f *repositoryFactory) NewProblemTypeRepository(ctx context.Context) service.ProblemTypeRepository {
	return NewProblemTypeRepository(f.db)
}

func (f *repositoryFactory) NewStudyTypeRepository(ctx context.Context) service.StudyTypeRepository {
	return NewStudyTypeRepository(f.db)
}

func (f *repositoryFactory) NewRecordbookRepository(ctx context.Context) (service.RecordbookRepository, error) {
	return NewRecordbookRepository(ctx, f, f.db, f.problemTypes)
}

func (f *repositoryFactory) NewAudioRepository(ctx context.Context) service.AudioRepository {
	return NewAudioRepository(f.db)
}

func (f *repositoryFactory) NewUserQuotaRepository(ctx context.Context) service.UserQuotaRepository {
	return NewUserQuotaRepository(f.db)
}

type audioRepositoryFactory struct {
	db *gorm.DB
}

func NewAudioRepositoryFactory(ctx context.Context, db *gorm.DB) service.AudioRepositoryFactory {
	return &audioRepositoryFactory{db: db}
}

func (f *audioRepositoryFactory) NewAudioRepository(ctx context.Context) service.AudioRepository {
	return NewAudioRepository(f.db)
}
