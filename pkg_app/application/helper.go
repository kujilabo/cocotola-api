package application

import (
	"context"

	"github.com/kujilabo/cocotola-api/pkg_app/domain"
	user "github.com/kujilabo/cocotola-api/pkg_user/domain"
	"golang.org/x/xerrors"
)

func findStudent(ctx context.Context, pf domain.ProcessorFactory, repo domain.RepositoryFactory, userRepo user.RepositoryFactory, organizationID user.OrganizationID, operatorID user.AppUserID) (domain.Student, error) {
	systemAdmin := user.SystemAdminInstance()
	systemOwner, err := systemAdmin.FindSystemOwnerByOrganizationID(ctx, organizationID)
	if err != nil {
		return nil, xerrors.Errorf("failed to FindSystemOwnerByOrganizationID. err: %w", err)
	}

	appUser, err := systemOwner.FindAppUserByID(ctx, operatorID)
	if err != nil {
		return nil, err
	}

	return domain.NewStudent(pf, repo, userRepo, appUser)
}
