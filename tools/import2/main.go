package main

import (
	"context"
	"database/sql"
	"encoding/csv"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path"
	"strconv"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"github.com/kujilabo/cocotola-api/pkg_app/config"
	appD "github.com/kujilabo/cocotola-api/pkg_app/domain"
	appG "github.com/kujilabo/cocotola-api/pkg_app/gateway"
	libD "github.com/kujilabo/cocotola-api/pkg_lib/domain"
	libG "github.com/kujilabo/cocotola-api/pkg_lib/gateway"
	"github.com/kujilabo/cocotola-api/pkg_plugin/common/domain"
	pluginCommonGateway "github.com/kujilabo/cocotola-api/pkg_plugin/common/gateway"
	pluginEnglishDomain "github.com/kujilabo/cocotola-api/pkg_plugin/english/domain"
	pluginEnglishGateway "github.com/kujilabo/cocotola-api/pkg_plugin/english/gateway"
	userD "github.com/kujilabo/cocotola-api/pkg_user/domain"
	userG "github.com/kujilabo/cocotola-api/pkg_user/gateway"
)

var defaultPageNo = 1
var defaultPageSize = 1000
var columnLength = 3

func initDB(cfg *config.DBConfig) (*gorm.DB, *sql.DB, error) {
	switch cfg.DriverName {
	case "sqlite3":
		db, err := libG.OpenSQLite("./" + cfg.SQLite3.File)
		if err != nil {
			return nil, nil, err
		}

		sqlDB, err := db.DB()
		if err != nil {
			return nil, nil, err
		}

		if err := sqlDB.Ping(); err != nil {
			return nil, nil, err
		}

		if err := appG.MigrateSQLiteDB(db); err != nil {
			return nil, nil, err
		}

		return db, sqlDB, nil
	default:
		return nil, nil, libD.ErrInvalidArgument
	}
}

func importDir() string {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	return path.Join(wd, "tools", "import2")
}

