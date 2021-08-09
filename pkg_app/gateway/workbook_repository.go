package gateway

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/casbin/casbin/v2"
	"golang.org/x/xerrors"
	"gorm.io/gorm"

	"github.com/kujilabo/cocotola-api/pkg_app/domain"
	libD "github.com/kujilabo/cocotola-api/pkg_lib/domain"
	libG "github.com/kujilabo/cocotola-api/pkg_lib/gateway"
	"github.com/kujilabo/cocotola-api/pkg_lib/log"
	user "github.com/kujilabo/cocotola-api/pkg_user/domain"
	casbinquery "github.com/pecolynx/casbin-query"
)

type workbookEntity struct {
	ID             uint
	Version        int
	CreatedAt      time.Time
	UpdatedAt      time.Time
	CreatedBy      uint
	UpdatedBy      uint
	OrganizationID uint
	SpaceID        uint
	OwnerID        uint
	Name           string
	ProblemTypeID  uint `gorm:"column:problem_type_id"`
	QuestionText   string
}

func (e *workbookEntity) TableName() string {
	return "workbook"
}

func (e *workbookEntity) toWorkbook(rf domain.RepositoryFactory, pf domain.ProcessorFactory, operator user.AppUser, problemType string, privs user.Privileges) (domain.Workbook, error) {
	model, err := user.NewModel(e.ID, e.Version, e.CreatedAt, e.UpdatedAt, e.CreatedBy, e.UpdatedBy)
	if err != nil {
		return nil, err
	}
	workbook, err := domain.NewWorkbook(rf, pf, model, user.SpaceID(e.SpaceID), user.AppUserID(e.OwnerID), privs, e.Name, problemType, e.QuestionText)
	if err != nil {
		return nil, xerrors.Errorf("failed to NewWorkbook. entity: %+v, err: %w", e, err)
	}
	return workbook, nil
}

type workbookRepository struct {
	driverName   string
	db           *gorm.DB
	rf           domain.RepositoryFactory
	userRepo     user.RepositoryFactory
	pf           domain.ProcessorFactory
	problemTypes []domain.ProblemType
}

func NewWorkbookRepository(ctx context.Context, driverName string, rf domain.RepositoryFactory, userRepo user.RepositoryFactory, pf domain.ProcessorFactory, db *gorm.DB, problemTypes []domain.ProblemType) domain.WorkbookRepository {
	return &workbookRepository{
		driverName:   driverName,
		db:           db,
		rf:           rf,
		userRepo:     userRepo,
		pf:           pf,
		problemTypes: problemTypes,
	}
}

func (r *workbookRepository) toProblemType(problemTypeID uint) string {
	for _, m := range r.problemTypes {
		if m.GetID() == problemTypeID {
			return m.GetName()
		}
	}
	return ""
}

func (r *workbookRepository) toProblemTypeID(problemType string) uint {
	for _, m := range r.problemTypes {
		if m.GetName() == problemType {
			return m.GetID()
		}
	}
	return 0
}

func (r *workbookRepository) FindPersonalWorkbooks(ctx context.Context, operator domain.Student, param domain.WorkbookSearchCondition) (*domain.WorkbookSearchResult, error) {
	logger := log.FromContext(ctx)
	logger.Infof("workbookRepository.FindWorkbooks %v", operator)
	if param == nil {
		return nil, libD.ErrInvalidArgument
	}

	limit := param.GetPageSize()
	offset := (param.GetPageNo() - 1) * param.GetPageSize()
	workbooks := []workbookEntity{}

	objectColumnName := "name"
	subQuery, err := casbinquery.QueryObject(r.db, r.driverName, objectColumnName, "user_"+strconv.Itoa(int(operator.GetID())), "read")
	if err != nil {
		return nil, err
	}

	if result := r.db.Model(&workbookEntity{}).
		Joins("inner join (?) AS t3 ON `workbook`.`id`= t3."+objectColumnName, subQuery).
		Order("`workbook`.`name`").Limit(limit).Offset(offset).
		Scan(&workbooks); result.Error != nil {
		return nil, result.Error
	}

	results := make([]domain.Workbook, len(workbooks))
	priv := user.NewPrivileges([]user.RBACAction{domain.PrivilegeRead})
	for i, e := range workbooks {
		w, err := e.toWorkbook(r.rf, r.pf, operator, r.toProblemType(e.ProblemTypeID), priv)
		if err != nil {
			return nil, xerrors.Errorf("failed to toWorkbook. err: %w", err)
		}
		results[i] = w
	}

	var count int64
	rows, err := r.db.Raw("select count(*) from workbook inner join (?) AS t3 ON `workbook`.`id`= t3."+objectColumnName, subQuery).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var c int64
		if err := rows.Scan(&c); err != nil {
			return nil, err
		}
		count += c
	}

	return &domain.WorkbookSearchResult{
		TotalCount: count,
		Results:    results,
	}, nil
}

