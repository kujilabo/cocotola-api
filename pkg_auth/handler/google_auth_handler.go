package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/kujilabo/cocotola-api/pkg_auth/application"
	"github.com/kujilabo/cocotola-api/pkg_auth/handler/entity"
	"github.com/kujilabo/cocotola-api/pkg_lib/log"
)

type GoogleAuthHandler interface {
	Authorize(c *gin.Context)
}

type googleAuthHandler struct {
	googleAuthService application.GoogleAuthService
}

func NewGoogleAuthHandler(googleAuthService application.GoogleAuthService) GoogleAuthHandler {
	return &googleAuthHandler{
		googleAuthService: googleAuthService,
	}
}

func (h *googleAuthHandler) Authorize(c *gin.Context) {
	ctx := c.Request.Context()
	logger := log.FromContext(ctx)
	logger.Info("Authorize")

	googleAuthParameter := entity.GoogleAuthParameter{}
	if err := c.BindJSON(&googleAuthParameter); err != nil {
		logger.Info("Invalid parameter")
		c.JSON(http.StatusBadRequest, gin.H{"message": http.StatusText(http.StatusBadRequest)})
		return
	}

	logger.Infof("RetrieveAccessToken. code: %s", googleAuthParameter)
	googleAuthResponse, err := h.googleAuthService.RetrieveAccessToken(ctx, googleAuthParameter.Code)
	if err != nil {
		logger.Warnf("Failed to RetrieveAccessToken. err: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"message": http.StatusText(http.StatusBadRequest)})
		return
	}

	logger.Infof("RetrieveUserInfo. googleResponse: %+v", googleAuthResponse)
	userInfo, err := h.googleAuthService.RetrieveUserInfo(ctx, googleAuthResponse)
	if err != nil {
		logger.Warnf("Failed to RetrieveUserInfo. error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"message": http.StatusText(http.StatusBadRequest)})
		return
	}

	logger.Info("RegisterStudent")
	authResult, err := h.googleAuthService.RegisterStudent(ctx, userInfo, googleAuthResponse, googleAuthParameter.OrganizationName)
	if err != nil {
		logger.Warnf("Failed to RegisterStudent. err: %+v", err)
		c.JSON(http.StatusBadRequest, gin.H{"message": http.StatusText(http.StatusBadRequest)})
		return
	}

	logger.Info("Authorize OK")
	c.JSON(http.StatusOK, entity.AuthResponse{
		AccessToken:  authResult.AccessToken,
		RefreshToken: authResult.RefreshToken,
	})
}
