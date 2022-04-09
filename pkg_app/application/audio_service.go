package application

import (
	"context"
	"errors"
	"strconv"

	"golang.org/x/xerrors"
	"gorm.io/gorm"

	"github.com/kujilabo/cocotola-api/pkg_app/domain"
	"github.com/kujilabo/cocotola-api/pkg_app/service"
	user "github.com/kujilabo/cocotola-api/pkg_user/domain"
	userS "github.com/kujilabo/cocotola-api/pkg_user/service"
)

type AudioService interface {
	FindAudioByID(ctx context.Context, organizationID user.OrganizationID, operatorID user.AppUserID, workbookID domain.WorkbookID, problemID domain.ProblemID, audioID domain.AudioID) (domain.AudioModel, error)
}

type audioService struct {
	db         *gorm.DB
	pf         service.ProcessorFactory
	rfFunc     service.RepositoryFactoryFunc
	userRfFunc userS.RepositoryFactoryFunc
}

func NewAudioService(db *gorm.DB, pf service.ProcessorFactory, rfFunc service.RepositoryFactoryFunc, userRfFunc userS.RepositoryFactoryFunc) AudioService {
	return &audioService{
		db:         db,
		pf:         pf,
		rfFunc:     rfFunc,
		userRfFunc: userRfFunc,
	}
}

func (s *audioService) FindAudioByID(ctx context.Context, organizationID user.OrganizationID, operatorID user.AppUserID, workbookID domain.WorkbookID, problemID domain.ProblemID, audioID domain.AudioID) (domain.AudioModel, error) {
	var result domain.AudioModel
	if err := s.db.Transaction(func(tx *gorm.DB) error {
		student, workbookService, err := s.findStudentAndWorkbook(ctx, tx, organizationID, operatorID, workbookID)
		if err != nil {
			return err
		}
		problem, err := workbookService.FindProblemByID(ctx, student, problemID)
		if err != nil {
			return err
		}
		if strconv.Itoa(int(audioID)) != problem.GetProperties(ctx)["audioId"] {
			return errors.New("invalid audio id")
		}
		// tmpResult, err := problem.FindAudioByID(ctx, audioID)
		// if err != nil {
		// 	return err
		// }
		rf, err := s.rfFunc(ctx, tx)
		if err != nil {
			return err
		}
		audioRepo, err := rf.NewAudioRepository(ctx)
		if err != nil {
			return err
		}
		tmpResult, err := audioRepo.FindAudioByAudioID(ctx, audioID)
		if err != nil {
			return err
		}
		result = tmpResult.GetAudioModel()
		return nil
	}); err != nil {
		return nil, err
	}

	return result, nil
}

func (s *audioService) findStudentAndWorkbook(ctx context.Context, tx *gorm.DB, organizationID user.OrganizationID, operatorID user.AppUserID, workbookID domain.WorkbookID) (service.Student, service.Workbook, error) {
	repo, err := s.rfFunc(ctx, tx)
	if err != nil {
		return nil, nil, err
	}
	userRepo, err := s.userRfFunc(ctx, tx)
	if err != nil {
		return nil, nil, err
	}
	studentService, err := findStudent(ctx, s.pf, repo, userRepo, organizationID, operatorID)
	if err != nil {
		return nil, nil, xerrors.Errorf("failed to findStudent. err: %w", err)
	}
	workbookService, err := studentService.FindWorkbookByID(ctx, workbookID)
	if err != nil {
		return nil, nil, err
	}
	return studentService, workbookService, nil
}
