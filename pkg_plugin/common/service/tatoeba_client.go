//go:generate mockery --output mock --name TatoebaClient
package service

import (
	"context"
	"io"
)

type TatoebaClient interface {
	FindSentencePairs(ctx context.Context, param TatoebaSentenceSearchCondition) (*TatoebaSentencePairSearchResult, error)

	FindSentenceBySentenceNumber(ctx context.Context, sentenceNumber int) (TatoebaSentence, error)

	ImportSentences(ctx context.Context, reader io.Reader) error

	ImportLinks(ctx context.Context, reader io.Reader) error
}
