package service

import (
	"github.com/kujilabo/cocotola-api/src/app/service"
	libD "github.com/kujilabo/cocotola-api/src/lib/domain"
	"github.com/kujilabo/cocotola-api/src/plugin/english/domain"
)

const EnglishWordProblemType = "english_word"

type EnglishWordProblem interface {
	domain.EnglishWordProblemModel
	service.ProblemFeature
}

type englishWordProblem struct {
	domain.EnglishWordProblemModel
	service.ProblemFeature
}

func NewEnglishWordProblem(problemModel domain.EnglishWordProblemModel, problem service.ProblemFeature) (EnglishWordProblem, error) {
	m := &englishWordProblem{
		EnglishWordProblemModel: problemModel,
		ProblemFeature:          problem,
	}

	return m, libD.Validator.Struct(m)
}
