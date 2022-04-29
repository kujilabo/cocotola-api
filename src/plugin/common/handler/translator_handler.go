package handler

import (
	"bytes"
	"encoding/csv"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	appD "github.com/kujilabo/cocotola-api/src/app/domain"
	"github.com/kujilabo/cocotola-api/src/lib/ginhelper"
	"github.com/kujilabo/cocotola-api/src/lib/log"
	"github.com/kujilabo/cocotola-api/src/plugin/common/domain"
	"github.com/kujilabo/cocotola-api/src/plugin/common/handler/converter"
	"github.com/kujilabo/cocotola-api/src/plugin/common/handler/entity"
	"github.com/kujilabo/cocotola-api/src/plugin/common/service"
	userD "github.com/kujilabo/cocotola-api/src/user/domain"
	"github.com/kujilabo/cocotola-api/src/user/handlerhelper"
)

type TranslationHandler interface {
	FindTranslations(c *gin.Context)
	FindTranslationByTextAndPos(c *gin.Context)
	FindTranslationsByText(c *gin.Context)
	AddTranslation(c *gin.Context)
	UpdateTranslation(c *gin.Context)
	RemoveTranslation(c *gin.Context)
	ExportTranslations(c *gin.Context)
}

type translationHandler struct {
	translatorClient service.TranslatorClient
}

func NewTranslationHandler(translatorClient service.TranslatorClient) TranslationHandler {
	return &translationHandler{translatorClient: translatorClient}
}

func (h *translationHandler) FindTranslations(c *gin.Context) {
	ctx := c.Request.Context()

	handlerhelper.HandleRoleFunction(c, "Owner", func(organizationID userD.OrganizationID, operatorID userD.AppUserID) error {

		param := entity.TranslationFindParameter{}
		if err := c.ShouldBindJSON(&param); err != nil {
			c.Status(http.StatusBadRequest)
			return nil
		}

		lang2, err := appD.NewLang2(param.Lang2)
		if err != nil {
			c.Status(http.StatusBadRequest)
			return nil
		}

		result, err := h.translatorClient.FindTranslationsByFirstLetter(ctx, lang2, param.Letter)
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

	handlerhelper.HandleRoleFunction(c, "Owner", func(organizationID userD.OrganizationID, operatorID userD.AppUserID) error {

		text := ginhelper.GetStringFromPath(c, "text")

		pos, err := ginhelper.GetIntFromPath(c, "pos")
		if err != nil {
			return err
		}

		wordPos, err := domain.NewWordPos(pos)
		if err != nil {
			return err
		}
		result, err := h.translatorClient.FindTranslationByTextAndPos(ctx, appD.Lang2JA, text, wordPos)
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

func (h *translationHandler) FindTranslationsByText(c *gin.Context) {
	ctx := c.Request.Context()

	handlerhelper.HandleRoleFunction(c, "Owner", func(organizationID userD.OrganizationID, operatorID userD.AppUserID) error {

		text := ginhelper.GetStringFromPath(c, "text")
		results, err := h.translatorClient.FindTranslationsByText(ctx, appD.Lang2JA, text)
		if err != nil {
			return err
		}

		response, err := converter.ToTranslationListResposne(ctx, results)
		if err != nil {
			return err
		}

		c.JSON(http.StatusOK, response)
		return nil
	}, h.errorHandle)
}

func (h *translationHandler) AddTranslation(c *gin.Context) {
	ctx := c.Request.Context()

	handlerhelper.HandleRoleFunction(c, "Owner", func(organizationID userD.OrganizationID, operatorID userD.AppUserID) error {
		param := entity.TranslationAddParameter{}
		if err := c.ShouldBindJSON(&param); err != nil {
			c.Status(http.StatusBadRequest)
			return nil
		}
		parameter, err := converter.ToTranslationAddParameter(ctx, &param)
		if err != nil {
			return err
		}

		if err := h.translatorClient.AddTranslation(ctx, parameter); err != nil {
			return err
		}

		c.Status(http.StatusOK)
		return nil
	}, h.errorHandle)
}

func (h *translationHandler) UpdateTranslation(c *gin.Context) {
	ctx := c.Request.Context()

	handlerhelper.HandleRoleFunction(c, "Owner", func(organizationID userD.OrganizationID, operatorID userD.AppUserID) error {
		text := ginhelper.GetStringFromPath(c, "text")

		pos, err := ginhelper.GetIntFromPath(c, "pos")
		if err != nil {
			return err
		}
		wordPos, err := domain.NewWordPos(pos)
		if err != nil {
			return err
		}

		param := entity.TranslationUpdateParameter{}
		if err := c.ShouldBindJSON(&param); err != nil {
			c.Status(http.StatusBadRequest)
			return nil
		}
		parameter, err := converter.ToTranslationUpdateParameter(ctx, &param)
		if err != nil {
			return err
		}

		if err := h.translatorClient.UpdateTranslation(ctx, appD.Lang2JA, text, wordPos, parameter); err != nil {
			return err
		}

		c.Status(http.StatusOK)
		return nil
	}, h.errorHandle)
}

func (h *translationHandler) RemoveTranslation(c *gin.Context) {
	ctx := c.Request.Context()

	handlerhelper.HandleRoleFunction(c, "Owner", func(organizationID userD.OrganizationID, operatorID userD.AppUserID) error {
		text := ginhelper.GetStringFromPath(c, "text")

		pos, err := ginhelper.GetIntFromPath(c, "pos")
		if err != nil {
			return err
		}
		wordPos, err := domain.NewWordPos(pos)
		if err != nil {
			return err
		}

		if err := h.translatorClient.RemoveTranslation(ctx, appD.Lang2JA, text, wordPos); err != nil {
			return err
		}

		c.Status(http.StatusOK)
		return nil
	}, h.errorHandle)
}

func (h *translationHandler) ExportTranslations(c *gin.Context) {
	handlerhelper.HandleRoleFunction(c, "Owner", func(organizationID userD.OrganizationID, operatorID userD.AppUserID) error {
		csvStruct := [][]string{
			{"name", "address", "phone"},
			{"Ram", "Tokyo", "1236524"},
			{"Shaym", "Beijing", "8575675484"},
		}
		b := new(bytes.Buffer)
		w := csv.NewWriter(b)
		if err := w.WriteAll(csvStruct); err != nil {
			return err
		}
		if _, err := c.Writer.Write(b.Bytes()); err != nil {
			return err
		}
		c.Status(http.StatusOK)
		return nil
	}, h.errorHandle)
}

func (h *translationHandler) errorHandle(c *gin.Context, err error) bool {
	ctx := c.Request.Context()
	logger := log.FromContext(ctx)

	if errors.Is(err, service.ErrTranslationAlreadyExists) {
		logger.Warnf("translationHandler. err: %v", err)
		c.JSON(http.StatusConflict, gin.H{"message": "Translation already exists"})
		return true
	}
	logger.Errorf("translationHandler. err: %v", err)
	return false
}
