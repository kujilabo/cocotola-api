package gateway

import (
	"context"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"

	"github.com/kujilabo/cocotola-api/pkg_auth/domain"
	"github.com/kujilabo/cocotola-api/pkg_lib/log"
	user "github.com/kujilabo/cocotola-api/pkg_user/domain"
)

type AppUserClaims struct {
	LoginID          string `json:"loginId"`
	AppUserID        uint   `json:"appUserId"`
	Username         string `json:"username"`
	OrganizationID   uint   `json:"organizationId"`
	OrganizationName string `json:"organizationName"`
	Role             string `json:"role"`
	TokenType        string `json:"tokenType"`
	jwt.StandardClaims
}

type authTokenManager struct {
	signingKey     []byte
	signingMethod  jwt.SigningMethod
	tokenTimeout   time.Duration
	refreshTimeout time.Duration
}

func NewAuthTokenManager(signingKey []byte, signingMethod jwt.SigningMethod, tokenTimeout, refreshTimeout time.Duration) domain.AuthTokenManager {
	return &authTokenManager{
		signingKey:     signingKey,
		signingMethod:  signingMethod,
		tokenTimeout:   tokenTimeout,
		refreshTimeout: refreshTimeout,
	}
}

func (m *authTokenManager) CreateTokenSet(ctx context.Context, appUser user.AppUser, organization user.Organization) (*domain.TokenSet, error) {
	accessToken, err := m.createJWT(ctx, appUser, organization, m.tokenTimeout, "access")
	if err != nil {
		return nil, err
	}

	refreshToken, err := m.createJWT(ctx, appUser, organization, m.refreshTimeout, "refresh")
	if err != nil {
		return nil, err
	}

	return &domain.TokenSet{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (m *authTokenManager) createJWT(ctx context.Context, appUser user.AppUser, organization user.Organization, duration time.Duration, tokenType string) (string, error) {
	logger := log.FromContext(ctx)
	now := time.Now()
	claims := AppUserClaims{
		LoginID:          appUser.GetLoginID(),
		AppUserID:        appUser.GetID(),
		Username:         appUser.GetUsername(),
		OrganizationID:   organization.GetID(),
		OrganizationName: organization.GetName(),
		Role:             appUser.GetRoles()[0],
		TokenType:        tokenType,
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  now.Unix(),
			ExpiresAt: now.Add(duration).Unix(),
		},
	}

	logger.Debugf("claims: %+v", claims)

	token := jwt.NewWithClaims(m.signingMethod, claims)
	return token.SignedString(m.signingKey)
}

func (m *authTokenManager) RefreshToken(ctx context.Context, tokenString string) (string, error) {
	logger := log.FromContext(ctx)
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		return m.signingKey, nil
	}

	currentToken, err := jwt.ParseWithClaims(tokenString, &AppUserClaims{}, keyFunc)
	if err != nil {
		logger.WithError(err).Infof("%v", err)
		return "", domain.NewUnauthorizedError(fmt.Sprintf("failed to ParseWithClaims. err: %v", err))
	}

	currentClaims, ok := currentToken.Claims.(*AppUserClaims)
	if !ok || !currentToken.Valid {
		return "", domain.NewUnauthorizedError("Invalid token. err: %v")
	}

	if currentClaims.TokenType != "refresh" {
		return "", domain.NewUnauthorizedError("Invalid token type")
	}

	now := time.Now()
	tmpID := uint(1)
	userModel, err := user.NewModel(currentClaims.AppUserID, 1, now, now, tmpID, tmpID)
	if err != nil {
		return "", err
	}

	appUser, err := user.NewAppUser(nil, userModel, user.OrganizationID(currentClaims.OrganizationID), currentClaims.LoginID, currentClaims.Username, []string{currentClaims.Role}, map[string]string{})
	if err != nil {
		return "", err
	}

	orgModel, err := user.NewModel(currentClaims.OrganizationID, 1, now, now, tmpID, tmpID)
	if err != nil {
		return "", err
	}

	organization, err := user.NewOrganization(orgModel, currentClaims.OrganizationName)
	if err != nil {
		return "", err
	}

	accessToken, err := m.createJWT(ctx, appUser, organization, m.tokenTimeout, "access")
	if err != nil {
		return "", err
	}

	return accessToken, nil
}
