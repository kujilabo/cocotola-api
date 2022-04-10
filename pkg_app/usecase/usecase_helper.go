package usecase

import (
	"context"

	"github.com/kujilabo/cocotola-api/pkg_app/domain"
	"github.com/kujilabo/cocotola-api/pkg_app/service"
	user "github.com/kujilabo/cocotola-api/pkg_user/domain"
	userS "github.com/kujilabo/cocotola-api/pkg_user/service"
	"golang.org/x/xerrors"
)

func FindStudent(ctx context.Context, pf service.ProcessorFactory, rf service.RepositoryFactory, userRf userS.RepositoryFactory, organizationID user.OrganizationID, operatorID user.AppUserID) (service.Student, error) {
	systemAdmin := userS.NewSystemAdmin(userRf)
	systemOwner, err := systemAdmin.FindSystemOwnerByOrganizationID(ctx, organizationID)
	if err != nil {
		return nil, xerrors.Errorf("failed to FindSystemOwnerByOrganizationID. err: %w", err)
	}

	appUser, err := systemOwner.FindAppUserByID(ctx, operatorID)
	if err != nil {
		return nil, err
	}

	studentModel, err := domain.NewStudentModel(appUser)
	if err != nil {
		return nil, err
	}

	return service.NewStudent(pf, rf, userRf, studentModel)
}
