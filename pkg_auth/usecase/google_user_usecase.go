package usecase

import (
	"context"
	"errors"

	"golang.org/x/xerrors"
	"gorm.io/gorm"

	"github.com/kujilabo/cocotola-api/pkg_auth/service"
	"github.com/kujilabo/cocotola-api/pkg_lib/log"
	user "github.com/kujilabo/cocotola-api/pkg_user/domain"
	userS "github.com/kujilabo/cocotola-api/pkg_user/service"
)

type GoogleUserUsecase interface {
	RetrieveAccessToken(ctx context.Context, code string) (*service.GoogleAuthResponse, error)

	RetrieveUserInfo(ctx context.Context, GoogleAuthResponse *service.GoogleAuthResponse) (*service.GoogleUserInfo, error)

	RegisterAppUser(ctx context.Context, googleUserInfo *service.GoogleUserInfo, googleAuthResponse *service.GoogleAuthResponse, organizationName string) (*service.TokenSet, error)
}

type googleUserUsecase struct {
	db                      *gorm.DB
	googleAuthClient        service.GoogleAuthClient
	authTokenManager        service.AuthTokenManager
	registerAppUserCallback func(ctx context.Context, db *gorm.DB, organizationName string, appUser user.AppUserModel) error
}

func NewGoogleUserUsecase(db *gorm.DB, googleAuthClient service.GoogleAuthClient, authTokenManager service.AuthTokenManager, registerAppUserCallback func(ctx context.Context, db *gorm.DB, organizationName string, appUser user.AppUserModel) error) GoogleUserUsecase {
	return &googleUserUsecase{
		db:                      db,
		googleAuthClient:        googleAuthClient,
		authTokenManager:        authTokenManager,
		registerAppUserCallback: registerAppUserCallback,
	}
}

func (s *googleUserUsecase) RetrieveAccessToken(ctx context.Context, code string) (*service.GoogleAuthResponse, error) {
	return s.googleAuthClient.RetrieveAccessToken(ctx, code)
}

func (s *googleUserUsecase) RetrieveUserInfo(ctx context.Context, googleAuthResponse *service.GoogleAuthResponse) (*service.GoogleUserInfo, error) {
	return s.googleAuthClient.RetrieveUserInfo(ctx, googleAuthResponse)
}

func (s *googleUserUsecase) RegisterAppUser(ctx context.Context, googleUserInfo *service.GoogleUserInfo, googleAuthResponse *service.GoogleAuthResponse, organizationName string) (*service.TokenSet, error) {
	logger := log.FromContext(ctx)
	var tokenSet *service.TokenSet

	if err := s.db.Transaction(func(tx *gorm.DB) error {
		systemAdmin, err := userS.NewSystemAdminFromDB(ctx, tx)
		if err != nil {
			return err
		}

		systemOwner, err := systemAdmin.FindSystemOwnerByOrganizationName(ctx, organizationName)
		if err != nil {
			return xerrors.Errorf("failed to FindSystemOwnerByOrganizationName. err: %w", err)
		}

		organization, err := systemOwner.GetOrganization(ctx)
		if err != nil {
			return xerrors.Errorf("failed to FindOrganization. err: %w", err)
		}

		loginID := googleUserInfo.Email
		logger.Infof("googleuserIndo: %+v", googleUserInfo)

		appUser, err := systemOwner.FindAppUserByLoginID(ctx, loginID)
		if err == nil {
			logger.Infof("user already exists. student: %+v", appUser)
			tokenSetTmp, err := s.authTokenManager.CreateTokenSet(ctx, appUser, organization)
			if err != nil {
				return err
			}

			tokenSet = tokenSetTmp
			return nil
		}

		if !errors.Is(err, userS.ErrAppUserNotFound) {
			logger.Infof("Unsupported %v", err)
			return err
		}

		logger.Infof("Add student. %+v", appUser)
		parameter, err := userS.NewAppUserAddParameter(
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
			return xerrors.Errorf("invalid AppUserAddParameter. err: %w", err)
		}

		studentID, err := systemOwner.AddAppUser(ctx, parameter)
		if err != nil {
			return xerrors.Errorf("failed to AddStudent. err: %w", err)
		}

		student2, err := systemOwner.FindAppUserByID(ctx, studentID)
		if err != nil {
			return xerrors.Errorf("failed to FindStudentByID. err: %w", err)
		}

		if err := s.registerAppUserCallback(ctx, tx, organizationName, student2); err != nil {
			return xerrors.Errorf("failed to registerStudentCallback. err: %w", err)
		}

		tokenSetTmp, err := s.authTokenManager.CreateTokenSet(ctx, student2, organization)
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
