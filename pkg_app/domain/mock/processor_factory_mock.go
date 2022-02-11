package domain_mock

import (
	"github.com/kujilabo/cocotola-api/pkg_app/domain"
	"github.com/stretchr/testify/mock"
)

type ProcessorFactoryMock struct {
	mock.Mock
}

func (m *ProcessorFactoryMock) NewProblemAddProcessor(processorType string) (domain.ProblemAddProcessor, error) {
	args := m.Called(processorType)
	return args.Get(0).(domain.ProblemAddProcessor), args.Error(1)
}
func (m *ProcessorFactoryMock) NewProblemUpdateProcessor(processorType string) (domain.ProblemUpdateProcessor, error) {
	args := m.Called(processorType)
	return args.Get(0).(domain.ProblemUpdateProcessor), args.Error(1)
}
func (m *ProcessorFactoryMock) NewProblemRemoveProcessor(processorType string) (domain.ProblemRemoveProcessor, error) {
	args := m.Called(processorType)
	return args.Get(0).(domain.ProblemRemoveProcessor), args.Error(1)
}
func (m *ProcessorFactoryMock) NewProblemImportProcessor(processorType string) (domain.ProblemImportProcessor, error) {
	args := m.Called(processorType)
	return args.Get(0).(domain.ProblemImportProcessor), args.Error(1)
}
func (m *ProcessorFactoryMock) NewProblemQuotaProcessor(processorType string) (domain.ProblemQuotaProcessor, error) {
	args := m.Called(processorType)
	return args.Get(0).(domain.ProblemQuotaProcessor), args.Error(1)
}
