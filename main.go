package main

import (
	"context"
	"database/sql"
	"errors"
	"flag"
	"io"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	ginlog "github.com/onrik/logrus/gin"
	"github.com/sirupsen/logrus"
	"golang.org/x/xerrors"
	"gorm.io/gorm"

	swaggerFiles "github.com/swaggo/files"     // swagger embed files
	ginSwagger "github.com/swaggo/gin-swagger" // gin-swagger middleware

	"github.com/kujilabo/cocotola-api/docs"
	"github.com/kujilabo/cocotola-api/pkg_app/application"
	"github.com/kujilabo/cocotola-api/pkg_app/config"
	appD "github.com/kujilabo/cocotola-api/pkg_app/domain"
	appG "github.com/kujilabo/cocotola-api/pkg_app/gateway"
	appH "github.com/kujilabo/cocotola-api/pkg_app/handler"
	authA "github.com/kujilabo/cocotola-api/pkg_auth/application"
	authG "github.com/kujilabo/cocotola-api/pkg_auth/gateway"
	authH "github.com/kujilabo/cocotola-api/pkg_auth/handler"
	authM "github.com/kujilabo/cocotola-api/pkg_auth/handler/middleware"
	english_word "github.com/kujilabo/cocotola-api/pkg_data/english_word"
	libD "github.com/kujilabo/cocotola-api/pkg_lib/domain"
	libG "github.com/kujilabo/cocotola-api/pkg_lib/gateway"
	"github.com/kujilabo/cocotola-api/pkg_lib/handler/middleware"
	"github.com/kujilabo/cocotola-api/pkg_lib/log"
	pluginApplication "github.com/kujilabo/cocotola-api/pkg_plugin/common/application"
	pluginCommonDomain "github.com/kujilabo/cocotola-api/pkg_plugin/common/domain"
	pluginCommonGateway "github.com/kujilabo/cocotola-api/pkg_plugin/common/gateway"
	pluginCommonHandler "github.com/kujilabo/cocotola-api/pkg_plugin/common/handler"
	pluginEnglishDomain "github.com/kujilabo/cocotola-api/pkg_plugin/english/domain"
	pluginEnglishGateway "github.com/kujilabo/cocotola-api/pkg_plugin/english/gateway"
	userD "github.com/kujilabo/cocotola-api/pkg_user/domain"
	userG "github.com/kujilabo/cocotola-api/pkg_user/gateway"
)

