package controller

import (
	"context"
	"io"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	ginlog "github.com/onrik/logrus/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"

	"github.com/kujilabo/cocotola-api/src/app/config"
	appD "github.com/kujilabo/cocotola-api/src/app/domain"
	appS "github.com/kujilabo/cocotola-api/src/app/service"
	studentU "github.com/kujilabo/cocotola-api/src/app/usecase/student"
	authH "github.com/kujilabo/cocotola-api/src/auth/controller"
	authM "github.com/kujilabo/cocotola-api/src/auth/controller/middleware"
	authG "github.com/kujilabo/cocotola-api/src/auth/gateway"
	authU "github.com/kujilabo/cocotola-api/src/auth/usecase"
	ginmiddleware "github.com/kujilabo/cocotola-api/src/lib/controller/middleware"
	pluginCommonController "github.com/kujilabo/cocotola-api/src/plugin/common/controller"
	pluginCommonService "github.com/kujilabo/cocotola-api/src/plugin/common/service"
)

type NewIteratorFunc func(ctx context.Context, workbookID appD.WorkbookID, problemType string, reader io.Reader) (appS.ProblemAddParameterIterator, error)

func NewRouter(googleUserUsecase authU.GoogleUserUsecase, guestUserUsecase authU.GuestUserUsecase, studentUsecaseWorkbook studentU.StudentUsecaseWorkbook, studentUsecaseProblem studentU.StudentUsecaseProblem, studentUsecaseAudio studentU.StudentUsecaseAudio, studentUsecaseStudy studentU.StudentUsecaseStudy, translatorClient pluginCommonService.TranslatorClient, tatoebaClient pluginCommonService.TatoebaClient, newIteratorFunc NewIteratorFunc, corsConfig cors.Config, appConfig *config.AppConfig, authConfig *config.AuthConfig, debugConfig *config.DebugConfig) *gin.Engine {
	if !debugConfig.GinMode {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()
	router.Use(cors.New(corsConfig))
	router.Use(gin.Recovery())

	if debugConfig.GinMode {
		router.Use(ginlog.Middleware(ginlog.DefaultConfig))
	}

	if debugConfig.Wait {
		router.Use(ginmiddleware.NewWaitMiddleware())
	}

	router.GET("/healthcheck", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	signingKey := []byte(authConfig.SigningKey)
	signingMethod := jwt.SigningMethodHS256
	authTokenManager := authG.NewAuthTokenManager(signingKey, signingMethod, time.Duration(authConfig.AccessTokenTTLMin)*time.Minute, time.Duration(authConfig.RefreshTokenTTLHour)*time.Hour)
	authMiddleware := authM.NewAuthMiddleware(signingKey)

	v1 := router.Group("v1")
	{
		v1.Use(otelgin.Middleware(appConfig.Name))
		v1.Use(ginmiddleware.NewTraceLogMiddleware(appConfig.Name))
		v1auth := v1.Group("auth")
		// googleUserUsecase := authU.NewGoogleUserUsecase(db, googleAuthClient, authTokenManager, registerAppUserCallback)
		// guestUserUsecase := authU.NewGuestUserUsecase(authTokenManager)
		authHandler := authH.NewAuthHandler(authTokenManager)
		googleAuthHandler := authH.NewGoogleAuthHandler(googleUserUsecase)
		guestAuthHandler := authH.NewGuestAuthHandler(guestUserUsecase)
		v1auth.POST("google/authorize", googleAuthHandler.Authorize)
		v1auth.POST("guest/authorize", guestAuthHandler.Authorize)
		v1auth.POST("refresh_token", authHandler.RefreshToken)

		v1Workbook := v1.Group("private/workbook")
		privateWorkbookHandler := NewPrivateWorkbookHandler(studentUsecaseWorkbook)
		v1Workbook.Use(authMiddleware)
		v1Workbook.POST(":workbookID", privateWorkbookHandler.FindWorkbooks)
		v1Workbook.GET(":workbookID", privateWorkbookHandler.FindWorkbookByID)
		v1Workbook.PUT(":workbookID", privateWorkbookHandler.UpdateWorkbook)
		v1Workbook.DELETE(":workbookID", privateWorkbookHandler.RemoveWorkbook)
		v1Workbook.POST("", privateWorkbookHandler.AddWorkbook)

		v1Problem := v1.Group("workbook/:workbookID/problem")
		problemHandler := NewProblemHandler(studentUsecaseProblem, newIteratorFunc)
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
		recordbookHandler := NewRecordbookHandler(studentUsecaseStudy)
		v1Study.Use(authMiddleware)
		v1Study.GET("study_type/:studyType", recordbookHandler.FindRecordbook)
		v1Study.POST("study_type/:studyType/problem/:problemID/record", recordbookHandler.SetStudyResult)
		v1Study.GET("completion_rate", recordbookHandler.GetCompletionRate)

		v1Audio := v1.Group("workbook/:workbookID/problem/:problemID/audio")

		audioHandler := NewAudioHandler(studentUsecaseAudio)
		v1Audio.Use(authMiddleware)
		v1Audio.GET(":audioID", audioHandler.FindAudioByID)
	}

	plugin := router.Group("plugin")
	{
		plugin.Use(otelgin.Middleware(appConfig.Name))
		plugin.Use(ginmiddleware.NewTraceLogMiddleware(appConfig.Name))
		plugin.Use(authMiddleware)

		InitTranslatorPluginRouter(plugin, translatorClient)
		InitTatoebaPluginRouter(plugin, tatoebaClient)
	}

	return router
}

func InitTranslatorPluginRouter(plugin *gin.RouterGroup, translatorClient pluginCommonService.TranslatorClient) {

	pluginTranslation := plugin.Group("translation")
	translationHandler := pluginCommonController.NewTranslationHandler(translatorClient)
	pluginTranslation.POST("find", translationHandler.FindTranslations)
	pluginTranslation.GET("text/:text/pos/:pos", translationHandler.FindTranslationByTextAndPos)
	pluginTranslation.GET("text/:text", translationHandler.FindTranslationsByText)
	pluginTranslation.PUT("text/:text/pos/:pos", translationHandler.UpdateTranslation)
	pluginTranslation.DELETE("text/:text/pos/:pos", translationHandler.RemoveTranslation)
	pluginTranslation.POST("", translationHandler.AddTranslation)
	pluginTranslation.POST("export", translationHandler.ExportTranslations)
}

func InitTatoebaPluginRouter(plugin *gin.RouterGroup, tatoebaClient pluginCommonService.TatoebaClient) {
	pluginTatoeba := plugin.Group("tatoeba")
	tatoebaHandler := pluginCommonController.NewTatoebaHandler(tatoebaClient)
	pluginTatoeba.POST("find", tatoebaHandler.FindSentencePairs)
	pluginTatoeba.POST("sentence/import", tatoebaHandler.ImportSentences)
	pluginTatoeba.POST("link/import", tatoebaHandler.ImportLinks)
}
