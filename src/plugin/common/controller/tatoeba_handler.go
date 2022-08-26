package controller

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/xerrors"

	"github.com/kujilabo/cocotola-api/src/lib/log"
	"github.com/kujilabo/cocotola-api/src/plugin/common/controller/converter"
	"github.com/kujilabo/cocotola-api/src/plugin/common/controller/entity"
	"github.com/kujilabo/cocotola-api/src/plugin/common/service"
	controllerhelper "github.com/kujilabo/cocotola-api/src/user/controller/helper"
	userD "github.com/kujilabo/cocotola-api/src/user/domain"
)

type TatoebaHandler interface {
	FindSentencePairs(c *gin.Context)
	ImportSentences(c *gin.Context)
	ImportLinks(c *gin.Context)
}

type tatoebaHandler struct {
	tatoebaClient service.TatoebaClient
}

func NewTatoebaHandler(tatoebaClient service.TatoebaClient) TatoebaHandler {
	return &tatoebaHandler{
		tatoebaClient: tatoebaClient,
	}
}

func (h *tatoebaHandler) FindSentencePairs(c *gin.Context) {
	ctx := c.Request.Context()
	logger := log.FromContext(ctx)

	controllerhelper.HandleFunction(c, func(organizationID userD.OrganizationID, operatorID userD.AppUserID) error {
		param := entity.TatoebaSentenceFindParameter{}
		if err := c.ShouldBindJSON(&param); err != nil {
			logger.Warnf("err: %+v", err)
			c.Status(http.StatusBadRequest)
			return nil
		}
		parameter, err := converter.ToTatoebaSentenceSearchCondition(ctx, &param)
		if err != nil {
			return xerrors.Errorf("failed to ToTatoebaSentenceSearchCondition. err: %w", err)
		}

		result, err := h.tatoebaClient.FindSentencePairs(ctx, parameter)
		if err != nil {
			return xerrors.Errorf("failed to FindSentencePairs. err: %w", err)
		}

		response, err := converter.ToTatoebaSentenceResponse(ctx, result)
		if err != nil {
			return xerrors.Errorf("failed to ToTatoebaSentenceResponse. err: %w", err)
		}

		c.JSON(http.StatusOK, response)
		return nil
	}, h.errorHandle)
}

func (h *tatoebaHandler) ImportSentences(c *gin.Context) {
	ctx := c.Request.Context()
	logger := log.FromContext(ctx)
	controllerhelper.HandleRoleFunction(c, "Owner", func(organizationID userD.OrganizationID, operatorID userD.AppUserID) error {

		file, err := c.FormFile("file")
		if err != nil {
			if errors.Is(err, http.ErrMissingFile) {
				logger.Warnf("err: %+v", err)
				c.Status(http.StatusBadRequest)
				return nil
			}
			return err
		}

		multipartFile, err := file.Open()
		if err != nil {
			return xerrors.Errorf("failed to file.Open. err: %w", err)
		}
		defer multipartFile.Close()

		if err := h.tatoebaClient.ImportSentences(ctx, multipartFile); err != nil {
			return xerrors.Errorf("failed to ImportSentences. err: %w", err)
		}

		c.Status(http.StatusOK)
		return nil
	}, h.errorHandle)
}

func (h *tatoebaHandler) ImportLinks(c *gin.Context) {
	ctx := c.Request.Context()
	logger := log.FromContext(ctx)
	controllerhelper.HandleRoleFunction(c, "Owner", func(organizationID userD.OrganizationID, operatorID userD.AppUserID) error {

		file, err := c.FormFile("file")
		if err != nil {
			if errors.Is(err, http.ErrMissingFile) {
				logger.Warnf("err: %+v", err)
				c.Status(http.StatusBadRequest)
				return nil
			}
			return err
		}

		multipartFile, err := file.Open()
		if err != nil {
			return xerrors.Errorf("failed to file.Open. err: %w", err)
		}
		defer multipartFile.Close()

		if err := h.tatoebaClient.ImportSentences(ctx, multipartFile); err != nil {
			return err
		}

		c.Status(http.StatusOK)
		return nil
	}, h.errorHandle)
}

func (h *tatoebaHandler) errorHandle(c *gin.Context, err error) bool {
	ctx := c.Request.Context()
	logger := log.FromContext(ctx)
	logger.Errorf("tatoebaHandler. err: %v", err)
	return false
}