func main() {
	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	ctx := context.Background()
	env := flag.String("env", "", "environment")
	flag.Parse()
	if len(*env) == 0 {
		appEnv := os.Getenv("APP_ENV")
		if len(appEnv) == 0 {
			*env = "local"
		} else {
			*env = appEnv
		}
	}

	logrus.Infof("env: %s", *env)

	go func() {
		sig := <-sigs
		logrus.Info()
		logrus.Info(sig)
		done <- true
	}()

	cfg, db, sqlDB, router, err := initialize(ctx, *env)
	if err != nil {
		panic(err)
	}
	defer sqlDB.Close()

	if err := initApp1(ctx, db, cfg.App.OwnerPassword); err != nil {
		panic(err)
	}

	router.GET("/healthcheck", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	synthesizer := pluginCommonGateway.NewSynthesizer(cfg.Google.SynthesizerKey, time.Duration(cfg.Google.SynthesizerTimeoutSec)*time.Minute)

	azureTranslationClient := pluginCommonGateway.NewAzureTranslationClient(cfg.Azure.SubscriptionKey)

	pluginRepo, err := pluginCommonGateway.NewRepositoryFactory(context.Background(), db, cfg.DB.DriverName)
	if err != nil {
		panic(err)
	}
	pluginRfFunc := func(db *gorm.DB) (pluginCommonDomain.RepositoryFactory, error) {
		return pluginCommonGateway.NewRepositoryFactory(context.Background(), db, cfg.DB.DriverName)
	}

	translator, err := pluginCommonDomain.NewTranslatior(pluginRepo, azureTranslationClient)
	if err != nil {
		panic(err)
	}

	tatoebaSentenceRepo, err := pluginRepo.NewTatoebaSentenceRepository(ctx)
	if err != nil {
		panic(err)
	}

	pf, problemRepositories, problemImportProcessor := initPf(synthesizer, translator, tatoebaSentenceRepo)

	newIterator := func(ctx context.Context, workbookID appD.WorkbookID, problemType string, reader io.Reader) (appD.ProblemAddParameterIterator, error) {
		processor, ok := problemImportProcessor[problemType]
		if ok {
			return processor.CreateCSVReader(ctx, workbookID, reader)
		}
		return nil, xerrors.Errorf("processor not found. problemType: %s", problemType)
	}

	userRfFunc := func(db *gorm.DB) (userD.RepositoryFactory, error) {
		return userG.NewRepositoryFactory(db)
	}
	appD.UserRfFunc = userRfFunc

	rfFunc := func(db *gorm.DB) (appD.RepositoryFactory, error) {
		return appG.NewRepositoryFactory(context.Background(), db, cfg.DB.DriverName, userRfFunc, pf, problemRepositories)
	}
	appD.RfFunc = rfFunc

	if err := initApp2(ctx, db, rfFunc, userRfFunc); err != nil {
		panic(err)
	}

	signingKey := []byte(cfg.Auth.SigningKey)
	signingMethod := jwt.SigningMethodHS256
	authTokenManager := authG.NewAuthTokenManager(signingKey, signingMethod, time.Duration(cfg.Auth.AccessTokenTTLMin)*time.Minute, time.Duration(cfg.Auth.RefreshTokenTTLHour)*time.Hour)

	googleAuthClient := authG.NewGoogleAuthClient(cfg.Auth.GoogleClientID, cfg.Auth.GoogleClientSecret, cfg.Auth.GoogleCallbackURL)
	authMiddleware := authM.NewAuthMiddleware(signingKey)

	registerAppUserCallback := func(ctx context.Context, organizationName string, appUser userD.AppUser) error {
		rf, err := rfFunc(db)
		if err != nil {
			return err
		}
		userRf, err := userRfFunc(db)
		if err != nil {
			return err
		}
		return callback(ctx, cfg.App.TestUserEmail, pf, rf, userRf, organizationName, appUser)
	}

	v1 := router.Group("v1")
	{
		v1auth := v1.Group("auth")
		googleAuthService := authA.NewGoogleAuthService(db, googleAuthClient, authTokenManager, registerAppUserCallback)
		guestAuthService := authA.NewGuestAuthService(authTokenManager)
		authHandler := authH.NewAuthHandler(authTokenManager)
		googleAuthHandler := authH.NewGoogleAuthHandler(googleAuthService)
		guestAuthHandler := authH.NewGuestAuthHandler(guestAuthService)
		v1auth.POST("google/authorize", googleAuthHandler.Authorize)
		v1auth.POST("guest/authorize", guestAuthHandler.Authorize)
		v1auth.POST("refresh_token", authHandler.RefreshToken)

		privateWorkbookService := application.NewPrivateWorkbookService(db, pf, rfFunc, userRfFunc)
		privateWorkbookHandler := appH.NewPrivateWorkbookHandler(privateWorkbookService)
		v1Workbook := v1.Group("private/workbook")
		v1Workbook.Use(authMiddleware)
		v1Workbook.POST(":workbookID", privateWorkbookHandler.FindWorkbooks)
		v1Workbook.GET(":workbookID", privateWorkbookHandler.FindWorkbookByID)
		v1Workbook.PUT(":workbookID", privateWorkbookHandler.UpdateWorkbook)
		v1Workbook.DELETE(":workbookID", privateWorkbookHandler.RemoveWorkbook)
		v1Workbook.POST("", privateWorkbookHandler.AddWorkbook)

		problemService := application.NewProblemService(db, pf, rfFunc, userRfFunc)
		problemHandler := appH.NewProblemHandler(problemService, newIterator)
		v1Problem := v1.Group("workbook/:workbookID")
		v1Problem.Use(authMiddleware)
		v1Problem.POST("problem", problemHandler.AddProblem)
		v1Problem.GET("problem/:problemID", problemHandler.FindProblemByID)
		v1Problem.DELETE("problem/:problemID", problemHandler.RemoveProblem)
		v1Problem.PUT("problem/:problemID", problemHandler.UpdateProblem)
		// v1Problem.GET("problem_ids", problemHandler.FindProblemIDs)
		v1Problem.POST("problem/find", problemHandler.FindProblems)
		v1Problem.POST("problem/find_all", problemHandler.FindAllProblems)
		v1Problem.POST("problem/find_by_ids", problemHandler.FindProblemsByProblemIDs)
		v1Problem.POST("problem/import", problemHandler.ImportProblems)

		studyService := application.NewStudyService(db, pf, rfFunc, userRfFunc)
		recordbookHandler := appH.NewRecordbookHandler(studyService)
		v1Study := v1.Group("study/workbook/:workbookID")
		v1Study.Use(authMiddleware)
		v1Study.GET("study_type/:studyType", recordbookHandler.FindRecordbook)
		v1Study.POST("study_type/:studyType/problem/:problemID/record", recordbookHandler.SetStudyResult)

		audioService := application.NewAudioService(db, rfFunc)
		audioHandler := appH.NewAudioHandler(audioService)
		v1Audio := v1.Group("audio")
		v1Audio.Use(authMiddleware)
		v1Audio.GET(":audioID", audioHandler.FindAudioByID)
	}

	plugin := router.Group("plugin")
	{
		plugin.Use(authMiddleware)
		{
			pluginTranslation := plugin.Group("translation")
			translationHandler := pluginCommonHandler.NewTranslationHandler(translator)
			pluginTranslation.POST("find", translationHandler.FindTranslations)
			pluginTranslation.GET("text/:text/pos/:pos", translationHandler.FindTranslationByTextAndPos)
			pluginTranslation.GET("text/:text", translationHandler.FindTranslationByText)
			pluginTranslation.PUT("text/:text/pos/:pos", translationHandler.UpdateTranslation)
			pluginTranslation.DELETE("text/:text/pos/:pos", translationHandler.RemoveTranslation)
			pluginTranslation.POST("", translationHandler.AddTranslation)
			pluginTranslation.POST("export", translationHandler.ExportTranslations)
		}
		{
			newSentenceReader := func(reader io.Reader) pluginCommonDomain.TatoebaSentenceAddParameterIterator {
				return pluginCommonGateway.NewTatoebaSentenceAddParameterReader(reader)
			}
			newLinkReader := func(reader io.Reader) pluginCommonDomain.TatoebaLinkAddParameterIterator {
				return pluginCommonGateway.NewTatoebaLinkAddParameterReader(reader)
			}
			pluginTatoeba := plugin.Group("tatoeba")
			tatoebaService := pluginApplication.NewTatoebaService(db, pluginRfFunc)
			tatoebaHandler := pluginCommonHandler.NewTatoebaHandler(tatoebaService, newSentenceReader, newLinkReader)
			pluginTatoeba.POST("find", tatoebaHandler.FindSentences)
			pluginTatoeba.POST("sentence/import", tatoebaHandler.ImportSentences)
			pluginTatoeba.POST("link/import", tatoebaHandler.ImportLinks)

		}
	}

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	docs.SwaggerInfo.Title = "Swagger Example API"
	docs.SwaggerInfo.Description = "This is a sample server Petstore server."
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "cocotola.com"
	docs.SwaggerInfo.BasePath = "/v1"
	docs.SwaggerInfo.Schemes = []string{"https"}

	gracefulShutdownTime1 := time.Duration(cfg.Shutdown.TimeSec1) * time.Second
	gracefulShutdownTime2 := time.Duration(cfg.Shutdown.TimeSec2) * time.Second
	server := http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil {
			logrus.Infof("failed to ListenAndServe. err: %v", err)
			done <- true
		}
	}()

	logrus.Info("awaiting signal")
	<-done
	logrus.Info("exiting")

	ctx, cancel := context.WithTimeout(context.Background(), gracefulShutdownTime1)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		logrus.Infof("Server forced to shutdown. err: %v", err)
	}
	time.Sleep(gracefulShutdownTime2)
	logrus.Info("exited")
}

