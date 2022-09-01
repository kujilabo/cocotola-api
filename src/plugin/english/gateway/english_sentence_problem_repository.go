package gateway

import (
	"context"
	"encoding/json"
	"errors"
	"math"
	"time"

	"gorm.io/gorm"

	appD "github.com/kujilabo/cocotola-api/src/app/domain"
	appS "github.com/kujilabo/cocotola-api/src/app/service"
	libD "github.com/kujilabo/cocotola-api/src/lib/domain"
	liberrors "github.com/kujilabo/cocotola-api/src/lib/errors"
	libG "github.com/kujilabo/cocotola-api/src/lib/gateway"
	"github.com/kujilabo/cocotola-api/src/lib/log"
	"github.com/kujilabo/cocotola-api/src/plugin/english/domain"
	"github.com/kujilabo/cocotola-api/src/plugin/english/service"
	userD "github.com/kujilabo/cocotola-api/src/user/domain"
)

type englishSentenceProblemEntity struct {
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
	Lang2          string
	Translated     string
	Note           string
}

func (e *englishSentenceProblemEntity) TableName() string {
	return "english_sentence_problem"
}

func (e *englishSentenceProblemEntity) toProblem(ctx context.Context, synthesizerClient appS.SynthesizerClient) (service.EnglishSentenceProblem, error) {
	model, err := userD.NewModel(e.ID, e.Version, e.CreatedAt, e.UpdatedAt, e.CreatedBy, e.UpdatedBy)
	if err != nil {
		return nil, err
	}

	properties := make(map[string]interface{})
	problemModel, err := appD.NewProblemModel(model, e.Number, domain.EnglishSentenceProblemType, properties)
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

	englishSentenceProblemModel, err := domain.NewEnglishSentenceProblemModel(problemModel, appD.AudioID(e.AudioID), "", e.Text, lang2, e.Translated, e.Note)
	if err != nil {
		return nil, err
	}

	englishSentenceProblem, err := service.NewEnglishSentenceProblem(englishSentenceProblemModel, problem)
	if err != nil {
		return nil, err
	}

	logger := log.FromContext(ctx)
	logger.Infof("properties: %+v", englishSentenceProblem.GetProperties(ctx))

	return englishSentenceProblem, nil
}

type englishSentenceProblemAddParameter struct {
	AudioID    uint
	Text       string `validate:"required"`
	Lang2      string `validate:"required"`
	Translated string
	Note       string
}

func makeTatoebaNote(param appS.ProblemAddParameter) (string, error) {
	provider, ok := param.GetProperties()[service.EnglishSentenceProblemAddPropertyProvider]
	if !ok {
		return "{}", nil
	}

	if provider == "tatoeba" {
		noteMap := map[string]string{
			"provider": "tatoeba",
		}
		for _, key := range []string{
			service.EnglishSentenceProblemAddPropertyTatoebaSentenceNumber1,
			service.EnglishSentenceProblemAddPropertyTatoebaSentenceNumber2,
			service.EnglishSentenceProblemAddPropertyTatoebaAuthor1,
			service.EnglishSentenceProblemAddPropertyTatoebaAuthor2} {
			if _, ok := param.GetProperties()[key]; !ok {
				return "", liberrors.Errorf("%s is not defined. err: %w", key, libD.ErrInvalidArgument)
			}
			noteMap[key] = param.GetProperties()[key]
		}

		noteBytes, err := json.Marshal(noteMap)
		if err != nil {
			return "", err
		}

		return string(noteBytes), nil
	}

	return "{}", nil
}

// func toEnglishSentenceProblemAddParameter(param appS.ProblemAddParameter) (*englishSentenceProblemAddParameter, error) {
// 	for _, key := range []string{
// 		service.EnglishSentenceProblemAddParemeterAudioID,
// 		service.EnglishSentenceProblemAddParemeterLang,
// 		service.EnglishSentenceProblemAddParemeterText} {

