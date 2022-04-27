package service

import (
	"github.com/kujilabo/cocotola-api/pkg_app/service"
	libD "github.com/kujilabo/cocotola-api/pkg_lib/domain"
	"github.com/kujilabo/cocotola-api/pkg_plugin/english/domain"
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
