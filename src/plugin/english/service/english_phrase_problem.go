package service

import (
	"github.com/kujilabo/cocotola-api/src/app/service"
	libD "github.com/kujilabo/cocotola-api/src/lib/domain"
	"github.com/kujilabo/cocotola-api/src/plugin/english/domain"
)

type EnglishPhraseProblem interface {
	domain.EnglishPhraseProblemModel
	service.ProblemFeature
}

type englishPhraseProblem struct {
	domain.EnglishPhraseProblemModel
	service.ProblemFeature
}

func NewEnglishPhraseProblem(problemModel domain.EnglishPhraseProblemModel, problem service.ProblemFeature) (EnglishPhraseProblem, error) {
	m := &englishPhraseProblem{
		EnglishPhraseProblemModel: problemModel,
		ProblemFeature:            problem,
	}

	return m, libD.Validator.Struct(m)
}
