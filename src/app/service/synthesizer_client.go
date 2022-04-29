//go:generate mockery --output mock --name SynthesizerClient
package service

import (
	"context"
	"errors"

	"github.com/kujilabo/cocotola-api/src/app/domain"
)

var ErrAudioNotFound = errors.New("audio not found")

type SynthesizerClient interface {
	Synthesize(ctx context.Context, lang2 domain.Lang2, text string) (Audio, error)

	FindAudioByAudioID(ctx context.Context, audioID domain.AudioID) (Audio, error)
}
