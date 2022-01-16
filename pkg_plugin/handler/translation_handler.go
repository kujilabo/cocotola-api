package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	app "github.com/kujilabo/cocotola-api/pkg_app/domain"
	"github.com/kujilabo/cocotola-api/pkg_lib/ginhelper"
	"github.com/kujilabo/cocotola-api/pkg_lib/log"
	"github.com/kujilabo/cocotola-api/pkg_plugin/common/domain"
	"github.com/kujilabo/cocotola-api/pkg_plugin/handler/converter"
	"github.com/kujilabo/cocotola-api/pkg_plugin/handler/entity"
	user "github.com/kujilabo/cocotola-api/pkg_user/domain"
	"github.com/kujilabo/cocotola-api/pkg_user/handlerhelper"
)

type TranslationHandler interface {
	FindTranslations(c *gin.Context)
	FindTranslationByTextAndPos(c *gin.Context)
}

type translationHandler struct {
	translator domain.Translator
}

func NewTranslationHandler(translator domain.Translator) TranslationHandler {
	return &translationHandler{translator: translator}
}

func (h *translationHandler) FindTranslations(c *gin.Context) {
	ctx := c.Request.Context()
	logger := log.FromContext(ctx)
	logger.Infof("FindTranslations")

	handlerhelper.HandleRoleFunction(c, "Owner", func(organizationID user.OrganizationID, operatorID user.AppUserID) error {

		param := entity.TranslationFindParameter{}
		if err := c.ShouldBindJSON(&param); err != nil {
			c.Status(http.StatusBadRequest)
			return nil
		}

		result, err := h.translator.FindTranslationsByFirstLetter(ctx, app.Lang2JA, param.Letter)
		if err != nil {
			return err
		}

		response, err := converter.ToTranslationFindResposne(ctx, result)
		if err != nil {
			return err
		}

		c.JSON(http.StatusOK, response)
		return nil
	}, h.errorHandle)
}

func (h *translationHandler) FindTranslationByTextAndPos(c *gin.Context) {
	ctx := c.Request.Context()
	logger := log.FromContext(ctx)
	logger.Infof("FindTranslationByTextAndPos")

	handlerhelper.HandleRoleFunction(c, "Owner", func(organizationID user.OrganizationID, operatorID user.AppUserID) error {

		text := ginhelper.GetString(c, "text")

		pos, err := ginhelper.GetInt(c, "pos")
		if err != nil {
			return err
		}

		wordPos, err := domain.NewWordPos(pos)
		if err != nil {
			return err
		}
		result, err := h.translator.FindTranslationByTextAndPos(ctx, app.Lang2JA, text, wordPos)
		if err != nil {
			return err
		}

		response, err := converter.ToTranslationResposne(ctx, result)
		if err != nil {
			return err
		}

		c.JSON(http.StatusOK, response)
		return nil
	}, h.errorHandle)
}

func (h *translationHandler) errorHandle(c *gin.Context, err error) bool {
	ctx := c.Request.Context()
	logger := log.FromContext(ctx)
	logger.Errorf("translationHandler error:%v", err)
	return false
}
