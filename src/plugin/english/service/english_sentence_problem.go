package service

import (
	"github.com/kujilabo/cocotola-api/src/app/service"
	libD "github.com/kujilabo/cocotola-api/src/lib/domain"
	"github.com/kujilabo/cocotola-api/src/plugin/english/domain"
)

const EnglishSentenceProblemType = "english_sentence"

type EnglishSentenceProblem interface {
	domain.EnglishSentenceProblemModel
	service.ProblemFeature
}

type englishSentenceProblem struct {
	domain.EnglishSentenceProblemModel
	service.ProblemFeature
}

func NewEnglishSentenceProblem(problemModel domain.EnglishSentenceProblemModel, problem service.ProblemFeature) (EnglishSentenceProblem, error) {
	m := &englishSentenceProblem{
		EnglishSentenceProblemModel: problemModel,
		ProblemFeature:              problem,
	}

	return m, libD.Validator.Struct(m)
}
