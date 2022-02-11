package application

import (
	"context"
	"time"

	"github.com/kujilabo/cocotola-api/pkg_auth/domain"
	user "github.com/kujilabo/cocotola-api/pkg_user/domain"
	"golang.org/x/xerrors"
)

type GuestAuthService interface {
	RetrieveGuestToken(ctx context.Context, organizationName string) (*domain.TokenSet, error)
}

type guestAuthService struct {
	authTokenManager domain.AuthTokenManager
}

func NewGuestAuthService(authTokenManager domain.AuthTokenManager) GuestAuthService {
	return &guestAuthService{
		authTokenManager: authTokenManager,
	}
}

func (s *guestAuthService) RetrieveGuestToken(ctx context.Context, organizationName string) (*domain.TokenSet, error) {
	systemAdmin := user.SystemAdminInstance()
	systemOwner, err := systemAdmin.FindSystemOwnerByOrganizationName(ctx, organizationName)
	if err != nil {
		return nil, xerrors.Errorf("failed to FindSystemOwnerByOrganizationName. err: %w", err)
	}

	// guest, err := systemOwner.FindAppUserByLoginID(ctx, "guest")
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to FindAppUserByLoginID. err: %w", err)
	// }

	organization, err := systemOwner.GetOrganization(ctx)
	if err != nil {
		return nil, xerrors.Errorf("failed to GetOrganization. err: %w", err)
	}

	model, err := user.NewModel(0, 1, time.Now(), time.Now(), 0, 0)
	if err != nil {
		return nil, xerrors.Errorf("failed to FindAppUserByLoginID. err: %w", err)
	}

	guest, err := user.NewAppUser(nil, model, user.OrganizationID(organization.GetID()), "guest", "Guest", []string{}, map[string]string{})
	if err != nil {
		return nil, xerrors.Errorf("failed to FindAppUserByLoginID. err: %w", err)
	}

	return s.authTokenManager.CreateTokenSet(ctx, guest, organization)
}
