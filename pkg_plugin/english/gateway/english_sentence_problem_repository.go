package gateway

import (
	"context"
	"encoding/json"
	"errors"
	"math"
	"time"

	"golang.org/x/xerrors"
	"gorm.io/gorm"

	app "github.com/kujilabo/cocotola-api/pkg_app/domain"
	appS "github.com/kujilabo/cocotola-api/pkg_app/service"
	lib "github.com/kujilabo/cocotola-api/pkg_lib/domain"
	libG "github.com/kujilabo/cocotola-api/pkg_lib/gateway"
	"github.com/kujilabo/cocotola-api/pkg_lib/log"
	"github.com/kujilabo/cocotola-api/pkg_plugin/english/domain"
	"github.com/kujilabo/cocotola-api/pkg_plugin/english/service"
	user "github.com/kujilabo/cocotola-api/pkg_user/domain"
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
	Lang           string
	Translated     string
	Note           string
}

func (e *englishSentenceProblemEntity) TableName() string {
	return "english_sentence_problem"
}

func (e *englishSentenceProblemEntity) toProblem(rf appS.AudioRepositoryFactory) (service.EnglishSentenceProblem, error) {
	model, err := user.NewModel(e.ID, e.Version, e.CreatedAt, e.UpdatedAt, e.CreatedBy, e.UpdatedBy)
	if err != nil {
		return nil, err
	}

	properties := make(map[string]interface{})
	problemModel, err := app.NewProblemModel(model, e.Number, domain.EnglishSentenceProblemType, properties)
	if err != nil {
		return nil, err
	}

	problem, err := appS.NewProblem(rf, problemModel)
	if err != nil {
		return nil, err
	}

	lang, err := app.NewLang2(e.Lang)
	if err != nil {
		return nil, err
	}

	englishSentenceProblemModel, err := domain.NewEnglishSentenceProblemModel(problemModel, app.AudioID(e.AudioID), "", e.Text, lang, e.Translated, e.Note)
	if err != nil {
		return nil, err
	}
	return service.NewEnglishSentenceProblem(englishSentenceProblemModel, problem)
}

type englishSentenceProblemAddParameter struct {
	AudioID    uint
	Text       string `validate:"required"`
	Lang       string `validate:"required"`
	Translated string
	Note       string
}

func toEnglishSentenceProblemAddParameter(param appS.ProblemAddParameter) (*englishSentenceProblemAddParameter, error) {
	for _, key := range []string{
		service.EnglishSentenceProblemAddParemeterAudioID,
		service.EnglishSentenceProblemAddParemeterLang,
		service.EnglishSentenceProblemAddParemeterText} {

		if _, ok := param.GetProperties()[key]; !ok {
			return nil, xerrors.Errorf("%s is not defined. err: %w", key, lib.ErrInvalidArgument)
		}
	}

	var note string
	if provider, ok := param.GetProperties()[service.EnglishSentenceProblemAddParemeterProvider]; ok {
		if provider == "tatoeba" {
			noteMap := map[string]string{}

			for _, key := range []string{
				service.EnglishSentenceProblemAddParemeterTatoebaSentenceNumber1,
				service.EnglishSentenceProblemAddParemeterTatoebaSentenceNumber2,
				service.EnglishSentenceProblemAddParemeterTatoebaAuthor1,
				service.EnglishSentenceProblemAddParemeterTatoebaAuthor2} {

				if _, ok := param.GetProperties()[key]; !ok {
					return nil, xerrors.Errorf("%s is not defined. err: %w", key, lib.ErrInvalidArgument)
				}

				noteMap[key] = param.GetProperties()[key]
			}

			nodeBytes, err := json.Marshal(noteMap)
			if err != nil {
				return nil, err
			}
			note = string(nodeBytes)
		}
	}
	// audioID, err := strconv.Atoi(param.GetProperties()["audioId"])
	// if err != nil {
	// 	return nil, xerrors.Errorf("audioId is not integer. err: %w", lib.ErrInvalidArgument)
	// }

	m := &englishSentenceProblemAddParameter{
		// AudioID:    uint(audioID),
		Lang:       param.GetProperties()[service.EnglishSentenceProblemAddParemeterLang],
		Text:       param.GetProperties()[service.EnglishSentenceProblemAddParemeterText],
		Translated: param.GetProperties()[service.EnglishSentenceProblemAddParemeterTranslated],
		Note:       note,
	}
	return m, lib.Validator.Struct(m)
}

type englishSentenceProblemRepository struct {
	db          *gorm.DB
	rf          appS.AudioRepositoryFactory
	problemType string
}

func NewEnglishSentenceProblemRepository(db *gorm.DB, rf appS.AudioRepositoryFactory, problemType string) (appS.ProblemRepository, error) {
	return &englishSentenceProblemRepository{
		db:          db,
		rf:          rf,
		problemType: problemType,
	}, nil
}

func (r *englishSentenceProblemRepository) FindProblems(ctx context.Context, operator app.StudentModel, param appS.ProblemSearchCondition) (appS.ProblemSearchResult, error) {
	logger := log.FromContext(ctx)
	logger.Debugf("englishSentenceProblemRepository.FindProblems")
	limit := param.GetPageSize()
	offset := (param.GetPageNo() - 1) * param.GetPageSize()
	var problemEntities []englishSentenceProblemEntity

	where := r.db.Where("organization_id = ? and workbook_id = ?", uint(operator.GetOrganizationID()), uint(param.GetWorkbookID()))
	if result := where.Order("workbook_id, number, created_at").
		Limit(limit).Offset(offset).Find(&problemEntities); result.Error != nil {
		return nil, result.Error
	}

	problems := make([]app.ProblemModel, len(problemEntities))
	for i, e := range problemEntities {
		p, err := e.toProblem(r.rf)
		if err != nil {
			return nil, err
		}
		problems[i] = p
	}

	var count int64
	if result := where.Model(&englishSentenceProblemEntity{}).Count(&count); result.Error != nil {
		return nil, result.Error
	}

	logger.Debugf("englishSentenceProblemRepository.FindProblems, problems: %d, count: %d", len(problems), count)

	if count > math.MaxInt32 {
		return nil, errors.New("overflow")
	}

	return appS.NewProblemSearchResult(int(count), problems)
}