func initPf(synthesizer pluginCommonDomain.Synthesizer, translator pluginCommonDomain.Translator, tatoebaSentenceRepository pluginCommonDomain.TatoebaSentenceRepositoryReadOnly) (appD.ProcessorFactory, map[string]func(*gorm.DB) (appD.ProblemRepository, error), map[string]appD.ProblemImportProcessor) {

	englishWordProblemProcessor := pluginEnglishDomain.NewEnglishWordProblemProcessor(synthesizer, translator, tatoebaSentenceRepository, pluginEnglishGateway.NewEnglishWordProblemAddParameterCSVReader)
	englishPhraseProblemProcessor := pluginEnglishDomain.NewEnglishPhraseProblemProcessor(synthesizer, translator)
	englishSentenceProblemProcessor := pluginEnglishDomain.NewEnglishSentenceProblemProcessor(synthesizer, translator, pluginEnglishGateway.NewEnglishSentenceProblemAddParameterCSVReader)

	problemAddProcessor := map[string]appD.ProblemAddProcessor{
		pluginEnglishDomain.EnglishWordProblemType:     englishWordProblemProcessor,
		pluginEnglishDomain.EnglishPhraseProblemType:   englishPhraseProblemProcessor,
		pluginEnglishDomain.EnglishSentenceProblemType: englishSentenceProblemProcessor,
	}
	problemUpdateProcessor := map[string]appD.ProblemUpdateProcessor{
		pluginEnglishDomain.EnglishWordProblemType: englishWordProblemProcessor,
	}
	problemRemoveProcessor := map[string]appD.ProblemRemoveProcessor{
		pluginEnglishDomain.EnglishWordProblemType:     englishWordProblemProcessor,
		pluginEnglishDomain.EnglishPhraseProblemType:   englishPhraseProblemProcessor,
		pluginEnglishDomain.EnglishSentenceProblemType: englishSentenceProblemProcessor,
	}
	problemImportProcessor := map[string]appD.ProblemImportProcessor{
		pluginEnglishDomain.EnglishWordProblemType: englishWordProblemProcessor,
	}
	problemQuotaProcessor := map[string]appD.ProblemQuotaProcessor{
		pluginEnglishDomain.EnglishWordProblemType: englishWordProblemProcessor,
	}

	englishWordProblemRepositoryFunc := func(db *gorm.DB) (appD.ProblemRepository, error) {
		return pluginEnglishGateway.NewEnglishWordProblemRepository(db, pluginEnglishDomain.EnglishWordProblemType)
	}
	englishPhraseProblemRepositoryFunc := func(db *gorm.DB) (appD.ProblemRepository, error) {
		return pluginEnglishGateway.NewEnglishPhraseProblemRepository(db, pluginEnglishDomain.EnglishPhraseProblemType)
	}
	englishSentenceProblemRepositoryFunc := func(db *gorm.DB) (appD.ProblemRepository, error) {
		return pluginEnglishGateway.NewEnglishSentenceProblemRepository(db, pluginEnglishDomain.EnglishSentenceProblemType)
	}

	pf := appD.NewProcessorFactory(problemAddProcessor, problemUpdateProcessor, problemRemoveProcessor, problemImportProcessor, problemQuotaProcessor)
	problemRepositories := map[string]func(*gorm.DB) (appD.ProblemRepository, error){
		pluginEnglishDomain.EnglishWordProblemType:     englishWordProblemRepositoryFunc,
		pluginEnglishDomain.EnglishPhraseProblemType:   englishPhraseProblemRepositoryFunc,
		pluginEnglishDomain.EnglishSentenceProblemType: englishSentenceProblemRepositoryFunc,
	}
	return pf, problemRepositories, problemImportProcessor
}

