package main

import (
	"context"
	"database/sql"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"os"
	"path"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"github.com/kujilabo/cocotola-api/pkg_app/config"
	app "github.com/kujilabo/cocotola-api/pkg_app/domain"
	appD "github.com/kujilabo/cocotola-api/pkg_app/domain"
	appG "github.com/kujilabo/cocotola-api/pkg_app/gateway"
	libD "github.com/kujilabo/cocotola-api/pkg_lib/domain"
	libG "github.com/kujilabo/cocotola-api/pkg_lib/gateway"
	pluginCommonGateway "github.com/kujilabo/cocotola-api/pkg_plugin/common/gateway"
	pluginEnglishDomain "github.com/kujilabo/cocotola-api/pkg_plugin/english/domain"
	pluginEnglishGateway "github.com/kujilabo/cocotola-api/pkg_plugin/english/gateway"
	userD "github.com/kujilabo/cocotola-api/pkg_user/domain"
	userG "github.com/kujilabo/cocotola-api/pkg_user/gateway"
)

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
	return path.Join(wd, "tools", "import1")
}

func registerEnglishPhraseProblemsFlushSentence(ctx context.Context, operator app.Student, repo app.RepositoryFactory, processor app.ProblemAddProcessor) error {

	fmt.Println("registerEnglishPhraseProblems")
	csvFilePath := importDir() + "/flush_sentence.csv"
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
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		if len(line) != 2 {
			return fmt.Errorf("invalid umber of column. row: %d", i)
		}
		i++
	}

	// workbookRepo, err := repo.NewWorkbookRepository(ctx)
	// if err != nil {
	// 	return err
	// }

	workbookName := "flush sentence"
	workbook, err := operator.FindWorkbookByName(ctx, workbookName)
	if err != nil {
		if err != app.ErrWorkbookNotFound {
			return err
		}
		param, err := app.NewWorkbookAddParameter(pluginEnglishDomain.EnglishPhraseProblemType, workbookName, "")
		if err != nil {
			return err
		}

		workbookID, err := operator.AddWorkbookToPersonalSpace(ctx, param)
		if err != nil {
			return err
		}

		wb2, err := operator.FindWorkbookByID(ctx, workbookID)
		if err != nil {
			return err
		}
		workbook = wb2
	}

	searchCondition, err := app.NewProblemSearchCondition(app.WorkbookID(workbook.GetID()), 1, 1000, "")
	if err != nil {
		return err
	}

	problems, err := workbook.FindProblems(ctx, operator, searchCondition)
	if err != nil {
		return err
	}
	problemMap := make(map[string]app.Problem)
	for _, p := range problems.Results {
		m := p.GetProperties(ctx)
		textObj, ok := m["text"]
		if !ok {
			return fmt.Errorf("text not found. %v", m)
		}
		text, ok := textObj.(string)
		if !ok {
			return fmt.Errorf("text is not string. %v", m)
		}

		problemMap[text] = p
	}
	{
		file, err := os.Open(csvFilePath)
		if err != nil {
			panic(err)
		}
		defer file.Close()

		reader := csv.NewReader(file)
		var i = 0
		for {
			var line []string
			line, err = reader.Read()
			if err != nil {
				break
			}

			if line[0] == "#" {
				continue
			}

			properties := map[string]string{
				"text":       line[0],
				"translated": line[1],
			}
			param, err := app.NewProblemAddParameter(app.WorkbookID(workbook.GetID()), 1, workbook.GetProblemType(), properties)
			if err != nil {
				return err
			}

			if _, ok := problemMap[line[0]]; !ok {
				processor.AddProblem(ctx, repo, operator, param)
			}

			i++
		}
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

	englishPhraseProblemProcessor := pluginEnglishDomain.NewEnglishPhraseProblemProcessor(synthesizer, translator)
	problemAddProcessor := map[string]appD.ProblemAddProcessor{
		pluginEnglishDomain.EnglishPhraseProblemType: englishPhraseProblemProcessor,
	}
	problemRemoveProcessor := map[string]appD.ProblemRemoveProcessor{
		pluginEnglishDomain.EnglishPhraseProblemType: englishPhraseProblemProcessor,
	}
	problemImportProcessor := map[string]appD.ProblemImportProcessor{}
	englishPhraseProblemRepository := func(db *gorm.DB) (appD.ProblemRepository, error) {
		return pluginEnglishGateway.NewEnglishPhraseProblemRepository(db, pluginEnglishDomain.EnglishPhraseProblemType)
	}

	pf := appD.NewProcessorFactory(problemAddProcessor, problemRemoveProcessor, problemImportProcessor)
	problemRepositories := map[string]func(*gorm.DB) (appD.ProblemRepository, error){
		pluginEnglishDomain.EnglishPhraseProblemType: englishPhraseProblemRepository,
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

	if err := registerEnglishPhraseProblemsFlushSentence(ctx, student, repoFunc(db), englishPhraseProblemProcessor); err != nil {
		panic(err)
	}
}
