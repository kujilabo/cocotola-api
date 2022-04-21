//go:generate mockery --output mock --name ProblemModel
package domain

import (
	"context"

	lib "github.com/kujilabo/cocotola-api/pkg_lib/domain"
	user "github.com/kujilabo/cocotola-api/pkg_user/domain"
)

type ProblemID uint

type ProblemModel interface {
	user.Model
	GetNumber() int
	GetProblemType() string
	GetProperties(ctx context.Context) map[string]interface{}
}

type problemModel struct {
	user.Model
	Number      int                    `validate:"required"`
	ProblemType string                 `validate:"required"`
	Properties  map[string]interface{} `validate:"required"`
}

func NewProblemModel(model user.Model, number int, problemType string, properties map[string]interface{}) (ProblemModel, error) {
	m := &problemModel{
		Model:       model,
		Number:      number,
		ProblemType: problemType,
		Properties:  properties,
	}

	return m, lib.Validator.Struct(m)
}

func (m *problemModel) GetNumber() int {
	return m.Number
}

func (m *problemModel) GetProblemType() string {
	return m.ProblemType
}

func (m *problemModel) GetProperties(ctx context.Context) map[string]interface{} {
	return m.Properties
}
