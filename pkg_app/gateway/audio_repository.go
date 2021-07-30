package gateway

import (
	"context"
	"errors"

	"golang.org/x/xerrors"
	"gorm.io/gorm"

	"github.com/kujilabo/cocotola-api/pkg_app/domain"
)

type audioEntity struct {
	ID           uint   `validate:"required"`
	Lang         string `validate:"required"`
	Text         string `validate:"required"`
	AudioContent string `validate:"required"`
}

func (e *audioEntity) TableName() string {
	return "audio"
}

func (e *audioEntity) toAudio() (domain.Audio, error) {
	return domain.NewAudio(e.ID, e.Lang, e.Text, e.AudioContent)
}

type audioRepository struct {
	db *gorm.DB
}

func NewAudioRepository(db *gorm.DB) domain.AudioRepository {
	return &audioRepository{
		db: db,
	}
}

func (r *audioRepository) AddAudio(ctx context.Context, lang domain.Lang5, text, audioContent string) (domain.AudioID, error) {
	entity := audioEntity{
		Lang:         string(lang),
		Text:         text,
		AudioContent: audioContent,
	}
	if result := r.db.Create(&entity); result.Error != nil {
		return 0, result.Error
	}
	return domain.AudioID(entity.ID), nil
}

func (r *audioRepository) FindAudioByAudioID(ctx context.Context, audioID domain.AudioID) (domain.Audio, error) {
	entity := audioEntity{}
	if result := r.db.Where("id = ?", uint(audioID)).First(&entity); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, domain.ErrAudioNotFound
		}
		return nil, result.Error
	}
	return entity.toAudio()
}

func (r *audioRepository) FindByLangAndText(ctx context.Context, lang domain.Lang5, text string) (domain.Audio, error) {
	entity := audioEntity{}
	if result := r.db.Where("lang = ? and text = ?", string(lang), text).First(&entity); result.Error != nil {
		return nil, result.Error
	}
	return entity.toAudio()
}

func (r *audioRepository) FindAudioID(ctx context.Context, lang domain.Lang5, text string) (domain.AudioID, error) {
	entity := audioEntity{}
	if result := r.db.Where("lang = ? and text = ?", string(lang), text).First(&entity); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return 0, domain.ErrAudioNotFound
		}
		return 0, result.Error
	}
	model, err := entity.toAudio()
	if err != nil {
		return 0, xerrors.Errorf("faield to toAudio. entity: %v, err: %w", entity, err)
	}
	return domain.AudioID(model.GetID()), nil
}