// 	if provider == "tatoeba" {
// 		noteMap := map[string]string{
// 			"provider": "tatoeba",
// 		}
// 	}

// 	var note string
// 	if provider, ok := param.GetProperties()[service.EnglishSentenceProblemAddParemeterProvider]; ok {
// 		if provider == "tatoeba" {
// 			noteMap := map[string]string{}

// 			for _, key := range []string{
// 				service.EnglishSentenceProblemAddParemeterTatoebaSentenceNumber1,
// 				service.EnglishSentenceProblemAddParemeterTatoebaSentenceNumber2,
// 				service.EnglishSentenceProblemAddParemeterTatoebaAuthor1,
// 				service.EnglishSentenceProblemAddParemeterTatoebaAuthor2} {

// 				noteMap[key] = param.GetProperties()[key]
// 			}
// 		}

// 		noteBytes, err := json.Marshal(noteMap)
// 		if err != nil {
// 			return "", err
// 		}
// 		return string(noteBytes), nil
// 	}
// 	return "", nil
// }

func toEnglishSentenceProblemAddParameter(param appS.ProblemAddParameter) (*englishSentenceProblemAddParameter, error) {
	for _, key := range []string{
		service.EnglishSentenceProblemAddPropertyAudioID,
		service.EnglishSentenceProblemAddPropertyLang2,
		service.EnglishSentenceProblemAddPropertyText} {

		if _, ok := param.GetProperties()[key]; !ok {
			return nil, liberrors.Errorf("%s is not defined. err: %w", key, libD.ErrInvalidArgument)
		}
	}

	note, err := makeTatoebaNote(param)
	if err != nil {
		return nil, err
	}
	// audioID, err := strconv.Atoi(param.GetProperties()["audioId"])
	// if err != nil {
	// 	return nil, liberrors.Errorf("audioId is not integer. err: %w", lib.ErrInvalidArgument)
	// }

	m := &englishSentenceProblemAddParameter{
		// AudioID:    uint(audioID),
		Lang2:      param.GetProperties()[service.EnglishSentenceProblemAddPropertyLang2],
		Text:       param.GetProperties()[service.EnglishSentenceProblemAddPropertyText],
		Translated: param.GetProperties()[service.EnglishSentenceProblemAddPropertyTranslated],
		Note:       note,
	}
	return m, libD.Validator.Struct(m)
}

type englishSentenceProblemRepository struct {
	db                *gorm.DB
	synthesizerClient appS.SynthesizerClient
	problemType       string
}

func NewEnglishSentenceProblemRepository(db *gorm.DB, synthesizerClient appS.SynthesizerClient, problemType string) (appS.ProblemRepository, error) {
	return &englishSentenceProblemRepository{
		db:                db,
		synthesizerClient: synthesizerClient,
		problemType:       problemType,
	}, nil
}

func (r *englishSentenceProblemRepository) FindProblems(ctx context.Context, operator appD.StudentModel, param appS.ProblemSearchCondition) (appS.ProblemSearchResult, error) {
	_, span := tracer.Start(ctx, "englishSentenceProblemRepository.FindProblems")
	defer span.End()

	limit := param.GetPageSize()
	offset := (param.GetPageNo() - 1) * param.GetPageSize()

	where := func() *gorm.DB {
		return r.db.
			Where("organization_id = ?", uint(operator.GetOrganizationID())).
			Where("workbook_id = ?", uint(param.GetWorkbookID()))
	}

	var problemEntities []englishSentenceProblemEntity
	if result := where().Order("workbook_id, number, text, created_at").
		Limit(limit).Offset(offset).Find(&problemEntities); result.Error != nil {
		return nil, result.Error
	}

	var count int64
	if result := where().Model(&englishSentenceProblemEntity{}).Count(&count); result.Error != nil {
		return nil, result.Error
	}

	return r.toProblemSearchResult(ctx, count, problemEntities)
}

