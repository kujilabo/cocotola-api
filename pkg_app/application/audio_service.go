package application

import (
	"context"
	"errors"
	"strconv"

	"github.com/kujilabo/cocotola-api/pkg_app/domain"
	user "github.com/kujilabo/cocotola-api/pkg_user/domain"
	"golang.org/x/xerrors"
	"gorm.io/gorm"
)

type AudioService interface {
	FindAudioByID(ctx context.Context, organizationID user.OrganizationID, operatorID user.AppUserID, workbookID domain.WorkbookID, problemID domain.ProblemID, audioID domain.AudioID) (domain.Audio, error)
}

type audioService struct {
	db         *gorm.DB
	pf         domain.ProcessorFactory
	rfFunc     domain.RepositoryFactoryFunc
	userRfFunc user.RepositoryFactoryFunc
}

func NewAudioService(db *gorm.DB, pf domain.ProcessorFactory, rfFunc domain.RepositoryFactoryFunc, userRfFunc user.RepositoryFactoryFunc) AudioService {
	return &audioService{
		db:         db,
		pf:         pf,
		rfFunc:     rfFunc,
		userRfFunc: userRfFunc,
	}
}

func (s *audioService) FindAudioByID(ctx context.Context, organizationID user.OrganizationID, operatorID user.AppUserID, workbookID domain.WorkbookID, problemID domain.ProblemID, audioID domain.AudioID) (domain.Audio, error) {
	var result domain.Audio
	if err := s.db.Transaction(func(tx *gorm.DB) error {
		student, workbook, err := s.findStudentAndWorkbook(ctx, tx, organizationID, operatorID, workbookID)
		if err != nil {
			return err
		}
		problem, err := workbook.FindProblemByID(ctx, student, problemID)
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
		result = tmpResult
		return nil
	}); err != nil {
		return nil, err
	}

	return result, nil
}

func (s *audioService) findStudentAndWorkbook(ctx context.Context, tx *gorm.DB, organizationID user.OrganizationID, operatorID user.AppUserID, workbookID domain.WorkbookID) (domain.Student, domain.Workbook, error) {
	repo, err := s.rfFunc(ctx, tx)
	if err != nil {
		return nil, nil, err
	}
	userRepo, err := s.userRfFunc(ctx, tx)
	if err != nil {
		return nil, nil, err
	}
	student, err := findStudent(ctx, s.pf, repo, userRepo, organizationID, operatorID)
	if err != nil {
		return nil, nil, xerrors.Errorf("failed to findStudent. err: %w", err)
	}
	workbook, err := student.FindWorkbookByID(ctx, workbookID)
	if err != nil {
		return nil, nil, err
	}
	return student, workbook, nil
}
