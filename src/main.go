package main

import (
	"context"
	"database/sql"
	"errors"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"golang.org/x/sync/errgroup"
	"golang.org/x/xerrors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"gorm.io/gorm"

	swaggerFiles "github.com/swaggo/files"     // swagger embed files
	ginSwagger "github.com/swaggo/gin-swagger" // gin-swagger middleware

	"github.com/kujilabo/cocotola-api/docs"
	"github.com/kujilabo/cocotola-api/src/app/config"
	"github.com/kujilabo/cocotola-api/src/app/controller"
	appD "github.com/kujilabo/cocotola-api/src/app/domain"
	appG "github.com/kujilabo/cocotola-api/src/app/gateway"
	appS "github.com/kujilabo/cocotola-api/src/app/service"
	studentU "github.com/kujilabo/cocotola-api/src/app/usecase/student"
	authG "github.com/kujilabo/cocotola-api/src/auth/gateway"
	authU "github.com/kujilabo/cocotola-api/src/auth/usecase"
	english_word "github.com/kujilabo/cocotola-api/src/data/english_word"
	"github.com/kujilabo/cocotola-api/src/lib/log"
	pluginCommonGateway "github.com/kujilabo/cocotola-api/src/plugin/common/gateway"
	pluginCommonS "github.com/kujilabo/cocotola-api/src/plugin/common/service"
	pluginEnglishDomain "github.com/kujilabo/cocotola-api/src/plugin/english/domain"
	pluginEnglishGateway "github.com/kujilabo/cocotola-api/src/plugin/english/gateway"
	pluginEnglishS "github.com/kujilabo/cocotola-api/src/plugin/english/service"
	userD "github.com/kujilabo/cocotola-api/src/user/domain"
	userG "github.com/kujilabo/cocotola-api/src/user/gateway"
	userS "github.com/kujilabo/cocotola-api/src/user/service"
)

// type newIteratorFunc func(ctx context.Context, workbookID appD.WorkbookID, problemType string, reader io.Reader) (appS.ProblemAddParameterIterator, error)

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

	cfg, db, sqlDB, tp, err := initialize(ctx, *env)
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

	synthesizer := appG.NewSynthesizerClient(cfg.Synthesizer.Endpoint, cfg.Synthesizer.Username, cfg.Synthesizer.Password, time.Duration(cfg.Synthesizer.TimeoutSec)*time.Second)

	// translator
	connTranslator, err := grpc.Dial(cfg.Translator.GRPCAddr, grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(otelgrpc.UnaryClientInterceptor()),
		grpc.WithStreamInterceptor(otelgrpc.StreamClientInterceptor()))
	if err != nil {
		panic(err)
	}
	defer connTranslator.Close()
	// translatorClient := pluginCommonGateway.NewTranslatorHTTPClient(cfg.Translator.Endpoint, cfg.Translator.Username, cfg.Translator.Password, time.Duration(cfg.Translator.TimeoutSec)*time.Second)
	translatorClient := pluginCommonGateway.NewTranslatorGRPCClient(connTranslator, cfg.Translator.Username, cfg.Translator.Password, time.Duration(cfg.Translator.TimeoutSec)*time.Second)

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

	gracefulShutdownTime2 := time.Duration(cfg.Shutdown.TimeSec2) * time.Second

	// {
	// 	conn, err := grpc.Dial(cfg.Translator.GRPCAddr, grpc.WithTransportCredentials(insecure.NewCredentials()),
	// 		grpc.WithUnaryInterceptor(otelgrpc.UnaryClientInterceptor()),
	// 		grpc.WithStreamInterceptor(otelgrpc.StreamClientInterceptor()))
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// 	defer conn.Close()

	// 	x := pluginCommonGateway.NewTranslatorGRPCClient(conn, "", "", time.Duration(cfg.Translator.TimeoutSec)*time.Second)
	// 	y, err := x.DictionaryLookup(ctx, appD.Lang2EN, appD.Lang2JA, "book")
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// 	logrus.Info("-----------------------------")
	// 	logrus.Info(y)
	// }

	result := run(context.Background(), cfg, db, pf, rfFunc, userRfFunc, synthesizer, translatorClient, tatoebaClient, newIterator)

	time.Sleep(gracefulShutdownTime2)
	logrus.Info("exited")
	os.Exit(result)
}

