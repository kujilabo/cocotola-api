//go:generate mockery --output mock --name SynthesizerClient
package service

import (
	"context"

	"github.com/kujilabo/cocotola-api/pkg_app/domain"
)

type SynthesizerClient interface {
	Synthesize(ctx context.Context, lang domain.Lang5, text string) (string, error)
}
