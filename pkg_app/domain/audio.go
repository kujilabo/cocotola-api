package domain

import libD "github.com/kujilabo/cocotola-api/pkg_lib/domain"

type AudioID uint

type Audio interface {
	GetID() uint
	GetLang() Lang5
	GetText() string
	GetAudioContent() string
}

type audio struct {
	ID           uint   `validate:"required"`
	Lang         Lang5  `validate:"required,len=5"`
	Text         string `validate:"required"`
	AudioContent string `validate:"required"`
}

func NewAudio(id uint, lang Lang5, text, audioContent string) (Audio, error) {
	m := &audio{
		ID:           id,
		Lang:         lang,
		Text:         text,
		AudioContent: audioContent,
	}

	return m, libD.Validator.Struct(m)
}

func (a *audio) GetID() uint {
	return a.ID
}

func (a *audio) GetLang() Lang5 {
	return a.Lang
}

func (a *audio) GetText() string {
	return a.Text
}

func (a *audio) GetAudioContent() string {
	return a.AudioContent
}