func run(ctx context.Context, cfg *config.Config, db *gorm.DB, pf appS.ProcessorFactory, rfFunc appS.RepositoryFactoryFunc, userRfFunc userS.RepositoryFactoryFunc, synthesizerClient appS.SynthesizerClient, translatorClient pluginCommonS.TranslatorClient, tatoebaClient pluginCommonS.TatoebaClient, newIteratorFunc controller.NewIteratorFunc) int {
	var eg *errgroup.Group
	eg, ctx = errgroup.WithContext(ctx)

	eg.Go(func() error {
		return httpServer(ctx, cfg, db, pf, rfFunc, userRfFunc, synthesizerClient, translatorClient, tatoebaClient, newIteratorFunc)
	})
	eg.Go(func() error {
		return signalNotify(ctx)
	})
	eg.Go(func() error {
		<-ctx.Done()
		return ctx.Err()
	})

	if err := eg.Wait(); err != nil {
		logrus.Error(err)
		return 1
	}
	return 0
}

func signalNotify(ctx context.Context) error {
	sigs := make(chan os.Signal, 1)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-ctx.Done():
		signal.Reset()
		return nil
	case sig := <-sigs:
		return fmt.Errorf("signal received: %v", sig.String())
	}
}

func httpServer(ctx context.Context, cfg *config.Config, db *gorm.DB, pf appS.ProcessorFactory, rfFunc appS.RepositoryFactoryFunc, userRfFunc userS.RepositoryFactoryFunc, synthesizerClient appS.SynthesizerClient, translatorClient pluginCommonS.TranslatorClient, tatoebaClient pluginCommonS.TatoebaClient, newIteratorFunc controller.NewIteratorFunc) error {
	// cors
	corsConfig := config.InitCORS(cfg.CORS)
	logrus.Infof("cors: %+v", corsConfig)

	if err := corsConfig.Validate(); err != nil {
		return err
	}

	if !cfg.Debug.GinMode {
		gin.SetMode(gin.ReleaseMode)
	}

	signingKey := []byte(cfg.Auth.SigningKey)
	signingMethod := jwt.SigningMethodHS256
	authTokenManager := authG.NewAuthTokenManager(signingKey, signingMethod, time.Duration(cfg.Auth.AccessTokenTTLMin)*time.Minute, time.Duration(cfg.Auth.RefreshTokenTTLHour)*time.Hour)

	googleAuthClient := authG.NewGoogleAuthClient(cfg.Auth.GoogleClientID, cfg.Auth.GoogleClientSecret, cfg.Auth.GoogleCallbackURL, time.Duration(cfg.Auth.APITimeoutSec)*time.Second)

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

	googleUserUsecase := authU.NewGoogleUserUsecase(db, googleAuthClient, authTokenManager, registerAppUserCallback)
	guestUserUsecase := authU.NewGuestUserUsecase(authTokenManager)
	studentUsecaseWorkbook := studentU.NewStudentUsecaseWorkbook(db, pf, rfFunc, userRfFunc)
	studentUsecaseProblem := studentU.NewStudentUsecaseProblem(db, pf, rfFunc, userRfFunc)
	studentUseCaseStudy := studentU.NewStudentUsecaseStudy(db, pf, rfFunc, userRfFunc)
	studentUsecaseAudio := studentU.NewStudentUsecaseAudio(db, pf, rfFunc, userRfFunc, synthesizerClient)

	router := controller.NewRouter(googleUserUsecase, guestUserUsecase, studentUsecaseWorkbook, studentUsecaseProblem, studentUsecaseAudio, studentUseCaseStudy, translatorClient, tatoebaClient, newIteratorFunc, corsConfig, cfg.App, cfg.Auth, cfg.Debug)

	// router := gin.New()
	// router.Use(cors.New(corsConfig))
	// router.Use(gin.Recovery())

	// if cfg.Debug.GinMode {
	// 	router.Use(ginlog.Middleware(ginlog.DefaultConfig))
	// }

	// if cfg.Debug.Wait {
	// 	router.Use(middleware.NewWaitMiddleware())
	// }

	// router.GET("/healthcheck", func(c *gin.Context) {
	// 	c.Status(http.StatusOK)
	// })

	// v1 := router.Group("v1")
	// {
	// 	v1.Use(otelgin.Middleware(cfg.App.Name))
	// 	v1.Use(middleware.NewTraceLogMiddleware(cfg.App.Name))
	// 	v1auth := v1.Group("auth")
	// 	googleUserUsecase := authU.NewGoogleUserUsecase(db, googleAuthClient, authTokenManager, registerAppUserCallback)
	// 	guestUserUsecase := authU.NewGuestUserUsecase(authTokenManager)
	// 	authHandler := authH.NewAuthHandler(authTokenManager)
	// 	googleAuthHandler := authH.NewGoogleAuthHandler(googleUserUsecase)
	// 	guestAuthHandler := authH.NewGuestAuthHandler(guestUserUsecase)
	// 	v1auth.POST("google/authorize", googleAuthHandler.Authorize)
	// 	v1auth.POST("guest/authorize", guestAuthHandler.Authorize)
	// 	v1auth.POST("refresh_token", authHandler.RefreshToken)

	// 	v1Workbook := v1.Group("private/workbook")
	// 	studentUsecaseWorkbook := studentU.NewStudentUsecaseWorkbook(db, pf, rfFunc, userRfFunc)
	// 	privateWorkbookHandler := appC.NewPrivateWorkbookHandler(studentUsecaseWorkbook)
	// 	v1Workbook.Use(authMiddleware)
	// 	v1Workbook.POST(":workbookID", privateWorkbookHandler.FindWorkbooks)
	// 	v1Workbook.GET(":workbookID", privateWorkbookHandler.FindWorkbookByID)
	// 	v1Workbook.PUT(":workbookID", privateWorkbookHandler.UpdateWorkbook)
	// 	v1Workbook.DELETE(":workbookID", privateWorkbookHandler.RemoveWorkbook)
	// 	v1Workbook.POST("", privateWorkbookHandler.AddWorkbook)

	// 	v1Problem := v1.Group("workbook/:workbookID/problem")
	// 	studentUsecaseProblem := studentU.NewStudentUsecaseProblem(db, pf, rfFunc, userRfFunc)
	// 	problemHandler := appC.NewProblemHandler(studentUsecaseProblem, newIteratorFunc)
	// 	v1Problem.Use(authMiddleware)
	// 	v1Problem.POST("", problemHandler.AddProblem)
	// 	v1Problem.GET(":problemID", problemHandler.FindProblemByID)
	// 	v1Problem.DELETE(":problemID", problemHandler.RemoveProblem)
	// 	v1Problem.PUT(":problemID", problemHandler.UpdateProblem)
	// 	// v1Problem.GET("problem_ids", problemHandler.FindProblemIDs)
	// 	v1Problem.POST("find", problemHandler.FindProblems)
	// 	v1Problem.POST("find_all", problemHandler.FindAllProblems)
	// 	v1Problem.POST("find_by_ids", problemHandler.FindProblemsByProblemIDs)
	// 	v1Problem.POST("import", problemHandler.ImportProblems)

	// 	v1Study := v1.Group("study/workbook/:workbookID")
	// 	studentUseCaseStudy := studentU.NewStudentUsecaseStudy(db, pf, rfFunc, userRfFunc)
	// 	recordbookHandler := appC.NewRecordbookHandler(studentUseCaseStudy)
	// 	v1Study.Use(authMiddleware)
	// 	v1Study.GET("study_type/:studyType", recordbookHandler.FindRecordbook)
	// 	v1Study.POST("study_type/:studyType/problem/:problemID/record", recordbookHandler.SetStudyResult)
	// 	v1Study.GET("completion_rate", recordbookHandler.GetCompletionRate)

	// 	v1Audio := v1.Group("workbook/:workbookID/problem/:problemID/audio")
	// 	studentUsecaseAudio := studentU.NewStudentUsecaseAudio(db, pf, rfFunc, userRfFunc, synthesizerClient)
	// 	audioHandler := appC.NewAudioHandler(studentUsecaseAudio)
	// 	v1Audio.Use(authMiddleware)
	// 	v1Audio.GET(":audioID", audioHandler.FindAudioByID)
	// }

	// plugin := router.Group("plugin")
	// {
	// 	plugin.Use(otelgin.Middleware(cfg.App.Name))
	// 	plugin.Use(middleware.NewTraceLogMiddleware(cfg.App.Name))
	// 	plugin.Use(authMiddleware)
	// 	{
	// 		pluginTranslation := plugin.Group("translation")
	// 		translationHandler := pluginCommonHandler.NewTranslationHandler(translatorClient)
	// 		pluginTranslation.POST("find", translationHandler.FindTranslations)
	// 		pluginTranslation.GET("text/:text/pos/:pos", translationHandler.FindTranslationByTextAndPos)
	// 		pluginTranslation.GET("text/:text", translationHandler.FindTranslationsByText)
	// 		pluginTranslation.PUT("text/:text/pos/:pos", translationHandler.UpdateTranslation)
	// 		pluginTranslation.DELETE("text/:text/pos/:pos", translationHandler.RemoveTranslation)
	// 		pluginTranslation.POST("", translationHandler.AddTranslation)
	// 		pluginTranslation.POST("export", translationHandler.ExportTranslations)
	// 	}
	// 	{
	// 		pluginTatoeba := plugin.Group("tatoeba")
	// 		tatoebaHandler := pluginCommonHandler.NewTatoebaHandler(tatoebaClient)
	// 		pluginTatoeba.POST("find", tatoebaHandler.FindSentencePairs)
	// 		pluginTatoeba.POST("sentence/import", tatoebaHandler.ImportSentences)
	// 		pluginTatoeba.POST("link/import", tatoebaHandler.ImportLinks)
	// 	}
	// }

	if cfg.Swagger.Enabled {
		router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
		docs.SwaggerInfo.Title = cfg.App.Name
		docs.SwaggerInfo.Version = "1.0"
		docs.SwaggerInfo.Host = cfg.Swagger.Host
		docs.SwaggerInfo.Schemes = []string{cfg.Swagger.Schema}
	}

	httpServer := http.Server{
		Addr:    ":" + strconv.Itoa(cfg.App.HTTPPort),
		Handler: router,
	}

	logrus.Printf("http server listening at %v", httpServer.Addr)

	errCh := make(chan error)
	go func() {
		defer close(errCh)
		if err := httpServer.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			logrus.Infof("failed to ListenAndServe. err: %v", err)
			errCh <- err
		}
	}()

	select {
	case <-ctx.Done():
		gracefulShutdownTime1 := time.Duration(cfg.Shutdown.TimeSec1) * time.Second
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), gracefulShutdownTime1)
		defer shutdownCancel()
		if err := httpServer.Shutdown(shutdownCtx); err != nil {
			logrus.Infof("Server forced to shutdown. err: %v", err)
			return err
		}
		return nil
	case err := <-errCh:
		return err
	}
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

func initialize(ctx context.Context, env string) (*config.Config, *gorm.DB, *sql.DB, *sdktrace.TracerProvider, error) {
	cfg, err := config.LoadConfig(env)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	// init log
	if err := config.InitLog(env, cfg.Log); err != nil {
		return nil, nil, nil, nil, err
	}

	// tracer
	tp, err := config.InitTracerProvider(cfg)
	if err != nil {
		return nil, nil, nil, nil, xerrors.Errorf("failed to InitTracerProvider. err: %w", err)
	}
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

	// init db
	db, sqlDB, err := config.InitDB(cfg.DB)
	if err != nil {
		return nil, nil, nil, nil, xerrors.Errorf("failed to InitDB. err: %w", err)
	}

	userRfFunc := func(ctx context.Context, db *gorm.DB) (userS.RepositoryFactory, error) {
		return userG.NewRepositoryFactory(db)
	}
	userS.InitSystemAdmin(userRfFunc)

	return cfg, db, sqlDB, tp, nil
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
