package domain_mock

import (
	"github.com/kujilabo/cocotola-api/pkg_user/domain"
	"github.com/stretchr/testify/mock"
)

type AppUserMock struct {
	mock.Mock
}

func (m *AppUserMock) GetID() uint {
	args := m.Called()
	return args.Get(0).(uint)
}

func (m *AppUserMock) GetOrganizationID() domain.OrganizationID {
	args := m.Called()
	return domain.OrganizationID(args.Get(0).(uint))
}

func (m *AppUserMock) GetLoginID() string {
	args := m.Called()
	return args.String(0)
}

func (m *AppUserMock) GetUsername() string {
	args := m.Called()
	return args.String(0)
}

func (m *AppUserMock) GetRoles() []string {
	args := m.Called()
	return args.Get(0).([]string)
}

func (m *AppUserMock) GetProperties() map[string]string {
	args := m.Called()
	return args.Get(0).(map[string]string)
}

func (m *AppUserMock) GetDefaultSpace() (domain.Space, error) {
	args := m.Called()
	return args.Get(0).(domain.Space), args.Error(1)
}

func (m *AppUserMock) GetPersonalSpace() (domain.Space, error) {
	args := m.Called()
	return args.Get(0).(domain.Space), args.Error(1)
}

// func Test_appUser_GetDefaultSpace(t *testing.T) {
// 	type fields struct {
// 		rf             RepositoryFactory
// 		Model          Model
// 		OrganizationID OrganizationID
// 		LoginID        string
// 		Username       string
// 		Roles          []string
// 		Properties     map[string]string
// 	}
// 	type args struct {
// 		ctx context.Context
// 	}
// 	tests := []struct {
// 		name    string
// 		fields  fields
// 		args    args
// 		want    Space
// 		wantErr bool
// 	}{
// 		// TODO: Add test cases.
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			a := &appUser{
// 				rf:             tt.fields.rf,
// 				Model:          tt.fields.Model,
// 				OrganizationID: tt.fields.OrganizationID,
// 				LoginID:        tt.fields.LoginID,
// 				Username:       tt.fields.Username,
// 				Roles:          tt.fields.Roles,
// 				Properties:     tt.fields.Properties,
// 			}
// 			got, err := a.GetDefaultSpace(tt.args.ctx)
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("appUser.GetDefaultSpace() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}
// 			if !reflect.DeepEqual(got, tt.want) {
// 				t.Errorf("appUser.GetDefaultSpace() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }
