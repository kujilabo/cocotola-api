package application

import (
	"context"
	"fmt"

	"gorm.io/gorm"

	"github.com/kujilabo/cocotola-api/pkg_app/domain"
	"github.com/kujilabo/cocotola-api/pkg_lib/log"
)

type AudioService interface {
	FindAudioByID(ctx context.Context, audioID domain.AudioID) (domain.Audio, error)
}

type audioService struct {
	db   *gorm.DB
	repo func(db *gorm.DB) (domain.RepositoryFactory, error)
}

func NewAudioService(db *gorm.DB, repo func(db *gorm.DB) (domain.RepositoryFactory, error)) AudioService {
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
		repo, err := s.repo(tx)
		if err != nil {
			return fmt.Errorf("failed to repo. err: %w", err)
		}
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
