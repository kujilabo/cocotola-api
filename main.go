package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"golang.org/x/xerrors"
	"gorm.io/gorm"

	"github.com/kujilabo/cocotola-api/pkg_app/config"
	"github.com/kujilabo/cocotola-api/pkg_lib/handler/middleware"
	"github.com/kujilabo/cocotola-api/pkg_lib/log"
	userD "github.com/kujilabo/cocotola-api/pkg_user/domain"
)

func main() {
	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	// bg := context.Background()
	// ctx := log.With(bg)
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

	// cors
	corsConfig := config.InitCORS(cfg.CORS)
	logrus.Infof("cors: %+v", corsConfig)

	if err := corsConfig.Validate(); err != nil {
		panic(err)
	}

	router := gin.New()
	router.Use(cors.New(corsConfig))
	router.Use(middleware.NewLogMiddleware())
	if !cfg.Debug.GinMode {
		gin.SetMode(gin.ReleaseMode)
	}
	if cfg.Debug.Wait {
		router.Use(middleware.NewWaitMiddleware())
	}

	// signingKey := []byte(cfg.Auth.SigningKey)
	// signingMethod := jwt.SigningMethodHS256
	// authTokenManager := authG.NewAuthTokenManager(signingKey, signingMethod, time.Duration(5)*time.Minute, time.Duration(24*30)*time.Hour)
	// authMiddleware := authM.NewAuthMiddleware(signingKey)

	router.GET("/healthcheck", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

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

func initialize(ctx context.Context, db *gorm.DB, password string) error {
	logger := log.FromContext(ctx)
	systemAdmin := userD.SystemAdminInstance()
	// repository := gateway.NewRepository(db)
	if err := db.Transaction(func(tx *gorm.DB) error {
		// repositoryFactory := gateway.NewRepositoryFactory(db, gh)
		organization, err := systemAdmin.FindOrganizationByName(ctx, "cocotola")
		if err != nil {
			if xerrors.Is(err, userD.ErrOrganizationNotFound) {
				organizationAddParameter := &userD.OrganizationAddParameter{
					Name: "cocotola",
					FirstOwner: &userD.FirstOwnerAddParameter{
						LoginID:  "cocotola-owner",
						Password: password,
						Username: "Owner(cocotola)",
					},
				}
				organizationID, err := systemAdmin.AddOrganization(ctx, organizationAddParameter)
				if err != nil {
					return fmt.Errorf("failed to AddOrganization: %w", err)
				}
				logger.Infof("organizationID: %d", organizationID)
				return nil
			}
			logger.Errorf("failed to AddOrganization: %w", err)
			return fmt.Errorf("failed to AddOrganization: %w", err)
		}
		logger.Infof("organization: %d", organization)
		return nil
	}); err != nil {
		return err
	}

	return nil
}
