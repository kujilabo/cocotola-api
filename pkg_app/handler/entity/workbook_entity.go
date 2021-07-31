package entity

import (
	"time"

	"github.com/go-playground/validator"
	user "github.com/kujilabo/cocotola-api/pkg_user/domain"
)

type Model struct {
	ID        uint      `validate:"required,gte=1" json:"id"`
	Version   int       `validate:"required,gte=1" json:"version"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	CreatedBy uint      `validate:"required,gte=1" json:"createdBy"`
	UpdatedBy uint      `validate:"required,gte=1" json:"updatedBy"`
}

func NewModel(model user.Model) (Model, error) {
	m := Model{
		ID:        model.GetID(),
		Version:   model.GetVersion(),
		CreatedAt: model.GetCreatedAt(),
		UpdatedAt: model.GetUpdatedAt(),
		CreatedBy: model.GetCreatedBy(),
		UpdatedBy: model.GetUpdatedBy(),
	}

	v := validator.New()
	return m, v.Struct(&m)
}

type Workbook struct {
	Model
	Name         string `json:"name"`
	ProblemType  string `json:"problemType"`
	QuestionText string `json:"questionText"`
}

type WorkbookSearchResponse struct {
	TotalCount int64      `json:"totalCount"`
	Results    []Workbook `json:"results"`
}

type WorkbookAddParameter struct {
	Name         string `json:"name" binding:"required"`
	ProblemType  string `json:"problemType" enums:"english_word,english_phrase,template"`
	QuestionText string `json:"questionText"`
}

type WorkbookUpdateParameter struct {
	Name         string `json:"name" binding:"required"`
	QuestionText string `json:"questionText"`
}

type WorkbookWithProblems struct {
	Model
	Name         string    `json:"name"`
	ProblemType  string    `json:"problemType"`
	QuestionText string    `json:"questionText"`
	Problems     []Problem `json:"problems"`
	Subscribed   bool      `json:"subscribed"`
}