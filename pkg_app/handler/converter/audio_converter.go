package converter

import (
	"context"

	"github.com/kujilabo/cocotola-api/pkg_app/domain"
	"github.com/kujilabo/cocotola-api/pkg_app/handler/entity"
)

func ToAudioResponse(ctx context.Context, audio domain.Audio) (*entity.Audio, error) {
	return &entity.Audio{
		ID:           audio.GetID(),
		Lang:         audio.GetLang().String(),
		Text:         audio.GetText(),
		AudioContent: audio.GetAudioContent(),
	}, nil
}
