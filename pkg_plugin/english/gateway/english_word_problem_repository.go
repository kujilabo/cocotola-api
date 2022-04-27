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
	libD "github.com/kujilabo/cocotola-api/pkg_lib/domain"
	libG "github.com/kujilabo/cocotola-api/pkg_lib/gateway"
	"github.com/kujilabo/cocotola-api/pkg_lib/log"
	"github.com/kujilabo/cocotola-api/pkg_plugin/english/domain"
	"github.com/kujilabo/cocotola-api/pkg_plugin/english/service"
	userD "github.com/kujilabo/cocotola-api/pkg_user/domain"
)

var (
	MaxNumberOfProblemsToFindAllProblems  = 1000
	MaxNumberOfProblemIDsToFindProblemIDs = 100
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
	Lang2             string
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

func (e *englishWordProblemEntity) toProblem(synthesizerClient appS.SynthesizerClient) (service.EnglishWordProblem, error) {
	model, err := userD.NewModel(e.ID, e.Version, e.CreatedAt, e.UpdatedAt, e.CreatedBy, e.UpdatedBy)
	if err != nil {
		return nil, err
	}

	properties := make(map[string]interface{})
	problemModel, err := appD.NewProblemModel(model, e.Number, domain.EnglishWordProblemType, properties)
	if err != nil {
		return nil, err
	}

	problem, err := appS.NewProblem(synthesizerClient, problemModel)
	if err != nil {
		return nil, err
	}

	lang2, err := appD.NewLang2(e.Lang2)
	if err != nil {
		return nil, err
	}

	phrases := make([]domain.EnglishPhraseProblemModel, 0)
	sentences := make([]domain.EnglishWordSentenceProblemModel, 0)
	if e.SentenceID1 != 0 {
		sentence, err := domain.NewEnglishWordProblemSentenceModel(appD.AudioID(0), e.SentenceText1, lang2, e.SentenceTranslated1, e.SentenceNote1)
		if err != nil {
			return nil, err
		}
		sentences = append(sentences, sentence)
	}
	englishWordProblemModel, err := domain.NewEnglishWordProblemModel(problemModel, appD.AudioID(e.AudioID), e.Text, e.Pos, e.Phonetic, e.PresentThird, e.PresentParticiple, e.PastTense, e.PastParticiple, lang2, e.Translated, phrases, sentences)
	if err != nil {
		return nil, err
	}

	return service.NewEnglishWordProblem(englishWordProblemModel, problem)
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
	Lang2             string `validate:"required"`
	Translated        string
	PhraseID1         uint
	PhraseID2         uint
	SentenceID1       uint
	SentenceID2       uint
}

func toEnglishWordProblemAddParameter(param appS.ProblemAddParameter) (*englishWordProblemAddParemeter, error) {
	if _, ok := param.GetProperties()["audioId"]; !ok {
		return nil, xerrors.Errorf("audioId is not defined. err: %w", libD.ErrInvalidArgument)
	}

	if _, ok := param.GetProperties()["pos"]; !ok {
		return nil, xerrors.Errorf("pos is not defined. err: %w", libD.ErrInvalidArgument)
	}

	if _, ok := param.GetProperties()["lang2"]; !ok {
		return nil, xerrors.Errorf("lang2 is not defined. err: %w", libD.ErrInvalidArgument)
	}

	if _, ok := param.GetProperties()["text"]; !ok {
		return nil, xerrors.Errorf("text is not defined. err: %w", libD.ErrInvalidArgument)
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
		Lang2:      param.GetProperties()["lang2"],
		Text:       param.GetProperties()["text"],
		Pos:        pos,
		Translated: param.GetProperties()["translated"],
	}
	return m, libD.Validator.Struct(m)
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

func toEnglishWordProblemUpdateParameter(param appS.ProblemUpdateParameter) (*englishWordProblemUpdateParemeter, error) {
	if _, ok := param.GetProperties()[service.EnglishWordProblemUpdatePropertyAudioID]; !ok {
		return nil, xerrors.Errorf("audioId is not defined. err: %w", libD.ErrInvalidArgument)
	}

	text, err := param.GetStringProperty(service.EnglishWordProblemUpdatePropertyText)
	if err != nil {
		return nil, xerrors.Errorf("text is not defined. err: %w", libD.ErrInvalidArgument)
	}

	audioID, err := param.GetIntProperty(service.EnglishWordProblemUpdatePropertyAudioID)
	if err != nil {
		return nil, err
	}

	sentenceID, err := param.GetIntProperty(service.EnglishWordProblemUpdatePropertySentenceID1)
	if err != nil {
		return nil, err
	}

	m := &englishWordProblemUpdateParemeter{
		AudioID:     uint(audioID),
		Text:        text,
		Translated:  param.GetProperties()[service.EnglishWordProblemUpdatePropertyTranslated],
		SentenceID1: uint(sentenceID),
	}
	return m, libD.Validator.Struct(m)
}

type englishWordProblemRepository struct {
	db                *gorm.DB
	synthesizerClient appS.SynthesizerClient
	problemType       string
}

func NewEnglishWordProblemRepository(db *gorm.DB, synthesizerClient appS.SynthesizerClient, problemType string) (appS.ProblemRepository, error) {
	return &englishWordProblemRepository{
		db:                db,
		synthesizerClient: synthesizerClient,
		problemType:       problemType,
	}, nil
}

func (r *englishWordProblemRepository) FindProblems(ctx context.Context, operator appD.StudentModel, param appS.ProblemSearchCondition) (appS.ProblemSearchResult, error) {
	_, span := tracer.Start(ctx, "englishWordProblemRepository.FindProblems")
	defer span.End()

	limit := param.GetPageSize()
	offset := (param.GetPageNo() - 1) * param.GetPageSize()

	where := func() *gorm.DB {
		return r.db.
			Where("organization_id = ?", uint(operator.GetOrganizationID())).
			Where("workbook_id = ?", uint(param.GetWorkbookID()))
	}

	var problemEntities []englishWordProblemEntity
	if result := where().Order("text, pos").
		Limit(limit).Offset(offset).Find(&problemEntities); result.Error != nil {
		return nil, result.Error
	}

	var count int64
	if result := where().Model(&englishWordProblemEntity{}).Count(&count); result.Error != nil {
		return nil, result.Error
	}

	return r.toProblemSearchResult(count, problemEntities)
}

func (r *englishWordProblemRepository) FindAllProblems(ctx context.Context, operator appD.StudentModel, workbookID appD.WorkbookID) (appS.ProblemSearchResult, error) {
	_, span := tracer.Start(ctx, "englishWordProblemRepository.FindAllProblems")
	defer span.End()

	limit := MaxNumberOfProblemsToFindAllProblems

	where := func() *gorm.DB {
		return r.db.
			Where("organization_id = ?", uint(operator.GetOrganizationID())).
			Where("workbook_id = ?", uint(workbookID))
	}

	var problemEntities []englishWordProblemEntity
	if result := where().Order("text, pos").
		Limit(limit).Find(&problemEntities); result.Error != nil {
		return nil, result.Error
	}

	var count int64
	if result := where().Model(&englishWordProblemEntity{}).Count(&count); result.Error != nil {
		return nil, result.Error
	}

	return r.toProblemSearchResult(count, problemEntities)
}

func (r *englishWordProblemRepository) FindProblemsByProblemIDs(ctx context.Context, operator appD.StudentModel, param appS.ProblemIDsCondition) (appS.ProblemSearchResult, error) {
	_, span := tracer.Start(ctx, "englishWordProblemRepository.FindProblemsByProblemIDs")
	defer span.End()

	if len(param.GetIDs()) > MaxNumberOfProblemIDsToFindProblemIDs {
		return nil, libD.ErrInvalidArgument
	}

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

func (r *englishWordProblemRepository) toProblemSearchResult(count int64, problemEntities []englishWordProblemEntity) (appS.ProblemSearchResult, error) {
	problems := make([]appD.ProblemModel, len(problemEntities))
	for i, e := range problemEntities {
		p, err := e.toProblem(r.synthesizerClient)
		if err != nil {
			return nil, err
		}
		problems[i] = p
	}

	if count > math.MaxInt32 {
		return nil, errors.New("overflow")
	}

	return appS.NewProblemSearchResult(int(count), problems)
}

func (r *englishWordProblemRepository) FindProblemByID(ctx context.Context, operator appD.StudentModel, id appS.ProblemSelectParameter1) (appS.Problem, error) {
	_, span := tracer.Start(ctx, "englishWordProblemRepository.FindProblemByID")
	defer span.End()

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
			return nil, appS.ErrProblemNotFound
		}
		return nil, result.Error
	}

	return problemEntity.toProblem(r.synthesizerClient)
}

func (r *englishWordProblemRepository) FindProblemIDs(ctx context.Context, operator appD.StudentModel, workbookID appD.WorkbookID) ([]appD.ProblemID, error) {
	_, span := tracer.Start(ctx, "englishWordProblemRepository.FindProblemIDs")
	defer span.End()

	pageNo := 1
	pageSize := 1000
	limit := pageSize

	ids := make([]appD.ProblemID, 0)
	for {
		offset := (pageNo - 1) * pageSize

		where := r.db.
			Where("organization_id = ?", uint(operator.GetOrganizationID())).
			Where("workbook_id = ?", uint(workbookID))

		var problemEntities []englishWordProblemEntity
		if result := where.Order("text, pos").
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

func (r *englishWordProblemRepository) FindProblemsByCustomCondition(ctx context.Context, operator appD.StudentModel, condition interface{}) ([]appD.ProblemModel, error) {
	return nil, errors.New("not implement")
}

func (r *englishWordProblemRepository) AddProblem(ctx context.Context, operator appD.StudentModel, param appS.ProblemAddParameter) (appD.ProblemID, error) {
	ctx, span := tracer.Start(ctx, "englishWordProblemRepository.AddProblem")
	defer span.End()

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
		Lang2:             problemParam.Lang2,
		Translated:        problemParam.Translated,
	}

	logger.Infof("englishWordProblemRepository.AddProblem. text: %s", problemParam.Text)
	if result := r.db.Create(&englishWordProblem); result.Error != nil {
		return 0, xerrors.Errorf("failed to Create. param: %+v, err: %w", param, libG.ConvertDuplicatedError(result.Error, appS.ErrProblemAlreadyExists))
	}

	return appD.ProblemID(englishWordProblem.ID), nil
}

func (r *englishWordProblemRepository) UpdateProblem(ctx context.Context, operator appD.StudentModel, id appS.ProblemSelectParameter2, param appS.ProblemUpdateParameter) error {
	ctx, span := tracer.Start(ctx, "englishWordProblemRepository.UpdateProblem")
	defer span.End()

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
		return appS.ErrProblemNotFound
	} else if result.RowsAffected != 1 {
		return appS.ErrProblemOtherError
	}

	return nil
}

func (r *englishWordProblemRepository) RemoveProblem(ctx context.Context, operator appD.StudentModel, id appS.ProblemSelectParameter2) error {
	_, span := tracer.Start(ctx, "englishWordProblemRepository.RemoveProblem")
	defer span.End()

	result := r.db.
		Where("organization_id = ?", uint(operator.GetOrganizationID())).
		Where("workbook_id = ?", uint(id.GetWorkbookID())).
		Where("id = ?", uint(id.GetProblemID())).
		Where("version = ?", id.GetVersion()).
		Delete(&englishWordProblemEntity{})

	if result.Error != nil {
		return result.Error
	} else if result.RowsAffected == 0 {
		return appS.ErrProblemNotFound
	} else if result.RowsAffected != 1 {
		return appS.ErrProblemOtherError
	}

	return nil
}

func (r *englishWordProblemRepository) CountProblems(ctx context.Context, operator appD.StudentModel, workbookID appD.WorkbookID) (int, error) {
	_, span := tracer.Start(ctx, "englishWordProblemRepository.CountProblems")
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
