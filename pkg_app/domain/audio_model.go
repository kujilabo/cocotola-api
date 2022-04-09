package domain

import libD "github.com/kujilabo/cocotola-api/pkg_lib/domain"

type AudioID uint

type AudioModel interface {
	GetID() uint
	GetLang() Lang5
	GetText() string
	GetAudioContent() string
}

type audioModel struct {
	ID           uint   `validate:"required"`
	Lang         Lang5  `validate:"required,len=5"`
	Text         string `validate:"required"`
	AudioContent string `validate:"required"`
}

func NewAudioModel(id uint, lang Lang5, text, audioContent string) (AudioModel, error) {
	m := &audioModel{
		ID:           id,
		Lang:         lang,
		Text:         text,
		AudioContent: audioContent,
	}

	return m, libD.Validator.Struct(m)
}

func (a *audioModel) GetID() uint {
	return a.ID
}

func (a *audioModel) GetLang() Lang5 {
	return a.Lang
}

func (a *audioModel) GetText() string {
	return a.Text
}

func (a *audioModel) GetAudioContent() string {
	return a.AudioContent
}
