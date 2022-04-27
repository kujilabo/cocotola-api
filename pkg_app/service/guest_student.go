//go:generate mockery --output mock --name GuestStudent
package service

import (
	"context"
	"errors"

	"github.com/kujilabo/cocotola-api/pkg_app/domain"
	libD "github.com/kujilabo/cocotola-api/pkg_lib/domain"
	userS "github.com/kujilabo/cocotola-api/pkg_user/service"
)

type GuestStudent interface {
	domain.StudentModel

	GetDefaultSpace(ctx context.Context) (userS.Space, error)

	FindWorkbooksFromPublicSpace(ctx context.Context, condition WorkbookSearchCondition) (WorkbookSearchResult, error)
}

type guestStudent struct {
	domain.StudentModel
	rf     RepositoryFactory
	pf     ProcessorFactory
	userRf userS.RepositoryFactory
}

func NewGuestStudent(pf ProcessorFactory, rf RepositoryFactory, userRf userS.RepositoryFactory, studentModel domain.StudentModel) (GuestStudent, error) {
	m := &guestStudent{
		StudentModel: studentModel,
		pf:           pf,
		rf:           rf,
		userRf:       userRf,
	}

	return m, libD.Validator.Struct(m)
}

func (s *guestStudent) GetDefaultSpace(ctx context.Context) (userS.Space, error) {
	return s.userRf.NewSpaceRepository().FindDefaultSpace(ctx, s)
}

func (s *guestStudent) FindWorkbooksFromPublicSpace(ctx context.Context, condition WorkbookSearchCondition) (WorkbookSearchResult, error) {
	return nil, errors.New("aaa")
	// space, err := s.GetPersonalSpace(ctx)
	// if err != nil {
	// 	return nil, xerrors.Errorf("failed to GetPersonalSpace. err: %w", err)
	// }

	// // specify space
	// newCondition, err := NewWorkbookSearchCondition(condition.GetPageNo(), condition.GetPageSize(), []userD.SpaceID{userD.SpaceID(space.GetID())})
	// if err != nil {
	// 	return nil, xerrors.Errorf("failed to NewWorkbookSearchCondition. err: %w", err)
	// }

	// workbookRepo, err := s.rf.NewWorkbookRepository(ctx)
	// if err != nil {
	// 	return nil, xerrors.Errorf("failed to NewWorkbookRepository. err: %w", err)
	// }

	// return workbookRepo.FindPersonalWorkbooks(ctx, s, newCondition)
}
