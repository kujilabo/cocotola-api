package domain

import "fmt"

type ProcessorFactory interface {
	NewProblemAddProcessor(processorType string) (ProblemAddProcessor, error)

	NewProblemRemoveProcessor(processorType string) (ProblemRemoveProcessor, error)
}

type processorFactrory struct {
	processors       map[string]ProblemAddProcessor
	removeProcessors map[string]ProblemRemoveProcessor
}

func NewProcessorFactory(processors map[string]ProblemAddProcessor, removeProcessors map[string]ProblemRemoveProcessor) ProcessorFactory {
	return &processorFactrory{
		processors:       processors,
		removeProcessors: removeProcessors,
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
