package gateway

import (
	"context"
	"math/rand"
	"strings"
	"time"

	"golang.org/x/xerrors"
	"gorm.io/gorm"

	app "github.com/kujilabo/cocotola-api/pkg_app/domain"
	libG "github.com/kujilabo/cocotola-api/pkg_lib/gateway"
	"github.com/kujilabo/cocotola-api/pkg_lib/log"
	"github.com/kujilabo/cocotola-api/pkg_plugin/common/domain"
)

const (
	shuffleBufferRate = 10
)

type tatoebaSentenceRepository struct {
	db *gorm.DB
}

type tatoebaSentenceEntity struct {
	SentenceNumber int
	Lang           string
	Text           string
	Author         string
	UpdatedAt      time.Time
}
type tatoebaSentencePairEntity struct {
	SrcSentenceNumber int
	SrcLang           string
	SrcText           string
	SrcAuthor         string
	SrcUpdatedAt      time.Time
	DstSentenceNumber int
	DstLang           string
	DstText           string
	DstAuthor         string
	DstUpdatedAt      time.Time
}

func (e *tatoebaSentenceEntity) toModel() (domain.TatoebaSentence, error) {
	lang, err := app.NewLang3(e.Lang)
	if err != nil {
		return nil, xerrors.Errorf("failed to NewLang3. err: %w", err)
	}
	author := e.Author
	if author == "\\N" {
		author = ""
	}
	return domain.NewTatoebaSentence(e.SentenceNumber, lang, e.Text, author, e.UpdatedAt)
}

func (e *tatoebaSentencePairEntity) toModel() (domain.TatoebaSentencePair, error) {
	srcE := tatoebaSentenceEntity{
		SentenceNumber: e.SrcSentenceNumber,
		Lang:           e.SrcLang,
		Text:           e.SrcText,
		Author:         e.SrcAuthor,
		UpdatedAt:      e.SrcUpdatedAt,
	}
	srcM, err := srcE.toModel()
	if err != nil {
		return nil, err
	}

	dstE := tatoebaSentenceEntity{
		SentenceNumber: e.DstSentenceNumber,
		Lang:           e.DstLang,
		Text:           e.DstText,
		Author:         e.DstAuthor,
		UpdatedAt:      e.DstUpdatedAt,
	}
	dstM, err := dstE.toModel()
	if err != nil {
		return nil, err
	}

	return domain.NewTatoebaSentencePair(srcM, dstM)
}

func (e *tatoebaSentenceEntity) TableName() string {
	return "tatoeba_sentence"
}

func NewTatoebaSentenceRepository(db *gorm.DB) domain.TatoebaSentenceRepository {
	return &tatoebaSentenceRepository{
		db: db,
	}
}

// func (r *tatoebaSentenceRepository) FindTatoebaSentences(ctx context.Context, param domain.TatoebaSentenceSearchCondition) (*domain.TatoebaSentenceSearchResult, error) {
// 	logger := log.FromContext(ctx)
// 	logger.Debug("tatoebaSentenceRepository.FindTatoebaSentences")
// 	limit := param.GetPageSize()
// 	offset := (param.GetPageNo() - 1) * param.GetPageSize()

// 	where := func() *gorm.DB {
// 		db := r.db.Where("lang = 'eng'")
// 		if param.GetKeyword() != "" {
// 			keyword := "%" + param.GetKeyword() + "%"
// 			db = db.Where("text like ?", keyword)
// 		}
// 		return db
// 	}

// 	entities := []tatoebaSentenceEntity{}
// 	if result := where().Limit(limit).Offset(offset).Find(&entities); result.Error != nil {
// 		return nil, result.Error
// 	}

// 	results := make([]domain.TatoebaSentence, 0)
// 	for _, e := range entities {
// 		m, err := e.toModel()
// 		if err != nil {
// 			return nil, err
// 		}
// 		results = append(results, m)
// 	}

// 	var count int64
// 	if result := where().Model(&azureTranslationEntity{}).Count(&count); result.Error != nil {
// 		return nil, result.Error
// 	}

// 	return &domain.TatoebaSentenceSearchResult{
// 		TotalCount: count,
// 		Results:    results,
// 	}, nil
// }
//SELECT *
// FROM development.tatoeba_sentence t1
// inner join development.tatoeba_link t2
// on t1.sentence_number= t2.`from`

// inner join development.tatoeba_sentence t3
// on t3.sentence_number= t2.`to`

// where t1.lang='eng' and t3.lang='jpn';

func (r *tatoebaSentenceRepository) FindTatoebaSentences(ctx context.Context, param domain.TatoebaSentenceSearchCondition) (*domain.TatoebaSentenceSearchResult, error) {
	if param.IsRandom() {
		return r.findTatoebaSentencesByRandom(ctx, param)
	}
	return r.findTatoebaSentences(ctx, param)
}

