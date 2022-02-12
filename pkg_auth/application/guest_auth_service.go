package application

import (
	"context"
	"time"

	"golang.org/x/xerrors"
	"gorm.io/gorm"

	"github.com/kujilabo/cocotola-api/pkg_auth/domain"
	user "github.com/kujilabo/cocotola-api/pkg_user/domain"
)

type GuestAuthService interface {
	RetrieveGuestToken(ctx context.Context, organizationName string) (*domain.TokenSet, error)
}

type guestAuthService struct {
	db               *gorm.DB
	authTokenManager domain.AuthTokenManager
}

func NewGuestAuthService(authTokenManager domain.AuthTokenManager) GuestAuthService {
	return &guestAuthService{
		authTokenManager: authTokenManager,
	}
}

func (s *guestAuthService) RetrieveGuestToken(ctx context.Context, organizationName string) (*domain.TokenSet, error) {
	var tokenSet *domain.TokenSet
	if err := s.db.Transaction(func(tx *gorm.DB) error {
		systemAdmin, err := user.NewSystemAdminFromDB(tx)
		if err != nil {
			return err
		}

		systemOwner, err := systemAdmin.FindSystemOwnerByOrganizationName(ctx, organizationName)
		if err != nil {
			return xerrors.Errorf("failed to FindSystemOwnerByOrganizationName. err: %w", err)
		}

		// guest, err := systemOwner.FindAppUserByLoginID(ctx, "guest")
		// if err != nil {
		// 	return nil, fmt.Errorf("failed to FindAppUserByLoginID. err: %w", err)
		// }

		organization, err := systemOwner.GetOrganization(ctx)
		if err != nil {
			return xerrors.Errorf("failed to GetOrganization. err: %w", err)
		}

		model, err := user.NewModel(0, 1, time.Now(), time.Now(), 0, 0)
		if err != nil {
			return xerrors.Errorf("failed to FindAppUserByLoginID. err: %w", err)
		}

		guest, err := user.NewAppUser(nil, model, user.OrganizationID(organization.GetID()), "guest", "Guest", []string{}, map[string]string{})
		if err != nil {
			return xerrors.Errorf("failed to FindAppUserByLoginID. err: %w", err)
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
