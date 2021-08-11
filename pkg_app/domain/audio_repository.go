package domain

import (
	"context"
	"errors"
)

var ErrAudioNotFound = errors.New("Audio not found")

type AudioRepository interface {
	AddAudio(ctx context.Context, lang Lang5, text, audioContent string) (AudioID, error)

	FindAudioByAudioID(ctx context.Context, audioID AudioID) (Audio, error)

	FindByLangAndText(ctx context.Context, lang Lang5, text string) (Audio, error)

	FindAudioIDByText(ctx context.Context, lang Lang5, text string) (AudioID, error)
}
