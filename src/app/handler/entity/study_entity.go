package entity

import "time"

type StudyResultParameter struct {
	Result    bool `json:"result"`
	Memorized bool `json:"memorized"`
}

type StudyRecord struct {
	ProblemID      uint       `json:"problemId"`
	Level          int        `json:"level"`
	ResultPrev1    bool       `json:"resultPrev1"`
	Memorized      bool       `json:"memorized"`
	LastAnsweredAt *time.Time `json:"lastAnsweredAt"`
}

type StudyRecords struct {
	Records []*StudyRecord `json:"records" validate:"dive"`
}

type IntValue struct {
	Value int `json:"value"`
}
