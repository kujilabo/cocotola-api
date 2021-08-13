package gateway

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/go-playground/validator"
	"gorm.io/gorm"

	app "github.com/kujilabo/cocotola-api/pkg_app/domain"
	libG "github.com/kujilabo/cocotola-api/pkg_lib/gateway"
	"github.com/kujilabo/cocotola-api/pkg_lib/log"
	"github.com/kujilabo/cocotola-api/pkg_plugin/english/domain"
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

func (e *englishPhraseProblemEntity) toProblem() (domain.EnglishPhraseProblem, error) {
	model, err := user.NewModel(e.ID, e.Version, e.CreatedAt, e.UpdatedAt, e.CreatedBy, e.UpdatedBy)
	if err != nil {
		return nil, err
	}

	properties := make(map[string]interface{})
	problem, err := app.NewProblem(model, e.Number, domain.EnglishPhraseProblemType, properties)
	if err != nil {
		return nil, err
	}

	lang, err := app.NewLang2(e.Lang)
	if err != nil {
		return nil, err
	}

	return domain.NewEnglishPhraseProblem(problem, app.AudioID(e.AudioID), e.Text, lang, e.Translated)
}

type newEnglishPhraseProblemParam struct {
	AudioID    uint
	Text       string
	Lang       string
	Translated string
}

func toNewEnglishPhraseProblemParam(param app.ProblemAddParameter) (*newEnglishPhraseProblemParam, error) {
	audioID, err := strconv.Atoi(param.GetProperties()["audioId"])
	if err != nil {
		return nil, err
	}

	m := &newEnglishPhraseProblemParam{
		AudioID:    uint(audioID),
		Text:       param.GetProperties()["text"],
		Translated: param.GetProperties()["translated"],
	}
	v := validator.New()
	return m, v.Struct(m)
}

type englishPhraseProblemRepository struct {
	db          *gorm.DB
	problemType string
}

func NewEnglishPhraseProblemRepository(db *gorm.DB, problemType string) (app.ProblemRepository, error) {
	return &englishPhraseProblemRepository{
		db:          db,
		problemType: problemType,
	}, nil
}

func (r *englishPhraseProblemRepository) FindProblems(ctx context.Context, operator app.Student, param app.ProblemSearchCondition) (*app.ProblemSearchResult, error) {
	limit := param.GetPageSize()
	offset := (param.GetPageNo() - 1) * param.GetPageSize()
	var problemEntities []englishPhraseProblemEntity

	where := r.db.Where("organization_id = ? and workbook_id = ?", uint(operator.GetOrganizationID()), uint(param.GetWorkbookID()))
	if result := where.Order("workbook_id, number, created_at").
		Limit(limit).Offset(offset).Find(&problemEntities); result.Error != nil {
		return nil, result.Error
	}

	problems := make([]app.Problem, len(problemEntities))
	for i, e := range problemEntities {
		p, err := e.toProblem()
		if err != nil {
			return nil, err
		}
		problems[i] = p
	}

	var count int64
	if result := where.Model(&englishPhraseProblemEntity{}).Count(&count); result.Error != nil {
		return nil, result.Error
	}

	return &app.ProblemSearchResult{
		TotalCount: count,
		Results:    problems,
	}, nil
}

func (r *englishPhraseProblemRepository) FindProblemsByProblemIDs(ctx context.Context, operator app.Student, param app.ProblemIDsCondition) (*app.ProblemSearchResult, error) {
	var problemEntities []englishPhraseProblemEntity

	ids := make([]uint, 0)
	for _, id := range param.GetIDs() {
		ids = append(ids, uint(id))
	}

	db := r.db.Where("organization_id = ?", uint(operator.GetOrganizationID()))
	db = db.Where("workbook_id = ?", uint(param.GetWorkbookID()))
	db = db.Where("id in ?", ids)
	if result := db.Find(&problemEntities); result.Error != nil {
		return nil, result.Error
	}

	problems := make([]app.Problem, len(problemEntities))
	for i, e := range problemEntities {
		p, err := e.toProblem()
		if err != nil {
			return nil, err
		}
		problems[i] = p
	}

	return &app.ProblemSearchResult{
		Results: problems,
	}, nil
}

func (r *englishPhraseProblemRepository) FindProblemByID(ctx context.Context, operator app.Student, workbookID app.WorkbookID, problemID app.ProblemID) (app.Problem, error) {
	var problemEntity englishPhraseProblemEntity

	db := r.db.Where("organization_id = ?", uint(operator.GetOrganizationID()))
	db = db.Where("workbook_id = ?", uint(workbookID))
	db = db.Where("id = ?", uint(problemID))
	if result := db.First(&problemEntity); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, app.ErrProblemNotFound
		}
		return nil, result.Error
	}

	return problemEntity.toProblem()
}

func (r *englishPhraseProblemRepository) FindProblemIDs(ctx context.Context, operator app.Student, workbookID app.WorkbookID) ([]app.ProblemID, error) {
	pageNo := 1
	pageSize := 1000
	ids := make([]app.ProblemID, 0)
	for {
		limit := pageSize
		offset := (pageNo - 1) * pageSize
		var problemEntities []englishPhraseProblemEntity

		where := r.db.Where("organization_id = ? and workbook_id = ?", uint(operator.GetOrganizationID()), uint(workbookID))
		if result := where.Order("workbook_id, number, created_at").
			Limit(limit).Offset(offset).Find(&problemEntities); result.Error != nil {
			return nil, result.Error
		}

		if len(problemEntities) == 0 {
			break
		}

		for _, r := range problemEntities {
			ids = append(ids, app.ProblemID(r.ID))
		}

		pageNo++
	}

	return ids, nil
}

func (r *englishPhraseProblemRepository) AddProblem(ctx context.Context, operator app.Student, param app.ProblemAddParameter) (app.ProblemID, error) {
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

	logger.Infof("englishPhraseProblemRepository.AddProblem. text: %s", problemParam.Text)
	if result := r.db.Create(&englishPhraseProblem); result.Error != nil {
		return 0, libG.ConvertDuplicatedError(result.Error, app.ErrProblemAlreadyExists)
	}

	return app.ProblemID(englishPhraseProblem.ID), nil
}

func (r *englishPhraseProblemRepository) RemoveProblem(ctx context.Context, operator app.Student, problemID app.ProblemID, version int) error {
	logger := log.FromContext(ctx)

	logger.Infof("englishPhraseProblemRepository.RemoveProblem. text: %d", problemID)
	result := r.db.Where("id = ? and version = ?", uint(problemID), version).Delete(&englishPhraseProblemEntity{})
	if result.Error != nil {
		return result.Error
	} else if result.RowsAffected == 0 {
		return app.ErrProblemNotFound
	} else if result.RowsAffected != 1 {
		return app.ErrProblemOtherError
	}

	return nil
}
