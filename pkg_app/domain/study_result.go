package domain

import "time"

const StudyMaxLevel = 10
const StudyMinLevel = 0

type StudyResultParameter struct {
	Result bool
}

type StudyRecordWithProblemID struct {
	ProblemID   ProblemID
	StudyRecord StudyRecord
}

type StudyRecord struct {
	Level          int
	ResultPrev1    bool
	Memorized      bool
	LastAnsweredAt *time.Time
}
