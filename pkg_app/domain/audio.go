package domain

import "github.com/go-playground/validator"

type AudioID uint

type Audio interface {
	GetID() uint
	GetLang() string
	GetText() string
	GetAudioContent() string
}

type audio struct {
	ID           uint   `validate:"required"`
	Lang         string `validate:"required"`
	Text         string `validate:"required"`
	AudioContent string `validate:"required"`
}

func NewAudio(id uint, lang, text, audioContent string) (Audio, error) {
	m := &audio{
		ID:           id,
		Lang:         lang,
		Text:         text,
		AudioContent: audioContent,
	}

	v := validator.New()
	return m, v.Struct(m)
}

func (a *audio) GetID() uint {
	return a.ID
}

func (a *audio) GetLang() string {
	return a.Lang
}

func (a *audio) GetText() string {
	return a.Text
}

func (a *audio) GetAudioContent() string {
	return a.AudioContent
}