func (r *englishSentenceProblemRepository) FindAllProblems(ctx context.Context, operator appD.StudentModel, workbookID appD.WorkbookID) (appS.ProblemSearchResult, error) {
	_, span := tracer.Start(ctx, "englishSentenceProblemRepository.FindAllProblems")
	defer span.End()

	limit := 1000

	where := func() *gorm.DB {
		return r.db.
			Where("organization_id = ?", uint(operator.GetOrganizationID())).
			Where("workbook_id = ?", uint(workbookID))
	}

	var problemEntities []englishSentenceProblemEntity
	if result := where().Order("workbook_id, number, text, created_at").
		Limit(limit).Find(&problemEntities); result.Error != nil {
		return nil, result.Error
	}

	var count int64
	if result := where().Model(&englishSentenceProblemEntity{}).Count(&count); result.Error != nil {
		return nil, result.Error
	}

	return r.toProblemSearchResult(ctx, count, problemEntities)
}

func (r *englishSentenceProblemRepository) FindProblemsByProblemIDs(ctx context.Context, operator appD.StudentModel, param appS.ProblemIDsCondition) (appS.ProblemSearchResult, error) {
	_, span := tracer.Start(ctx, "englishSentenceProblemRepository.FindProblemsByProblemIDs")
	defer span.End()

	ids := make([]uint, 0)
	for _, id := range param.GetIDs() {
		ids = append(ids, uint(id))
	}

	db := r.db.
		Where("organization_id = ?", uint(operator.GetOrganizationID())).
		Where("workbook_id = ?", uint(param.GetWorkbookID())).
		Where("id in ?", ids)

	var problemEntities []englishSentenceProblemEntity
	if result := db.Find(&problemEntities); result.Error != nil {
		return nil, result.Error
	}

	problems := make([]appD.ProblemModel, len(problemEntities))
	for i, e := range problemEntities {
		p, err := e.toProblem(ctx, r.synthesizerClient)
		if err != nil {
			return nil, liberrors.Errorf("failed to toProblem. err: %w", err)
		}
		problems[i] = p
	}

	return r.toProblemSearchResult(ctx, 0, problemEntities)
}
func (r *englishSentenceProblemRepository) FindProblemsByCustomCondition(ctx context.Context, operator appD.StudentModel, condition interface{}) ([]appD.ProblemModel, error) {
	_, span := tracer.Start(ctx, "englishSentenceProblemRepository.FindProblemsByCustomCondition")
	defer span.End()

	limit := 1000

	condition1, ok := condition.(map[string]interface{})
	if !ok {
		return nil, libD.ErrInvalidArgument
	}

	conditionWorkbookID, ok := condition1["workbookId"].(uint)
	if !ok {
		return nil, liberrors.Errorf("workbookId is not defined. err: %w", libD.ErrInvalidArgument)
	}

	conditionText, ok := condition1["text"].(string)
	if !ok {
		return nil, liberrors.Errorf("text is not defined. err: %w", libD.ErrInvalidArgument)
	}

	conditionTranslated, ok := condition1["translated"].(string)
	if !ok {
		return nil, liberrors.Errorf("translated is not defined. err: %w", libD.ErrInvalidArgument)
	}

	var problemEntity englishSentenceProblemEntity

	db := r.db.
		Where("organization_id = ?", uint(operator.GetOrganizationID())).
		Where("workbook_id = ?", conditionWorkbookID).
		Where("text = ?", conditionText).
		Where("translated = ?", conditionTranslated)

	if result := db.Limit(limit).First(&problemEntity); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return []appD.ProblemModel{}, nil
		}
		return nil, result.Error
	}

	problem, err := problemEntity.toProblem(ctx, r.synthesizerClient)
	if err != nil {
		return nil, err
	}

	return []appD.ProblemModel{problem}, nil
}

