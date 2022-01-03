package entity

import "encoding/json"

type ProblemSearchParameter struct {
	PageNo   int    `json:"pageNo" binding:"required,gte=1"`
	PageSize int    `json:"pageSize" binding:"required,gte=1"`
	Keyword  string `json:"keyword"`
}

type ProblemIDsParameter struct {
	IDs []uint `json:"ids"`
}

type Problem struct {
	Model
	Number      int             `json:"number"`
	ProblemType string          `json:"problemType"`
	Properties  json.RawMessage `json:"properties"`
}

type ProblemSearchResponse struct {
	TotalCount int64     `json:"totalCount"`
	Results    []Problem `json:"results"`
}

type SimpleProblem struct {
	ID          uint            `validate:"required,gte=1" json:"id"`
	Number      int             `json:"number"`
	ProblemType string          `json:"problemType"`
	Properties  json.RawMessage `json:"properties"`
}

type ProblemFindAllResponse struct {
	TotalCount int64           `json:"totalCount"`
	Results    []SimpleProblem `json:"results"`
}

type ProblemAddParameter struct {
	Number      int             `json:"number" binding:"required"`
	ProblemType string          `json:"problemType"`
	Properties  json.RawMessage `json:"properties"`
}

type ProblemIDs struct {
	Results []uint `json:"results"`
}
