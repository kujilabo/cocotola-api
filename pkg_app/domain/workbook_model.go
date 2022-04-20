//go:generate mockery --output mock --name WorkbookModel
package domain

import (
	lib "github.com/kujilabo/cocotola-api/pkg_lib/domain"
	user "github.com/kujilabo/cocotola-api/pkg_user/domain"
)

type WorkbookID uint

type WorkbookModel interface {
	user.Model
	GetSpaceID() user.SpaceID
	GetOwnerID() user.AppUserID
	GetName() string
	GetProblemType() string
	GetQuestionText() string
	GetProperties() map[string]string
	HasPrivilege(privilege user.RBACAction) bool
}

type workbookModel struct {
	user.Model
	spaceID      user.SpaceID    `validate:"required"`
	ownerID      user.AppUserID  `validate:"required"`
	privileges   user.Privileges `validate:"required"`
	Name         string          `validate:"required"`
	ProblemType  string          `validate:"required"`
	QuestionText string
	Properties   map[string]string
}

func NewWorkbookModel(model user.Model, spaceID user.SpaceID, ownerID user.AppUserID, privileges user.Privileges, name string, problemType string, questsionText string, properties map[string]string) (WorkbookModel, error) {
	m := &workbookModel{
		privileges: privileges,
		Model:      model,
		spaceID:    spaceID,
		ownerID:    ownerID,
		// Properties: workbookProperties{
		// 	Name:         name,
		// 	ProblemType:  problemType,
		// 	QuestionText: questsionText,
		// },

		Name:         name,
		ProblemType:  problemType,
		QuestionText: questsionText,
		Properties:   properties,
	}

	return m, lib.Validator.Struct(m)
}

func (m *workbookModel) GetSpaceID() user.SpaceID {
	return m.spaceID
}

func (m *workbookModel) GetOwnerID() user.AppUserID {
	return m.ownerID
}

func (m *workbookModel) GetName() string {
	return m.Name
}

func (m *workbookModel) GetProblemType() string {
	return m.ProblemType
}

func (m *workbookModel) GetQuestionText() string {
	return m.QuestionText
}

func (m *workbookModel) GetProperties() map[string]string {
	return m.Properties
}

func (m *workbookModel) HasPrivilege(privilege user.RBACAction) bool {
	return m.privileges.HasPrivilege(privilege)
}
