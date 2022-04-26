//go:generate mockery --output mock --name WorkbookModel
package domain

import (
	libD "github.com/kujilabo/cocotola-api/pkg_lib/domain"
	userD "github.com/kujilabo/cocotola-api/pkg_user/domain"
)

type WorkbookID uint

type WorkbookModel interface {
	userD.Model
	GetSpaceID() userD.SpaceID
	GetOwnerID() userD.AppUserID
	GetName() string
	GetLang() Lang2
	GetProblemType() string
	GetQuestionText() string
	GetProperties() map[string]string
	HasPrivilege(privilege userD.RBACAction) bool
}

type workbookModel struct {
	userD.Model
	spaceID      userD.SpaceID    `validate:"required"`
	ownerID      userD.AppUserID  `validate:"required"`
	privileges   userD.Privileges `validate:"required"`
	Name         string           `validate:"required"`
	Lang2        Lang2
	ProblemType  string `validate:"required"`
	QuestionText string
	Properties   map[string]string
}

func NewWorkbookModel(model userD.Model, spaceID userD.SpaceID, ownerID userD.AppUserID, privileges userD.Privileges, name string, lang Lang2, problemType string, questsionText string, properties map[string]string) (WorkbookModel, error) {
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
		Lang2:        lang,
		ProblemType:  problemType,
		QuestionText: questsionText,
		Properties:   properties,
	}

	return m, libD.Validator.Struct(m)
}

func (m *workbookModel) GetSpaceID() userD.SpaceID {
	return m.spaceID
}

func (m *workbookModel) GetOwnerID() userD.AppUserID {
	return m.ownerID
}

func (m *workbookModel) GetName() string {
	return m.Name
}

func (m *workbookModel) GetLang() Lang2 {
	return m.Lang2
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

func (m *workbookModel) HasPrivilege(privilege userD.RBACAction) bool {
	return m.privileges.HasPrivilege(privilege)
}
