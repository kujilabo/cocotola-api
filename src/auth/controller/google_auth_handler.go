package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/kujilabo/cocotola-api/src/auth/controller/entity"
	"github.com/kujilabo/cocotola-api/src/auth/usecase"
	"github.com/kujilabo/cocotola-api/src/lib/log"
)

type GoogleUserHandler interface {
	Authorize(c *gin.Context)
}

type googleUserHandler struct {
	googleUserUsecase usecase.GoogleUserUsecase
}

func NewGoogleAuthHandler(googleUserUsecase usecase.GoogleUserUsecase) GoogleUserHandler {
	return &googleUserHandler{
		googleUserUsecase: googleUserUsecase,
	}
}

func (h *googleUserHandler) Authorize(c *gin.Context) {
	ctx := c.Request.Context()
	logger := log.FromContext(ctx)
	logger.Info("Authorize")

	googleAuthParameter := entity.GoogleAuthParameter{}
	if err := c.BindJSON(&googleAuthParameter); err != nil {
		logger.Warnf("invalid parameter. err: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"message": http.StatusText(http.StatusBadRequest)})
		return
	}

	logger.Infof("RetrieveAccessToken. code: %s", googleAuthParameter)
	googleAuthResponse, err := h.googleUserUsecase.RetrieveAccessToken(ctx, googleAuthParameter.Code)
	if err != nil {
		logger.Warnf("failed to RetrieveAccessToken. err: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"message": http.StatusText(http.StatusBadRequest)})
		return
	}

	logger.Infof("RetrieveUserInfo. googleResponse: %+v", googleAuthResponse)
	userInfo, err := h.googleUserUsecase.RetrieveUserInfo(ctx, googleAuthResponse)
	if err != nil {
		logger.Warnf("failed to RetrieveUserInfo. error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"message": http.StatusText(http.StatusBadRequest)})
		return
	}

	logger.Info("RegisterAppUser")
	authResult, err := h.googleUserUsecase.RegisterAppUser(ctx, userInfo, googleAuthResponse, googleAuthParameter.OrganizationName)
	if err != nil {
		logger.Warnf("failed to RegisterStudent. err: %+v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": http.StatusText(http.StatusBadRequest)})
		return
	}

	logger.Info("Authorize OK")
	c.JSON(http.StatusOK, entity.AuthResponse{
		AccessToken:  authResult.AccessToken,
		RefreshToken: authResult.RefreshToken,
	})
}
