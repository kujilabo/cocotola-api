package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"strconv"
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
	libD "github.com/kujilabo/cocotola-api/pkg_lib/domain"
	libG "github.com/kujilabo/cocotola-api/pkg_lib/gateway"
	"github.com/kujilabo/cocotola-api/pkg_lib/handler/middleware"
	"github.com/kujilabo/cocotola-api/pkg_lib/log"
	pluginCommonDomain "github.com/kujilabo/cocotola-api/pkg_plugin/common/domain"
	pluginCommonGateway "github.com/kujilabo/cocotola-api/pkg_plugin/common/gateway"
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

	// cfg, err := config.LoadConfig(*env)
	// if err != nil {
	// 	panic(err)
	// }

	// // init log
	// if err := config.InitLog(*env, cfg.Log); err != nil {
	// 	panic(err)
	// }

	// // cors
	// corsConfig := config.InitCORS(cfg.CORS)
	// logrus.Infof("cors: %+v", corsConfig)

	// if err := corsConfig.Validate(); err != nil {
	// 	panic(err)
	// }

	// // init db
	// db, sqlDB, err := initDB(cfg.DB)
	// if err != nil {
	// 	fmt.Printf("failed to InitDB. err: %+v", err)
	// 	panic(err)
	// }
	// defer sqlDB.Close()

	// rf, err := userG.NewRepositoryFactory(db)
	// if err != nil {
	// 	panic(err)
	// }

	// userD.InitSystemAdmin(rf)

	// if err := initApp(ctx, db, cfg.App.OwnerPassword); err != nil {
	// 	panic(err)
	// }

	// if !cfg.Debug.GinMode {
	// 	gin.SetMode(gin.ReleaseMode)
	// }

	// router := gin.New()
	// router.Use(cors.New(corsConfig))
	// router.Use(ginlog.Middleware(ginlog.DefaultConfig))
	// router.Use(middleware.NewLogMiddleware(), gin.Recovery())

	// if cfg.Debug.Wait {
	// 	router.Use(middleware.NewWaitMiddleware())
	// }

	cfg, db, sqlDB, router, err := initialize(ctx, *env)
	if err != nil {
		panic(err)
	}
	defer sqlDB.Close()

	router.GET("/healthcheck", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	synthesizer := pluginCommonGateway.NewSynthesizer(cfg.Google.SynthesizerKey, time.Duration(cfg.Google.SynthesizerTimeoutSec)*time.Minute)

	translatorClient := pluginCommonGateway.NewAzureTranslatorClient(cfg.Azure.SubscriptionKey)
	translatorRepository := pluginCommonGateway.NewAzureTranslationRepository(db)
	translator := pluginCommonGateway.NewAzureCachedTranslatorClient(translatorClient, translatorRepository)

	englishWordProblemProcessor := pluginEnglishDomain.NewEnglishWordProblemProcessor(synthesizer, translator, pluginEnglishGateway.NewEnglishWordProblemAddParameterCSVReader)
	englishPhraseProblemProcessor := pluginEnglishDomain.NewEnglishPhraseProblemProcessor(synthesizer, translator)
	englishSentenceProblemProcessor := pluginEnglishDomain.NewEnglishSentenceProblemProcessor(synthesizer, translator, pluginEnglishGateway.NewEnglishSentenceProblemAddParameterCSVReader)

	problemAddProcessor := map[string]appD.ProblemAddProcessor{
		pluginEnglishDomain.EnglishWordProblemType:     englishWordProblemProcessor,
		pluginEnglishDomain.EnglishPhraseProblemType:   englishPhraseProblemProcessor,
		pluginEnglishDomain.EnglishSentenceProblemType: englishSentenceProblemProcessor,
	}
	problemRemoveProcessor := map[string]appD.ProblemRemoveProcessor{
		pluginEnglishDomain.EnglishWordProblemType:     englishWordProblemProcessor,
		pluginEnglishDomain.EnglishPhraseProblemType:   englishPhraseProblemProcessor,
		pluginEnglishDomain.EnglishSentenceProblemType: englishSentenceProblemProcessor,
	}
	problemImportProcessor := map[string]appD.ProblemImportProcessor{
		pluginEnglishDomain.EnglishWordProblemType: englishWordProblemProcessor,
	}

	englishWordProblemRepository := func(db *gorm.DB) (appD.ProblemRepository, error) {
		return pluginEnglishGateway.NewEnglishWordProblemRepository(db, pluginEnglishDomain.EnglishWordProblemType)
	}
	englishPhraseProblemRepository := func(db *gorm.DB) (appD.ProblemRepository, error) {
		return pluginEnglishGateway.NewEnglishPhraseProblemRepository(db, pluginEnglishDomain.EnglishPhraseProblemType)
	}
	englishSentenceProblemRepository := func(db *gorm.DB) (appD.ProblemRepository, error) {
		return pluginEnglishGateway.NewEnglishSentenceProblemRepository(db, pluginEnglishDomain.EnglishSentenceProblemType)
	}

	pf := appD.NewProcessorFactory(problemAddProcessor, problemRemoveProcessor, problemImportProcessor)
	problemRepositories := map[string]func(*gorm.DB) (appD.ProblemRepository, error){
		pluginEnglishDomain.EnglishWordProblemType:     englishWordProblemRepository,
		pluginEnglishDomain.EnglishPhraseProblemType:   englishPhraseProblemRepository,
		pluginEnglishDomain.EnglishSentenceProblemType: englishSentenceProblemRepository,
	}

	newIterator := func(ctx context.Context, workbookID appD.WorkbookID, problemType string, reader io.Reader) (appD.ProblemAddParameterIterator, error) {
		processor, ok := problemImportProcessor[problemType]
		if ok {
			return processor.CreateCSVReader(ctx, workbookID, problemType, reader)
		}
		return nil, fmt.Errorf("processor not found. problemType: %s", problemType)
	}

	userRepoFunc := func(db *gorm.DB) (userD.RepositoryFactory, error) {
		return userG.NewRepositoryFactory(db)
	}
	repoFunc := func(db *gorm.DB) (appD.RepositoryFactory, error) {
		return appG.NewRepositoryFactory(context.Background(), db, cfg.DB.DriverName, userRepoFunc, pf, problemRepositories)
	}

	signingKey := []byte(cfg.Auth.SigningKey)
	signingMethod := jwt.SigningMethodHS256
	authTokenManager := authG.NewAuthTokenManager(signingKey, signingMethod, time.Duration(cfg.Auth.AccessTokenTTLMin)*time.Minute, time.Duration(cfg.Auth.RefreshTokenTTLHour)*time.Hour)

	googleAuthClient := authG.NewGoogleAuthClient(cfg.Auth.GoogleClientID, cfg.Auth.GoogleClientSecret, cfg.Auth.GoogleCallbackURL)
	authMiddleware := authM.NewAuthMiddleware(signingKey)

	registerAppUserCallback := func(ctx context.Context, organizationName string, appUser userD.AppUser) error {
		repo, err := repoFunc(db)
		if err != nil {
			return err
		}
		userRepo, err := userRepoFunc(db)
		if err != nil {
			return err
		}
		return callback(ctx, cfg.App.TestUserEmail, repo, userRepo, organizationName, appUser)
	}

	v1 := router.Group("v1")
	{
		v1auth := v1.Group("auth")
		googleAuthService := authA.NewGoogleAuthService(userRepoFunc, googleAuthClient, authTokenManager, registerAppUserCallback)
		authHandler := authH.NewAuthHandler(authTokenManager)
		googleAuthHandler := authH.NewGoogleAuthHandler(googleAuthService)
		v1auth.POST("google/authorize", googleAuthHandler.Authorize)
		v1auth.POST("refresh_token", authHandler.RefreshToken)

		privateWorkbookService := application.NewPrivateWorkbookService(db, repoFunc, userRepoFunc)
		privateWorkbookHandler := appH.NewPrivateWorkbookHandler(privateWorkbookService)
		v1Workbook := v1.Group("private/workbook")
		v1Workbook.Use(authMiddleware)
		v1Workbook.POST(":workbookID", privateWorkbookHandler.FindWorkbooks)
		v1Workbook.GET(":workbookID", privateWorkbookHandler.FindWorkbookByID)
		v1Workbook.POST("", privateWorkbookHandler.AddWorkbook)

		problemService := application.NewProblemService(db, repoFunc, userRepoFunc)
		problemHandler := appH.NewProblemHandler(problemService, newIterator)
		v1Problem := v1.Group("workbook/:workbookID")
		v1Problem.Use(authMiddleware)
		v1Problem.POST("problem", problemHandler.AddProblem)
		v1Problem.GET("problem/:problemID", problemHandler.FindProblemByID)
		v1Problem.DELETE("problem/:problemID", problemHandler.RemoveProblem)
		// v1Problem.GET("problem_ids", problemHandler.FindProblemIDs)
		v1Problem.POST("problem/search", problemHandler.FindProblems)
		v1Problem.POST("problem/find_all", problemHandler.FindAllProblems)
		v1Problem.POST("problem/search_by_ids", problemHandler.FindProblemsByProblemIDs)
		v1Problem.POST("problem/import", problemHandler.ImportProblems)

		studyService := application.NewStudyService(db, repoFunc, userRepoFunc)
		studyHandler := appH.NewStudyHandler(studyService)
		v1Study := v1.Group("study/workbook/:workbookID")
		v1Study.Use(authMiddleware)
		v1Study.GET("study_type/:studyType", studyHandler.FindRecordbook)
		v1Study.POST("study_type/:studyType/problem/:problemID/result", studyHandler.SetStudyResult)

		audioService := application.NewAudioService(db, repoFunc)
		audioHandler := appH.NewAudioHandler(audioService)
		v1Audio := v1.Group("audio")
		v1Audio.Use(authMiddleware)
		v1Audio.GET(":audioID", audioHandler.FindAudioByID)
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

	rf, err := userG.NewRepositoryFactory(db)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	userD.InitSystemAdmin(rf)

	if err := initApp(ctx, db, cfg.App.OwnerPassword); err != nil {
		return nil, nil, nil, nil, err
	}

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

func initApp(ctx context.Context, db *gorm.DB, password string) error {
	logger := log.FromContext(ctx)
	systemAdmin := userD.SystemAdminInstance()
	if err := db.Transaction(func(tx *gorm.DB) error {
		organization, err := systemAdmin.FindOrganizationByName(ctx, "cocotola")
		if err != nil {
			if !xerrors.Is(err, userD.ErrOrganizationNotFound) {
				return fmt.Errorf("failed to AddOrganization: %w", err)
			}

			firstOwnerAddParam, err := userD.NewFirstOwnerAddParameter("cocotola-owner", password, "Owner(cocotola)")
			if err != nil {
				return fmt.Errorf("failed to AddOrganization: %w", err)
			}

			organizationAddParameter, err := userD.NewOrganizationAddParameter("cocotola", firstOwnerAddParam)
			if err != nil {
				return fmt.Errorf("failed to AddOrganization: %w", err)
			}

			organizationID, err := systemAdmin.AddOrganization(ctx, organizationAddParameter)
			if err != nil {
				return fmt.Errorf("failed to AddOrganization: %w", err)
			}

			logger.Infof("organizationID: %d", organizationID)
			return nil
		}
		logger.Infof("organization: %d", organization)
		return nil
	}); err != nil {
		return err
	}

	return nil
}

func callback(ctx context.Context, testUserEmail string, repo appD.RepositoryFactory, userRepo userD.RepositoryFactory, organizationName string, appUser userD.AppUser) error {
	logger := log.FromContext(ctx)
	logger.Infof("callback. loginID: %s", appUser.GetLoginID())

	if appUser.GetLoginID() == testUserEmail {
		logger.Info("%s", appUser.GetLoginID())
		student, err := appD.NewStudent(repo, userRepo, appUser)
		if err != nil {
			return xerrors.Errorf("failed to NewStudent. err: %w", err)
		}

		if err := createDemoWordWorkbook(ctx, student, "Example", pluginCommonDomain.PosNoun, []string{"butcher", "bakery", "library", "bookstore", "drugstore", "restaurant", "garage", "barbershop", "bank", "market"}); err != nil {
			return err
		}

		if err := createDemoWordWorkbook(ctx, student, "NGSL-100", pluginCommonDomain.PosOther, []string{
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
	}

	return nil
}

func createDemoWordWorkbook(ctx context.Context, student appD.Student, workbookName string, pos pluginCommonDomain.WordPos, words []string) error {
	logger := log.FromContext(ctx)

	workbookProperties := map[string]string{
		"audioEnabled": "false",
	}
	param, err := appD.NewWorkbookAddParameter(pluginEnglishDomain.EnglishWordProblemType, workbookName, "", workbookProperties)
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
		param, err := appD.NewProblemAddParameter(workbookID, i+1, pluginEnglishDomain.EnglishWordProblemType, properties)
		if err != nil {
			return xerrors.Errorf("failed to NewProblemAddParameter. err: %w", err)
		}

		problemID, err := workbook.AddProblem(ctx, student, param)
		if err != nil {
			return xerrors.Errorf("failed to NewProblemAddParameter. err: %w", err)
		}
		logger.Infof("problemID: %d", problemID)
	}

	logger.Infof("Example %d", workbookID)
	return nil
}
