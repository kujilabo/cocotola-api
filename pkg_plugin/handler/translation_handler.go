package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kujilabo/cocotola-api/pkg_lib/log"
	user "github.com/kujilabo/cocotola-api/pkg_user/domain"
	"github.com/kujilabo/cocotola-api/pkg_user/handlerhelper"
)

type TranslationHandler interface {
	FindTranslations(c *gin.Context)
}

type translationHandler struct {
}

func NewTranslationHandler() TranslationHandler {
	return &translationHandler{}
}

func (h *translationHandler) FindTranslations(c *gin.Context) {

	ctx := c.Request.Context()
	logger := log.FromContext(ctx)
	logger.Infof("FindTranslations")

	handlerhelper.HandleRoleFunction(c, "Owner", func(organizationID user.OrganizationID, operatorID user.AppUserID) error {
		c.Status(http.StatusOK)
		return nil
	}, h.errorHandle)
}

func (h *translationHandler) errorHandle(c *gin.Context, err error) bool {
	ctx := c.Request.Context()
	logger := log.FromContext(ctx)
	logger.Errorf("translationHandler error:%v", err)
	return false
}
