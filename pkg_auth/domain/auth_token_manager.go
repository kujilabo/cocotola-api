package domain

import (
	"context"

	user "github.com/kujilabo/cocotola-api/pkg_user/domain"
)

type AuthTokenManager interface {
	CreateTokenSet(ctx context.Context, appUser user.AppUser, organization user.Organization) (*TokenSet, error)
	RefreshToken(ctx context.Context, tokenString string) (string, error)
}
