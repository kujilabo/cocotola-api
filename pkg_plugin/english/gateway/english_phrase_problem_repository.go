package gateway

import (
	"context"
	"errors"
	"math"
	"strconv"
	"time"

	"golang.org/x/xerrors"
	"gorm.io/gorm"

	appD "github.com/kujilabo/cocotola-api/pkg_app/domain"
	appS "github.com/kujilabo/cocotola-api/pkg_app/service"
	lib "github.com/kujilabo/cocotola-api/pkg_lib/domain"
	libG "github.com/kujilabo/cocotola-api/pkg_lib/gateway"
	"github.com/kujilabo/cocotola-api/pkg_lib/log"
	"github.com/kujilabo/cocotola-api/pkg_plugin/english/domain"
	"github.com/kujilabo/cocotola-api/pkg_plugin/english/service"
	user "github.com/kujilabo/cocotola-api/pkg_user/domain"
)

type englishPhraseProblemEntity struct {
	ID             uint
	Version        int
	CreatedAt      time.Time
	UpdatedAt      time.Time
	CreatedBy      uint
	UpdatedBy      uint
	OrganizationID uint
	WorkbookID     uint
	Number         int
	AudioID        uint
	Text           string
	Lang           string
	Translated     string
}

func (e *englishPhraseProblemEntity) TableName() string {
	return "english_phrase_problem"
}

func (e *englishPhraseProblemEntity) toProblem(rf appS.AudioRepositoryFactory) (service.EnglishPhraseProblem, error) {
	model, err := user.NewModel(e.ID, e.Version, e.CreatedAt, e.UpdatedAt, e.CreatedBy, e.UpdatedBy)
	if err != nil {
		return nil, err
	}

	properties := make(map[string]interface{})
	for k, v := range toEnglishPhraseProblemProperties(e.Lang, e.Text, e.Translated) {
		properties[k] = v
	}

	problemModel, err := appD.NewProblemModel(model, e.Number, domain.EnglishPhraseProblemType, properties)
	if err != nil {
		return nil, err
	}

	problem, err := appS.NewProblem(rf, problemModel)
	if err != nil {
		return nil, err
	}

	lang, err := appD.NewLang2(e.Lang)
	if err != nil {
		return nil, err
	}

	englishPhraseProblemModel, err := domain.NewEnglishPhraseProblemModel(problemModel, appD.AudioID(e.AudioID), e.Text, lang, e.Translated)
	if err != nil {
		return nil, err
	}

	return service.NewEnglishPhraseProblem(englishPhraseProblemModel, problem)
}

func fromEnglishPhraseProblemProperties(properties map[string]string) (string, string, string) {
	return properties["lang"], properties["text"], properties["translated"]
}

func toEnglishPhraseProblemProperties(lang, text, translated string) map[string]string {
	return map[string]string{
		"lang":       lang,
		"text":       text,
		"translated": translated,
	}
}

type newEnglishPhraseProblemParam struct {
	AudioID    uint
	Lang       string
	Text       string
	Translated string
}

func toNewEnglishPhraseProblemParam(param appS.ProblemAddParameter) (*newEnglishPhraseProblemParam, error) {
	audioID, err := strconv.Atoi(param.GetProperties()["audioId"])
	if err != nil {
		return nil, err
	}

	lang, text, translated := fromEnglishPhraseProblemProperties(param.GetProperties())
	m := &newEnglishPhraseProblemParam{
		AudioID:    uint(audioID),
		Lang:       lang,
		Text:       text,
		Translated: translated,
	}
	return m, lib.Validator.Struct(m)
}

type englishPhraseProblemRepository struct {
	db          *gorm.DB
	rf          appS.AudioRepositoryFactory
	problemType string
}

func NewEnglishPhraseProblemRepository(db *gorm.DB, rf appS.AudioRepositoryFactory, problemType string) (appS.ProblemRepository, error) {
	return &englishPhraseProblemRepository{
		db:          db,
		rf:          rf,
		problemType: problemType,
	}, nil
}

