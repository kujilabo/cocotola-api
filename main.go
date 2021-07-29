package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"golang.org/x/xerrors"
	"gorm.io/gorm"

	swaggerFiles "github.com/swaggo/files"     // swagger embed files
	ginSwagger "github.com/swaggo/gin-swagger" // gin-swagger middleware

	"github.com/kujilabo/cocotola-api/docs"
	"github.com/kujilabo/cocotola-api/pkg_app/config"
	"github.com/kujilabo/cocotola-api/pkg_app/gateway"
	authA "github.com/kujilabo/cocotola-api/pkg_auth/application"
	authG "github.com/kujilabo/cocotola-api/pkg_auth/gateway"
	authH "github.com/kujilabo/cocotola-api/pkg_auth/handler"
	libG "github.com/kujilabo/cocotola-api/pkg_lib/gateway"
	"github.com/kujilabo/cocotola-api/pkg_lib/handler/middleware"
	"github.com/kujilabo/cocotola-api/pkg_lib/log"
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
			*env = "development"
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

	cfg, err := config.LoadConfig(*env)
	if err != nil {
		panic(err)
	}

	// init log
	if err := config.InitLog(*env, cfg.Log); err != nil {
		panic(err)
	}

	// cors
	corsConfig := config.InitCORS(cfg.CORS)
	logrus.Infof("cors: %+v", corsConfig)

	if err := corsConfig.Validate(); err != nil {
		panic(err)
	}

	// init db
	db, sqlDB, err := initDB()
	if err != nil {
		fmt.Printf("Failed to InitDB. err: %+v1", err)
		panic(err)
	}
	defer sqlDB.Close()

	rf := userG.NewRepositoryFactory(db)
	userD.InitSystemAdmin(rf)

	if err := initApp(ctx, db, cfg.App.OwnerPassword); err != nil {
		panic(err)
	}

	if !cfg.Debug.GinMode {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()
	router.Use(cors.New(corsConfig))
	router.Use(middleware.NewLogMiddleware())

	if cfg.Debug.Wait {
		router.Use(middleware.NewWaitMiddleware())
	}

	router.GET("/healthcheck", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	userRepoFunc := func(db *gorm.DB) userD.RepositoryFactory {
		return userG.NewRepositoryFactory(db)
	}
	signingKey := []byte(cfg.Auth.SigningKey)
	signingMethod := jwt.SigningMethodHS256
	authTokenManager := authG.NewAuthTokenManager(signingKey, signingMethod, time.Duration(cfg.Auth.AccessTokenTTLMin)*time.Minute, time.Duration(cfg.Auth.RefreshTokenTTLHour)*time.Hour)

	googleAuthClient := authG.NewGoogleAuthClient(cfg.Auth.GoogleClientID, cfg.Auth.GoogleClientSecret, cfg.Auth.GoogleCallbackURL)
	// authMiddleware := authM.NewAuthMiddleware(signingKey)

	registerAppUsedrCallback := func(ctx context.Context, organizationName string, appUser userD.AppUser) error {
		logger := log.FromContext(ctx)
		logger.Infof("%s", appUser.GetLoginID())

		if appUser.GetLoginID() == cfg.App.TestUserEmail {
			logger.Info("%s", appUser.GetLoginID())
		}
		return nil
	}
	googleAuthService := authA.NewGoogleAuthService(userRepoFunc, googleAuthClient, authTokenManager, registerAppUsedrCallback)
	authHandler := authH.NewAuthHandler(authTokenManager)
	googleAuthHandler := authH.NewGoogleAuthHandler(googleAuthService)
	v1 := router.Group("v1")
	{
		v1auth := v1.Group("auth")
		v1auth.POST("google/authorize", googleAuthHandler.Authorize)
		v1auth.POST("refresh_token", authHandler.RefreshToken)
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

func initDB() (*gorm.DB, *sql.DB, error) {
	// init db
	db, err := libG.OpenSQLite("./app.db")
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

	if err := gateway.MigrateSQLiteDB(db); err != nil {
		return nil, nil, err
	}

	return db, sqlDB, nil
}

func initApp(ctx context.Context, db *gorm.DB, password string) error {
	logger := log.FromContext(ctx)
	systemAdmin := userD.SystemAdminInstance()
	// repository := gateway.NewRepository(db)
	if err := db.Transaction(func(tx *gorm.DB) error {
		// repositoryFactory := gateway.NewRepositoryFactory(db, gh)
		organization, err := systemAdmin.FindOrganizationByName(ctx, "cocotola")
		if err != nil {
			if !xerrors.Is(err, userD.ErrOrganizationNotFound) {
				return fmt.Errorf("failed to AddOrganization: %w", err)
			}

			firstOwnerAddParam, err := userD.NewFirstOwnerAddParameter("cocotola-owner", password, "Owner(cocotola)")
			if err != nil {
				return fmt.Errorf("failed to AddOrganization: %w", err)
			}
			organizationAddParameter, err := userD.NewOrganizationAddParameter(
				"cocotola", firstOwnerAddParam)
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
