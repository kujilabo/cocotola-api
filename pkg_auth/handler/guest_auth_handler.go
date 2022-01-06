package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/kujilabo/cocotola-api/pkg_auth/application"
	"github.com/kujilabo/cocotola-api/pkg_auth/handler/entity"
	"github.com/kujilabo/cocotola-api/pkg_lib/log"
)

type GuestAuthHandler interface {
	Authorize(c *gin.Context)
}

type guestAuthHandler struct {
	guestAuthService application.GuestAuthService
}

func NewGuestAuthHandler(guestAuthService application.GuestAuthService) GuestAuthHandler {
	return &guestAuthHandler{
		guestAuthService: guestAuthService,
	}
}

func (h *guestAuthHandler) Authorize(c *gin.Context) {
	ctx := c.Request.Context()
	logger := log.FromContext(ctx)
	logger.Info("Authorize")

	guestAuthParameter := entity.GuestAuthParameter{}
	if err := c.BindJSON(&guestAuthParameter); err != nil {
		logger.Warnf("invalid parameter. err: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"message": http.StatusText(http.StatusBadRequest)})
		return
	}

	authResult, err := h.guestAuthService.RetrieveGuestToken(ctx, guestAuthParameter.OrganizationName)
	if err != nil {
		logger.Warnf("failed to RetrieveGuestToken. err: %+v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": http.StatusText(http.StatusBadRequest)})
		return
	}

	logger.Info("Authorize OK")
	c.JSON(http.StatusOK, entity.AuthResponse{
		AccessToken:  authResult.AccessToken,
		RefreshToken: authResult.RefreshToken,
	})
}