func initialize(ctx context.Context, env string) (*config.Config, *gorm.DB, *sql.DB, *gin.Engine, error) {
	cfg, err := config.LoadConfig(env)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	// init log
	if err := config.InitLog(env, cfg.Log); err != nil {
		return nil, nil, nil, nil, err
	}

	// cors
	corsConfig := config.InitCORS(cfg.CORS)
	logrus.Infof("cors: %+v", corsConfig)

	if err := corsConfig.Validate(); err != nil {
		return nil, nil, nil, nil, err
	}

	// init db
	db, sqlDB, err := initDB(cfg.DB)
	if err != nil {
		return nil, nil, nil, nil, xerrors.Errorf("failed to InitDB. err: %w", err)
	}

	userRfFunc := func(db *gorm.DB) (userD.RepositoryFactory, error) {
		return userG.NewRepositoryFactory(db)
	}

	userD.InitSystemAdmin(userRfFunc)

	router := gin.New()
	router.Use(cors.New(corsConfig))
	router.Use(middleware.NewLogMiddleware())
	router.Use(gin.Recovery())

	if cfg.Debug.GinMode {
		router.Use(ginlog.Middleware(ginlog.DefaultConfig))
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	if cfg.Debug.Wait {
		router.Use(middleware.NewWaitMiddleware())
	}

	return cfg, db, sqlDB, router, nil
}

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
	case "mysql":
		db, err := libG.OpenMySQL(cfg.MySQL.Username, cfg.MySQL.Password, cfg.MySQL.Host, cfg.MySQL.Port, cfg.MySQL.Database)
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

		if err := appG.MigrateMySQLDB(db); err != nil {
			return nil, nil, xerrors.Errorf("failed to MigrateMySQLDB. err: %w", err)
		}

		return db, sqlDB, nil
	default:
		return nil, nil, libD.ErrInvalidArgument
	}
}

