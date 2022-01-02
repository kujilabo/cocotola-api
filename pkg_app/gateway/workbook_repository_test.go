package gateway

import (
	"context"
	"testing"

	appD "github.com/kujilabo/cocotola-api/pkg_app/domain"
	userD "github.com/kujilabo/cocotola-api/pkg_user/domain"
	userG "github.com/kujilabo/cocotola-api/pkg_user/gateway"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func Test_workbookRepository_FindPersonalWorkbooks(t *testing.T) {
	logrus.SetLevel(logrus.DebugLevel)
	bg := context.Background()

	userD.InitSystemAdmin(nil)
	for driverName, db := range dbList() {
		logrus.Println(driverName)
		sqlDB, err := db.DB()
		assert.NoError(t, err)
		defer sqlDB.Close()
		userRepo, err := userG.NewRepositoryFactory(db)
		assert.NoError(t, err)
		_, sysOwner, owner := testInitOrganization(t, db)
		appUserRepo := userG.NewAppUserRepository(nil, db)

		rbacRepo := userG.NewRBACRepository(db)
		err = rbacRepo.Init()
		assert.NoError(t, err)

		userID1, err := appUserRepo.AddAppUser(bg, owner, testNewAppUserAddParameter(t, "LOGIN_ID_1", "USERNAME_1"))
		assert.NoError(t, err)
		user1, err := appUserRepo.FindAppUserByID(bg, owner, userID1)
		assert.NoError(t, err)
		assert.Equal(t, "LOGIN_ID_1", user1.GetLoginID())
		userID2, err := appUserRepo.AddAppUser(bg, owner, testNewAppUserAddParameter(t, "LOGIN_ID_2", "USERNAME_2"))
		assert.NoError(t, err)
		user2, err := appUserRepo.FindAppUserByID(bg, owner, userID2)
		assert.NoError(t, err)
		assert.Equal(t, "LOGIN_ID_2", user2.GetLoginID())

		englishWord := testNewProblemType(t, "english_word_problem")
		workbookRepo := NewWorkbookRepository(bg, driverName, nil, userRepo, nil, db, []appD.ProblemType{englishWord})
		spaceRepo := userG.NewSpaceRepository(db)

		// user1 has two workbooks
		student1 := testNewStudent(t, user1)
		spaceID1, err := spaceRepo.AddPersonalSpace(bg, sysOwner, user1)
		assert.NoError(t, err)
		workbookID11, err := workbookRepo.AddWorkbook(bg, student1, spaceID1, testNewWorkbookAddParameter(t, "WB11"))
		assert.NoError(t, err)
		workbookID12, err := workbookRepo.AddWorkbook(bg, student1, spaceID1, testNewWorkbookAddParameter(t, "WB12"))
		assert.NoError(t, err)

		// user2 has one workbook
		student2 := testNewStudent(t, user2)
		spaceID2, err := spaceRepo.AddPersonalSpace(bg, sysOwner, user2)
		assert.NoError(t, err)
		workbookID21, err := workbookRepo.AddWorkbook(bg, student2, spaceID2, testNewWorkbookAddParameter(t, "WB21"))
		assert.NoError(t, err)
		assert.GreaterOrEqual(t, uint(workbookID21), uint(1))

		type args struct {
			operator appD.Student
			param    appD.WorkbookSearchCondition
		}
		type want struct {
			workbookID   appD.WorkbookID
			workbookName string
		}
		tests := []struct {
			name    string
			args    args
			want    []want
			wantErr bool
		}{
			{
				name: "user1",
				args: args{
					operator: student1,
					param:    testNewWorkbookSearchCondition(t),
				},
				want: []want{
					{
						workbookID:   workbookID11,
						workbookName: "WB11",
					},
					{
						workbookID:   workbookID12,
						workbookName: "WB12",
					},
				},
			},
			{
				name: "user2",
				args: args{
					operator: student2,
					param:    testNewWorkbookSearchCondition(t),
				},
				want: []want{
					{
						workbookID:   workbookID21,
						workbookName: "WB21",
					},
				},
			},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				got, err := workbookRepo.FindPersonalWorkbooks(bg, tt.args.operator, tt.args.param)
				if (err != nil) != tt.wantErr {
					t.Errorf("workbookRepository.FindPersonalWorkbooks() error = %v, wantErr %v", err, tt.wantErr)
					return
				}

				if err == nil {
					assert.Len(t, got.Results, len(tt.want))
					for i, want := range tt.want {
						assert.Equal(t, uint(want.workbookID), got.Results[i].GetID())
						assert.Equal(t, want.workbookName, got.Results[i].GetName())
					}
					assert.Equal(t, int64(len(tt.want)), got.TotalCount)
				}
			})
		}
	}
}