func (r *englishPhraseProblemRepository) FindProblems(ctx context.Context, operator appD.StudentModel, param appS.ProblemSearchCondition) (appS.ProblemSearchResult, error) {
	limit := param.GetPageSize()
	offset := (param.GetPageNo() - 1) * param.GetPageSize()
	var problemEntities []englishPhraseProblemEntity

	where := r.db.
		Where("organization_id = ?", uint(operator.GetOrganizationID())).
		Where("workbook_id = ?", uint(param.GetWorkbookID()))

	if result := where.Order("workbook_id, number, created_at").
		Limit(limit).Offset(offset).Find(&problemEntities); result.Error != nil {
		return nil, xerrors.Errorf("failed to Find. err: %w", result.Error)
	}

	problems := make([]appD.ProblemModel, len(problemEntities))
	for i, e := range problemEntities {
		p, err := e.toProblem(r.rf)
		if err != nil {
			return nil, xerrors.Errorf("failed to toProblem. err: %w", err)
		}
		problems[i] = p
	}

	var count int64
	if result := where.Model(&englishPhraseProblemEntity{}).Count(&count); result.Error != nil {
		return nil, xerrors.Errorf("failed to Count. err: %w", result.Error)
	}

	if count > math.MaxInt32 {
		return nil, errors.New("overflow")
	}

	return appS.NewProblemSearchResult(int(count), problems)
}

func (r *englishPhraseProblemRepository) FindAllProblems(ctx context.Context, operator appD.StudentModel, workbookID appD.WorkbookID) (appS.ProblemSearchResult, error) {
	limit := 1000
	var problemEntities []englishPhraseProblemEntity

	where := func() *gorm.DB {
		return r.db.
			Where("organization_id = ?", uint(operator.GetOrganizationID())).
			Where("workbook_id = ?", uint(workbookID))
	}
	if result := where().Order("workbook_id, number, created_at").
		Limit(limit).Find(&problemEntities); result.Error != nil {
		return nil, xerrors.Errorf("failed to Find. err: %w", result.Error)
	}

	problems := make([]appD.ProblemModel, len(problemEntities))
	for i, e := range problemEntities {
		p, err := e.toProblem(r.rf)
		if err != nil {
			return nil, xerrors.Errorf("failed to toProblem. err: %w", err)
		}
		problems[i] = p
	}

	var count int64
	if result := where().Model(&englishPhraseProblemEntity{}).Count(&count); result.Error != nil {
		return nil, xerrors.Errorf("failed to Count. err: %w", result.Error)
	}

	if count > math.MaxInt32 {
		return nil, errors.New("overflow")
	}

	return appS.NewProblemSearchResult(int(count), problems)
}

func (r *englishPhraseProblemRepository) FindProblemsByProblemIDs(ctx context.Context, operator appD.StudentModel, param appS.ProblemIDsCondition) (appS.ProblemSearchResult, error) {
	var problemEntities []englishPhraseProblemEntity

	ids := make([]uint, 0)
	for _, id := range param.GetIDs() {
		ids = append(ids, uint(id))
	}

	db := r.db.
		Where("organization_id = ?", uint(operator.GetOrganizationID())).
		Where("workbook_id = ?", uint(param.GetWorkbookID())).
		Where("id in ?", ids)

	if result := db.Find(&problemEntities); result.Error != nil {
		return nil, result.Error
	}

	problems := make([]appD.ProblemModel, len(problemEntities))
	for i, e := range problemEntities {
		p, err := e.toProblem(r.rf)
		if err != nil {
			return nil, err
		}
		problems[i] = p
	}

	return appS.NewProblemSearchResult(0, problems)
}

func (r *englishPhraseProblemRepository) FindProblemByID(ctx context.Context, operator appD.StudentModel, id appS.ProblemSelectParameter1) (appS.Problem, error) {
	var problemEntity englishPhraseProblemEntity

	db := r.db.
		Where("organization_id = ?", uint(operator.GetOrganizationID())).
		Where("workbook_id = ?", uint(id.GetWorkbookID())).
		Where("id = ?", uint(id.GetProblemID()))

	if result := db.First(&problemEntity); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, appS.ErrProblemNotFound
		}
		return nil, result.Error
	}

	return problemEntity.toProblem(r.rf)
}

