package application

import (
	"context"
	"errors"
	"io"

	"github.com/kujilabo/cocotola-api/pkg_lib/log"
	"github.com/kujilabo/cocotola-api/pkg_plugin/common/domain"
	"gorm.io/gorm"
)

const (
	commitSize = 1000
	logSize    = 100000
)

type TatoebaService interface {
	ImportSentences(ctx context.Context, iterator domain.TatoebaSentenceAddParameterIterator) error
	ImportLinks(ctx context.Context, iterator domain.TatoebaLinkAddParameterIterator) error
}

type tatoebaService struct {
	db *gorm.DB
	rf func(db *gorm.DB) (domain.RepositoryFactory, error)
}

func NewTatoebaService(db *gorm.DB, rf func(db *gorm.DB) (domain.RepositoryFactory, error)) TatoebaService {
	return &tatoebaService{
		db: db,
		rf: rf,
	}
}

func (s *tatoebaService) ImportSentences(ctx context.Context, iterator domain.TatoebaSentenceAddParameterIterator) error {
	logger := log.FromContext(ctx)

	var count = 0
	var loop = true
	for loop {
		if err := s.db.Transaction(func(tx *gorm.DB) error {
			rf, err := s.rf(tx)
			if err != nil {
				return err
			}

			repo, err := rf.NewTatoebaSentenceRepository(ctx)
			if err != nil {
				return err
			}

			i := 0
			for {
				param, err := iterator.Next(ctx)
				if errors.Is(err, io.EOF) {
					loop = false
					break
				}
				if err != nil {
					return err
				}
				if param == nil {
					// logger.Infof("skip count: %d", count)
					continue
				}

				if err := repo.Add(ctx, param); err != nil {
					logger.Warnf("failed to add .commit i: %d, err: %v", i, err)
					continue
				}
				i++
				count++
				if i >= commitSize {
					if count%logSize == 0 {
						logger.Infof("commit i: %d", i)
					}
					break
				}
			}

			return nil
		}); err != nil {
			return err
		}
	}

	return nil
}

func (s *tatoebaService) ImportLinks(ctx context.Context, iterator domain.TatoebaLinkAddParameterIterator) error {
	logger := log.FromContext(ctx)

	var count = 0
	var loop = true
	for loop {
		if err := s.db.Transaction(func(tx *gorm.DB) error {
			rf, err := s.rf(tx)
			if err != nil {
				return err
			}

			repo, err := rf.NewTatoebaLinkRepository(ctx)
			if err != nil {
				return err
			}

			i := 0
			for {
				param, err := iterator.Next(ctx)
				if errors.Is(err, io.EOF) {
					loop = false
					break
				}
				if err != nil {
					return err
				}
				if param == nil {
					logger.Infof("skip count: %d", count)
					continue
				}

				if err := repo.Add(ctx, param); err != nil {
					logger.Warnf("commit i: %d", i)
					continue
				}
				i++
				count++
				if i >= commitSize {
					if count%logSize == 0 {
						logger.Infof("commit i: %d", i)
					}
					break
				}
			}

			return nil
		}); err != nil {
			return err
		}
	}

	return nil
}
