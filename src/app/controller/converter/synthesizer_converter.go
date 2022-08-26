package converter

import (
	"context"

	"github.com/kujilabo/cocotola-api/src/app/controller/entity"
	"github.com/kujilabo/cocotola-api/src/app/service"
	libD "github.com/kujilabo/cocotola-api/src/lib/domain"
)

func ToAudioResponse(ctx context.Context, audio service.Audio) (*entity.AudioResponse, error) {
	audioModel := audio.GetAudioModel()
	e := &entity.AudioResponse{
		ID:      int(audioModel.GetID()),
		Lang2:   audioModel.GetLang2().String(),
		Text:    audioModel.GetText(),
		Content: audioModel.GetContent(),
	}

	return e, libD.Validator.Struct(e)
}
