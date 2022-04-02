package service

import (
	"context"

	user "github.com/kujilabo/cocotola-api/pkg_user/domain"
)

type AuthTokenManager interface {
	CreateTokenSet(ctx context.Context, appUser user.AppUserModel, organization user.OrganizationModel) (*TokenSet, error)
	RefreshToken(ctx context.Context, tokenString string) (string, error)
}
