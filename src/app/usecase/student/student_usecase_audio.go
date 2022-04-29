package student

import (
	"context"

	"golang.org/x/xerrors"
	"gorm.io/gorm"

	"github.com/kujilabo/cocotola-api/src/app/domain"
	"github.com/kujilabo/cocotola-api/src/app/service"
	"github.com/kujilabo/cocotola-api/src/app/usecase"
	userD "github.com/kujilabo/cocotola-api/src/user/domain"
	userS "github.com/kujilabo/cocotola-api/src/user/service"
)

type StudentUsecaseAudio interface {
	FindAudioByID(ctx context.Context, organizationID userD.OrganizationID, operatorID userD.AppUserID, workbookID domain.WorkbookID, problemID domain.ProblemID, audioID domain.AudioID) (service.Audio, error)
}

type studentUsecaseAudio struct {
	db         *gorm.DB
	pf         service.ProcessorFactory
	rfFunc     service.RepositoryFactoryFunc
	userRfFunc userS.RepositoryFactoryFunc
}

func NewStudentUsecaseAudio(db *gorm.DB, pf service.ProcessorFactory, rfFunc service.RepositoryFactoryFunc, userRfFunc userS.RepositoryFactoryFunc) StudentUsecaseAudio {
	return &studentUsecaseAudio{
		db:         db,
		pf:         pf,
		rfFunc:     rfFunc,
		userRfFunc: userRfFunc,
	}
}

func (s *studentUsecaseAudio) FindAudioByID(ctx context.Context, organizationID userD.OrganizationID, operatorID userD.AppUserID, workbookID domain.WorkbookID, problemID domain.ProblemID, audioID domain.AudioID) (service.Audio, error) {
	var result service.Audio
	if err := s.db.Transaction(func(tx *gorm.DB) error {
		student, workbookService, err := s.findStudentAndWorkbook(ctx, tx, organizationID, operatorID, workbookID)
		if err != nil {
			return err
		}

		problem, err := workbookService.FindProblemByID(ctx, student, problemID)
		if err != nil {
			return err
		}

		tmpResult, err := problem.FindAudioByAudioID(ctx, audioID)
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

func (s *studentUsecaseAudio) findStudentAndWorkbook(ctx context.Context, tx *gorm.DB, organizationID userD.OrganizationID, operatorID userD.AppUserID, workbookID domain.WorkbookID) (service.Student, service.Workbook, error) {
	repo, err := s.rfFunc(ctx, tx)
	if err != nil {
		return nil, nil, err
	}
	userRepo, err := s.userRfFunc(ctx, tx)
	if err != nil {
		return nil, nil, err
	}
	studentService, err := usecase.FindStudent(ctx, s.pf, repo, userRepo, organizationID, operatorID)
	if err != nil {
		return nil, nil, xerrors.Errorf("failed to findStudent. err: %w", err)
	}
	workbookService, err := studentService.FindWorkbookByID(ctx, workbookID)
	if err != nil {
		return nil, nil, err
	}
	return studentService, workbookService, nil
}
