package domain

import "github.com/kujilabo/cocotola-api/pkg_app/domain"

type Synthesizer interface {
	Synthesize(lang domain.Lang5, text string) (string, error)
}