func testNewProblemType(t *testing.T, name string) appD.ProblemType {
	p, err := appD.NewProblemType(1, name)
	assert.NoError(t, err)
	return p
}

func testNewStudent(t *testing.T, appUser userD.AppUser) appD.Student {
	s, err := appD.NewStudent(nil, nil, appUser)
	assert.NoError(t, err)
	return s
}

func testNewWorkbookSearchCondition(t *testing.T) appD.WorkbookSearchCondition {
	p, err := appD.NewWorkbookSearchCondition(1, 10, []userD.SpaceID{})
	assert.NoError(t, err)
	return p
}

func testNewWorkbookAddParameter(t *testing.T, name string) appD.WorkbookAddParameter {
	p, err := appD.NewWorkbookAddParameter("english_word_problem", name, "", map[string]string{"audioEnabled": "false"})
	assert.NoError(t, err)
	return p
}

func Test_workbookRepository_FindWorkbookByName(t *testing.T) {
	logrus.SetLevel(logrus.DebugLevel)
	bg := context.Background()

	userD.InitSystemAdmin(nil)
	for driverName, db := range dbList() {
		logrus.Println(driverName)
		sqlDB, err := db.DB()
		assert.NoError(t, err)
		defer sqlDB.Close()
		userRepo, err := userG.NewRepositoryFactory(db)
		assert.NoError(t, err)
		_, sysOwner, owner := testInitOrganization(t, db)
		appUserRepo := userG.NewAppUserRepository(nil, db)

		rbacRepo := userG.NewRBACRepository(db)
		err = rbacRepo.Init()
		assert.NoError(t, err)

		userID1, err := appUserRepo.AddAppUser(bg, owner, testNewAppUserAddParameter(t, "LOGIN_ID_1", "USERNAME_1"))
		assert.NoError(t, err)
		user1, err := appUserRepo.FindAppUserByID(bg, owner, userID1)
		assert.NoError(t, err)
		assert.Equal(t, "LOGIN_ID_1", user1.GetLoginID())
		userID2, err := appUserRepo.AddAppUser(bg, owner, testNewAppUserAddParameter(t, "LOGIN_ID_2", "USERNAME_2"))
		assert.NoError(t, err)
		user2, err := appUserRepo.FindAppUserByID(bg, owner, userID2)
		assert.NoError(t, err)
		assert.Equal(t, "LOGIN_ID_2", user2.GetLoginID())

		englishWord := testNewProblemType(t, "english_word_problem")
		workbookRepo := NewWorkbookRepository(bg, driverName, nil, userRepo, nil, db, []appD.ProblemType{englishWord})
		spaceRepo := userG.NewSpaceRepository(db)

		// user1 has two workbooks
		student1 := testNewStudent(t, user1)
		spaceID1, err := spaceRepo.AddPersonalSpace(bg, sysOwner, user1)
		assert.NoError(t, err)
		workbookID11, err := workbookRepo.AddWorkbook(bg, student1, spaceID1, testNewWorkbookAddParameter(t, "WB11"))
		if err != nil {
			panic(err)
		}
		assert.NoError(t, err)
		assert.Greater(t, int(workbookID11), 0)
		workbookID12, err := workbookRepo.AddWorkbook(bg, student1, spaceID1, testNewWorkbookAddParameter(t, "WB12"))
		assert.NoError(t, err)
		assert.Greater(t, int(workbookID12), 0)

		type args struct {
			operator appD.Student
			param    string
		}
		type want struct {
			workbookID   appD.WorkbookID
			workbookName string
			audioEnabled string
		}
		tests := []struct {
			name    string
			args    args
			want    want
			wantErr bool
		}{
			{
				name: "user1",
				args: args{
					operator: student1,
					param:    "WB11",
				},
				want: want{
					workbookID:   workbookID11,
					workbookName: "WB11",
					audioEnabled: "false",
				},
			},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				got, err := workbookRepo.FindWorkbookByName(bg, tt.args.operator, tt.args.param)
				if (err != nil) != tt.wantErr {
					t.Errorf("workbookRepository.FindWorkbookByName() error = %v, wantErr %v", err, tt.wantErr)
					return
				}

				if err == nil {
					assert.Equal(t, uint(tt.want.workbookID), got.GetID())
					assert.Equal(t, tt.want.workbookName, got.GetName())
					assert.Equal(t, tt.want.audioEnabled, got.GetProperties()["audioEnabled"])
				}
			})
		}
	}

}