func initApp1(ctx context.Context, db *gorm.DB, password string) error {
	logger := log.FromContext(ctx)

	if err := db.Transaction(func(tx *gorm.DB) error {
		systemAdmin, err := userD.NewSystemAdminFromDB(tx)
		if err != nil {
			return err
		}

		organization, err := systemAdmin.FindOrganizationByName(ctx, "cocotola")
		if err != nil {
			if !errors.Is(err, userD.ErrOrganizationNotFound) {
				return xerrors.Errorf("failed to AddOrganization. err: %w", err)
			}

			firstOwnerAddParam, err := userD.NewFirstOwnerAddParameter("cocotola-owner", "Owner(cocotola)", password)
			if err != nil {
				return xerrors.Errorf("failed to AddOrganization. err: %w", err)
			}

			organizationAddParameter, err := userD.NewOrganizationAddParameter("cocotola", firstOwnerAddParam)
			if err != nil {
				return xerrors.Errorf("failed to AddOrganization. err: %w", err)
			}

			organizationID, err := systemAdmin.AddOrganization(ctx, organizationAddParameter)
			if err != nil {
				return xerrors.Errorf("failed to AddOrganization. err: %w", err)
			}

			logger.Infof("organizationID: %d", organizationID)
			return nil
		}
		logger.Infof("organization: %d", organization.GetID())
		return nil
	}); err != nil {
		return err
	}
	return nil
}

func initApp2(ctx context.Context, db *gorm.DB, rfFunc func(db *gorm.DB) (appD.RepositoryFactory, error), userRfFunc func(db *gorm.DB) (userD.RepositoryFactory, error)) error {
	if err := initApp2_1(ctx, db, rfFunc, userRfFunc); err != nil {
		return err
	}

	if err := initApp2_2(ctx, db, rfFunc, userRfFunc); err != nil {
		return err
	}

	if err := initApp2_3(ctx, db, rfFunc, userRfFunc); err != nil {
		return err
	}

	return nil
}

func initApp2_1(ctx context.Context, db *gorm.DB, rfFunc func(db *gorm.DB) (appD.RepositoryFactory, error), userRfFunc func(db *gorm.DB) (userD.RepositoryFactory, error)) error {
	var propertiesSystemStudentID userD.AppUserID

	if err := db.Transaction(func(tx *gorm.DB) error {
		userRf, err := userRfFunc(tx)
		if err != nil {
			return err
		}

		systemAdmin := userD.NewSystemAdmin(userRf)

		systemOwner, err := systemAdmin.FindSystemOwnerByOrganizationName(ctx, "cocotola")
		if err != nil {
			return xerrors.Errorf("failed to FindSystemOwnerByOrganizationName. err: %w", err)
		}

		systemStudent, err := systemOwner.FindAppUserByLoginID(ctx, appD.SystemStudentLoginID)
		if err != nil {
			if !errors.Is(err, userD.ErrAppUserNotFound) {
				return xerrors.Errorf("failed to FindAppUserByLoginID. err: %w", err)
			}

			param, err := userD.NewAppUserAddParameter(appD.SystemStudentLoginID, "SystemStudent(cocotola)", []string{}, map[string]string{})
			if err != nil {
				return xerrors.Errorf("failed to NewAppUserAddParameter. err: %w", err)
			}

			systemStudentID, err := systemOwner.AddAppUser(ctx, param)
			if err != nil {
				return xerrors.Errorf("failed to AddAppUser. err: %w", err)
			}

			propertiesSystemStudentID = systemStudentID
		} else {
			propertiesSystemStudentID = userD.AppUserID(systemStudent.GetID())
		}
		return nil
	}); err != nil {
		return err
	}

	appD.SetSystemStudentID(propertiesSystemStudentID)

	return nil
}

