package student

import (
	"context"
	"errors"

	"golang.org/x/xerrors"
	"gorm.io/gorm"

	"github.com/kujilabo/cocotola-api/src/app/domain"
	"github.com/kujilabo/cocotola-api/src/app/service"
	"github.com/kujilabo/cocotola-api/src/app/usecase"
	libD "github.com/kujilabo/cocotola-api/src/lib/domain"
	"github.com/kujilabo/cocotola-api/src/lib/log"
	userD "github.com/kujilabo/cocotola-api/src/user/domain"
	userS "github.com/kujilabo/cocotola-api/src/user/service"
)

type StudentUsecaseAudio interface {
	FindAudioByID(ctx context.Context, organizationID userD.OrganizationID, operatorID userD.AppUserID, workbookID domain.WorkbookID, problemID domain.ProblemID, audioID domain.AudioID) (service.Audio, error)
}

type studentUsecaseAudio struct {
	db                *gorm.DB
	pf                service.ProcessorFactory
	rfFunc            service.RepositoryFactoryFunc
	userRfFunc        userS.RepositoryFactoryFunc
	synthesizerClient service.SynthesizerClient
}

func NewStudentUsecaseAudio(db *gorm.DB, pf service.ProcessorFactory, rfFunc service.RepositoryFactoryFunc, userRfFunc userS.RepositoryFactoryFunc, synthesizerClient service.SynthesizerClient) StudentUsecaseAudio {
	return &studentUsecaseAudio{
		db:                db,
		pf:                pf,
		rfFunc:            rfFunc,
		userRfFunc:        userRfFunc,
		synthesizerClient: synthesizerClient,
	}
}

func (s *studentUsecaseAudio) FindAudioByID(ctx context.Context, organizationID userD.OrganizationID, operatorID userD.AppUserID, workbookID domain.WorkbookID, problemID domain.ProblemID, audioID domain.AudioID) (service.Audio, error) {
	var result service.Audio
	if err := s.db.Transaction(func(tx *gorm.DB) error {
		student, workbook, err := s.findStudentAndWorkbook(ctx, tx, organizationID, operatorID, workbookID)
		if err != nil {
			return err
		}

		problem, err := workbook.FindProblemByID(ctx, student, problemID)
		if err != nil {
			return err
		}

		savedAudioID, ok := (problem.GetProperties(ctx)["audioId"]).(domain.AudioID)
		if !ok {
			return errors.New("mismatch")
		}

		logger := log.FromContext(ctx)
		if audioID != savedAudioID {
			logger.Debugf("properties: %+v", problem.GetProperties(ctx))
			logger.Warnf("audioID: %d, %s", audioID, problem.GetProperties(ctx)["audioId"])
			message := "invalid audio id"
			return domain.NewPluginError(domain.ErrorType(domain.ErrorTypeClient), message, []string{message}, libD.ErrInvalidArgument)
		}

		tmpResult, err := s.synthesizerClient.FindAudioByAudioID(ctx, audioID)
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
