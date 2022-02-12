package application

import (
	"context"

	"golang.org/x/xerrors"
	"gorm.io/gorm"

	"github.com/kujilabo/cocotola-api/pkg_app/domain"
	"github.com/kujilabo/cocotola-api/pkg_lib/log"
)

type AudioService interface {
	FindAudioByID(ctx context.Context, audioID domain.AudioID) (domain.Audio, error)
}

type audioService struct {
	db     *gorm.DB
	rfFunc func(db *gorm.DB) (domain.RepositoryFactory, error)
}

func NewAudioService(db *gorm.DB, rfFunc func(db *gorm.DB) (domain.RepositoryFactory, error)) AudioService {
	return &audioService{
		db:     db,
		rfFunc: rfFunc,
	}
}

func (s *audioService) FindAudioByID(ctx context.Context, audioID domain.AudioID) (domain.Audio, error) {
	logger := log.FromContext(ctx)
	logger.Infof("audioID: %d", audioID)
	var audio domain.Audio
	if err := s.db.Transaction(func(tx *gorm.DB) error {
		rf, err := s.rfFunc(tx)
		if err != nil {
			return xerrors.Errorf("failed to rf. err: %w", err)
		}
		audioRepo, err := rf.NewAudioRepository(ctx)
		if err != nil {
			return err
		}
		model, err := audioRepo.FindAudioByAudioID(ctx, audioID)
		if err != nil {
			return err
		}

		audio = model
		return nil
	}); err != nil {
		return nil, err
	}
	return audio, nil
}
