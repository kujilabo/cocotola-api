//go:generate mockery --output mock --name SynthesizerClient
package service

import (
	"context"
	"errors"

	"github.com/kujilabo/cocotola-api/pkg_app/domain"
)

var ErrAudioNotFound = errors.New("audio not found")

type SynthesizerClient interface {
	Synthesize(ctx context.Context, lang domain.Lang2, text string) (Audio, error)

	FindAudioByAudioID(ctx context.Context, audioID domain.AudioID) (Audio, error)
}