func initApp2_2(ctx context.Context, db *gorm.DB, rfFunc func(db *gorm.DB) (appD.RepositoryFactory, error), userRfFunc func(db *gorm.DB) (userD.RepositoryFactory, error)) error {

	var propertiesSystemSpaceID userD.SpaceID

	if err := db.Transaction(func(tx *gorm.DB) error {
		userRf, err := userRfFunc(tx)
		if err != nil {
			return err
		}

		systemAdmin := userD.NewSystemAdmin(userRf)

		systemOwner, err := systemAdmin.FindSystemOwnerByOrganizationName(ctx, appD.OrganizationName)
		if err != nil {
			return xerrors.Errorf("failed to FindSystemOwnerByOrganizationName. err: %w", err)
		}

		systemSpace, err := systemOwner.FindSystemSpace(ctx)
		if err != nil {
			if !errors.Is(err, userD.ErrSpaceNotFound) {
				return xerrors.Errorf("failed to FindSystemSpace. err: %w", err)
			}

			spaceID, err := systemOwner.AddSystemSpace(ctx)
			if err != nil {
				return xerrors.Errorf("failed to AddSystemSpace. err: %w", err)
			}

			propertiesSystemSpaceID = spaceID
		} else {
			propertiesSystemSpaceID = userD.SpaceID(systemSpace.GetID())
		}

		return nil
	}); err != nil {
		return err
	}

	appD.SetSystemSpaceID(propertiesSystemSpaceID)

	return nil
}

func initApp2_3(ctx context.Context, db *gorm.DB, rfFunc func(db *gorm.DB) (appD.RepositoryFactory, error), userRfFunc func(db *gorm.DB) (userD.RepositoryFactory, error)) error {

	var propertiesTatoebaWorkbookID appD.WorkbookID
	if err := db.Transaction(func(tx *gorm.DB) error {
		userRf, err := userRfFunc(tx)
		if err != nil {
			return err
		}

		systemAdmin := userD.NewSystemAdmin(userRf)

		systemOwner, err := systemAdmin.FindSystemOwnerByOrganizationName(ctx, appD.OrganizationName)
		if err != nil {
			return xerrors.Errorf("failed to FindSystemOwnerByOrganizationName. err: %w", err)
		}

		systemStudentAppUser, err := systemOwner.FindAppUserByLoginID(ctx, appD.SystemStudentLoginID)
		if err != nil {
			return xerrors.Errorf("failed to FindAppUserByLoginID. err: %w", err)
		}

		rf, err := rfFunc(tx)
		if err != nil {
			return err
		}

		systemStudent, err := appD.NewSystemStudent(rf, systemStudentAppUser)
		if err != nil {
			return err
		}

		tatoebaWorkbook, err := systemStudent.FindWorkbookFromSystemSpace(ctx, appD.TatoebaWorkbookName)
		if err != nil {
			if !errors.Is(err, appD.ErrWorkbookNotFound) {
				return err
			}

			paramToAddWorkbook, err := appD.NewWorkbookAddParameter(pluginEnglishDomain.EnglishSentenceProblemType, appD.TatoebaWorkbookName, "", map[string]string{})
			if err != nil {
				return err
			}

			tatoebaWorkbookID, err := systemStudent.AddWorkbookToSystemSpace(ctx, paramToAddWorkbook)
			if err != nil {
				return err
			}

			propertiesTatoebaWorkbookID = tatoebaWorkbookID
		} else {
			propertiesTatoebaWorkbookID = appD.WorkbookID(tatoebaWorkbook.GetID())
		}

		return nil
	}); err != nil {
		return err
	}

	appD.SetTatoebaWorkbookID(propertiesTatoebaWorkbookID)

	return nil
}

func callback(ctx context.Context, testUserEmail string, pf appD.ProcessorFactory, repo appD.RepositoryFactory, userRepo userD.RepositoryFactory, organizationName string, appUser userD.AppUser) error {
	logger := log.FromContext(ctx)
	logger.Infof("callback. loginID: %s", appUser.GetLoginID())

	if appUser.GetLoginID() == testUserEmail {
		student, err := appD.NewStudent(pf, repo, userRepo, appUser)
		if err != nil {
			return xerrors.Errorf("failed to NewStudent. err: %w", err)
		}

		if err := english_word.CreateDemoWorkbook(ctx, student); err != nil {
			return err
		}

		if err := english_word.Create20NGSLWorkbook(ctx, student); err != nil {
			return err
		}

		if err := english_word.Create300NGSLWorkbook(ctx, student); err != nil {
			return err
		}
	}

	return nil
}
