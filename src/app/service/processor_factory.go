//go:generate mockery --output mock --name ProcessorFactory
package service

import liberrors "github.com/kujilabo/cocotola-api/src/lib/errors"

type ProcessorFactory interface {
	NewProblemAddProcessor(processorType string) (ProblemAddProcessor, error)

	NewProblemUpdateProcessor(processorType string) (ProblemUpdateProcessor, error)

	NewProblemRemoveProcessor(processorType string) (ProblemRemoveProcessor, error)

	NewProblemImportProcessor(processorType string) (ProblemImportProcessor, error)

	NewProblemQuotaProcessor(processorType string) (ProblemQuotaProcessor, error)
}

type processorFactrory struct {
	addProcessors    map[string]ProblemAddProcessor
	updateProcessors map[string]ProblemUpdateProcessor
	removeProcessors map[string]ProblemRemoveProcessor
	importProcessors map[string]ProblemImportProcessor
	quotaProcessors  map[string]ProblemQuotaProcessor
}

func NewProcessorFactory(addProcessors map[string]ProblemAddProcessor, updateProcessors map[string]ProblemUpdateProcessor, removeProcessors map[string]ProblemRemoveProcessor, importProcessors map[string]ProblemImportProcessor, quotaProcessors map[string]ProblemQuotaProcessor) ProcessorFactory {
	return &processorFactrory{
		addProcessors:    addProcessors,
		updateProcessors: updateProcessors,
		removeProcessors: removeProcessors,
		importProcessors: importProcessors,
		quotaProcessors:  quotaProcessors,
	}
}

func (f *processorFactrory) NewProblemAddProcessor(processorType string) (ProblemAddProcessor, error) {
	processor, ok := f.addProcessors[processorType]
	if !ok {
		return nil, liberrors.Errorf("newProblemProcessor not found. processorType: %s", processorType)
	}
	return processor, nil
}

func (f *processorFactrory) NewProblemUpdateProcessor(processorType string) (ProblemUpdateProcessor, error) {
	processor, ok := f.updateProcessors[processorType]
	if !ok {
		return nil, liberrors.Errorf("newProblemProcessor not found. processorType: %s", processorType)
	}
	return processor, nil
}

func (f *processorFactrory) NewProblemRemoveProcessor(processorType string) (ProblemRemoveProcessor, error) {
	processor, ok := f.removeProcessors[processorType]
	if !ok {
		return nil, liberrors.Errorf("newProblemRemoveProcessor not found. processorType: %s", processorType)
	}
	return processor, nil
}

func (f *processorFactrory) NewProblemImportProcessor(processorType string) (ProblemImportProcessor, error) {
	processor, ok := f.importProcessors[processorType]
	if !ok {
		return nil, liberrors.Errorf("NewProblemImportProcessor not found. processorType: %s", processorType)
	}
	return processor, nil
}

func (f *processorFactrory) NewProblemQuotaProcessor(processorType string) (ProblemQuotaProcessor, error) {
	processor, ok := f.quotaProcessors[processorType]
	if !ok {
		return nil, liberrors.Errorf("NewProblemQuotaProcessor not found. processorType: %s", processorType)
	}
	return processor, nil
}
