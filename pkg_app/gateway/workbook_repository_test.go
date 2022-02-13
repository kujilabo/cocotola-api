package gateway_test

import (
	"context"
	"errors"
	"testing"

	"github.com/kujilabo/cocotola-api/pkg_app/domain"
	appD "github.com/kujilabo/cocotola-api/pkg_app/domain"
	"github.com/kujilabo/cocotola-api/pkg_app/gateway"
	userD "github.com/kujilabo/cocotola-api/pkg_user/domain"
	userG "github.com/kujilabo/cocotola-api/pkg_user/gateway"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func Test_workbookRepository_FindPersonalWorkbooks(t *testing.T) {
	logrus.SetLevel(logrus.DebugLevel)
	bg := context.Background()
	userRfFunc := func(db *gorm.DB) (userD.RepositoryFactory, error) {
		return userG.NewRepositoryFactory(db)
	}

	userD.InitSystemAdmin(userRfFunc)
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
		workbookRepo := gateway.NewWorkbookRepository(bg, driverName, nil, userRepo, nil, db, []appD.ProblemType{englishWord})
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
					assert.Len(t, got.GetResults(), len(tt.want))
					for i, want := range tt.want {
						assert.Equal(t, uint(want.workbookID), got.GetResults()[i].GetID())
						assert.Equal(t, want.workbookName, got.GetResults()[i].GetName())
					}
					assert.Equal(t, len(tt.want), got.GetTotalCount())
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
	s, err := appD.NewStudent(nil, nil, nil, appUser)
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
func testNewAppUser(t *testing.T, ctx context.Context, db *gorm.DB, owner userD.Owner, loginID, username string) userD.AppUser {
	appUserRepo := userG.NewAppUserRepository(nil, db)
	userID1, err := appUserRepo.AddAppUser(ctx, owner, testNewAppUserAddParameter(t, loginID, username))
	assert.NoError(t, err)
	user1, err := appUserRepo.FindAppUserByID(ctx, owner, userID1)
	assert.NoError(t, err)
	assert.Equal(t, loginID, user1.GetLoginID())
	return user1
}
func testNewWorkbook(t *testing.T, ctx context.Context, db *gorm.DB, workbookRepo appD.WorkbookRepository, student appD.Student, spaceID userD.SpaceID, workbookName string) appD.Workbook {
	workbookID11, err := workbookRepo.AddWorkbook(ctx, student, spaceID, testNewWorkbookAddParameter(t, workbookName))
	assert.NoError(t, err)
	assert.Greater(t, int(workbookID11), 0)
	workbook, err := workbookRepo.FindWorkbookByID(ctx, student, workbookID11)
	assert.NoError(t, err)
	return workbook
}
func Test_workbookRepository_FindWorkbookByName(t *testing.T) {
	logrus.SetLevel(logrus.DebugLevel)
	bg := context.Background()

	userRfFunc := func(db *gorm.DB) (userD.RepositoryFactory, error) {
		return userG.NewRepositoryFactory(db)
	}

	userD.InitSystemAdmin(userRfFunc)
	for driverName, db := range dbList() {
		logrus.Println(driverName)
		sqlDB, err := db.DB()
		assert.NoError(t, err)
		defer sqlDB.Close()
		userRepo, err := userG.NewRepositoryFactory(db)
		assert.NoError(t, err)
		_, sysOwner, owner := testInitOrganization(t, db)

		rbacRepo := userG.NewRBACRepository(db)
		err = rbacRepo.Init()
		assert.NoError(t, err)

		user1 := testNewAppUser(t, bg, db, owner, "LOGIN_ID_1", "USERNAME_1")
		testNewAppUser(t, bg, db, owner, "LOGIN_ID_2", "USERNAME_2")

		englishWord := testNewProblemType(t, "english_word_problem")
		workbookRepo := gateway.NewWorkbookRepository(bg, driverName, nil, userRepo, nil, db, []appD.ProblemType{englishWord})
		spaceRepo := userG.NewSpaceRepository(db)

		// user1 has two workbooks
		student1 := testNewStudent(t, user1)
		spaceID1, err := spaceRepo.AddPersonalSpace(bg, sysOwner, user1)
		assert.NoError(t, err)

		workbook11 := testNewWorkbook(t, bg, db, workbookRepo, student1, spaceID1, "WB11")
		testNewWorkbook(t, bg, db, workbookRepo, student1, spaceID1, "WB12")

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
					workbookID:   appD.WorkbookID(workbook11.GetID()),
					workbookName: "WB11",
					audioEnabled: "false",
				},
			},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				got, err := workbookRepo.FindWorkbookByName(bg, tt.args.operator, spaceID1, tt.args.param)
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

func Test_workbookRepository_FindWorkbookByID_priv(t *testing.T) {
	logrus.SetLevel(logrus.DebugLevel)
	bg := context.Background()

	userRfFunc := func(db *gorm.DB) (userD.RepositoryFactory, error) {
		return userG.NewRepositoryFactory(db)
	}

	userD.InitSystemAdmin(userRfFunc)
	for driverName, db := range dbList() {
		logrus.Println(driverName)
		sqlDB, err := db.DB()
		assert.NoError(t, err)
		defer sqlDB.Close()
		userRepo, err := userG.NewRepositoryFactory(db)
		assert.NoError(t, err)
		_, sysOwner, owner := testInitOrganization(t, db)

		rbacRepo := userG.NewRBACRepository(db)
		err = rbacRepo.Init()
		assert.NoError(t, err)

		user1 := testNewAppUser(t, bg, db, owner, "LOGIN_ID_1", "USERNAME_1")
		user2 := testNewAppUser(t, bg, db, owner, "LOGIN_ID_2", "USERNAME_2")

		englishWord := testNewProblemType(t, "english_word_problem")
		workbookRepo := gateway.NewWorkbookRepository(bg, driverName, nil, userRepo, nil, db, []appD.ProblemType{englishWord})
		spaceRepo := userG.NewSpaceRepository(db)

		// user1 has two workbooks(WB11, WB12)
		student1 := testNewStudent(t, user1)
		spaceID1, err := spaceRepo.AddPersonalSpace(bg, sysOwner, user1)
		assert.NoError(t, err)

		workbook11 := testNewWorkbook(t, bg, db, workbookRepo, student1, spaceID1, "WB11")
		workbook12 := testNewWorkbook(t, bg, db, workbookRepo, student1, spaceID1, "WB12")

		// user2 has two workbooks(WB11, WB12)
		student2 := testNewStudent(t, user2)
		spaceID2, err := spaceRepo.AddPersonalSpace(bg, sysOwner, user2)
		assert.NoError(t, err)

		workbook21 := testNewWorkbook(t, bg, db, workbookRepo, student2, spaceID2, "WB21")
		workbook22 := testNewWorkbook(t, bg, db, workbookRepo, student2, spaceID2, "WB22")

		// user1 can read user1's workbooks(WB11, WB12)
		workbook11Tmp, err := workbookRepo.FindWorkbookByID(bg, student1, appD.WorkbookID(workbook11.GetID()))
		assert.NoError(t, err)
		assert.Equal(t, workbook11Tmp.GetID(), workbook11.GetID())
		workbook12Tmp, err := workbookRepo.FindWorkbookByID(bg, student1, appD.WorkbookID(workbook12.GetID()))
		assert.NoError(t, err)
		assert.Equal(t, workbook12Tmp.GetID(), workbook12.GetID())

		// user1 cannot read user2's workbooks(WB21, WB22)
		if _, err := workbookRepo.FindWorkbookByID(bg, student1, appD.WorkbookID(workbook21.GetID())); err != nil {
			assert.True(t, errors.Is(err, domain.ErrWorkbookPermissionDenied))
		} else {
			assert.Fail(t, "err is nil")
		}
		if _, err := workbookRepo.FindWorkbookByID(bg, student1, appD.WorkbookID(workbook22.GetID())); err != nil {
			assert.True(t, errors.Is(err, domain.ErrWorkbookPermissionDenied))
		} else {
			assert.Fail(t, "err is nil")
		}

	}

}