func (r *workbookRepository) getAllWorkbookRoles(workbookID domain.WorkbookID) []user.RBACRole {
	return []user.RBACRole{domain.NewWorkbookWriter(workbookID), domain.NewWorkbookReader(workbookID)}
}

func (r *workbookRepository) getAllWorkbookPrivileges() []user.RBACAction {
	return []user.RBACAction{domain.PrivilegeRead, domain.PrivilegeUpdate, domain.PrivilegeRemove}
}

func (r *workbookRepository) checkPrivileges(e *casbin.Enforcer, userObject user.RBACUser, workbookObject user.RBACObject, privs []user.RBACAction) (user.Privileges, error) {
	actions := make([]user.RBACAction, 0)
	for _, priv := range privs {
		ok, err := e.Enforce(string(userObject), string(workbookObject), string(priv))
		if err != nil {
			return nil, err
		}
		if ok {
			actions = append(actions, priv)
		}
	}
	return user.NewPrivileges(actions), nil
}

func (r *workbookRepository) FindWorkbookByID(ctx context.Context, operator domain.Student, workbookID domain.WorkbookID) (domain.Workbook, error) {
	workbook := workbookEntity{}
	if result := r.db.Where("id = ?", uint(workbookID)).First(&workbook); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, domain.ErrWorkbookNotFound
		}
		return nil, result.Error
	}

	rbacRepo := r.userRepo.NewRBACRepository()
	workbookRoles := r.getAllWorkbookRoles(domain.WorkbookID(workbook.ID))
	userObject := user.NewUserObject(user.AppUserID(operator.GetID()))
	e, err := rbacRepo.NewEnforcerWithRolesAndUsers(workbookRoles, []user.RBACUser{userObject})
	if err != nil {
		return nil, xerrors.Errorf("failed to NewEnforcerWithRolesAndUsers. err: %w", err)
	}
	workbookObject := domain.NewWorkbookObject(domain.WorkbookID(workbook.ID))
	privs := r.getAllWorkbookPrivileges() // TODO

	priv, err := r.checkPrivileges(e, userObject, workbookObject, privs)
	if err != nil {
		return nil, xerrors.Errorf("failed to checkPrivileges. err: %w", err)
	}

	// defaultSpace, err := operator.GetDefaultSpace(ctx)
	// if err != nil {
	// 	return nil, err
	// }

	logger := log.FromContext(ctx)
	logger.Infof("ownerId: %d, operatorId: %d", workbook.OwnerID, operator.GetID())

	return workbook.toWorkbook(r.rf, r.pf, operator, r.toProblemType(workbook.ProblemTypeID), priv)
}

