package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/kujilabo/cocotola-api/src/auth/handler/entity"
	"github.com/kujilabo/cocotola-api/src/auth/usecase"
	"github.com/kujilabo/cocotola-api/src/lib/log"
)

type GuestUserHandler interface {
	Authorize(c *gin.Context)
}

type guestUserHandler struct {
	guestUserUsecase usecase.GuestUserUsecase
}

func NewGuestAuthHandler(guestUserUsecase usecase.GuestUserUsecase) GuestUserHandler {
	return &guestUserHandler{
		guestUserUsecase: guestUserUsecase,
	}
}

func (h *guestUserHandler) Authorize(c *gin.Context) {
	ctx := c.Request.Context()
	logger := log.FromContext(ctx)
	logger.Info("Authorize")

	guestAuthParameter := entity.GuestAuthParameter{}
	if err := c.BindJSON(&guestAuthParameter); err != nil {
		logger.Warnf("invalid parameter. err: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"message": http.StatusText(http.StatusBadRequest)})
		return
	}

	authResult, err := h.guestUserUsecase.RetrieveGuestToken(ctx, guestAuthParameter.OrganizationName)
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
