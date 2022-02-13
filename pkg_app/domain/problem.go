package domain

import (
	"context"

	lib "github.com/kujilabo/cocotola-api/pkg_lib/domain"
	user "github.com/kujilabo/cocotola-api/pkg_user/domain"
)

type ProblemID uint

type Problem interface {
	user.Model
	GetNumber() int
	GetProblemType() string
	GetProperties(ctx context.Context) map[string]interface{}

	FindAudioByID(ctx context.Context, audioID AudioID) (Audio, error)
}

type problem struct {
	rf AudioRepositoryFactory
	user.Model
	Number      int                    `validate:"required"`
	ProblemType string                 `validate:"required"`
	Properties  map[string]interface{} `validate:"required"`
}

func NewProblem(rf AudioRepositoryFactory, model user.Model, number int, problemType string, properties map[string]interface{}) (Problem, error) {
	m := &problem{
		rf:          rf,
		Model:       model,
		Number:      number,
		ProblemType: problemType,
		Properties:  properties,
	}

	return m, lib.Validator.Struct(m)
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

func (m *problem) FindAudioByID(ctx context.Context, audioID AudioID) (Audio, error) {
	audioRepo, err := m.rf.NewAudioRepository(ctx)
	if err != nil {
		return nil, err
	}
	return audioRepo.FindAudioByAudioID(ctx, audioID)
}

type ProblemWithResults interface {
	Problem
	GetResults() []bool
	GetLevel() int
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

func (m *problemWithResults) GetResults() []bool {
	return m.results
}

func (m *problemWithResults) GetLevel() int {
	return m.level
}
