package gateway

import (
	"context"
	"fmt"
	"time"

	"gorm.io/gorm"

	libG "github.com/kujilabo/cocotola-api/pkg_lib/gateway"
	"github.com/kujilabo/cocotola-api/pkg_plugin/common/domain"
)

type tatoebaSentenceRepository struct {
	db *gorm.DB
}

type tatoebaSentenceEntity struct {
	SentenceNumber int
	Lang           string
	Text           string
	Author         string
	UpdatedAt      time.Time
}

func (e *tatoebaSentenceEntity) TableName() string {
	return "tatoeba_sentence"
}

func NewTatoebaSentenceRepository(db *gorm.DB) domain.TatoebaSentenceRepository {
	return &tatoebaSentenceRepository{
		db: db,
	}
}

func (r *tatoebaSentenceRepository) Add(ctx context.Context, param domain.TatoebaSentenceAddParameter) error {
	entity := tatoebaSentenceEntity{
		SentenceNumber: param.GetSentenceNumber(),
		Lang:           param.GetLang().String(),
		Text:           param.GetText(),
		Author:         param.GetAuthor(),
		UpdatedAt:      param.GetUpdatedAt(),
	}

	if result := r.db.Create(&entity); result.Error != nil {
		err := libG.ConvertDuplicatedError(result.Error, domain.ErrTatoebaSentenceAlreadyExists)
		return fmt.Errorf("failed to Add tatoebaSentence. err: %w", err)
	}

	return nil
}
