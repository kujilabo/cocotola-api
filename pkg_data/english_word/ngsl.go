package english_word

import (
	"context"
	"strconv"

	"golang.org/x/xerrors"

	appD "github.com/kujilabo/cocotola-api/pkg_app/domain"
	appS "github.com/kujilabo/cocotola-api/pkg_app/service"
	"github.com/kujilabo/cocotola-api/pkg_lib/log"
	pluginCommonDomain "github.com/kujilabo/cocotola-api/pkg_plugin/common/domain"
	pluginEnglishDomain "github.com/kujilabo/cocotola-api/pkg_plugin/english/domain"
)

func CreateDemoWorkbook(ctx context.Context, studentService appS.Student) error {
	if err := CreateWorkbook(ctx, studentService, "Example", pluginCommonDomain.PosOther, []string{"butcher", "bakery", "library", "bookstore", "drugstore", "restaurant", "garage", "barbershop", "bank", "market"}); err != nil {
		return err
	}
	return nil
}

func Create20NGSLWorkbook(ctx context.Context, studentService appS.Student) error {
	if err := CreateWorkbook(ctx, studentService, "NGSL-20", pluginCommonDomain.PosOther, []string{
		"know",
		"more",
		"get",
		"who",
		"like",
		"when",
		"think",
		"make",
		"time",
		"see",
		"what",
		"up",
		"some",
		"other",
		"out",
		"good",
		"people",
		"year",
		"take",
		"no",
		"well",
		"because",
		"very",
		"just",
		"come",
		"could",
		"work",
		"use",
		"than",
		"now",
	}); err != nil {
		return err
	}
	return nil
}

func Create300NGSLWorkbook(ctx context.Context, studentService appS.Student) error {
	if err := CreateWorkbook(ctx, studentService, "NGSL-300", pluginCommonDomain.PosOther, []string{
		"know",
		"more",
		"get",
		"who",
		"like",
		"when",
		"think",
		"make",
		"time",
		"see",
		"what",
		"up",
		"some",
		"other",
		"out",
		"good",
		"people",
		"year",
		"take",
		"no",
		"well",
		"because",
		"very",
		"just",
		"come",
		"could",
		"work",
		"use",
		"than",
		"now",
		"then",
		"also",
		"into",
		"only",
		"look",
		"want",
		"give",
		"first",
		"new",
		"way",
		"find",
		"over",
		"any",
		"after",
		"day",
		"where",
		"thing",
		"most",
		"should",
		"need",
		"much",
		"right",
		"how",
		"back",
		"mean",
		"even",
		"may",
		"here",
		"many",
		"such",
		"last",
		"child",
		"tell",
		"really",
		"call",
		"before",
		"company",
		"through",
		"down",
		"show",
		"life",
		"man",
		"change",
		"place",
		"long",
		"between",
		"feel",
		"too",
		"still",
		"problem",
		"write",
		"same",
		"lot",
		"great",
		"try",
		"leave",
		"number",
		"both",
		"own",
		"part",
		"point",
		"little",
		"help",
		"ask",
		"meet",
		"start",
		"talk",
		"something",
		"put",
		"another",
		"become",
		"interest",
		"country",
		"old",
		"each",
		"school",
		"late",
		"high",
		"different",
		"off",
		"next",
		"end",
		"live",
		"why",
		"while",
		"world",
		"week",
		"play",
		"might",
		"must",
		"home",
		"never",
		"include",
		"course",
		"house",
		"report",
		"group",
		"case",
		"woman",
		"around",
		"book",
		"family",
		"seem",
		"let",
		"again",
		"kind",
		"keep",
		"hear",
		"system",
		"every",
		"question",
		"during",
		"always",
		"big",
		"set",
		"small",
		"study",
		"follow",
		"begin",
		"important",
		"since",
		"run",
		"under",
		"turn",
		"few",
		"bring",
		"early",
		"hand",
		"state",
		"move",
		"money",
		"fact",
		"however",
		"area",
		"provide",
		"name",
		"read",
		"friend",
		"month",
		"large",
		"business",
		"without",
		"information",
		"open",
		"order",
		"government",
		"word",
		"issue",
		"market",
		"pay",
		"build",
		"hold",
		"service",
		"against",
		"believe",
		"second",
		"though",
		"yes",
		"love",
		"increase",
		"job",
		"plan",
		"result",
		"away",
		"example",
		"happen",
		"offer",
		"young",
		"close",
		"program",
		"lead",
		"buy",
		"understand",
		"thank",
		"far",
		"today",
		"hour",
		"student",
		"face",
		"hope",
		"idea",
		"cost",
		"less",
		"room",
		"until",
		"reason",
		"form",
		"spend",
		"head",
		"car",
		"learn",
		"level",
		"person",
		"experience",
		"once",
		"member",
		"enough",
		"bad",
		"city",
		"night",
		"able",
		"support",
		"whether",
		"line",
		"present",
		"side",
		"quite",
		"although",
		"sure",
		"term",
		"least",
		"age",
		"low",
		"speak",
		"within",
		"process",
		"public",
		"often",
		"train",
		"possible",
		"actually",
		"rather",
		"view",
		"together",
		"consider",
		"price",
		"parent",
		"hard",
		"party",
		"local",
		"control",
		"already",
		"concern",
		"product",
		"lose",
		"story",
		"almost",
		"continue",
		"stand",
		"whole",
		"yet",
		"rate",
		"care",
		"expect",
		"effect",
		"sort",
		"ever",
		"anything",
		"cause",
		"fall",
		"deal",
		"water",
		"send",
		"allow",
		"soon",
		"watch",
		"base",
		"probably",
		"suggest",
		"past",
		"power",
		"test",
		"visit",
		"center",
		"grow",
		"nothing",
		"return",
		"mother",
		"walk",
		"matter",
	}); err != nil {
		return err
	}
	return nil
}
func CreateWorkbook(ctx context.Context, student appS.Student, workbookName string, pos pluginCommonDomain.WordPos, words []string) error {
	logger := log.FromContext(ctx)

	workbookProperties := map[string]string{
		"audioEnabled": "false",
	}
	param, err := appS.NewWorkbookAddParameter(pluginEnglishDomain.EnglishWordProblemType, workbookName, appD.Lang2JA, "", workbookProperties)
	if err != nil {
		return xerrors.Errorf("failed to NewWorkbookAddParameter. err: %w", err)
	}

	workbookID, err := student.AddWorkbookToPersonalSpace(ctx, param)
	if err != nil {
		return xerrors.Errorf("failed to AddWorkbookToPersonalSpace. err: %w", err)
	}

	workbook, err := student.FindWorkbookByID(ctx, workbookID)
	if err != nil {
		return xerrors.Errorf("failed to FindWorkbookByID. err: %w", err)
	}

	for i, word := range words {
		properties := map[string]string{
			"text": word,
			"lang": "ja",
			"pos":  strconv.Itoa(int(pos)),
		}
		param, err := appS.NewProblemAddParameter(workbookID, i+1, properties)
		if err != nil {
			return xerrors.Errorf("failed to NewProblemAddParameter. err: %w", err)
		}

		problemIDs, err := workbook.AddProblem(ctx, student, param)
		if err != nil {
			return xerrors.Errorf("failed to NewProblemAddParameter. err: %w", err)
		}
		logger.Infof("problemIDs: %v", problemIDs)
	}

	logger.Infof("Example %d", workbookID)
	return nil
}
