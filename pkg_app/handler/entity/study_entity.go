package entity

type StudyResultParameter struct {
	Result bool `json:"result"`
}

type ProblemWithLevel struct {
	ProblemID uint `json:"problemId"`
	Level     int  `json:"level"`
}

type ProblemWithLevelList struct {
	Records []ProblemWithLevel `json:"records"`
}
