package gateway

import (
	"context"
	"errors"
	"math"
	"strconv"
	"time"

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
	// joined columns
	SentenceText1       string `gorm:"->"` // readonly
	SentenceTranslated1 string `gorm:"->"` // readonly
	SentenceNote1       string `gorm:"->"` // readonly
}

func (e *englishWordProblemEntity) TableName() string {
	return "english_word_problem"
}

func (e *englishWordProblemEntity) toProblem(rf app.AudioRepositoryFactory) (domain.EnglishWordProblem, error) {
	model, err := user.NewModel(e.ID, e.Version, e.CreatedAt, e.UpdatedAt, e.CreatedBy, e.UpdatedBy)
	if err != nil {
		return nil, err
	}

	properties := make(map[string]interface{})
	problem, err := app.NewProblem(rf, model, e.Number, domain.EnglishWordProblemType, properties)
	if err != nil {
		return nil, err
	}

	lang, err := app.NewLang2(e.Lang)
	if err != nil {
		return nil, err
	}

	phrases := make([]domain.EnglishPhraseProblem, 0)
	sentences := make([]domain.EnglishWordProblemSentence, 0)

	if e.SentenceID1 != 0 {
		sentence, err := domain.NewEnglishWordProblemSentence(app.AudioID(0), e.SentenceText1, app.Lang2JA, e.SentenceTranslated1, e.SentenceNote1)
		if err != nil {
			return nil, err
		}
		sentences = append(sentences, sentence)
	}
	return domain.NewEnglishWordProblem(problem, app.AudioID(e.AudioID), e.Text, e.Pos, e.Phonetic, e.PresentThird, e.PresentParticiple, e.PastTense, e.PastParticiple, lang, e.Translated, phrases, sentences)
}

type englishWordProblemAddParemeter struct {
	AudioID           uint
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
	return m, lib.Validator.Struct(m)
}

type englishWordProblemUpdateParemeter struct {
	AudioID           uint
	Text              string `validate:"required"`
	Phonetic          string
	PresentThird      string
	PresentParticiple string
	PastTense         string
	PastParticiple    string
	Translated        string
	PhraseID1         uint
	PhraseID2         uint
	SentenceID1       uint
	SentenceID2       uint
}

func toEnglishWordProblemUpdateParameter(param app.ProblemUpdateParameter) (*englishWordProblemUpdateParemeter, error) {
	if _, ok := param.GetProperties()[domain.EnglishWordProblemUpdatePropertyAudioID]; !ok {
		return nil, xerrors.Errorf("audioId is not defined. err: %w", lib.ErrInvalidArgument)
	}

	text, err := param.GetStringProperty(domain.EnglishWordProblemUpdatePropertyText)
	if err != nil {
		return nil, xerrors.Errorf("text is not defined. err: %w", lib.ErrInvalidArgument)
	}

	audioID, err := param.GetIntProperty(domain.EnglishWordProblemUpdatePropertyAudioID)
	if err != nil {
		return nil, err
	}

	sentenceID, err := param.GetIntProperty(domain.EnglishWordProblemUpdatePropertySentenceID1)
	if err != nil {
		return nil, err
	}

	m := &englishWordProblemUpdateParemeter{
		AudioID:     uint(audioID),
		Text:        text,
		Translated:  param.GetProperties()[domain.EnglishWordProblemUpdatePropertyTranslated],
		SentenceID1: uint(sentenceID),
	}
	return m, lib.Validator.Struct(m)
}

type englishWordProblemRepository struct {
	db          *gorm.DB
	rf          app.AudioRepositoryFactory
	problemType string
}

func NewEnglishWordProblemRepository(db *gorm.DB, rf app.AudioRepositoryFactory, problemType string) (app.ProblemRepository, error) {
	return &englishWordProblemRepository{
		db:          db,
		rf:          rf,
		problemType: problemType,
	}, nil
}

func (r *englishWordProblemRepository) FindProblems(ctx context.Context, operator app.Student, param app.ProblemSearchCondition) (app.ProblemSearchResult, error) {
	logger := log.FromContext(ctx)
	logger.Debug("englishWordProblemRepository.FindProblems")

	limit := param.GetPageSize()
	offset := (param.GetPageNo() - 1) * param.GetPageSize()

	where := func() *gorm.DB {
		return r.db.
			Where("organization_id = ?", uint(operator.GetOrganizationID())).
			Where("workbook_id = ?", uint(param.GetWorkbookID()))
	}

	var problemEntities []englishWordProblemEntity
	if result := where().Order("workbook_id, number, created_at").
		Limit(limit).Offset(offset).Find(&problemEntities); result.Error != nil {
		return nil, result.Error
	}

	var count int64
	if result := where().Model(&englishWordProblemEntity{}).Count(&count); result.Error != nil {
		return nil, result.Error
	}

	return r.toProblemSearchResult(count, problemEntities)
}

func (r *englishWordProblemRepository) FindAllProblems(ctx context.Context, operator app.Student, workbookID app.WorkbookID) (app.ProblemSearchResult, error) {
	logger := log.FromContext(ctx)
	logger.Debug("englishWordProblemRepository.FindProblems")
	limit := 1000

	where := func() *gorm.DB {
		return r.db.
			Where("organization_id = ?", uint(operator.GetOrganizationID())).
			Where("workbook_id = ?", uint(workbookID))
	}

	var problemEntities []englishWordProblemEntity
	if result := where().Order("workbook_id, number, created_at").
		Limit(limit).Find(&problemEntities); result.Error != nil {
		return nil, result.Error
	}

	var count int64
	if result := where().Model(&englishWordProblemEntity{}).Count(&count); result.Error != nil {
		return nil, result.Error
	}

	return r.toProblemSearchResult(count, problemEntities)
}

