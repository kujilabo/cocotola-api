package service

type ProblemAddParameterIterator interface {
	Next() (ProblemAddParameter, error)
}
