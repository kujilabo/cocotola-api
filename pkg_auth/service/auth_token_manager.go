package service

import (
	"context"

	userD "github.com/kujilabo/cocotola-api/pkg_user/domain"
)

type AuthTokenManager interface {
	CreateTokenSet(ctx context.Context, appUser userD.AppUserModel, organization userD.OrganizationModel) (*TokenSet, error)
	RefreshToken(ctx context.Context, tokenString string) (string, error)
}
