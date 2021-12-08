package gateway

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/go-playground/validator"
	"golang.org/x/xerrors"
	"gorm.io/gorm"

	app "github.com/kujilabo/cocotola-api/pkg_app/domain"
	lib "github.com/kujilabo/cocotola-api/pkg_lib/domain"
	libG "github.com/kujilabo/cocotola-api/pkg_lib/gateway"
	"github.com/kujilabo/cocotola-api/pkg_lib/log"
	"github.com/kujilabo/cocotola-api/pkg_plugin/english/domain"
	user "github.com/kujilabo/cocotola-api/pkg_user/domain"
)

type englishWordProblemEntity struct {
	ID                uint
	Version           int
	CreatedAt         time.Time
	UpdatedAt         time.Time
	CreatedBy         uint
	UpdatedBy         uint
	OrganizationID    uint
	WorkbookID        uint
	Number            int
	AudioID           uint
	Text              string
	Pos               int
	Phonetic          string
	PresentThird      string
	PresentParticiple string
	PastTense         string
	PastParticiple    string
	Lang              string
	Translated        string
	PhraseID1         uint
	PhraseID2         uint
	SentenceID1       uint
	SentenceID2       uint
}

func (e *englishWordProblemEntity) TableName() string {
	return "english_word_problem"
}

func (e *englishWordProblemEntity) toProblem() (domain.EnglishWordProblem, error) {
	model, err := user.NewModel(e.ID, e.Version, e.CreatedAt, e.UpdatedAt, e.CreatedBy, e.UpdatedBy)
	if err != nil {
		return nil, err
	}

	properties := make(map[string]interface{})
	problem, err := app.NewProblem(model, e.Number, domain.EnglishWordProblemType, properties)
	if err != nil {
		return nil, err
	}

	lang, err := app.NewLang2(e.Lang)
	if err != nil {
		return nil, err
	}

	phrases := make([]domain.EnglishPhraseProblem, 0)
	sentences := make([]domain.EnglishSentenceProblem, 0)

	return domain.NewEnglishWordProblem(problem, app.AudioID(e.AudioID), e.Text, e.Pos, e.Phonetic, e.PresentThird, e.PresentParticiple, e.PastTense, e.PastParticiple, lang, e.Translated, phrases, sentences)
}

type englishWordProblemAddParemeter struct {
	AudioID           uint   `validate:"required"`
	Text              string `validate:"required"`
	Pos               int    `validate:"required"`
	Phonetic          string
	PresentThird      string
	PresentParticiple string
	PastTense         string
	PastParticiple    string
	Lang              string `validate:"required"`
	Translated        string
	PhraseID1         uint
	PhraseID2         uint
	SentenceID1       uint
	SentenceID2       uint
}

func toEnglishWordProblemAddParameter(param app.ProblemAddParameter) (*englishWordProblemAddParemeter, error) {
	if _, ok := param.GetProperties()["audioId"]; !ok {
		return nil, xerrors.Errorf("audioId is not defined. err: %w", lib.ErrInvalidArgument)
	}

	if _, ok := param.GetProperties()["pos"]; !ok {
		return nil, xerrors.Errorf("pos is not defined. err: %w", lib.ErrInvalidArgument)
	}

	if _, ok := param.GetProperties()["lang"]; !ok {
		return nil, xerrors.Errorf("lang is not defined. err: %w", lib.ErrInvalidArgument)
	}

	if _, ok := param.GetProperties()["text"]; !ok {
		return nil, xerrors.Errorf("text is not defined. err: %w", lib.ErrInvalidArgument)
	}

	audioID, err := strconv.Atoi(param.GetProperties()["audioId"])
	if err != nil {
		return nil, err
	}

	pos, err := strconv.Atoi(param.GetProperties()["pos"])
	if err != nil {
		return nil, err
	}

	m := &englishWordProblemAddParemeter{
		AudioID:    uint(audioID),
		Lang:       param.GetProperties()["lang"],
		Text:       param.GetProperties()["text"],
		Pos:        pos,
		Translated: param.GetProperties()["translated"],
	}
	v := validator.New()
	return m, v.Struct(m)
}

