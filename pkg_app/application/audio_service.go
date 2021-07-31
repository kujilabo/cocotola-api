package application

import (
	"context"

	"gorm.io/gorm"

	"github.com/kujilabo/cocotola-api/pkg_app/domain"
	"github.com/kujilabo/cocotola-api/pkg_lib/log"
)

type AudioService interface {
	FindAudioByID(ctx context.Context, audioID domain.AudioID) (domain.Audio, error)
}

type audioService struct {
	db   *gorm.DB
	repo func(db *gorm.DB) domain.RepositoryFactory
}

func NewAudioService(db *gorm.DB, repo func(db *gorm.DB) domain.RepositoryFactory) AudioService {
	return &audioService{
		db:   db,
		repo: repo,
	}
}

func (s *audioService) FindAudioByID(ctx context.Context, audioID domain.AudioID) (domain.Audio, error) {
	logger := log.FromContext(ctx)
	logger.Infof("audioID: %d", audioID)
	var audio domain.Audio
	if err := s.db.Transaction(func(tx *gorm.DB) error {
		repo := s.repo(tx)
		audioRepo, err := repo.NewAudioRepository(ctx)
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
