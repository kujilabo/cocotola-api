package service

import (
	"context"

	"github.com/kujilabo/cocotola-api/pkg_app/domain"
	lib "github.com/kujilabo/cocotola-api/pkg_lib/domain"
)

type ProblemFeature interface {
	FindAudioByID(ctx context.Context, audioID domain.AudioID) (Audio, error)
}

type Problem interface {
	domain.ProblemModel
	// GetProblemModel() domain.ProblemModel
	// GetProblemType() string
	// GetProperties(ctx context.Context) map[string]interface{}
	ProblemFeature
}

type problem struct {
	domain.ProblemModel
	rf           AudioRepositoryFactory
	problemModel domain.ProblemModel
}

func NewProblem(rf AudioRepositoryFactory, problemModel domain.ProblemModel) (Problem, error) {
	s := &problem{
		ProblemModel: problemModel,
		rf:           rf,
		problemModel: problemModel,
	}

	return s, lib.Validator.Struct(s)
}

// func (s *problem) GetProblemModel() domain.ProblemModel {
// 	return s.problemModel
// }

func (s *problem) FindAudioByID(ctx context.Context, audioID domain.AudioID) (Audio, error) {
	audioRepo := s.rf.NewAudioRepository(ctx)
	return audioRepo.FindAudioByAudioID(ctx, audioID)
}

type ProblemWithResults interface {
	domain.ProblemModel
	GetResults() []bool
	GetLevel() int
}

type problemWithResults struct {
	domain.ProblemModel
	results []bool
	level   int
}

func NewProblemWithResults(problem domain.ProblemModel, results []bool, level int) ProblemWithResults {
	return &problemWithResults{
		ProblemModel: problem,
		results:      results,
		level:        level,
	}
}

func (m *problemWithResults) GetResults() []bool {
	return m.results
}

func (m *problemWithResults) GetLevel() int {
	return m.level
}
