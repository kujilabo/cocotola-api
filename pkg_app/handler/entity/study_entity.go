package entity

type StudyResultParameter struct {
	Result    bool `json:"result"`
	Memorized bool `json:"memorized"`
}

type ProblemWithLevel struct {
	ProblemID uint `json:"problemId"`
	Level     int  `json:"level"`
	Memorized bool `json:"memorized"`
}

type ProblemWithLevelList struct {
	Records []ProblemWithLevel `json:"records"`
}

type IntValue struct {
	Value int `json:"value"`
}