func (r *workbookRepository) FindWorkbookByName(ctx context.Context, operator domain.Student, name string) (domain.Workbook, error) {
	workbook := workbookEntity{}
	if result := r.db.Where("name = ?", name).First(&workbook); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, domain.ErrWorkbookNotFound
		}
		return nil, result.Error
	}

	rbacRepo := r.userRepo.NewRBACRepository()
	workbookRoles := r.getAllWorkbookRoles(domain.WorkbookID(workbook.ID))
	userObject := user.NewUserObject(user.AppUserID(operator.GetID()))
	e, err := rbacRepo.NewEnforcerWithRolesAndUsers(workbookRoles, []user.RBACUser{userObject})
	if err != nil {
		return nil, xerrors.Errorf("failed to NewEnforcerWithRolesAndUsers. err: %w", err)
	}
	workbookObject := domain.NewWorkbookObject(domain.WorkbookID(workbook.ID))
	privs := r.getAllWorkbookPrivileges() // TODO

	priv, err := r.checkPrivileges(e, userObject, workbookObject, privs)
	if err != nil {
		return nil, xerrors.Errorf("failed to checkPrivileges. err: %w", err)
	}

	logger := log.FromContext(ctx)
	logger.Infof("ownerId: %d, operatorId: %d", workbook.OwnerID, operator.GetID())

	return workbook.toWorkbook(r.rf, r.pf, operator, r.toProblemType(workbook.ProblemTypeID), priv)
}

// func (r *workbookRepository) FindWorkbookByName(ctx context.Context, operator domain.AbstractStudent, name string) (domain.Workbook, error) {
// 	workbook := workbookEntity{}
// 	if result := r.db.Where("name = ?", name).First(&workbook); result.Error != nil {
// 		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
// 			return nil, domain.NewWorkbookNotFoundError(0)
// 		}
// 		return nil, result.Error
// 	}
// 	priv := user.NewPrivileges(user.PrivRead, user.PrivUpdate, domain.PrivDelete)
// 	return workbook.toWorkbook(r.rf, r.gh, operator, priv)
// }

// func (r *workbookRepository) scan(rows *sql.Rows) *workbookEntity {
// 	var id uint
// 	var version int
// 	var createdAt time.Time
// 	var updatedAt time.Time
// 	var createdBy uint
// 	var updatedBy uint
// 	var organizationID uint
// 	var spaceID uint
// 	var ownerID uint
// 	var problemType string
// 	var name string
// 	var questionText string
// 	rows.Scan(&id, &version, &createdAt, &updatedAt, &createdBy, &updatedBy, &organizationID, &spaceID, &ownerID, &problemType, &name, &questionText)

// 	return &workbookEntity{
// 		ID:             id,
// 		Version:        version,
// 		CreatedAt:      createdAt,
// 		UpdatedAt:      updatedAt,
// 		CreatedBy:      createdBy,
// 		UpdatedBy:      updatedBy,
// 		OrganizationID: organizationID,
// 		SpaceID:        spaceID,
// 		OwnerID:        ownerID,
// 		ProblemTypeID:  r.toProblemTypeID(problemType),
// 		Name:           name,
// 		QuestionText:   questionText,
// 	}

// }

// func (r *workbookRepository) FindWorkbooksFromDefaultSpace(ctx context.Context, operator domain.Student, spaceID uint, param *domain.WorkbookSearchCondition) (*domain.WorkbookSearchResult, error) {
// 	limit := param.PageSize
// 	offset := (param.PageNo - 1) * param.PageSize
// 	var workbooks []workbookEntity
// 	db := r.db.Where("organization_id = ?", operator.OrganizationID())
// 	db = db.Where("space_id = ?", spaceID)
// 	if result := db.Limit(limit).Offset(offset).Find(&workbooks); result.Error != nil {
// 		return nil, result.Error
// 	}
// 	var count int64
// 	if result := db.Model(workbookEntity{}).Count(&count); result.Error != nil {
// 		return nil, result.Error
// 	}
// 	priv := user.NewPrivileges(user.PrivRead)
// 	results := make([]domain.Workbook, len(workbooks))
// 	for i, e := range workbooks {
// 		w, err := e.toWorkbook(r.rf, operator, priv)
// 		if err != nil {
// 			return nil, err
// 		}
// 		results[i] = w
// 	}
// 	return &domain.WorkbookSearchResult{
// 		TotalCount: count,
// 		Results:    results,
// 	}, nil
// }