type englishWordProblemRepository struct {
	db          *gorm.DB
	problemType string
}

func NewEnglishWordProblemRepository(db *gorm.DB, problemType string) (app.ProblemRepository, error) {
	return &englishWordProblemRepository{
		db:          db,
		problemType: problemType,
	}, nil
}

func (r *englishWordProblemRepository) FindProblems(ctx context.Context, operator app.Student, param app.ProblemSearchCondition) (*app.ProblemSearchResult, error) {
	logger := log.FromContext(ctx)
	logger.Debugf("englishWordProblemRepository.FindProblems")
	limit := param.GetPageSize()
	offset := (param.GetPageNo() - 1) * param.GetPageSize()
	var problemEntities []englishWordProblemEntity

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
	if result := where.Model(&englishWordProblemEntity{}).Count(&count); result.Error != nil {
		return nil, result.Error
	}

	logger.Debugf("englishWordProblemRepository.FindProblems, problems: %d, count: %d", len(problems), count)

	return &app.ProblemSearchResult{
		TotalCount: count,
		Results:    problems,
	}, nil
}

func (r *englishWordProblemRepository) FindProblemsByProblemIDs(ctx context.Context, operator app.Student, param app.ProblemIDsCondition) (*app.ProblemSearchResult, error) {
	var problemEntities []englishWordProblemEntity

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

func (r *englishWordProblemRepository) FindProblemByID(ctx context.Context, operator app.Student, workbookID app.WorkbookID, problemID app.ProblemID) (app.Problem, error) {
	var problemEntity englishWordProblemEntity

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

func (r *englishWordProblemRepository) FindProblemIDs(ctx context.Context, operator app.Student, workbookID app.WorkbookID) ([]app.ProblemID, error) {
	pageNo := 1
	pageSize := 1000
	ids := make([]app.ProblemID, 0)
	for {
		limit := pageSize
		offset := (pageNo - 1) * pageSize
		var problemEntities []englishWordProblemEntity

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

func (r *englishWordProblemRepository) AddProblem(ctx context.Context, operator app.Student, param app.ProblemAddParameter) (app.ProblemID, error) {
	logger := log.FromContext(ctx)

	problemParam, err := toEnglishWordProblemAddParameter(param)
	if err != nil {
		return 0, xerrors.Errorf("failed to toEnglishWordProblemAddParameter. param: %+v, err: %w", param, err)
	}
	englishWordProblem := englishWordProblemEntity{
		Version:           1,
		CreatedBy:         operator.GetID(),
		UpdatedBy:         operator.GetID(),
		OrganizationID:    uint(operator.GetOrganizationID()),
		WorkbookID:        uint(param.GetWorkbookID()),
		AudioID:           problemParam.AudioID,
		Number:            param.GetNumber(),
		Text:              problemParam.Text,
		Pos:               problemParam.Pos,
		Phonetic:          problemParam.Phonetic,
		PresentThird:      problemParam.PresentThird,
		PresentParticiple: problemParam.PresentParticiple,
		PastTense:         problemParam.PastTense,
		PastParticiple:    problemParam.PastParticiple,
		Lang:              problemParam.Lang,
		Translated:        problemParam.Translated,
	}

	logger.Infof("englishWordProblemRepository.AddProblem. text: %s", problemParam.Text)
	if result := r.db.Create(&englishWordProblem); result.Error != nil {
		return 0, libG.ConvertDuplicatedError(result.Error, app.ErrProblemAlreadyExists)
	}

	return app.ProblemID(englishWordProblem.ID), nil
}

func (r *englishWordProblemRepository) RemoveProblem(ctx context.Context, operator app.Student, problemID app.ProblemID, version int) error {
	logger := log.FromContext(ctx)

	logger.Infof("englishWordProblemRepository.RemoveProblem. text: %d", problemID)
	result := r.db.Where("id = ? and version = ?", uint(problemID), version).Delete(&englishWordProblemEntity{})
	if result.Error != nil {
		return result.Error
	} else if result.RowsAffected == 0 {
		return app.ErrProblemNotFound
	} else if result.RowsAffected != 1 {
		return app.ErrProblemOtherError
	}

	return nil
}
