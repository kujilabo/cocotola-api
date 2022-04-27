//go:generate mockery --output mock --name AudioModel
package domain

import libD "github.com/kujilabo/cocotola-api/pkg_lib/domain"

type AudioID uint

type AudioModel interface {
	GetID() uint
	GetLang() Lang2
	GetText() string
	GetContent() string
}

type audioModel struct {
	ID      uint `validate:"required"`
	Lang    Lang2
	Text    string `validate:"required"`
	Content string `validate:"required"`
}

func NewAudioModel(id uint, lang Lang2, text, content string) (AudioModel, error) {
	m := &audioModel{
		ID:      id,
		Lang:    lang,
		Text:    text,
		Content: content,
	}

	return m, libD.Validator.Struct(m)
}

func (a *audioModel) GetID() uint {
	return a.ID
}

func (a *audioModel) GetLang() Lang2 {
	return a.Lang
}

func (a *audioModel) GetText() string {
	return a.Text
}

func (a *audioModel) GetContent() string {
	return a.Content
}
