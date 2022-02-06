package domain

import (
	"context"
	"errors"

	libD "github.com/kujilabo/cocotola-api/pkg_lib/domain"
)

var ErrTatoebaLinkAlreadyExists = errors.New("tatoebaLink already exists")
var ErrTatoebaLinkSourceNotFound = errors.New("tatoebaLink source not found")

type TatoebaLinkAddParameter interface {
	GetFrom() int
	GetTo() int
}

type tatoebaLinkAddParameter struct {
	From int `validate:"required"`
	To   int `validate:"required"`
}

func NewTatoebaLinkAddParameter(from, to int) (TatoebaLinkAddParameter, error) {
	m := &tatoebaLinkAddParameter{
		From: from,
		To:   to,
	}
	return m, libD.Validator.Struct(m)
}

func (p *tatoebaLinkAddParameter) GetFrom() int {
	return p.From
}

func (p *tatoebaLinkAddParameter) GetTo() int {
	return p.To
}

type TatoebaLinkRepository interface {
	Add(ctx context.Context, param TatoebaLinkAddParameter) error
}
