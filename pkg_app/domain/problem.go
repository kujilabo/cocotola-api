package domain

import (
	"context"

	"github.com/go-playground/validator"

	user "github.com/kujilabo/cocotola-api/pkg_user/domain"
)

type ProblemID uint

type Problem interface {
	user.Model
	GetNumber() int
	GetProblemType() string
	GetProperties(ctx context.Context) map[string]interface{}
}

type problem struct {
	user.Model
	Number      int                    `validate:"required"`
	ProblemType string                 `validate:"required"`
	Properties  map[string]interface{} `validate:"required"`
}

func NewProblem(model user.Model, number int, problemType string, properties map[string]interface{}) (Problem, error) {
	m := &problem{
		Model:       model,
		Number:      number,
		ProblemType: problemType,
		Properties:  properties,
	}

	v := validator.New()
	return m, v.Struct(m)
}

func (m *problem) GetNumber() int {
	return m.Number
}

func (m *problem) GetProblemType() string {
	return m.ProblemType
}

func (m *problem) GetProperties(ctx context.Context) map[string]interface{} {
	return m.Properties
}

type ProblemWithResults interface {
	Problem
	Results() []bool
	Level() int
}

type problemWithResults struct {
	Problem
	results []bool
	level   int
}

func NewProblemWithResults(problem Problem, results []bool, level int) ProblemWithResults {
	return &problemWithResults{
		Problem: problem,
		results: results,
		level:   level,
	}
}

func (m *problemWithResults) Results() []bool {
	return m.results
}

func (m *problemWithResults) Level() int {
	return m.level
}
