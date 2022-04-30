//go:generate mockery --output mock --name ProblemModel
package domain

import (
	"context"
	"errors"

	libD "github.com/kujilabo/cocotola-api/src/lib/domain"
	userD "github.com/kujilabo/cocotola-api/src/user/domain"
)

type ProblemID uint

type ProblemModel interface {
	userD.Model
	GetNumber() int
	GetProblemType() string
	GetProperties(ctx context.Context) map[string]interface{}
}

type problemModel struct {
	userD.Model
	Number      int    `validate:"required"`
	ProblemType string `validate:"required"`
	// Properties  map[string]interface{} `validate:"required"`
}

func NewProblemModel(model userD.Model, number int, problemType string, properties map[string]interface{}) (ProblemModel, error) {
	m := &problemModel{
		Model:       model,
		Number:      number,
		ProblemType: problemType,
		// Properties:  properties,
	}

	return m, libD.Validator.Struct(m)
}

func (m *problemModel) GetNumber() int {
	return m.Number
}

func (m *problemModel) GetProblemType() string {
	return m.ProblemType
}

func (m *problemModel) GetProperties(ctx context.Context) map[string]interface{} {
	panic(errors.New("not implemented"))
}
