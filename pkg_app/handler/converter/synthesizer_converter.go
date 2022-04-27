package converter

import (
	"context"

	"github.com/kujilabo/cocotola-api/pkg_app/handler/entity"
	"github.com/kujilabo/cocotola-api/pkg_app/service"
)

func ToAudioResponse(ctx context.Context, audio service.Audio) (*entity.AudioResponse, error) {
	audioModel := audio.GetAudioModel()
	return &entity.AudioResponse{
		ID:      int(audioModel.GetID()),
		Lang2:   audioModel.GetLang2().String(),
		Text:    audioModel.GetText(),
		Content: audioModel.GetContent(),
	}, nil
}