func (r *workbookRepository) AddWorkbook(ctx context.Context, operator domain.Student, spaceID user.SpaceID, param domain.WorkbookAddParameter) (domain.WorkbookID, error) {
	problemTypeID := r.toProblemTypeID(param.GetProblemType())
	if problemTypeID == 0 {
		return 0, xerrors.Errorf("unsupported problemType. problemType: %s", param.GetProblemType())
	}
	workbook := workbookEntity{
		Version:        1,
		CreatedBy:      operator.GetID(),
		UpdatedBy:      operator.GetID(),
		OrganizationID: uint(operator.GetOrganizationID()),
		SpaceID:        uint(spaceID),
		OwnerID:        operator.GetID(),
		ProblemTypeID:  problemTypeID,
		Name:           param.GetName(),
		QuestionText:   param.GetQuestionText(),
	}
	if result := r.db.Create(&workbook); result.Error != nil {
		return 0, libG.ConvertDuplicatedError(result.Error, domain.ErrWorkbookAlreadyExists)
	}

	workbookID := domain.WorkbookID(workbook.ID)

	rbacRepo := r.userRepo.NewRBACRepository()
	userObject := user.NewUserObject(user.AppUserID(operator.GetID()))
	workbookObject := domain.NewWorkbookObject(workbookID)
	workbookWriter := domain.NewWorkbookWriter(workbookID)

	// wrokbookWriter role can read, update, remove
	if err := rbacRepo.AddNamedPolicy(workbookWriter, workbookObject, domain.PrivilegeRead); err != nil {
		return 0, err
	}
	if err := rbacRepo.AddNamedPolicy(workbookWriter, workbookObject, domain.PrivilegeUpdate); err != nil {
		return 0, err
	}
	if err := rbacRepo.AddNamedPolicy(workbookWriter, workbookObject, domain.PrivilegeRemove); err != nil {
		return 0, err
	}

	// user is assigned workbookWriter role
	if err := rbacRepo.AddNamedGroupingPolicy(userObject, workbookWriter); err != nil {
		return 0, err
	}

	// rbacRepo.NewEnforcerWithRolesAndUsers([]user.RBACRole{workbookWriter}, []user.RBACUser{userObject})

	return workbookID, nil
}

func (r *workbookRepository) RemoveWorkbook(ctx context.Context, operator domain.Student, id domain.WorkbookID, version int) error {
	workbook := workbookEntity{}
	if result := r.db.Where("organization_id = ? and id = ? and version = ?", operator.GetOrganizationID(), id, version).Delete(&workbook); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return domain.ErrWorkbookNotFound
		}

		return result.Error
	}

	return nil
}

func (r *workbookRepository) UpdateWorkbook(ctx context.Context, operator domain.Student, id domain.WorkbookID, version int, param domain.WorkbookUpdateParameter) error {
	workbook := workbookEntity{
		Name:         param.GetName(),
		QuestionText: param.GetQuestionText(),
	}
	if result := r.db.Model(&workbookEntity{}).
		Where("organization_id = ? and id = ? and version = ?",
			uint(operator.GetOrganizationID()), uint(id), version).
		Updates(&workbook); result.Error != nil {
		return libG.ConvertDuplicatedError(result.Error, domain.ErrWorkbookAlreadyExists)
	}

	return nil
}

// func (r *workbookRepository) ChangeSpace(ctx context.Context, operator domain.AbstractStudent, id uint, spaceID uint) error {
// 	result := r.db.Model(&workbookEntity{}).Where(workbookEntity{
// 		OrganizationID: operator.OrganizationID(),
// 		ID:             id,
// 	}).Update(workbookEntity{
// 		SpaceID: spaceID,
// 	})
// 	if result.Error != nil {
// 		return result.Error
// 	}
// 	if result.RowsAffected == 0 {
// 		return domain.NewWorkbookNotFoundError(id)
// 	}

// 	return nil
// }
