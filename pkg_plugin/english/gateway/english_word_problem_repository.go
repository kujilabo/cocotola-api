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

	phrases := make([]domain.EnglishPhraseProblem, 0)
	sentences := make([]domain.EnglishSentenceProblem, 0)

	return domain.NewEnglishWordProblem(problem, app.AudioID(e.AudioID), e.Text, e.Pos, e.Phonetic, e.PresentThird, e.PresentParticiple, e.PastTense, e.PastParticiple, e.Lang, e.Translated, phrases, sentences)
}

type newEnglishWordProblemParam struct {
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

func toNewEnglishWordProblemParam(param *app.ProblemAddParameter) (*newEnglishWordProblemParam, error) {
	audioID, err := strconv.Atoi(param.Properties["audioId"])
	if err != nil {
		return nil, err
	}
	pos, err := strconv.Atoi(param.Properties["pos"])
	if err != nil {
		return nil, err
	}

	m := &newEnglishWordProblemParam{
		AudioID:    uint(audioID),
		Text:       param.Properties["text"],
		Pos:        pos,
		Translated: param.Properties["translated"],
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

func (r *englishWordProblemRepository) FindProblems(ctx context.Context, operator app.Student, param *app.ProblemSearchCondition) (*app.ProblemSearchResult, error) {
	limit := param.PageSize
	offset := (param.PageNo - 1) * param.PageSize
	var problemEntities []englishWordProblemEntity

	where := r.db.Where("organization_id = ? and workbook_id = ?", uint(operator.GetOrganizationID()), uint(param.WorkbookID))
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

	return &app.ProblemSearchResult{
		TotalCount: count,
		Results:    problems,
	}, nil
}

func (r *englishWordProblemRepository) FindProblemsByProblemIDs(ctx context.Context, operator app.Student, param *app.ProblemIDsCondition) (*app.ProblemSearchResult, error) {
	var problemEntities []englishWordProblemEntity

	ids := make([]uint, 0)
	for _, id := range param.IDs {
		ids = append(ids, uint(id))
	}

	db := r.db.Where("organization_id = ?", uint(operator.GetOrganizationID()))
	db = db.Where("workbook_id = ?", uint(param.WorkbookID))
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

func (r *englishWordProblemRepository) AddProblem(ctx context.Context, operator app.Student, param *app.ProblemAddParameter) (app.ProblemID, error) {
	logger := log.FromContext(ctx)

	problemParam, err := toNewEnglishWordProblemParam(param)
	if err != nil {
		return 0, err
	}
	englishWordProblem := englishWordProblemEntity{
		Version:           1,
		CreatedBy:         operator.GetID(),
		UpdatedBy:         operator.GetID(),
		OrganizationID:    uint(operator.GetOrganizationID()),
		WorkbookID:        uint(param.WorkbookID),
		AudioID:           problemParam.AudioID,
		Number:            param.Number,
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
