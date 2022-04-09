package gateway_test

import (
	"context"
	"errors"
	"log"
	"strconv"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"

	"github.com/kujilabo/cocotola-api/pkg_user/domain"
	"github.com/kujilabo/cocotola-api/pkg_user/gateway"
	"github.com/kujilabo/cocotola-api/pkg_user/service"
)

func Test_spaceRepository_FindDefaultSpace(t *testing.T) {
	logrus.SetLevel(logrus.DebugLevel)
	bg := context.Background()

	userRfFunc := func(ctx context.Context, db *gorm.DB) (service.RepositoryFactory, error) {
		return gateway.NewRepositoryFactory(db)
	}

	service.InitSystemAdmin(userRfFunc)
	for i, db := range dbList() {
		log.Printf("%d", i)
		sqlDB, err := db.DB()
		assert.NoError(t, err)
		defer sqlDB.Close()

		orgID, owner := testInitOrganization(t, db)

		type args struct {
			operator domain.AppUserModel
		}

		model, err := domain.NewModel(1, 1, time.Now(), time.Now(), 1, 1)
		assert.NoError(t, err)
		spaceModel, err := domain.NewSpaceModel(model, orgID, 1, "default", "Default", "")
		assert.NoError(t, err)
		space, err := service.NewSpace(spaceModel)
		assert.NoError(t, err)
		tests := []struct {
			name string
			args args
			want service.Space
			err  error
		}{
			{
				name: "",
				args: args{
					operator: owner,
				},
				want: space,
				err:  nil,
			},
		}
		spaceRepo := gateway.NewSpaceRepository(db)
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				got, err := spaceRepo.FindDefaultSpace(bg, tt.args.operator)
				if err != nil && !errors.Is(err, tt.err) {
					t.Errorf("spaceRepository.FindDefaultSpace() error = %v, err %v", err, tt.err)
					return
				}
				if err == nil {
					assert.Equal(t, space.GetKey(), got.GetKey())
					assert.Equal(t, space.GetName(), got.GetName())
					assert.Equal(t, space.GetDescription(), got.GetDescription())
				}
			})
		}
	}
}

func Test_spaceRepository_FindPersonalSpace(t *testing.T) {
	logrus.SetLevel(logrus.DebugLevel)
	bg := context.Background()

	userRfFunc := func(ctx context.Context, db *gorm.DB) (service.RepositoryFactory, error) {
		return gateway.NewRepositoryFactory(db)
	}

	service.InitSystemAdmin(userRfFunc)
	for i, db := range dbList() {
		log.Printf("%d", i)
		sqlDB, err := db.DB()
		assert.NoError(t, err)
		defer sqlDB.Close()

		orgID, owner := testInitOrganization(t, db)

		type args struct {
			operator domain.AppUserModel
		}

		model, err := domain.NewModel(1, 1, time.Now(), time.Now(), 1, 1)
		assert.NoError(t, err)
		spaceModel, err := domain.NewSpaceModel(model, orgID, 1, strconv.Itoa(int(owner.GetID())), "Default", "")
		assert.NoError(t, err)
		space, err := service.NewSpace(spaceModel)
		assert.NoError(t, err)
		tests := []struct {
			name string
			args args
			want service.Space
			err  error
		}{
			{
				name: "",
				args: args{
					operator: owner,
				},
				want: space,
				err:  nil,
			},
		}
		spaceRepo := gateway.NewSpaceRepository(db)
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				got, err := spaceRepo.FindPersonalSpace(bg, tt.args.operator)
				if err != nil && !errors.Is(err, tt.err) {
					t.Errorf("spaceRepository.FindPersonalSpace() error = %v, err %v", err, tt.err)
					return
				}
				if err == nil {
					assert.Equal(t, space.GetKey(), got.GetKey())
					assert.Equal(t, space.GetName(), got.GetName())
					assert.Equal(t, space.GetDescription(), got.GetDescription())
				}
			})
		}
	}
}
