package entity

type StudyResultParameter struct {
	Result bool `json:"result"`
}

type ProblemWithLevel struct {
	ProblemID uint `json:"problemId"`
	Level     int  `json:"leve"`
}

type ProblemWithLevelList struct {
	Results []ProblemWithLevel `json:"results"`
}
