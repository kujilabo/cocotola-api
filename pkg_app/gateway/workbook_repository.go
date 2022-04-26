package gateway

import (
	"context"
	"encoding/json"
	"errors"
	"math"
	"strconv"
	"time"

	"github.com/casbin/casbin/v2"
	"golang.org/x/xerrors"
	"gorm.io/gorm"

	"github.com/kujilabo/cocotola-api/pkg_app/domain"
	"github.com/kujilabo/cocotola-api/pkg_app/gateway/casbinquery"
	"github.com/kujilabo/cocotola-api/pkg_app/service"
	libD "github.com/kujilabo/cocotola-api/pkg_lib/domain"
	libG "github.com/kujilabo/cocotola-api/pkg_lib/gateway"
	"github.com/kujilabo/cocotola-api/pkg_lib/log"
	user "github.com/kujilabo/cocotola-api/pkg_user/domain"
	userS "github.com/kujilabo/cocotola-api/pkg_user/service"
	// casbinquery "github.com/pecolynx/casbin-query"
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
	Lang           string
	ProblemTypeID  uint `gorm:"column:problem_type_id"`
	QuestionText   string
	Properties     string
}

func (e *workbookEntity) TableName() string {
	return "workbook"
}

func jsonToStringMap(s string) (map[string]string, error) {
	var m map[string]string
	err := json.Unmarshal([]byte(s), &m)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func stringMapToJSON(m map[string]string) (string, error) {
	b, err := json.Marshal(m)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func (e *workbookEntity) toWorkbookModel(rf service.RepositoryFactory, pf service.ProcessorFactory, operator user.AppUserModel, problemType string, privs user.Privileges) (domain.WorkbookModel, error) {
	model, err := user.NewModel(e.ID, e.Version, e.CreatedAt, e.UpdatedAt, e.CreatedBy, e.UpdatedBy)
	if err != nil {
		return nil, err
	}

	properties, err := jsonToStringMap(e.Properties)
	if err != nil {
		return nil, xerrors.Errorf("failed to jsonToStringMap. err: %w ", err)
	}

	lang, err := domain.NewLang2(e.Lang)
	if err != nil {
		return nil, xerrors.Errorf("invalid lang. lang: %s, err: %w", e.Lang, err)
	}

	workbook, err := domain.NewWorkbookModel(model, user.SpaceID(e.SpaceID), user.AppUserID(e.OwnerID), privs, e.Name, lang, problemType, e.QuestionText, properties)
	if err != nil {
		return nil, xerrors.Errorf("failed to NewWorkbook. entity: %+v, err: %w", e, err)
	}
	return workbook, nil
}

type workbookRepository struct {
	driverName   string
	db           *gorm.DB
	rf           service.RepositoryFactory
	userRf       userS.RepositoryFactory
	pf           service.ProcessorFactory
	problemTypes []domain.ProblemType
}

func NewWorkbookRepository(ctx context.Context, driverName string, rf service.RepositoryFactory, userRf userS.RepositoryFactory, pf service.ProcessorFactory, db *gorm.DB, problemTypes []domain.ProblemType) service.WorkbookRepository {
	return &workbookRepository{
		driverName:   driverName,
		db:           db,
		rf:           rf,
		userRf:       userRf,
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

func (r *workbookRepository) FindPersonalWorkbooks(ctx context.Context, operator domain.StudentModel, param service.WorkbookSearchCondition) (service.WorkbookSearchResult, error) {
	ctx, span := tracer.Start(ctx, "workbookRepository.FindWorkbooks")
	defer span.End()

	logger := log.FromContext(ctx)
	logger.Debugf("workbookRepository.FindWorkbooks. OperatorID: %d", operator.GetID())

	if param == nil {
		return nil, libD.ErrInvalidArgument
	}

	limit := param.GetPageSize()
	offset := (param.GetPageNo() - 1) * param.GetPageSize()
	workbooks := []workbookEntity{}

	objectColumnName := "name"
	subQuery, err := casbinquery.QueryObject(r.db, r.driverName, domain.WorkbookObjectPrefix, objectColumnName, "user_"+strconv.Itoa(int(operator.GetID())), "read")
	if err != nil {
		return nil, err
	}

	if result := r.db.Model(&workbookEntity{}).
		Joins("inner join (?) AS t3 ON `workbook`.`id`= t3."+objectColumnName, subQuery).
		Order("`workbook`.`name`").Limit(limit).Offset(offset).
		Scan(&workbooks); result.Error != nil {
		return nil, result.Error
	}

	results := make([]domain.WorkbookModel, len(workbooks))
	priv := user.NewPrivileges([]user.RBACAction{domain.PrivilegeRead})
	for i, e := range workbooks {
		w, err := e.toWorkbookModel(r.rf, r.pf, operator, r.toProblemType(e.ProblemTypeID), priv)
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

	if count > math.MaxInt32 {
		return nil, errors.New("overflow")
	}

	return service.NewWorkbookSearchResult(int(count), results)
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

// func (r *workbookRepository) canReadWorkbook(operator user.AppUser, workbookID domain.WorkbookID) error {
// 	objectColumnName := "name"
// 	object := domain.WorkbookObjectPrefix + strconv.Itoa(int(uint(workbookID)))
// 	subject := "user_" + strconv.Itoa(int(operator.GetID()))
// 	casbinQuery, err := casbinquery.FindObject(r.db, r.driverName, object, objectColumnName, subject, "read")
// 	if err != nil {
// 		return err
// 	}
// 	var name string
// 	if result := casbinQuery.First(&name); result.Error != nil {
// 		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
// 			return domain.ErrWorkbookPermissionDenied
// 		}
// 		return result.Error
// 	}
// 	return nil
// }

func (r *workbookRepository) FindWorkbookByID(ctx context.Context, operator domain.StudentModel, workbookID domain.WorkbookID) (service.Workbook, error) {
	ctx, span := tracer.Start(ctx, "workbookRepository.FindWorkbookByID")
	defer span.End()

	workbookEntity := workbookEntity{}
	if result := r.db.
		Where("organization_id = ?", uint(operator.GetOrganizationID())).
		Where("id = ?", uint(workbookID)).
		First(&workbookEntity); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, service.ErrWorkbookNotFound
		}
		return nil, result.Error
	}

	priv, err := r.getPrivileges(ctx, operator, domain.WorkbookID(workbookEntity.ID))
	if err != nil {
		return nil, xerrors.Errorf("failed to checkPrivileges. err: %w", err)
	}
	if !priv.HasPrivilege(domain.PrivilegeRead) {
		return nil, service.ErrWorkbookPermissionDenied
	}

	logger := log.FromContext(ctx)
	logger.Infof("ownerId: %d, operatorId: %d", workbookEntity.OwnerID, operator.GetID())

	workbookModel, err := workbookEntity.toWorkbookModel(r.rf, r.pf, operator, r.toProblemType(workbookEntity.ProblemTypeID), priv)
	if err != nil {
		return nil, err
	}
	return service.NewWorkbook(r.rf, r.pf, workbookModel)
}

func (r *workbookRepository) FindWorkbookByName(ctx context.Context, operator user.AppUserModel, spaceID user.SpaceID, name string) (service.Workbook, error) {
	ctx, span := tracer.Start(ctx, "workbookRepository.FindWorkbookByName")
	defer span.End()

	workbookEntity := workbookEntity{}
	if result := r.db.
		Where("organization_id = ?", uint(operator.GetOrganizationID())).
		Where("space_id = ?", uint(spaceID)).
		Where("name = ?", name).
		First(&workbookEntity); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, service.ErrWorkbookNotFound
		}
		return nil, result.Error
	}

	var priv user.Privileges
	if spaceID == service.GetSystemSpaceID() {
		priv = user.NewPrivileges([]user.RBACAction{domain.PrivilegeRead})
	} else {
		privTmp, err := r.getPrivileges(ctx, operator, domain.WorkbookID(workbookEntity.ID))
		if err != nil {
			return nil, xerrors.Errorf("failed to checkPrivileges. err: %w", err)
		}
		if !privTmp.HasPrivilege(domain.PrivilegeRead) {
			return nil, service.ErrWorkbookPermissionDenied
		}
		priv = privTmp
	}

	logger := log.FromContext(ctx)
	logger.Infof("ownerId: %d, operatorId: %d", workbookEntity.OwnerID, operator.GetID())

	workbookModel, err := workbookEntity.toWorkbookModel(r.rf, r.pf, operator, r.toProblemType(workbookEntity.ProblemTypeID), priv)
	if err != nil {
		return nil, err
	}
	return service.NewWorkbook(r.rf, r.pf, workbookModel)
}

func (r *workbookRepository) getPrivileges(ctx context.Context, operator user.AppUserModel, workbookID domain.WorkbookID) (user.Privileges, error) {
	rbacRepo := r.userRf.NewRBACRepository()
	workbookRoles := r.getAllWorkbookRoles(workbookID)
	userObject := user.NewUserObject(user.AppUserID(operator.GetID()))
	e, err := rbacRepo.NewEnforcerWithRolesAndUsers(workbookRoles, []user.RBACUser{userObject})
	if err != nil {
		return nil, xerrors.Errorf("failed to NewEnforcerWithRolesAndUsers. err: %w", err)
	}
	workbookObject := domain.NewWorkbookObject(workbookID)
	privs := r.getAllWorkbookPrivileges()
	return r.checkPrivileges(e, userObject, workbookObject, privs)
}

func (r *workbookRepository) AddWorkbook(ctx context.Context, operator user.AppUserModel, spaceID user.SpaceID, param service.WorkbookAddParameter) (domain.WorkbookID, error) {
	_, span := tracer.Start(ctx, "workbookRepository.AddWorkbook")
	defer span.End()

	problemTypeID := r.toProblemTypeID(param.GetProblemType())
	if problemTypeID == 0 {
		return 0, xerrors.Errorf("unsupported problemType. problemType: %s", param.GetProblemType())
	}
	propertiesJSON, err := stringMapToJSON(param.GetProperties())
	if err != nil {
		return 0, err
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
		Properties:     propertiesJSON,
	}
	if result := r.db.Create(&workbook); result.Error != nil {
		return 0, libG.ConvertDuplicatedError(result.Error, service.ErrWorkbookAlreadyExists)
	}

	workbookID := domain.WorkbookID(workbook.ID)

	rbacRepo := r.userRf.NewRBACRepository()
	userObject := user.NewUserObject(user.AppUserID(operator.GetID()))
	workbookObject := domain.NewWorkbookObject(workbookID)
	workbookWriter := domain.NewWorkbookWriter(workbookID)

	// the workbookWriter role can read, update, remove
	if err := rbacRepo.AddNamedPolicy(workbookWriter, workbookObject, domain.PrivilegeRead); err != nil {
		return 0, xerrors.Errorf("Failed to AddNamedPolicy. priv: read, err: %w", err)
	}
	if err := rbacRepo.AddNamedPolicy(workbookWriter, workbookObject, domain.PrivilegeUpdate); err != nil {
		return 0, xerrors.Errorf("Failed to AddNamedPolicy. priv: update, err: %w", err)
	}
	if err := rbacRepo.AddNamedPolicy(workbookWriter, workbookObject, domain.PrivilegeRemove); err != nil {
		return 0, xerrors.Errorf("Failed to AddNamedPolicy. priv: remove, err: %w", err)
	}

	// user is assigned the workbookWriter role
	if err := rbacRepo.AddNamedGroupingPolicy(userObject, workbookWriter); err != nil {
		return 0, xerrors.Errorf("Failed to AddNamedGroupingPolicy. err: %w", err)
	}

	// rbacRepo.NewEnforcerWithRolesAndUsers([]user.RBACRole{workbookWriter}, []user.RBACUser{userObject})

	return workbookID, nil
}

func (r *workbookRepository) RemoveWorkbook(ctx context.Context, operator domain.StudentModel, id domain.WorkbookID, version int) error {
	_, span := tracer.Start(ctx, "workbookRepository.RemoveWorkbook")
	defer span.End()

	workbook := workbookEntity{}
	if result := r.db.Where("organization_id = ? and id = ? and version = ?", operator.GetOrganizationID(), id, version).Delete(&workbook); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return service.ErrWorkbookNotFound
		}

		return result.Error
	}

	return nil
}

func (r *workbookRepository) UpdateWorkbook(ctx context.Context, operator domain.StudentModel, id domain.WorkbookID, version int, param service.WorkbookUpdateParameter) error {
	_, span := tracer.Start(ctx, "workbookRepository.UpdateWorkbook")
	defer span.End()

	if result := r.db.Model(&workbookEntity{}).
		Where("organization_id = ? and id = ? and version = ?",
			uint(operator.GetOrganizationID()), uint(id), version).
		Updates(map[string]interface{}{
			"name":          param.GetName(),
			"question_text": param.GetQuestionText(),
			"version":       gorm.Expr("version + 1"),
		}); result.Error != nil {
		return libG.ConvertDuplicatedError(result.Error, service.ErrWorkbookAlreadyExists)
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
