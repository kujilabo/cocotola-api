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
	"strconv"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	ginlog "github.com/onrik/logrus/gin"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"golang.org/x/xerrors"
	"gorm.io/gorm"

	swaggerFiles "github.com/swaggo/files"     // swagger embed files
	ginSwagger "github.com/swaggo/gin-swagger" // gin-swagger middleware

	"github.com/kujilabo/cocotola-api/docs"
	"github.com/kujilabo/cocotola-api/src/app/config"
	appD "github.com/kujilabo/cocotola-api/src/app/domain"
	appG "github.com/kujilabo/cocotola-api/src/app/gateway"
	appH "github.com/kujilabo/cocotola-api/src/app/handler"
	appS "github.com/kujilabo/cocotola-api/src/app/service"
	studentU "github.com/kujilabo/cocotola-api/src/app/usecase/student"
	authG "github.com/kujilabo/cocotola-api/src/auth/gateway"
	authH "github.com/kujilabo/cocotola-api/src/auth/handler"
	authM "github.com/kujilabo/cocotola-api/src/auth/handler/middleware"
	authU "github.com/kujilabo/cocotola-api/src/auth/usecase"
	english_word "github.com/kujilabo/cocotola-api/src/data/english_word"
	"github.com/kujilabo/cocotola-api/src/lib/handler/middleware"
	"github.com/kujilabo/cocotola-api/src/lib/log"
	pluginCommonGateway "github.com/kujilabo/cocotola-api/src/plugin/common/gateway"
	pluginCommonHandler "github.com/kujilabo/cocotola-api/src/plugin/common/handler"
	pluginCommonS "github.com/kujilabo/cocotola-api/src/plugin/common/service"
	pluginEnglishDomain "github.com/kujilabo/cocotola-api/src/plugin/english/domain"
	pluginEnglishGateway "github.com/kujilabo/cocotola-api/src/plugin/english/gateway"
	pluginEnglishS "github.com/kujilabo/cocotola-api/src/plugin/english/service"
	userD "github.com/kujilabo/cocotola-api/src/user/domain"
	userG "github.com/kujilabo/cocotola-api/src/user/gateway"
	userS "github.com/kujilabo/cocotola-api/src/user/service"
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

	cfg, db, sqlDB, router, tp, err := initialize(ctx, *env)
	if err != nil {
		panic(err)
	}
	defer sqlDB.Close()
	defer tp.ForceFlush(ctx) // flushes any pending spans

	userRfFunc := func(ctx context.Context, db *gorm.DB) (userS.RepositoryFactory, error) {
		return userG.NewRepositoryFactory(db)
	}
	appS.UserRfFunc = userRfFunc
	if err := initApp1(ctx, db, cfg.App.OwnerPassword); err != nil {
		panic(err)
	}

	router.GET("/healthcheck", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	synthesizer := appG.NewSynthesizerClient(cfg.Synthesizer.Endpoint, cfg.Synthesizer.Username, cfg.Synthesizer.Password, time.Duration(cfg.Synthesizer.TimeoutSec)*time.Second)

	translatorClient := pluginCommonGateway.NewTranslatorClient(cfg.Translator.Endpoint, cfg.Translator.Username, cfg.Translator.Password, time.Duration(cfg.Translator.TimeoutSec)*time.Second)
	tatoebaClient := pluginCommonGateway.NewTatoebaClient(cfg.Tatoeba.Endpoint, cfg.Tatoeba.Username, cfg.Tatoeba.Password, time.Duration(cfg.Tatoeba.TimeoutSec)*time.Second)

	pf, problemRepositories, problemImportProcessor := initPf(synthesizer, translatorClient, tatoebaClient)

	newIterator := func(ctx context.Context, workbookID appD.WorkbookID, problemType string, reader io.Reader) (appS.ProblemAddParameterIterator, error) {
		processor, ok := problemImportProcessor[problemType]
		if ok {
			return processor.CreateCSVReader(ctx, workbookID, reader)
		}
		return nil, xerrors.Errorf("processor not found. problemType: %s", problemType)
	}

	problemTypeRepo := appG.NewProblemTypeRepository(db)
	problemTypes, err := problemTypeRepo.FindAllProblemTypes(ctx)
	if err != nil {
		panic(err)
	}

	studyTypeRepo := appG.NewStudyTypeRepository(db)
	studyTypes, err := studyTypeRepo.FindAllStudyTypes(ctx)
	if err != nil {
		panic(err)
	}

	rfFunc := func(ctx context.Context, db *gorm.DB) (appS.RepositoryFactory, error) {
		return appG.NewRepositoryFactory(ctx, db, cfg.DB.DriverName, userRfFunc, pf, problemTypes, studyTypes, problemRepositories)
	}
	appS.RfFunc = rfFunc

	if err := initApp2(ctx, db, rfFunc, userRfFunc); err != nil {
		panic(err)
	}

	signingKey := []byte(cfg.Auth.SigningKey)
	signingMethod := jwt.SigningMethodHS256
	authTokenManager := authG.NewAuthTokenManager(signingKey, signingMethod, time.Duration(cfg.Auth.AccessTokenTTLMin)*time.Minute, time.Duration(cfg.Auth.RefreshTokenTTLHour)*time.Hour)

	googleAuthClient := authG.NewGoogleAuthClient(cfg.Auth.GoogleClientID, cfg.Auth.GoogleClientSecret, cfg.Auth.GoogleCallbackURL, time.Duration(cfg.Auth.APITimeoutSec)*time.Second)
	authMiddleware := authM.NewAuthMiddleware(signingKey)

	registerAppUserCallback := func(ctx context.Context, db *gorm.DB, organizationName string, appUser userD.AppUserModel) error {
		rf, err := rfFunc(ctx, db)
		if err != nil {
			return err
		}
		userRf, err := userRfFunc(ctx, db)
		if err != nil {
			return err
		}
		return callback(ctx, cfg.App.TestUserEmail, pf, rf, userRf, organizationName, appUser)
	}

	v1 := router.Group("v1")
	{
		v1.Use(otelgin.Middleware(cfg.App.Name))
		v1.Use(middleware.NewTraceLogMiddleware(cfg.App.Name))
		v1auth := v1.Group("auth")
		googleUserUsecase := authU.NewGoogleUserUsecase(db, googleAuthClient, authTokenManager, registerAppUserCallback)
		guestUserUsecase := authU.NewGuestUserUsecase(authTokenManager)
		authHandler := authH.NewAuthHandler(authTokenManager)
		googleAuthHandler := authH.NewGoogleAuthHandler(googleUserUsecase)
		guestAuthHandler := authH.NewGuestAuthHandler(guestUserUsecase)
		v1auth.POST("google/authorize", googleAuthHandler.Authorize)
		v1auth.POST("guest/authorize", guestAuthHandler.Authorize)
		v1auth.POST("refresh_token", authHandler.RefreshToken)

		v1Workbook := v1.Group("private/workbook")
		studentUsecaseWorkbook := studentU.NewStudentUsecaseWorkbook(db, pf, rfFunc, userRfFunc)
		privateWorkbookHandler := appH.NewPrivateWorkbookHandler(studentUsecaseWorkbook)
		v1Workbook.Use(authMiddleware)
		v1Workbook.POST(":workbookID", privateWorkbookHandler.FindWorkbooks)
		v1Workbook.GET(":workbookID", privateWorkbookHandler.FindWorkbookByID)
		v1Workbook.PUT(":workbookID", privateWorkbookHandler.UpdateWorkbook)
		v1Workbook.DELETE(":workbookID", privateWorkbookHandler.RemoveWorkbook)
		v1Workbook.POST("", privateWorkbookHandler.AddWorkbook)

		v1Problem := v1.Group("workbook/:workbookID/problem")
		studentUsecaseProblem := studentU.NewStudentUsecaseProblem(db, pf, rfFunc, userRfFunc)
		problemHandler := appH.NewProblemHandler(studentUsecaseProblem, newIterator)
		v1Problem.Use(authMiddleware)
		v1Problem.POST("", problemHandler.AddProblem)
		v1Problem.GET(":problemID", problemHandler.FindProblemByID)
		v1Problem.DELETE(":problemID", problemHandler.RemoveProblem)
		v1Problem.PUT(":problemID", problemHandler.UpdateProblem)
		// v1Problem.GET("problem_ids", problemHandler.FindProblemIDs)
		v1Problem.POST("find", problemHandler.FindProblems)
		v1Problem.POST("find_all", problemHandler.FindAllProblems)
		v1Problem.POST("find_by_ids", problemHandler.FindProblemsByProblemIDs)
		v1Problem.POST("import", problemHandler.ImportProblems)

		v1Study := v1.Group("study/workbook/:workbookID")
		studentUseCaseStudy := studentU.NewStudentUsecaseStudy(db, pf, rfFunc, userRfFunc)
		recordbookHandler := appH.NewRecordbookHandler(studentUseCaseStudy)
		v1Study.Use(authMiddleware)
		v1Study.GET("study_type/:studyType", recordbookHandler.FindRecordbook)
		v1Study.POST("study_type/:studyType/problem/:problemID/record", recordbookHandler.SetStudyResult)
		v1Study.GET("completion_rate", recordbookHandler.GetCompletionRate)

		v1Audio := v1.Group("workbook/:workbookID/problem/:problemID/audio")
		studentUsecaseAudio := studentU.NewStudentUsecaseAudio(db, pf, rfFunc, userRfFunc, synthesizer)
		audioHandler := appH.NewAudioHandler(studentUsecaseAudio)
		v1Audio.Use(authMiddleware)
		v1Audio.GET(":audioID", audioHandler.FindAudioByID)
	}

	plugin := router.Group("plugin")
	{
		plugin.Use(otelgin.Middleware(cfg.App.Name))
		plugin.Use(middleware.NewTraceLogMiddleware(cfg.App.Name))
		plugin.Use(authMiddleware)
		{
			pluginTranslation := plugin.Group("translation")
			translationHandler := pluginCommonHandler.NewTranslationHandler(translatorClient)
			pluginTranslation.POST("find", translationHandler.FindTranslations)
			pluginTranslation.GET("text/:text/pos/:pos", translationHandler.FindTranslationByTextAndPos)
			pluginTranslation.GET("text/:text", translationHandler.FindTranslationsByText)
			pluginTranslation.PUT("text/:text/pos/:pos", translationHandler.UpdateTranslation)
			pluginTranslation.DELETE("text/:text/pos/:pos", translationHandler.RemoveTranslation)
			pluginTranslation.POST("", translationHandler.AddTranslation)
			pluginTranslation.POST("export", translationHandler.ExportTranslations)
		}
		{
			pluginTatoeba := plugin.Group("tatoeba")
			tatoebaHandler := pluginCommonHandler.NewTatoebaHandler(tatoebaClient)
			pluginTatoeba.POST("find", tatoebaHandler.FindSentencePairs)
			pluginTatoeba.POST("sentence/import", tatoebaHandler.ImportSentences)
			pluginTatoeba.POST("link/import", tatoebaHandler.ImportLinks)
		}
	}

	if cfg.Swagger.Enabled {
		router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
		docs.SwaggerInfo.Title = cfg.App.Name
		docs.SwaggerInfo.Version = "1.0"
		docs.SwaggerInfo.Host = cfg.Swagger.Host
		docs.SwaggerInfo.Schemes = []string{cfg.Swagger.Schema}
	}

	gracefulShutdownTime1 := time.Duration(cfg.Shutdown.TimeSec1) * time.Second
	gracefulShutdownTime2 := time.Duration(cfg.Shutdown.TimeSec2) * time.Second
	server := http.Server{
		Addr:    ":" + strconv.Itoa(cfg.App.Port),
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

func initPf(synthesizerClient appS.SynthesizerClient, translatorClient pluginCommonS.TranslatorClient, tatoebaClient pluginCommonS.TatoebaClient) (appS.ProcessorFactory, map[string]func(context.Context, *gorm.DB) (appS.ProblemRepository, error), map[string]appS.ProblemImportProcessor) {

	englishWordProblemProcessor := pluginEnglishS.NewEnglishWordProblemProcessor(synthesizerClient, translatorClient, tatoebaClient, pluginEnglishGateway.NewEnglishWordProblemAddParameterCSVReader)
	englishPhraseProblemProcessor := pluginEnglishS.NewEnglishPhraseProblemProcessor(synthesizerClient, translatorClient)
	englishSentenceProblemProcessor := pluginEnglishS.NewEnglishSentenceProblemProcessor(synthesizerClient, translatorClient, pluginEnglishGateway.NewEnglishSentenceProblemAddParameterCSVReader)

	problemAddProcessor := map[string]appS.ProblemAddProcessor{
		pluginEnglishDomain.EnglishWordProblemType:     englishWordProblemProcessor,
		pluginEnglishDomain.EnglishPhraseProblemType:   englishPhraseProblemProcessor,
		pluginEnglishDomain.EnglishSentenceProblemType: englishSentenceProblemProcessor,
	}
	problemUpdateProcessor := map[string]appS.ProblemUpdateProcessor{
		pluginEnglishDomain.EnglishWordProblemType: englishWordProblemProcessor,
	}
	problemRemoveProcessor := map[string]appS.ProblemRemoveProcessor{
		pluginEnglishDomain.EnglishWordProblemType:     englishWordProblemProcessor,
		pluginEnglishDomain.EnglishPhraseProblemType:   englishPhraseProblemProcessor,
		pluginEnglishDomain.EnglishSentenceProblemType: englishSentenceProblemProcessor,
	}
	problemImportProcessor := map[string]appS.ProblemImportProcessor{
		pluginEnglishDomain.EnglishWordProblemType: englishWordProblemProcessor,
	}
	problemQuotaProcessor := map[string]appS.ProblemQuotaProcessor{
		pluginEnglishDomain.EnglishWordProblemType: englishWordProblemProcessor,
	}

	englishWordProblemRepositoryFunc := func(ctx context.Context, db *gorm.DB) (appS.ProblemRepository, error) {
		// fmt.Println("-------Word")
		return pluginEnglishGateway.NewEnglishWordProblemRepository(db, synthesizerClient, pluginEnglishDomain.EnglishWordProblemType)
	}
	englishPhraseProblemRepositoryFunc := func(ctx context.Context, db *gorm.DB) (appS.ProblemRepository, error) {
		return pluginEnglishGateway.NewEnglishPhraseProblemRepository(db, synthesizerClient, pluginEnglishDomain.EnglishPhraseProblemType)
	}
	englishSentenceProblemRepositoryFunc := func(ctx context.Context, db *gorm.DB) (appS.ProblemRepository, error) {
		// fmt.Println("-------Sentence")
		return pluginEnglishGateway.NewEnglishSentenceProblemRepository(db, synthesizerClient, pluginEnglishDomain.EnglishSentenceProblemType)
	}

	pf := appS.NewProcessorFactory(problemAddProcessor, problemUpdateProcessor, problemRemoveProcessor, problemImportProcessor, problemQuotaProcessor)

	problemRepositories := map[string]func(context.Context, *gorm.DB) (appS.ProblemRepository, error){
		pluginEnglishDomain.EnglishWordProblemType:     englishWordProblemRepositoryFunc,
		pluginEnglishDomain.EnglishPhraseProblemType:   englishPhraseProblemRepositoryFunc,
		pluginEnglishDomain.EnglishSentenceProblemType: englishSentenceProblemRepositoryFunc,
	}
	return pf, problemRepositories, problemImportProcessor
}

func initialize(ctx context.Context, env string) (*config.Config, *gorm.DB, *sql.DB, *gin.Engine, *sdktrace.TracerProvider, error) {
	cfg, err := config.LoadConfig(env)
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}

	// init log
	if err := config.InitLog(env, cfg.Log); err != nil {
		return nil, nil, nil, nil, nil, err
	}

	// cors
	corsConfig := config.InitCORS(cfg.CORS)
	logrus.Infof("cors: %+v", corsConfig)

	if err := corsConfig.Validate(); err != nil {
		return nil, nil, nil, nil, nil, err
	}

	// tracer
	tp, err := config.InitTracerProvider(cfg)
	if err != nil {
		return nil, nil, nil, nil, nil, xerrors.Errorf("failed to InitTracerProvider. err: %w", err)
	}
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

	// init db
	db, sqlDB, err := config.InitDB(cfg.DB)
	if err != nil {
		return nil, nil, nil, nil, nil, xerrors.Errorf("failed to InitDB. err: %w", err)
	}

	userRfFunc := func(ctx context.Context, db *gorm.DB) (userS.RepositoryFactory, error) {
		return userG.NewRepositoryFactory(db)
	}
	userS.InitSystemAdmin(userRfFunc)

	if !cfg.Debug.GinMode {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()
	router.Use(cors.New(corsConfig))
	router.Use(gin.Recovery())

	if cfg.Debug.GinMode {
		router.Use(ginlog.Middleware(ginlog.DefaultConfig))
	}

	if cfg.Debug.Wait {
		router.Use(middleware.NewWaitMiddleware())
	}

	return cfg, db, sqlDB, router, tp, nil
}

func initApp1(ctx context.Context, db *gorm.DB, password string) error {
	logger := log.FromContext(ctx)

	if err := db.Transaction(func(tx *gorm.DB) error {
		systemAdmin, err := userS.NewSystemAdminFromDB(ctx, tx)
		if err != nil {
			return err
		}

		organization, err := systemAdmin.FindOrganizationByName(ctx, "cocotola")
		if err != nil {
			if !errors.Is(err, userS.ErrOrganizationNotFound) {
				return xerrors.Errorf("failed to AddOrganization. err: %w", err)
			}

			firstOwnerAddParam, err := userS.NewFirstOwnerAddParameter("cocotola-owner", "Owner(cocotola)", password)
			if err != nil {
				return xerrors.Errorf("failed to AddOrganization. err: %w", err)
			}

			organizationAddParameter, err := userS.NewOrganizationAddParameter("cocotola", firstOwnerAddParam)
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

func initApp2(ctx context.Context, db *gorm.DB, rfFunc appS.RepositoryFactoryFunc, userRfFunc userS.RepositoryFactoryFunc) error {
	if err := initApp2_1(ctx, db, rfFunc, userRfFunc); err != nil {
		return xerrors.Errorf("failed to initApp2_1. err: %w", err)
	}

	if err := initApp2_2(ctx, db, rfFunc, userRfFunc); err != nil {
		return xerrors.Errorf("failed to initApp2_2. err: %w", err)
	}

	if err := initApp2_3(ctx, db, rfFunc, userRfFunc); err != nil {
		return xerrors.Errorf("failed to initApp2_3. err: %w", err)
	}

	return nil
}

func initApp2_1(ctx context.Context, db *gorm.DB, rfFunc appS.RepositoryFactoryFunc, userRfFunc userS.RepositoryFactoryFunc) error {
	var propertiesSystemStudentID userD.AppUserID

	if err := db.Transaction(func(tx *gorm.DB) error {
		userRf, err := userRfFunc(ctx, tx)
		if err != nil {
			return err
		}

		systemAdmin := userS.NewSystemAdmin(userRf)

		systemOwner, err := systemAdmin.FindSystemOwnerByOrganizationName(ctx, "cocotola")
		if err != nil {
			return xerrors.Errorf("failed to FindSystemOwnerByOrganizationName. err: %w", err)
		}

		systemStudent, err := systemOwner.FindAppUserByLoginID(ctx, appS.SystemStudentLoginID)
		if err != nil {
			if !errors.Is(err, userS.ErrAppUserNotFound) {
				return xerrors.Errorf("failed to FindAppUserByLoginID. err: %w", err)
			}

			param, err := userS.NewAppUserAddParameter(appS.SystemStudentLoginID, "SystemStudent(cocotola)", []string{}, map[string]string{})
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

	appS.SetSystemStudentID(propertiesSystemStudentID)

	return nil
}

func initApp2_2(ctx context.Context, db *gorm.DB, rfFunc appS.RepositoryFactoryFunc, userRfFunc userS.RepositoryFactoryFunc) error {

	var propertiesSystemSpaceID userD.SpaceID

	if err := db.Transaction(func(tx *gorm.DB) error {
		userRf, err := userRfFunc(ctx, tx)
		if err != nil {
			return err
		}

		systemAdmin := userS.NewSystemAdmin(userRf)

		systemOwner, err := systemAdmin.FindSystemOwnerByOrganizationName(ctx, appS.OrganizationName)
		if err != nil {
			return xerrors.Errorf("failed to FindSystemOwnerByOrganizationName. err: %w", err)
		}

		systemSpace, err := systemOwner.FindSystemSpace(ctx)
		if err != nil {
			if !errors.Is(err, userS.ErrSpaceNotFound) {
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

	appS.SetSystemSpaceID(propertiesSystemSpaceID)

	return nil
}

func initApp2_3(ctx context.Context, db *gorm.DB, rfFunc appS.RepositoryFactoryFunc, userRfFunc userS.RepositoryFactoryFunc) error {

	var propertiesTatoebaWorkbookID appD.WorkbookID
	if err := db.Transaction(func(tx *gorm.DB) error {
		userRf, err := userRfFunc(ctx, tx)
		if err != nil {
			return err
		}

		systemAdmin := userS.NewSystemAdmin(userRf)

		systemOwner, err := systemAdmin.FindSystemOwnerByOrganizationName(ctx, appS.OrganizationName)
		if err != nil {
			return xerrors.Errorf("failed to FindSystemOwnerByOrganizationName. err: %w", err)
		}

		systemStudentAppUser, err := systemOwner.FindAppUserByLoginID(ctx, appS.SystemStudentLoginID)
		if err != nil {
			return xerrors.Errorf("failed to FindAppUserByLoginID. err: %w", err)
		}

		rf, err := rfFunc(ctx, tx)
		if err != nil {
			return err
		}

		systemStudent, err := appS.NewSystemStudent(rf, systemStudentAppUser)
		if err != nil {
			return err
		}

		tatoebaWorkbook, err := systemStudent.FindWorkbookFromSystemSpace(ctx, appS.TatoebaWorkbookName)
		if err != nil {
			if !errors.Is(err, appS.ErrWorkbookNotFound) {
				return err
			}

			paramToAddWorkbook, err := appS.NewWorkbookAddParameter(pluginEnglishDomain.EnglishSentenceProblemType, appS.TatoebaWorkbookName, appD.Lang2JA, "", map[string]string{})
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

	appS.SetTatoebaWorkbookID(propertiesTatoebaWorkbookID)

	return nil
}

func callback(ctx context.Context, testUserEmail string, pf appS.ProcessorFactory, repo appS.RepositoryFactory, userRepo userS.RepositoryFactory, organizationName string, appUser userD.AppUserModel) error {
	logger := log.FromContext(ctx)
	logger.Infof("callback. loginID: %s", appUser.GetLoginID())

	if appUser.GetLoginID() == testUserEmail {
		student, err := appS.NewStudent(pf, repo, userRepo, appUser)
		if err != nil {
			return xerrors.Errorf("failed to NewStudent. err: %w", err)
		}

		if err := english_word.CreateDemoWorkbook(ctx, student); err != nil {
			return xerrors.Errorf("failed to CreateDemoWorkbook. err: %w", err)
		}

		if err := english_word.Create20NGSLWorkbook(ctx, student); err != nil {
			return xerrors.Errorf("failed to Create20NGSLWorkbook. err: %w", err)
		}

		// if err := english_word.Create300NGSLWorkbook(ctx, student); err != nil {
		// 	return xerrors.Errorf("failed to Create300NGSLWorkbook. err: %w", err)
		// }
	}

	return nil
}
