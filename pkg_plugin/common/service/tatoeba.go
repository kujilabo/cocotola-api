package service

import (
	"errors"
	"time"

	appD "github.com/kujilabo/cocotola-api/pkg_app/domain"
	libD "github.com/kujilabo/cocotola-api/pkg_lib/domain"
)

var ErrTatoebaSentenceAlreadyExists = errors.New("tatoebaSentence already exists")

type TatoebaSentence interface {
	GetSentenceNumber() int
	GetLang2() appD.Lang2
	GetText() string
	GetAuthor() string
	GetUpdatedAt() time.Time
}

type tatoebaSentence struct {
	SentenceNumber int
	Lang2          appD.Lang2
	Text           string
	Author         string
	UpdatedAt      time.Time
}

func NewTatoebaSentence(sentenceNumber int, lang2 appD.Lang2, text, author string, updatedAt time.Time) (TatoebaSentence, error) {
	m := &tatoebaSentence{
		SentenceNumber: sentenceNumber,
		Lang2:          lang2,
		Text:           text,
		Author:         author,
		UpdatedAt:      updatedAt,
	}

	return m, libD.Validator.Struct(m)
}

func (m *tatoebaSentence) GetSentenceNumber() int {
	return m.SentenceNumber
}

func (m *tatoebaSentence) GetLang2() appD.Lang2 {
	return m.Lang2
}

func (m *tatoebaSentence) GetText() string {
	return m.Text
}

func (m *tatoebaSentence) GetAuthor() string {
	return m.Author
}

func (m *tatoebaSentence) GetUpdatedAt() time.Time {
	return m.UpdatedAt
}

type TatoebaSentencePair interface {
	GetSrc() TatoebaSentence
	GetDst() TatoebaSentence
}

type tatoebaSentencePair struct {
	Src TatoebaSentence
	Dst TatoebaSentence
}

func NewTatoebaSentencePair(src, dst TatoebaSentence) (TatoebaSentencePair, error) {
	m := &tatoebaSentencePair{
		Src: src,
		Dst: dst,
	}

	return m, libD.Validator.Struct(m)
}

func (m *tatoebaSentencePair) GetSrc() TatoebaSentence {
	return m.Src
}

func (m *tatoebaSentencePair) GetDst() TatoebaSentence {
	return m.Dst
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

	return m, libD.Validator.Struct(m)
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

type TatoebaSentencePairSearchResult struct {
	TotalCount int64
	Results    []TatoebaSentencePair
}
