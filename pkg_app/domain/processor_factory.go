package domain

import "fmt"

type ProcessorFactory interface {
	NewProblemAddProcessor(processorType string) (ProblemAddProcessor, error)

	NewProblemRemoveProcessor(processorType string) (ProblemRemoveProcessor, error)

	NewProblemImportProcessor(processorType string) (ProblemImportProcessor, error)

	NewProblemQuotaProcessor(processorType string) (ProblemQuotaProcessor, error)
}

type processorFactrory struct {
	processors       map[string]ProblemAddProcessor
	removeProcessors map[string]ProblemRemoveProcessor
	importProcessors map[string]ProblemImportProcessor
	quotaProcessors  map[string]ProblemQuotaProcessor
}

func NewProcessorFactory(processors map[string]ProblemAddProcessor, removeProcessors map[string]ProblemRemoveProcessor, importProcessors map[string]ProblemImportProcessor, quotaProcessors map[string]ProblemQuotaProcessor) ProcessorFactory {
	return &processorFactrory{
		processors:       processors,
		removeProcessors: removeProcessors,
		importProcessors: importProcessors,
		quotaProcessors:  quotaProcessors,
	}
}

func (f *processorFactrory) NewProblemAddProcessor(processorType string) (ProblemAddProcessor, error) {
	processor, ok := f.processors[processorType]
	if !ok {
		return nil, fmt.Errorf("newProblemProcessor not found. processorType: %s", processorType)
	}
	return processor, nil
}

func (f *processorFactrory) NewProblemRemoveProcessor(processorType string) (ProblemRemoveProcessor, error) {
	processor, ok := f.removeProcessors[processorType]
	if !ok {
		return nil, fmt.Errorf("newProblemRemoveProcessor not found. processorType: %s", processorType)
	}
	return processor, nil
}

func (f *processorFactrory) NewProblemImportProcessor(processorType string) (ProblemImportProcessor, error) {
	processor, ok := f.importProcessors[processorType]
	if !ok {
		return nil, fmt.Errorf("NewProblemImportProcessor not found. processorType: %s", processorType)
	}
	return processor, nil
}

func (f *processorFactrory) NewProblemQuotaProcessor(processorType string) (ProblemQuotaProcessor, error) {
	processor, ok := f.quotaProcessors[processorType]
	if !ok {
		return nil, fmt.Errorf("NewProblemQuotaProcessor not found. processorType: %s", processorType)
	}
	return processor, nil
}
