package service

import (
	"github.com/kujilabo/cocotola-api/pkg_app/service"
	libD "github.com/kujilabo/cocotola-api/pkg_lib/domain"
	"github.com/kujilabo/cocotola-api/pkg_plugin/english/domain"
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
