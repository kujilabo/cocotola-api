package entity

import "encoding/json"

type ProblemFindParameter struct {
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
	ProblemType string          `json:"problemType" validate:"required"`
	Properties  json.RawMessage `json:"properties"`
}

type ProblemFindResponse struct {
	TotalCount int        `json:"totalCount" validate:"gte=0"`
	Results    []*Problem `json:"results" validate:"dive"`
}

type SimpleProblem struct {
	ID          uint            `json:"id" validate:"required,gte=1"`
	Number      int             `json:"number"`
	ProblemType string          `json:"problemType" validate:"required"`
	Properties  json.RawMessage `json:"properties"`
}

type ProblemFindAllResponse struct {
	TotalCount int              `json:"totalCount" validate:"gte=0"`
	Results    []*SimpleProblem `json:"results" validate:"dive"`
}

type ProblemAddParameter struct {
	Number     int             `json:"number" binding:"required"`
	Properties json.RawMessage `json:"properties"`
}

type ProblemUpdateParameter struct {
	Number     int             `json:"number" binding:"required"`
	Properties json.RawMessage `json:"properties"`
}

type ProblemIDs struct {
	Results []uint `json:"results"`
}
