package application

import (
	"context"
	"fmt"

	"github.com/kujilabo/cocotola-api/pkg_app/domain"
	user "github.com/kujilabo/cocotola-api/pkg_user/domain"
)

func findStudent(ctx context.Context, repo domain.RepositoryFactory, userRepo user.RepositoryFactory, organizationID user.OrganizationID, operatorID user.AppUserID) (domain.Student, error) {
	systemAdmin := user.SystemAdminInstance()
	systemOwner, err := systemAdmin.FindSystemOwnerByOrganizationID(ctx, organizationID)
	if err != nil {
		return nil, fmt.Errorf("failed to FindSystemOwnerByOrganizationID. err: %w", err)
	}

	appUser, err := systemOwner.FindAppUserByID(ctx, operatorID)
	if err != nil {
		return nil, err
	}

	return domain.NewStudent(repo, userRepo, appUser)
}
