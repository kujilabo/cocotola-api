package service

import "context"

type TatoebaLinkAddParameterIterator interface {
	Next(ctx context.Context) (TatoebaLinkAddParameter, error)
}

type TatoebaSentenceAddParameterIterator interface {
	Next(ctx context.Context) (TatoebaSentenceAddParameter, error)
}