func (r *englishWordProblemRepository) FindProblemsByProblemIDs(ctx context.Context, operator app.Student, param app.ProblemIDsCondition) (app.ProblemSearchResult, error) {

	ids := make([]uint, 0)
	for _, id := range param.GetIDs() {
		ids = append(ids, uint(id))
	}

	db := r.db.
		Where("organization_id = ?", uint(operator.GetOrganizationID())).
		Where("workbook_id = ?", uint(param.GetWorkbookID())).
		Where("id in ?", ids)

	var problemEntities []englishWordProblemEntity
	if result := db.Find(&problemEntities); result.Error != nil {
		return nil, result.Error
	}

	return r.toProblemSearchResult(0, problemEntities)
}

func (r *englishWordProblemRepository) FindProblemsByCustomCondition(ctx context.Context, operator app.Student, condition interface{}) ([]app.Problem, error) {
	return nil, errors.New("not implement")
}

func (r *englishWordProblemRepository) toProblemSearchResult(count int64, problemEntities []englishWordProblemEntity) (app.ProblemSearchResult, error) {
	problems := make([]app.Problem, len(problemEntities))
	for i, e := range problemEntities {
		p, err := e.toProblem(r.rf)
		if err != nil {
			return nil, err
		}
		problems[i] = p
	}

	if count > math.MaxInt32 {
		return nil, errors.New("overflow")
	}

	return app.NewProblemSearchResult(int(count), problems)
}

func (r *englishWordProblemRepository) FindProblemByID(ctx context.Context, operator app.Student, id app.ProblemSelectParameter1) (app.Problem, error) {
	var problemEntity englishWordProblemEntity

	db := r.db.Table("english_word_problem AS T1").
		Select("T1.*,"+
			"T2.text AS sentence_text1,"+
			"T2.translated AS sentence_translated1,"+
			"T2.note AS sentence_note1").
		Joins("LEFT JOIN english_sentence_problem as T2 ON T1.sentence_id1 = T2.id").
		Where("T1.organization_id = ?", uint(operator.GetOrganizationID())).
		Where("T1.workbook_id = ?", uint(id.GetWorkbookID())).
		Where("T1.id = ?", uint(id.GetProblemID()))

	if result := db.First(&problemEntity); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, app.ErrProblemNotFound
		}
		return nil, result.Error
	}

	return problemEntity.toProblem(r.rf)
}

func (r *englishWordProblemRepository) FindProblemIDs(ctx context.Context, operator app.Student, workbookID app.WorkbookID) ([]app.ProblemID, error) {
	pageNo := 1
	pageSize := 1000
	limit := pageSize

	ids := make([]app.ProblemID, 0)
	for {
		offset := (pageNo - 1) * pageSize

		where := r.db.
			Where("organization_id = ?", uint(operator.GetOrganizationID())).
			Where("workbook_id = ?", uint(workbookID))

		var problemEntities []englishWordProblemEntity
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
		return 0, xerrors.Errorf("failed to Create. param: %+v, err: %w", param, libG.ConvertDuplicatedError(result.Error, app.ErrProblemAlreadyExists))
	}

	return app.ProblemID(englishWordProblem.ID), nil
}

func (r *englishWordProblemRepository) UpdateProblem(ctx context.Context, operator app.Student, id app.ProblemSelectParameter2, param app.ProblemUpdateParameter) error {
	logger := log.FromContext(ctx)

	problemParam, err := toEnglishWordProblemUpdateParameter(param)
	if err != nil {
		return xerrors.Errorf("failed to toEnglishWordProblemUdateParameter. param: %+v, err: %w", param, err)
	}

	englishWordProblem := englishWordProblemEntity{
		Version:           id.GetVersion() + 1,
		UpdatedBy:         operator.GetID(),
		AudioID:           problemParam.AudioID,
		Number:            param.GetNumber(),
		Phonetic:          problemParam.Phonetic,
		PresentThird:      problemParam.PresentThird,
		PresentParticiple: problemParam.PresentParticiple,
		PastTense:         problemParam.PastTense,
		PastParticiple:    problemParam.PastParticiple,
		Translated:        problemParam.Translated,
		SentenceID1:       problemParam.SentenceID1,
	}

	logger.Infof("englishWordProblemRepository.UpdateProblem. text: %s", problemParam.Text)

	result := r.db.
		Where("organization_id = ?", uint(operator.GetOrganizationID())).
		Where("workbook_id = ?", uint(id.GetWorkbookID())).
		Where("id = ?", uint(id.GetProblemID())).
		Where("version = ?", id.GetVersion()).
		UpdateColumns(&englishWordProblem)

	if result.Error != nil {
		return result.Error
	} else if result.RowsAffected == 0 {
		return app.ErrProblemNotFound
	} else if result.RowsAffected != 1 {
		return app.ErrProblemOtherError
	}

	return nil
}

func (r *englishWordProblemRepository) RemoveProblem(ctx context.Context, operator app.Student, id app.ProblemSelectParameter2) error {
	logger := log.FromContext(ctx)

	logger.Infof("englishWordProblemRepository.RemoveProblem. ProblemID: %d", id.GetProblemID())

	result := r.db.
		Where("organization_id = ?", uint(operator.GetOrganizationID())).
		Where("workbook_id = ?", uint(id.GetWorkbookID())).
		Where("id = ?", uint(id.GetProblemID())).
		Where("version = ?", id.GetVersion()).
		Delete(&englishWordProblemEntity{})

	if result.Error != nil {
		return result.Error
	} else if result.RowsAffected == 0 {
		return app.ErrProblemNotFound
	} else if result.RowsAffected != 1 {
		return app.ErrProblemOtherError
	}

	return nil
}
