package service

import (
	"context"

	"github.com/kujilabo/cocotola-api/pkg_app/domain"
)

type Synthesizer interface {
	Synthesize(ctx context.Context, lang domain.Lang5, text string) (string, error)
}
