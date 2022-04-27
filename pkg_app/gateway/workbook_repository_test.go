package gateway_test

import (
	"context"
	"errors"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"

	"github.com/kujilabo/cocotola-api/pkg_app/domain"
	"github.com/kujilabo/cocotola-api/pkg_app/gateway"
	"github.com/kujilabo/cocotola-api/pkg_app/service"
	userD "github.com/kujilabo/cocotola-api/pkg_user/domain"
	userG "github.com/kujilabo/cocotola-api/pkg_user/gateway"
	userS "github.com/kujilabo/cocotola-api/pkg_user/service"
)

func Test_workbookRepository_FindPersonalWorkbooks(t *testing.T) {
	logrus.SetLevel(logrus.DebugLevel)
	bg := context.Background()
	userRfFunc := func(ctx context.Context, db *gorm.DB) (userS.RepositoryFactory, error) {
		return userG.NewRepositoryFactory(db)
	}

	userS.InitSystemAdmin(userRfFunc)
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
		workbookRepo := gateway.NewWorkbookRepository(bg, driverName, nil, userRepo, nil, db, []domain.ProblemType{englishWord})
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
			operator service.Student
			param    service.WorkbookSearchCondition
		}
		type want struct {
			workbookID   domain.WorkbookID
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

func testNewProblemType(t *testing.T, name string) domain.ProblemType {
	p, err := domain.NewProblemType(1, name)
	assert.NoError(t, err)
	return p
}

func testNewStudent(t *testing.T, appUser userS.AppUser) service.Student {
	s, err := service.NewStudent(nil, nil, nil, appUser)
	assert.NoError(t, err)
	return s
}

func testNewWorkbookSearchCondition(t *testing.T) service.WorkbookSearchCondition {
	p, err := service.NewWorkbookSearchCondition(1, 10, []userD.SpaceID{})
	assert.NoError(t, err)
	return p
}

func testNewWorkbookAddParameter(t *testing.T, name string) service.WorkbookAddParameter {
	p, err := service.NewWorkbookAddParameter("english_word_problem", name, domain.Lang2JA, "", map[string]string{"audioEnabled": "false"})
	assert.NoError(t, err)
	return p
}

func testNewAppUser(t *testing.T, ctx context.Context, db *gorm.DB, owner userS.Owner, loginID, username string) userS.AppUser {
	appUserRepo := userG.NewAppUserRepository(nil, db)
	userID1, err := appUserRepo.AddAppUser(ctx, owner, testNewAppUserAddParameter(t, loginID, username))
	assert.NoError(t, err)
	user1, err := appUserRepo.FindAppUserByID(ctx, owner, userID1)
	assert.NoError(t, err)
	assert.Equal(t, loginID, user1.GetLoginID())
	return user1
}

func testNewWorkbook(t *testing.T, ctx context.Context, db *gorm.DB, workbookRepo service.WorkbookRepository, student service.Student, spaceID userD.SpaceID, workbookName string) service.Workbook {
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

	userRfFunc := func(ctx context.Context, db *gorm.DB) (userS.RepositoryFactory, error) {
		return userG.NewRepositoryFactory(db)
	}

	userS.InitSystemAdmin(userRfFunc)
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
		workbookRepo := gateway.NewWorkbookRepository(bg, driverName, nil, userRepo, nil, db, []domain.ProblemType{englishWord})
		spaceRepo := userG.NewSpaceRepository(db)

		// user1 has two workbooks
		student1 := testNewStudent(t, user1)
		spaceID1, err := spaceRepo.AddPersonalSpace(bg, sysOwner, user1)
		assert.NoError(t, err)

		workbook11 := testNewWorkbook(t, bg, db, workbookRepo, student1, spaceID1, "WB11")
		testNewWorkbook(t, bg, db, workbookRepo, student1, spaceID1, "WB12")

		type args struct {
			operator service.Student
			param    string
		}
		type want struct {
			workbookID   domain.WorkbookID
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
					workbookID:   domain.WorkbookID(workbook11.GetID()),
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

	userRfFunc := func(ctx context.Context, db *gorm.DB) (userS.RepositoryFactory, error) {
		return userG.NewRepositoryFactory(db)
	}

	userS.InitSystemAdmin(userRfFunc)
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
		workbookRepo := gateway.NewWorkbookRepository(bg, driverName, nil, userRepo, nil, db, []domain.ProblemType{englishWord})
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
		workbook11Tmp, err := workbookRepo.FindWorkbookByID(bg, student1, domain.WorkbookID(workbook11.GetID()))
		assert.NoError(t, err)
		assert.Equal(t, workbook11Tmp.GetID(), workbook11.GetID())
		workbook12Tmp, err := workbookRepo.FindWorkbookByID(bg, student1, domain.WorkbookID(workbook12.GetID()))
		assert.NoError(t, err)
		assert.Equal(t, workbook12Tmp.GetID(), workbook12.GetID())

		// user1 cannot read user2's workbooks(WB21, WB22)
		if _, err := workbookRepo.FindWorkbookByID(bg, student1, domain.WorkbookID(workbook21.GetID())); err != nil {
			assert.True(t, errors.Is(err, service.ErrWorkbookPermissionDenied))
		} else {
			assert.Fail(t, "err is nil")
		}
		if _, err := workbookRepo.FindWorkbookByID(bg, student1, domain.WorkbookID(workbook22.GetID())); err != nil {
			assert.True(t, errors.Is(err, service.ErrWorkbookPermissionDenied))
		} else {
			assert.Fail(t, "err is nil")
		}

	}

}