func (r *tatoebaSentenceRepository) findTatoebaSentences(ctx context.Context, param domain.TatoebaSentenceSearchCondition) (*domain.TatoebaSentenceSearchResult, error) {
	logger := log.FromContext(ctx)
	logger.Debug("tatoebaSentenceRepository.FindTatoebaSentences")
	limit := param.GetPageSize()
	offset := (param.GetPageNo() - 1) * param.GetPageSize()

	//db.Model(&User{}).Select("users.name, emails.email").Joins("left join emails on emails.user_id = users.id").Scan(&result{})

	// 	FROM `sandbox` AS s
	// INNER JOIN (
	//   SELECT CEIL(RAND() * (SELECT MAX(`id`) FROM `sandbox`)) AS `id`
	// ) AS `tmp` ON s.id >= tmp.id
	// ORDER BY s.id

	where := func() *gorm.DB {
		db := r.db.Table("tatoeba_sentence AS T1").Select(
			// Src
			"T1.sentence_number AS src_sentence_number," +
				"T1.lang AS src_lang," +
				"T1.text AS src_text," +
				"T1.author AS src_author," +
				"T1.updated_at AS src_updated_at," +
				// Dst
				"T3.sentence_number AS dst_sentence_number," +
				"T3.lang AS dst_lang," +
				"T3.text AS dst_text," +
				"T3.author AS dst_author," +
				"T3.updated_at AS dst_updated_at").
			Joins("INNER JOIN tatoeba_link AS T2 ON T1.sentence_number = T2.`from`").
			Joins("INNER JOIN tatoeba_sentence AS T3 ON T3.sentence_number = T2.`to`").
			Where("T1.lang = 'eng' AND T3.lang = 'jpn'")
		if param.GetKeyword() != "" {
			keyword1 := strings.ReplaceAll(param.GetKeyword(), "%", "\\%")
			keyword2 := "%" + keyword1 + "%"
			db = db.Where("T1.text like ?", keyword2)
		}
		return db
	}

	entities := []tatoebaSentencePairEntity{}
	if result := where().Limit(limit).Offset(offset).Scan(&entities); result.Error != nil {
		return nil, result.Error
	}

	results := make([]domain.TatoebaSentencePair, len(entities))
	for i, e := range entities {
		m, err := e.toModel()
		if err != nil {
			return nil, err
		}
		results[i] = m
	}

	var count int64 = 0
	// if result := where().Count(&count); result.Error != nil {
	// 	return nil, result.Error
	// }

	return &domain.TatoebaSentenceSearchResult{
		TotalCount: count,
		Results:    results,
	}, nil
}

func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

func (r *tatoebaSentenceRepository) findTatoebaSentencesByRandom(ctx context.Context, param domain.TatoebaSentenceSearchCondition) (*domain.TatoebaSentenceSearchResult, error) {
	logger := log.FromContext(ctx)
	logger.Debug("tatoebaSentenceRepository.FindTatoebaSentences")
	limit := param.GetPageSize() * shuffleBufferRate
	offset := (param.GetPageNo() - 1) * param.GetPageSize()

	where := func() *gorm.DB {
		db := r.db.Table("tatoeba_sentence AS T1").Select(
			// Src
			"T1.sentence_number AS src_sentence_number," +
				"T1.lang AS src_lang," +
				"T1.text AS src_text," +
				"T1.author AS src_author," +
				"T1.updated_at AS src_updated_at," +
				// Dst
				"T3.sentence_number AS dst_sentence_number," +
				"T3.lang AS dst_lang," +
				"T3.text AS dst_text," +
				"T3.author AS dst_author," +
				"T3.updated_at AS dst_updated_at").
			Joins("INNER JOIN tatoeba_link AS T2 ON T1.sentence_number = T2.`from`").
			Joins("INNER JOIN tatoeba_sentence AS T3 ON T3.sentence_number = T2.`to`").
			Joins("INNER JOIN (SELECT CEIL(RAND() * (SELECT MAX(`sentence_number`) FROM `tatoeba_sentence`)) AS `sentence_number`) AS `tmp` ON T1.sentence_number >= tmp.sentence_number").
			Where("T1.lang = 'eng' AND T3.lang = 'jpn'")
		if param.GetKeyword() != "" {
			keyword1 := strings.ReplaceAll(param.GetKeyword(), "%", "\\%")
			keyword2 := "%" + keyword1 + "%"
			db = db.Where("T1.text like ?", keyword2)
		}
		return db
	}

	entities := []tatoebaSentencePairEntity{}
	if result := where().Limit(limit).Offset(offset).Scan(&entities); result.Error != nil {
		return nil, result.Error
	}

	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(entities), func(i, j int) { entities[i], entities[j] = entities[j], entities[i] })

	logger.Infof("len(entities): %d", len(entities))

	length := min(param.GetPageSize(), len(entities))
	results := make([]domain.TatoebaSentencePair, length)
	for i := 0; i < length; i++ {
		m, err := entities[i].toModel()
		if err != nil {
			return nil, err
		}
		results[i] = m
	}

	var count int64 = 0
	// if result := where().Count(&count); result.Error != nil {
	// 	return nil, result.Error
	// }

	return &domain.TatoebaSentenceSearchResult{
		TotalCount: count,
		Results:    results,
	}, nil
}

func (r *tatoebaSentenceRepository) FindTatoebaSentenceBySentenceNumber(ctx context.Context, sentenceNumber int) (domain.TatoebaSentence, error) {
	entity := &tatoebaSentenceEntity{}
	if result := r.db.Where("sentence_number = ?", r.db.Statement.TableExpr.Vars...).Find(&entity); result.Error != nil {
		return nil, result.Error
	}

	sentence, err := entity.toModel()
	if err != nil {
		return nil, err
	}

	return sentence, nil
}

func (r *tatoebaSentenceRepository) Add(ctx context.Context, param domain.TatoebaSentenceAddParameter) error {
	entity := tatoebaSentenceEntity{
		SentenceNumber: param.GetSentenceNumber(),
		Lang:           param.GetLang().String(),
		Text:           param.GetText(),
		Author:         param.GetAuthor(),
		UpdatedAt:      param.GetUpdatedAt(),
	}

	if result := r.db.Create(&entity); result.Error != nil {
		err := libG.ConvertDuplicatedError(result.Error, domain.ErrTatoebaSentenceAlreadyExists)
		return xerrors.Errorf("failed to Add tatoebaSentence. err: %w", err)
	}

	return nil
}
