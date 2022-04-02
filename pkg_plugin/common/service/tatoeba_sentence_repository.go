package service

import (
	"context"
	"errors"
	"time"

	appD "github.com/kujilabo/cocotola-api/pkg_app/domain"
	lib "github.com/kujilabo/cocotola-api/pkg_lib/domain"
)

var ErrTatoebaSentenceAlreadyExists = errors.New("tatoebaSentence already exists")

type TatoebaSentenceAddParameter interface {
	GetSentenceNumber() int
	GetLang() appD.Lang3
	GetText() string
	GetAuthor() string
	GetUpdatedAt() time.Time
}

type tatoebaSentenceAddParameter struct {
	SentenceNumber int `validate:"required"`
	Lang           appD.Lang3
	Text           string `validate:"required"`
	Author         string `validate:"required"`
	UpdatedAt      time.Time
}

func NewTatoebaSentenceAddParameter(sentenceNumber int, lang appD.Lang3, text, author string, updatedAt time.Time) (TatoebaSentenceAddParameter, error) {
	m := &tatoebaSentenceAddParameter{
		SentenceNumber: sentenceNumber,
		Lang:           lang,
		Text:           text,
		Author:         author,
		UpdatedAt:      updatedAt,
	}

	return m, lib.Validator.Struct(m)
}

func (p *tatoebaSentenceAddParameter) GetSentenceNumber() int {
	return p.SentenceNumber
}

func (p *tatoebaSentenceAddParameter) GetLang() appD.Lang3 {
	return p.Lang
}

func (p *tatoebaSentenceAddParameter) GetText() string {
	return p.Text
}

func (p *tatoebaSentenceAddParameter) GetAuthor() string {
	return p.Author
}

func (p *tatoebaSentenceAddParameter) GetUpdatedAt() time.Time {
	return p.UpdatedAt
}

type TatoebaSentenceSearchCondition interface {
	GetPageNo() int
	GetPageSize() int
	GetKeyword() string
	IsRandom() bool
}

type tatoebaSentenceSearchCondition struct {
	PageNo   int `validate:"required,gte=1"`
	PageSize int `validate:"required,gte=1,lte=100"`
	Keyword  string
	Random   bool
}

func NewTatoebaSentenceSearchCondition(pageNo, pageSize int, keyword string, random bool) (TatoebaSentenceSearchCondition, error) {
	m := &tatoebaSentenceSearchCondition{
		PageNo:   pageNo,
		PageSize: pageSize,
		Keyword:  keyword,
		Random:   random,
	}

	return m, lib.Validator.Struct(m)
}

func (c *tatoebaSentenceSearchCondition) GetPageNo() int {
	return c.PageNo
}

func (c *tatoebaSentenceSearchCondition) GetPageSize() int {
	return c.PageSize
}

func (c *tatoebaSentenceSearchCondition) GetKeyword() string {
	return c.Keyword
}

func (c *tatoebaSentenceSearchCondition) IsRandom() bool {
	return c.Random
}

type TatoebaSentenceSearchResult struct {
	TotalCount int64
	Results    []TatoebaSentencePair
}

type TatoebaSentenceRepository interface {
	FindTatoebaSentences(ctx context.Context, param TatoebaSentenceSearchCondition) (*TatoebaSentenceSearchResult, error)

	FindTatoebaSentenceBySentenceNumber(ctx context.Context, sentenceNumber int) (TatoebaSentence, error)

	Add(ctx context.Context, param TatoebaSentenceAddParameter) error
}

type TatoebaSentenceRepositoryReadOnly interface {
	FindTatoebaSentences(ctx context.Context, param TatoebaSentenceSearchCondition) (*TatoebaSentenceSearchResult, error)

	FindTatoebaSentenceBySentenceNumber(ctx context.Context, sentenceNumber int) (TatoebaSentence, error)
}
