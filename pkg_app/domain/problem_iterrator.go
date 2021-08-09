package domain

type ProblemAddParameterIterator interface {
	Next() (ProblemAddParameter, error)
}