func (r *englishSentenceProblemRepository) toProblemSearchResult(ctx context.Context, count int64, problemEntities []englishSentenceProblemEntity) (appS.ProblemSearchResult, error) {
	problems := make([]appD.ProblemModel, len(problemEntities))
	for i, e := range problemEntities {
		p, err := e.toProblem(ctx, r.synthesizerClient)
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

func (r *englishSentenceProblemRepository) FindProblemByID(ctx context.Context, operator appD.StudentModel, id appS.ProblemSelectParameter1) (appS.Problem, error) {
	_, span := tracer.Start(ctx, "englishSentenceProblemRepository.FindProblemByID")
	defer span.End()

	var problemEntity englishSentenceProblemEntity

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

	return problemEntity.toProblem(ctx, r.synthesizerClient)
}

func (r *englishSentenceProblemRepository) FindProblemIDs(ctx context.Context, operator appD.StudentModel, workbookID appD.WorkbookID) ([]appD.ProblemID, error) {
	_, span := tracer.Start(ctx, "englishSentenceProblemRepository.FindProblemIDs")
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

		var problemEntities []englishSentenceProblemEntity
		if result := where.Order("workbook_id, number, text, created_at").
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

func (r *englishSentenceProblemRepository) AddProblem(ctx context.Context, operator appD.StudentModel, param appS.ProblemAddParameter) (appD.ProblemID, error) {
	_, span := tracer.Start(ctx, "englishSentenceProblemRepository.AddProblem")
	defer span.End()

	logger := log.FromContext(ctx)

	problemParam, err := toEnglishSentenceProblemAddParameter(param)
	if err != nil {
		return 0, err
	}

	englishSentenceProblem := englishSentenceProblemEntity{
		Version:        1,
		CreatedBy:      operator.GetID(),
		UpdatedBy:      operator.GetID(),
		OrganizationID: uint(operator.GetOrganizationID()),
		WorkbookID:     uint(param.GetWorkbookID()),
		AudioID:        problemParam.AudioID,
		Number:         param.GetNumber(),
		Text:           problemParam.Text,
		Lang2:          problemParam.Lang2,
		Translated:     problemParam.Translated,
		Note:           problemParam.Note,
	}

	logger.Infof("englishSentenceProblemRepository.AddProblem. text: %s", problemParam.Text)

	if result := r.db.Create(&englishSentenceProblem); result.Error != nil {
		return 0, libG.ConvertDuplicatedError(liberrors.Errorf("failed to Create englishSentenceProblem. err: %w", result.Error), appS.ErrProblemAlreadyExists)
	}

	return appD.ProblemID(englishSentenceProblem.ID), nil
}

func (r *englishSentenceProblemRepository) UpdateProblem(ctx context.Context, operator appD.StudentModel, id appS.ProblemSelectParameter2, param appS.ProblemUpdateParameter) error {
	return errors.New("not implemented")
}

func (r *englishSentenceProblemRepository) RemoveProblem(ctx context.Context, operator appD.StudentModel, id appS.ProblemSelectParameter2) error {
	_, span := tracer.Start(ctx, "englishSentenceProblemRepository.RemoveProblem")
	defer span.End()

	logger := log.FromContext(ctx)

	logger.Infof("englishSentenceProblemRepository.RemoveProblem. problemID: %d", id.GetProblemID())
	result := r.db.Where("id = ? and version = ?", uint(id.GetProblemID()), id.GetVersion()).Delete(&englishSentenceProblemEntity{})
	if result.Error != nil {
		return result.Error
	} else if result.RowsAffected == 0 {
		return appS.ErrProblemNotFound
	} else if result.RowsAffected != 1 {
		return appS.ErrProblemOtherError
	}

	return nil
}

func (r *englishSentenceProblemRepository) CountProblems(ctx context.Context, operator appD.StudentModel, workbookID appD.WorkbookID) (int, error) {
	_, span := tracer.Start(ctx, "englishSentenceProblemRepository.CountProblems")
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
