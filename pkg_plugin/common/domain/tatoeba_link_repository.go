package domain

import (
	"context"
	"errors"

	"github.com/go-playground/validator"
)

var ErrTatoebaLinkAlreadyExists = errors.New("tatoebaLink already exists")

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

	v := validator.New()
	return m, v.Struct(m)
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
