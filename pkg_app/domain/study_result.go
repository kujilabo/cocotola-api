package domain

import "time"

const StudyMaxLevel = 10
const StudyMinLevel = 0

type StudyResultParameter struct {
	Result bool
}

type ProblemWithLevel struct {
	ProblemID      ProblemID
	Level          int
	ResultPrev1    bool
	Memorized      bool
	LastAnsweredAt *time.Time
}

type StudyStatus struct {
	Level          int
	ResultPrev1    bool
	Memorized      bool
	LastAnsweredAt *time.Time
}
