package usecase

import (
	"context"
	"time"

	"gorm.io/gorm"

	"github.com/kujilabo/cocotola-api/src/auth/service"
	liberrors "github.com/kujilabo/cocotola-api/src/lib/errors"
	userD "github.com/kujilabo/cocotola-api/src/user/domain"
	userS "github.com/kujilabo/cocotola-api/src/user/service"
)

type GuestUserUsecase interface {
	RetrieveGuestToken(ctx context.Context, organizationName string) (*service.TokenSet, error)
}

type guestUserUsecase struct {
	db               *gorm.DB
	authTokenManager service.AuthTokenManager
}

func NewGuestUserUsecase(authTokenManager service.AuthTokenManager) GuestUserUsecase {
	return &guestUserUsecase{
		authTokenManager: authTokenManager,
	}
}

func (s *guestUserUsecase) RetrieveGuestToken(ctx context.Context, organizationName string) (*service.TokenSet, error) {
	var tokenSet *service.TokenSet
	if err := s.db.Transaction(func(tx *gorm.DB) error {
		systemAdmin, err := userS.NewSystemAdminFromDB(ctx, tx)
		if err != nil {
			return err
		}

		systemOwner, err := systemAdmin.FindSystemOwnerByOrganizationName(ctx, organizationName)
		if err != nil {
			return liberrors.Errorf("failed to FindSystemOwnerByOrganizationName. err: %w", err)
		}

		// guest, err := systemOwner.FindAppUserByLoginID(ctx, "guest")
		// if err != nil {
		// 	return nil, fmt.Errorf("failed to FindAppUserByLoginID. err: %w", err)
		// }

		organization, err := systemOwner.GetOrganization(ctx)
		if err != nil {
			return liberrors.Errorf("failed to GetOrganization. err: %w", err)
		}

		model, err := userD.NewModel(0, 1, time.Now(), time.Now(), 0, 0)
		if err != nil {
			return liberrors.Errorf("failed to FindAppUserByLoginID. err: %w", err)
		}

		guest, err := userD.NewAppUserModel(model, userD.OrganizationID(organization.GetID()), "guest", "Guest", []string{}, map[string]string{})
		if err != nil {
			return liberrors.Errorf("failed to FindAppUserByLoginID. err: %w", err)
		}

		tokenSetTmp, err := s.authTokenManager.CreateTokenSet(ctx, guest, organization)
		if err != nil {
			return err
		}

		tokenSet = tokenSetTmp
		return nil
	}); err != nil {
		return nil, err
	}
	return tokenSet, nil
}
