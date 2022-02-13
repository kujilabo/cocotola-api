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
	FindSentences(ctx context.Context, param domain.TatoebaSentenceSearchCondition) (*domain.TatoebaSentenceSearchResult, error)

	ImportSentences(ctx context.Context, iterator domain.TatoebaSentenceAddParameterIterator) error

	ImportLinks(ctx context.Context, iterator domain.TatoebaLinkAddParameterIterator) error
}

type tatoebaService struct {
	db *gorm.DB
	rf domain.RepositoryFactoryFunc
}

func NewTatoebaService(db *gorm.DB, rf domain.RepositoryFactoryFunc) TatoebaService {
	return &tatoebaService{
		db: db,
		rf: rf,
	}
}

func (s *tatoebaService) FindSentences(ctx context.Context, param domain.TatoebaSentenceSearchCondition) (*domain.TatoebaSentenceSearchResult, error) {
	var result *domain.TatoebaSentenceSearchResult
	if err := s.db.Transaction(func(tx *gorm.DB) error {
		rf, err := s.rf(ctx, tx)
		if err != nil {
			return err
		}

		repo, err := rf.NewTatoebaSentenceRepository(ctx)
		if err != nil {
			return err
		}

		tmpResult, err := repo.FindTatoebaSentences(ctx, param)
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

func (s *tatoebaService) ImportSentences(ctx context.Context, iterator domain.TatoebaSentenceAddParameterIterator) error {
	logger := log.FromContext(ctx)

	var count = 0
	var loop = true
	for loop {
		if err := s.db.Transaction(func(tx *gorm.DB) error {
			rf, err := s.rf(ctx, tx)
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
					logger.Warnf("failed to Add. count: %d, err: %v", count, err)
					continue
				}
				i++
				count++
				if i >= commitSize {
					if count%logSize == 0 {
						logger.Infof("commit count: %d", count)
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
			rf, err := s.rf(ctx, tx)
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
					logger.Infof("skip to Add Link. count: %d", count)
					continue
				}

				if err := repo.Add(ctx, param); err != nil {
					logger.Warnf("failed to Add Link. count: %d", count)
					continue
				}
				i++
				count++
				if i >= commitSize {
					if count%logSize == 0 {
						logger.Infof("commit count: %d", count)
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
