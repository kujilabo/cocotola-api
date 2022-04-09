package service_mock

import (
	"github.com/stretchr/testify/mock"

	"github.com/kujilabo/cocotola-api/pkg_app/service"
)

type ProcessorFactoryMock struct {
	mock.Mock
}

func (m *ProcessorFactoryMock) NewProblemAddProcessor(processorType string) (service.ProblemAddProcessor, error) {
	args := m.Called(processorType)
	return args.Get(0).(service.ProblemAddProcessor), args.Error(1)
}
func (m *ProcessorFactoryMock) NewProblemUpdateProcessor(processorType string) (service.ProblemUpdateProcessor, error) {
	args := m.Called(processorType)
	return args.Get(0).(service.ProblemUpdateProcessor), args.Error(1)
}
func (m *ProcessorFactoryMock) NewProblemRemoveProcessor(processorType string) (service.ProblemRemoveProcessor, error) {
	args := m.Called(processorType)
	return args.Get(0).(service.ProblemRemoveProcessor), args.Error(1)
}
func (m *ProcessorFactoryMock) NewProblemImportProcessor(processorType string) (service.ProblemImportProcessor, error) {
	args := m.Called(processorType)
	return args.Get(0).(service.ProblemImportProcessor), args.Error(1)
}
func (m *ProcessorFactoryMock) NewProblemQuotaProcessor(processorType string) (service.ProblemQuotaProcessor, error) {
	args := m.Called(processorType)
	return args.Get(0).(service.ProblemQuotaProcessor), args.Error(1)
}
