package entity

import (
	"time"

	libD "github.com/kujilabo/cocotola-api/src/lib/domain"
	userD "github.com/kujilabo/cocotola-api/src/user/domain"
)

type Model struct {
	ID        uint      `json:"id" validate:"required,gte=1"`
	Version   int       `json:"version" validate:"required,gte=1"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	CreatedBy uint      `json:"createdBy" validate:"required,gte=1"`
	UpdatedBy uint      `json:"updatedBy" validate:"required,gte=1"`
}

func NewModel(model userD.Model) (Model, error) {
	m := Model{
		ID:        model.GetID(),
		Version:   model.GetVersion(),
		CreatedAt: model.GetCreatedAt(),
		UpdatedAt: model.GetUpdatedAt(),
		CreatedBy: model.GetCreatedBy(),
		UpdatedBy: model.GetUpdatedBy(),
	}

	return m, libD.Validator.Struct(m)
}

type Workbook struct {
	Model
	Name         string `json:"name" validate:"required"`
	Lang2        string `json:"lang2" validate:"required,len=2"`
	ProblemType  string `json:"problemType" validate:"required"`
	QuestionText string `json:"questionText"`
}

type WorkbookSearchResponse struct {
	TotalCount int         `json:"totalCount" validate:"gte=0"`
	Results    []*Workbook `json:"results" validate:"dive"`
}

type WorkbookAddParameter struct {
	Name         string `json:"name" binding:"required"`
	ProblemType  string `json:"problemType" binding:"required"`
	QuestionText string `json:"questionText"`
}

type WorkbookUpdateParameter struct {
	Name         string `json:"name" binding:"required"`
	QuestionText string `json:"questionText"`
}

type WorkbookWithProblems struct {
	Model
	Name         string     `json:"name" binding:"required"`
	Lang2        string     `json:"lang2" validate:"required,len=2"`
	ProblemType  string     `json:"problemType" binding:"required"`
	QuestionText string     `json:"questionText"`
	Problems     []*Problem `json:"problems" validate:"dive"`
	Subscribed   bool       `json:"subscribed"`
}
