package gateway

import (
	"context"
	"errors"

	"golang.org/x/xerrors"
	"gorm.io/gorm"

	"github.com/kujilabo/cocotola-api/pkg_app/domain"
	"github.com/kujilabo/cocotola-api/pkg_app/service"
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

func (e *audioEntity) toAudioModel() (domain.AudioModel, error) {
	lang5, err := domain.NewLang5(e.Lang)
	if err != nil {
		return nil, err
	}

	return domain.NewAudioModel(e.ID, lang5, e.Text, e.AudioContent)
}

type audioRepository struct {
	db *gorm.DB
}

func NewAudioRepository(db *gorm.DB) service.AudioRepository {
	return &audioRepository{
		db: db,
	}
}

func (r *audioRepository) AddAudio(ctx context.Context, lang domain.Lang5, text, audioContent string) (domain.AudioID, error) {
	entity := audioEntity{
		Lang:         lang.String(),
		Text:         text,
		AudioContent: audioContent,
	}
	if result := r.db.Create(&entity); result.Error != nil {
		return 0, result.Error
	}
	return domain.AudioID(entity.ID), nil
}

func (r *audioRepository) FindAudioByAudioID(ctx context.Context, audioID domain.AudioID) (service.Audio, error) {
	entity := audioEntity{}
	if result := r.db.Where("id = ?", uint(audioID)).First(&entity); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, service.ErrAudioNotFound
		}
		return nil, result.Error
	}
	audioModel, err := entity.toAudioModel()
	if err != nil {
		return nil, err
	}
	audio, err := service.NewAudio(audioModel)
	if err != nil {
		return nil, err
	}
	return audio, nil
}

func (r *audioRepository) FindByLangAndText(ctx context.Context, lang domain.Lang5, text string) (service.Audio, error) {
	entity := audioEntity{}
	if result := r.db.Where("lang = ? and text = ?", lang.String(), text).First(&entity); result.Error != nil {
		return nil, result.Error
	}
	audio, err := entity.toAudioModel()
	if err != nil {
		return nil, err
	}
	audioService, err := service.NewAudio(audio)
	if err != nil {
		return nil, err
	}
	return audioService, nil
}

func (r *audioRepository) FindAudioIDByText(ctx context.Context, lang domain.Lang5, text string) (domain.AudioID, error) {
	entity := audioEntity{}
	if result := r.db.Where("lang = ? and text = ?", lang.String(), text).First(&entity); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return 0, service.ErrAudioNotFound
		}
		return 0, result.Error
	}
	model, err := entity.toAudioModel()
	if err != nil {
		return 0, xerrors.Errorf("failed to toAudio. entity: %v, err: %w", entity, err)
	}
	return domain.AudioID(model.GetID()), nil
}
