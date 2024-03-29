package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/kujilabo/cocotola-api/src/auth/controller/entity"
	"github.com/kujilabo/cocotola-api/src/auth/service"
	"github.com/kujilabo/cocotola-api/src/lib/log"
)

type AuthHandler interface {
	RefreshToken(c *gin.Context)
}

type authHandler struct {
	authTokenManager service.AuthTokenManager
}

func NewAuthHandler(authTokenManager service.AuthTokenManager) AuthHandler {
	return &authHandler{
		authTokenManager: authTokenManager,
	}
}

func (h *authHandler) RefreshToken(c *gin.Context) {
	ctx := c.Request.Context()
	logger := log.FromContext(ctx)
	logger.Info("Authorize")
	refreshTokenParameter := entity.RefreshTokenParameter{}
	if err := c.BindJSON(&refreshTokenParameter); err != nil {
		return
	}

	token, err := h.authTokenManager.RefreshToken(ctx, refreshTokenParameter.RefreshToken)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	logger.Info("Authorize OK")
	c.JSON(http.StatusOK, entity.AuthResponse{
		AccessToken: token,
	})
}
