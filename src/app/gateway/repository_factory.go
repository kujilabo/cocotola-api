package gateway

import (
	"context"

	"gorm.io/gorm"

	"github.com/kujilabo/cocotola-api/src/app/domain"
	"github.com/kujilabo/cocotola-api/src/app/service"
	libD "github.com/kujilabo/cocotola-api/src/lib/domain"
	liberrors "github.com/kujilabo/cocotola-api/src/lib/errors"
	"github.com/kujilabo/cocotola-api/src/lib/log"
	userS "github.com/kujilabo/cocotola-api/src/user/service"
)

type repositoryFactory struct {
	db                  *gorm.DB
	driverName          string
	userRfFunc          userS.RepositoryFactoryFunc
	pf                  service.ProcessorFactory
	problemRepositories map[string]func(context.Context, *gorm.DB) (service.ProblemRepository, error)
	problemTypes        []domain.ProblemType
	studyTypes          []domain.StudyType
}

func NewRepositoryFactory(ctx context.Context, db *gorm.DB, driverName string, userRfFunc userS.RepositoryFactoryFunc, pf service.ProcessorFactory, problemTypes []domain.ProblemType, studyTypes []domain.StudyType, problemRepositories map[string]func(context.Context, *gorm.DB) (service.ProblemRepository, error)) (service.RepositoryFactory, error) {
	if db == nil {
		return nil, libD.ErrInvalidArgument
	}

	return &repositoryFactory{
		db:                  db,
		driverName:          driverName,
		userRfFunc:          userRfFunc,
		pf:                  pf,
		problemRepositories: problemRepositories,
		problemTypes:        problemTypes,
		studyTypes:          studyTypes,
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
		return nil, liberrors.Errorf("problem repository not found. problemType: %s", problemType)
	}
	return problemRepository(ctx, f.db)
}

func (f *repositoryFactory) NewProblemTypeRepository(ctx context.Context) service.ProblemTypeRepository {
	return NewProblemTypeRepository(f.db)
}

func (f *repositoryFactory) NewStudyTypeRepository(ctx context.Context) service.StudyTypeRepository {
	return NewStudyTypeRepository(f.db)
}

func (f *repositoryFactory) NewRecordbookRepository(ctx context.Context) service.RecordbookRepository {
	return NewRecordbookRepository(ctx, f, f.db, f.problemTypes, f.studyTypes)
}

func (f *repositoryFactory) NewUserQuotaRepository(ctx context.Context) service.UserQuotaRepository {
	return NewUserQuotaRepository(f.db)
}
