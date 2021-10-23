package domain

const StudyMaxLevel = 10
const StudyMinLevel = 0

type StudyResultParameter struct {
	Result bool
}

type ProblemWithLevel struct {
	ProblemID ProblemID
	Level     int
}