func (r *englishPhraseProblemRepository) FindProblemIDs(ctx context.Context, operator appD.StudentModel, workbookID appD.WorkbookID) ([]appD.ProblemID, error) {
	pageNo := 1
	pageSize := 1000
	ids := make([]appD.ProblemID, 0)
	for {
		limit := pageSize
		offset := (pageNo - 1) * pageSize
		var problemEntities []englishPhraseProblemEntity

		where := r.db.
			Where("organization_id = ?", uint(operator.GetOrganizationID())).
			Where("workbook_id = ?", uint(workbookID))

		if result := where.Order("workbook_id, number, created_at").
			Limit(limit).Offset(offset).Find(&problemEntities); result.Error != nil {
			return nil, result.Error
		}

		if len(problemEntities) == 0 {
			break
		}

		for _, r := range problemEntities {
			ids = append(ids, appD.ProblemID(r.ID))
		}

		pageNo++
	}

	return ids, nil
}

func (r *englishPhraseProblemRepository) FindProblemsByCustomCondition(ctx context.Context, operator appD.StudentModel, condition interface{}) ([]appD.ProblemModel, error) {
	return nil, errors.New("not implement")
}

func (r *englishPhraseProblemRepository) AddProblem(ctx context.Context, operator appD.StudentModel, param appS.ProblemAddParameter) (appD.ProblemID, error) {
	logger := log.FromContext(ctx)

	problemParam, err := toNewEnglishPhraseProblemParam(param)
	if err != nil {
		return 0, err
	}
	englishPhraseProblem := englishPhraseProblemEntity{
		Version:        1,
		CreatedBy:      operator.GetID(),
		UpdatedBy:      operator.GetID(),
		OrganizationID: uint(operator.GetOrganizationID()),
		WorkbookID:     uint(param.GetWorkbookID()),
		AudioID:        problemParam.AudioID,
		Number:         param.GetNumber(),
		Text:           problemParam.Text,
		Lang:           problemParam.Lang,
		Translated:     problemParam.Translated,
	}

	logger.Infof("englishPhraseProblemRepository.AddProblem. lang: %s, text: %s", problemParam.Lang, problemParam.Text)
	if result := r.db.Create(&englishPhraseProblem); result.Error != nil {
		return 0, libG.ConvertDuplicatedError(result.Error, appS.ErrProblemAlreadyExists)
	}

	return appD.ProblemID(englishPhraseProblem.ID), nil
}

func (r *englishPhraseProblemRepository) UpdateProblem(ctx context.Context, operator appD.StudentModel, id appS.ProblemSelectParameter2, param appS.ProblemUpdateParameter) error {
	return errors.New("not implemented")
}

func (r *englishPhraseProblemRepository) RemoveProblem(ctx context.Context, operator appD.StudentModel, id appS.ProblemSelectParameter2) error {
	logger := log.FromContext(ctx)

	logger.Infof("englishPhraseProblemRepository.RemoveProblem. text: %d", id.GetProblemID())

	result := r.db.
		Where("organization_id = ?", uint(operator.GetOrganizationID())).
		Where("workbook_id = ?", uint(id.GetWorkbookID())).
		Where("id = ?", uint(id.GetProblemID())).
		Where("version = ?", id.GetVersion()).
		Delete(&englishPhraseProblemEntity{})

	if result.Error != nil {
		return result.Error
	} else if result.RowsAffected == 0 {
		return appS.ErrProblemNotFound
	} else if result.RowsAffected != 1 {
		return appS.ErrProblemOtherError
	}

	return nil
}

func (r *englishPhraseProblemRepository) CountProblems(ctx context.Context, operator appD.StudentModel, workbookID appD.WorkbookID) (int, error) {
	ctx, span := tracer.Start(ctx, "englishSentenceProblemRepository.CountProblems")
	defer span.End()

	where := func() *gorm.DB {
		return r.db.
			Where("organization_id = ?", uint(operator.GetOrganizationID())).
			Where("workbook_id = ?", uint(workbookID))
	}

	var count int64
	if result := where().Model(&englishWordProblemEntity{}).Count(&count); result.Error != nil {
		return 0, result.Error
	}

	return int(count), nil
}