func (r *englishSentenceProblemRepository) FindAllProblems(ctx context.Context, operator app.StudentModel, workbookID app.WorkbookID) (appS.ProblemSearchResult, error) {
	logger := log.FromContext(ctx)
	logger.Debugf("englishSentenceProblemRepository.FindProblems")
	limit := 1000
	var problemEntities []englishSentenceProblemEntity

	where := func() *gorm.DB {
		return r.db.Where("organization_id = ? and workbook_id = ?", uint(operator.GetOrganizationID()), uint(workbookID))
	}
	if result := where().Order("workbook_id, number, created_at").
		Limit(limit).Find(&problemEntities); result.Error != nil {
		return nil, result.Error
	}

	problems := make([]app.ProblemModel, len(problemEntities))
	for i, e := range problemEntities {
		p, err := e.toProblem(r.rf)
		if err != nil {
			return nil, xerrors.Errorf("failed to toProblem. err: %w", err)
		}
		problems[i] = p
	}

	var count int64
	if result := where().Model(&englishSentenceProblemEntity{}).Count(&count); result.Error != nil {
		return nil, result.Error
	}

	logger.Debugf("englishSentenceProblemRepository.FindProblems, problems: %d, count: %d", len(problems), count)

	if count > math.MaxInt32 {
		return nil, errors.New("overflow")
	}

	return appS.NewProblemSearchResult(int(count), problems)
}

func (r *englishSentenceProblemRepository) FindProblemsByProblemIDs(ctx context.Context, operator app.StudentModel, param appS.ProblemIDsCondition) (appS.ProblemSearchResult, error) {
	var problemEntities []englishSentenceProblemEntity

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

	problems := make([]app.ProblemModel, len(problemEntities))
	for i, e := range problemEntities {
		p, err := e.toProblem(r.rf)
		if err != nil {
			return nil, err
		}
		problems[i] = p
	}

	return appS.NewProblemSearchResult(0, problems)
}

func (r *englishSentenceProblemRepository) FindProblemByID(ctx context.Context, operator app.StudentModel, id appS.ProblemSelectParameter1) (appS.Problem, error) {
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

	return problemEntity.toProblem(r.rf)
}

func (r *englishSentenceProblemRepository) FindProblemIDs(ctx context.Context, operator app.StudentModel, workbookID app.WorkbookID) ([]app.ProblemID, error) {
	pageNo := 1
	pageSize := 1000
	ids := make([]app.ProblemID, 0)
	for {
		limit := pageSize
		offset := (pageNo - 1) * pageSize
		var problemEntities []englishSentenceProblemEntity

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
			ids = append(ids, app.ProblemID(r.ID))
		}

		pageNo++
	}

	return ids, nil
}
func (r *englishSentenceProblemRepository) FindProblemsByCustomCondition(ctx context.Context, operator app.StudentModel, condition interface{}) ([]app.ProblemModel, error) {
	condition1, ok := condition.(map[string]interface{})
	if !ok {
		return nil, lib.ErrInvalidArgument
	}
	conditionWorkbookID, ok := condition1["workbookId"].(uint)
	if !ok {
		return nil, xerrors.Errorf("workbookId is not defined. err: %w", lib.ErrInvalidArgument)
	}
	conditionText, ok := condition1["text"].(string)
	if !ok {
		return nil, xerrors.Errorf("text is not defined. err: %w", lib.ErrInvalidArgument)
	}
	conditionTranslated, ok := condition1["translated"].(string)
	if !ok {
		return nil, xerrors.Errorf("translated is not defined. err: %w", lib.ErrInvalidArgument)
	}

	var problemEntity englishSentenceProblemEntity

	db := r.db.Where("organization_id = ?", uint(operator.GetOrganizationID())).
		Where("workbook_id = ?", conditionWorkbookID).
		Where("text = ?", conditionText).
		Where("translated = ?", conditionTranslated)

	if result := db.First(&problemEntity); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return []app.ProblemModel{}, nil
		}
		return nil, result.Error
	}

	problem, err := problemEntity.toProblem(r.rf)
	if err != nil {
		return nil, err
	}
	return []app.ProblemModel{problem}, nil
}

func (r *englishSentenceProblemRepository) AddProblem(ctx context.Context, operator app.StudentModel, param appS.ProblemAddParameter) (app.ProblemID, error) {
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
		Lang:           problemParam.Lang,
		Translated:     problemParam.Translated,
		Note:           problemParam.Note,
	}

	logger.Infof("englishSentenceProblemRepository.AddProblem. text: %s", problemParam.Text)
	if result := r.db.Create(&englishSentenceProblem); result.Error != nil {
		return 0, libG.ConvertDuplicatedError(result.Error, appS.ErrProblemAlreadyExists)
	}

	return app.ProblemID(englishSentenceProblem.ID), nil
}

func (r *englishSentenceProblemRepository) UpdateProblem(ctx context.Context, operator app.StudentModel, id appS.ProblemSelectParameter2, param appS.ProblemUpdateParameter) error {
	return errors.New("not implemented")
}

func (r *englishSentenceProblemRepository) RemoveProblem(ctx context.Context, operator app.StudentModel, id appS.ProblemSelectParameter2) error {
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
