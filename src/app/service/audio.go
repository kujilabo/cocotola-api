package service

import (
	"github.com/kujilabo/cocotola-api/src/app/domain"
	libD "github.com/kujilabo/cocotola-api/src/lib/domain"
)

type Audio interface {
	GetAudioModel() domain.AudioModel
}

type audio struct {
	AudioModel domain.AudioModel
}

func NewAudio(audioModel domain.AudioModel) (Audio, error) {
	m := &audio{
		AudioModel: audioModel,
	}

	return m, libD.Validator.Struct(m)
}

func (s *audio) GetAudioModel() domain.AudioModel {
	return s.AudioModel
}