func checkFile(csvFilePath string) error {
	file, err := os.Open(csvFilePath)
	if err != nil {
		return err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	var i = 1
	for {
		var line []string
		line, err = reader.Read()
		if errors.Is(err, io.EOF) {
			break
		}

		if err != nil {
			return err
		}

		if len(line) != columnLength {
			return fmt.Errorf("invalid umber of column. row: %d", i)
		}

		i++
	}

	return nil
}

func initWorkbook(ctx context.Context, operator appD.Student, workbookName string) (appD.Workbook, error) {
	workbook, err := operator.FindWorkbookByName(ctx, workbookName)
	if err != nil {
		if !errors.Is(err, appD.ErrWorkbookNotFound) {
			return nil, err
		}

		param, err := appD.NewWorkbookAddParameter(pluginEnglishDomain.EnglishWordProblemType, workbookName, "")
		if err != nil {
			return nil, err
		}

		workbookID, err := operator.AddWorkbookToPersonalSpace(ctx, param)
		if err != nil {
			return nil, err
		}

		wb2, err := operator.FindWorkbookByID(ctx, workbookID)
		if err != nil {
			return nil, err
		}

		workbook = wb2
	}

	return workbook, nil
}

func initProblems(ctx context.Context, operator appD.Student, workbook appD.Workbook) (map[string]bool, error) {
	searchCondition, err := appD.NewProblemSearchCondition(appD.WorkbookID(workbook.GetID()), defaultPageNo, defaultPageSize, "")
	if err != nil {
		return nil, err
	}

	problems, err := workbook.FindProblems(ctx, operator, searchCondition)
	if err != nil {
		return nil, err
	}

	problemMap := make(map[string]bool)
	for _, p := range problems.Results {
		m := p.GetProperties(ctx)
		textObj, ok := m["text"]
		if !ok {
			return nil, fmt.Errorf("text not found. %v", m)
		}
		text, ok := textObj.(string)
		if !ok {
			return nil, fmt.Errorf("text is not string. %v", m)
		}

		problemMap[text] = true
	}

	return problemMap, nil
}

func registerEnglishWordProblems(ctx context.Context, operator appD.Student, repo appD.RepositoryFactory, processor appD.ProblemAddProcessor) error {

	fmt.Println("registerEnglishWordProblems")
	csvFilePath := importDir() + "/kikutan.csv"

	if err := checkFile(csvFilePath); err != nil {
		return err
	}

	workbookName := "kikutan"
	workbook, err := initWorkbook(ctx, operator, workbookName)
	if err != nil {
		return err
	}

	problemMap, err := initProblems(ctx, operator, workbook)
	if err != nil {
		return err
	}
	file, err := os.Open(csvFilePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	var i = 1
	for {
		var line []string
		line, err = reader.Read()
		if err != nil {
			break
		}

		if line[0] == "#" {
			continue
		}

		pos, err := domain.ParsePos(line[0])
		if err != nil {
			fmt.Println("parsePos")
			return err
		}

		properties := map[string]string{
			"lang":       "ja",
			"pos":        strconv.Itoa(int(pos)),
			"text":       line[1],
			"translated": line[2],
		}
		param, err := appD.NewProblemAddParameter(appD.WorkbookID(workbook.GetID()), i, workbook.GetProblemType(), properties)
		if err != nil {
			return err
		}

		if _, ok := problemMap[line[0]]; !ok {
			if _, err := processor.AddProblem(ctx, repo, operator, param); err != nil {
				return err
			}
		}

		i++
	}
	return nil
}

func main() {
	ctx := context.Background()

	env := flag.String("env", "", "environment")
	flag.Parse()
	if len(*env) == 0 {
		appEnv := os.Getenv("APP_ENV")
		if len(appEnv) == 0 {
			*env = "development"
		} else {
			*env = appEnv
		}
	}

	logrus.Infof("env: %s", *env)

	cfg, err := config.LoadConfig(*env)
	if err != nil {
		panic(err)
	}

	// init log
	if err := config.InitLog(*env, cfg.Log); err != nil {
		panic(err)
	}

	// init db
	db, sqlDB, err := initDB(cfg.DB)
	if err != nil {
		fmt.Printf("failed to InitDB. err: %+v", err)
		panic(err)
	}
	defer sqlDB.Close()

	rf := userG.NewRepositoryFactory(db)
	userD.InitSystemAdmin(rf)
	systemAdmin := userD.SystemAdminInstance()
	systemOwner, err := systemAdmin.FindSystemOwnerByOrganizationName(ctx, "cocotola")
	if err != nil {
		panic(err)
	}

	synthesizer := pluginCommonGateway.NewSynthesizer(cfg.Google.SynthesizerKey, time.Duration(cfg.Google.SynthesizerTimeoutSec)*time.Minute)

	translatorClient := pluginCommonGateway.NewAzureTranslatorClient(cfg.Azure.SubscriptionKey)
	translatorRepository := pluginCommonGateway.NewAzureTranslationRepository(db)
	translator := pluginCommonGateway.NewAzureCachedTranslatorClient(translatorClient, translatorRepository)

	englishWordProblemProcessor := pluginEnglishDomain.NewEnglishWordProblemProcessor(synthesizer, translator, pluginEnglishGateway.NewEnglishWordProblemAddParameterCSVReader)
	problemAddProcessor := map[string]appD.ProblemAddProcessor{
		pluginEnglishDomain.EnglishWordProblemType: englishWordProblemProcessor,
	}
	problemRemoveProcessor := map[string]appD.ProblemRemoveProcessor{
		pluginEnglishDomain.EnglishWordProblemType: englishWordProblemProcessor,
	}
	problemImportProcessor := map[string]appD.ProblemImportProcessor{
		pluginEnglishDomain.EnglishWordProblemType: englishWordProblemProcessor,
	}
	englishWordProblemRepository := func(db *gorm.DB) (appD.ProblemRepository, error) {
		return pluginEnglishGateway.NewEnglishWordProblemRepository(db, pluginEnglishDomain.EnglishWordProblemType)
	}

	pf := appD.NewProcessorFactory(problemAddProcessor, problemRemoveProcessor, problemImportProcessor)
	problemRepositories := map[string]func(*gorm.DB) (appD.ProblemRepository, error){
		pluginEnglishDomain.EnglishWordProblemType: englishWordProblemRepository,
	}

	userRepoFunc := func(db *gorm.DB) userD.RepositoryFactory {
		return userG.NewRepositoryFactory(db)
	}
	repoFunc := func(db *gorm.DB) appD.RepositoryFactory {
		return appG.NewRepositoryFactory(db, cfg.DB.DriverName, userRepoFunc, pf, problemRepositories)
	}

	appUser, err := userRepoFunc(db).NewAppUserRepository().FindAppUserByLoginID(ctx, systemOwner, cfg.App.TestUserEmail)
	student, err := appD.NewStudent(repoFunc(db), userRepoFunc(db), appUser)
	if err != nil {
		panic(err)
	}

	if err := registerEnglishWordProblems(ctx, student, repoFunc(db), englishWordProblemProcessor); err != nil {
		panic(err)
	}
}
