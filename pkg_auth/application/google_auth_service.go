package application

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"

	"github.com/kujilabo/cocotola-api/pkg_auth/domain"
	"github.com/kujilabo/cocotola-api/pkg_lib/log"
	user "github.com/kujilabo/cocotola-api/pkg_user/domain"
)

type GoogleAuthService interface {
	RetrieveAccessToken(ctx context.Context, code string) (*domain.GoogleAuthResponse, error)
	RetrieveUserInfo(ctx context.Context, GoogleAuthResponse *domain.GoogleAuthResponse) (*domain.GoogleUserInfo, error)
	RegisterStudent(ctx context.Context, googleUserInfo *domain.GoogleUserInfo, googleAuthResponse *domain.GoogleAuthResponse, organizationName string) (*domain.TokenSet, error)
}

type googleAuthService struct {
	userRepo                func(db *gorm.DB) (user.RepositoryFactory, error)
	googleAuthClient        domain.GoogleAuthClient
	authTokenManager        domain.AuthTokenManager
	registerAppUserCallback func(ctx context.Context, organizationName string, appUser user.AppUser) error
}

func NewGoogleAuthService(userRepo func(db *gorm.DB) (user.RepositoryFactory, error), googleAuthClient domain.GoogleAuthClient, authTokenManager domain.AuthTokenManager, registerAppUserCallback func(ctx context.Context, organizationName string, appUser user.AppUser) error) GoogleAuthService {
	return &googleAuthService{
		userRepo:                userRepo,
		googleAuthClient:        googleAuthClient,
		authTokenManager:        authTokenManager,
		registerAppUserCallback: registerAppUserCallback,
	}
}

func (s *googleAuthService) RetrieveAccessToken(ctx context.Context, code string) (*domain.GoogleAuthResponse, error) {
	return s.googleAuthClient.RetrieveAccessToken(ctx, code)
}

func (s *googleAuthService) RetrieveUserInfo(ctx context.Context, googleAuthResponse *domain.GoogleAuthResponse) (*domain.GoogleUserInfo, error) {
	return s.googleAuthClient.RetrieveUserInfo(ctx, googleAuthResponse)
}

func (s *googleAuthService) RegisterStudent(ctx context.Context, googleUserInfo *domain.GoogleUserInfo, googleAuthResponse *domain.GoogleAuthResponse, organizationName string) (*domain.TokenSet, error) {
	logger := log.FromContext(ctx)

	systemAdmin := user.SystemAdminInstance()
	systemOwner, err := systemAdmin.FindSystemOwnerByOrganizationName(ctx, organizationName)
	if err != nil {
		return nil, fmt.Errorf("failed to FindSystemOwnerByOrganizationName. err: %w", err)
	}

	organization, err := systemOwner.GetOrganization(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to FindOrganization. err: %w", err)
	}

	loginID := googleUserInfo.Email
	logger.Infof("googleuserIndo: %+v", googleUserInfo)

	appUser, err := systemOwner.FindAppUserByLoginID(ctx, loginID)
	if err == nil {
		logger.Infof("user already exists. student: %+v", appUser)
		return s.authTokenManager.CreateTokenSet(ctx, appUser, organization)
	}

	if !errors.Is(err, user.ErrAppUserNotFound) {
		logger.Infof("Unsupported %v", err)
		return nil, err
	}

	logger.Infof("Add student. %+v", appUser)
	parameter, err := user.NewAppUserAddParameter(
		googleUserInfo.Email,
		googleUserInfo.Name,
		[]string{""},
		map[string]string{
			"password":             "----",
			"provider":             "google",
			"providerId":           googleUserInfo.Email,
			"providerAccessToken":  googleAuthResponse.AccessToken,
			"providerRefreshToken": googleAuthResponse.RefreshToken,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("invalid AppUserAddParameter. err: %w", err)
	}

	studentID, err := systemOwner.AddAppUser(ctx, parameter)
	if err != nil {
		return nil, fmt.Errorf("failed to AddStudent. err: %w", err)
	}

	student2, err := systemOwner.FindAppUserByID(ctx, studentID)
	if err != nil {
		return nil, fmt.Errorf("failed to FindStudentByID. err: %w", err)
	}

	if err := s.registerAppUserCallback(ctx, organizationName, student2); err != nil {
		return nil, fmt.Errorf("failed to registerStudentCallback. err: %w", err)
	}

	return s.authTokenManager.CreateTokenSet(ctx, student2, organization)
}
