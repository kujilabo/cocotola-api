//go:generate mockery --output mock --name WorkbookModel
package domain

import (
	libD "github.com/kujilabo/cocotola-api/src/lib/domain"
	userD "github.com/kujilabo/cocotola-api/src/user/domain"
)

type WorkbookID uint

type WorkbookModel interface {
	userD.Model
	GetSpaceID() userD.SpaceID
	GetOwnerID() userD.AppUserID
	GetName() string
	GetLang2() Lang2
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
	Lang2        Lang2            `validate:"required,len=2"`
	ProblemType  string           `validate:"required"`
	QuestionText string
	Properties   map[string]string
}

func NewWorkbookModel(model userD.Model, spaceID userD.SpaceID, ownerID userD.AppUserID, privileges userD.Privileges, name string, lang2 Lang2, problemType string, questsionText string, properties map[string]string) (WorkbookModel, error) {
	m := &workbookModel{
		Model:        model,
		spaceID:      spaceID,
		ownerID:      ownerID,
		privileges:   privileges,
		Name:         name,
		Lang2:        lang2,
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

func (m *workbookModel) GetLang2() Lang2 {
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
