package handler

import (
	"errors"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kujilabo/cocotola-api/pkg_lib/log"
	"github.com/kujilabo/cocotola-api/pkg_plugin/common/application"
	"github.com/kujilabo/cocotola-api/pkg_plugin/common/domain"
	"github.com/kujilabo/cocotola-api/pkg_plugin/common/handler/converter"
	"github.com/kujilabo/cocotola-api/pkg_plugin/common/handler/entity"
	user "github.com/kujilabo/cocotola-api/pkg_user/domain"
	"github.com/kujilabo/cocotola-api/pkg_user/handlerhelper"
	"golang.org/x/xerrors"
)

type TatoebaHandler interface {
	FindSentences(c *gin.Context)
	ImportSentences(c *gin.Context)
	ImportLinks(c *gin.Context)
}

type tatoebaHandler struct {
	tatoebaService                       application.TatoebaService
	newTatoebaSentenceAddParameterReader func(reader io.Reader) domain.TatoebaSentenceAddParameterIterator
	newTatoebaLinkAddParameterReader     func(reader io.Reader) domain.TatoebaLinkAddParameterIterator
}

func NewTatoebaHandler(tatoebaService application.TatoebaService, newTatoebaSentenceAddParameterReader func(reader io.Reader) domain.TatoebaSentenceAddParameterIterator, newTatoebaLinkAddParameterReader func(reader io.Reader) domain.TatoebaLinkAddParameterIterator) TatoebaHandler {
	return &tatoebaHandler{
		tatoebaService:                       tatoebaService,
		newTatoebaSentenceAddParameterReader: newTatoebaSentenceAddParameterReader,
		newTatoebaLinkAddParameterReader:     newTatoebaLinkAddParameterReader,
	}
}

func (h *tatoebaHandler) FindSentences(c *gin.Context) {
	ctx := c.Request.Context()
	handlerhelper.HandleRoleFunction(c, "Owner", func(organizationID user.OrganizationID, operatorID user.AppUserID) error {
		param := entity.TatoebaSentenceFindParameter{}
		if err := c.ShouldBindJSON(&param); err != nil {
			c.Status(http.StatusBadRequest)
			return nil
		}
		parameter, err := converter.ToTatoebaSentenceSearchCondition(ctx, &param)
		if err != nil {
			return err
		}
		result, err := h.tatoebaService.FindSentences(ctx, parameter)
		if err != nil {
			return xerrors.Errorf("failed to FindSentences. err: %w", err)
		}
		response, err := converter.ToTatoebaSentenceResponse(ctx, result)
		if err != nil {
			return err
		}

		c.JSON(http.StatusOK, response)
		return nil
	}, h.errorHandle)
}

func (h *tatoebaHandler) ImportSentences(c *gin.Context) {
	ctx := c.Request.Context()
	logger := log.FromContext(ctx)
	handlerhelper.HandleRoleFunction(c, "Owner", func(organizationID user.OrganizationID, operatorID user.AppUserID) error {

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

		iterator := h.newTatoebaSentenceAddParameterReader(multipartFile)

		if err := h.tatoebaService.ImportSentences(ctx, iterator); err != nil {
			return xerrors.Errorf("failed to ImportSentences. err: %w", err)
		}

		c.Status(http.StatusOK)
		return nil
	}, h.errorHandle)
}

func (h *tatoebaHandler) ImportLinks(c *gin.Context) {
	ctx := c.Request.Context()
	logger := log.FromContext(ctx)
	handlerhelper.HandleRoleFunction(c, "Owner", func(organizationID user.OrganizationID, operatorID user.AppUserID) error {

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

		iterator := h.newTatoebaLinkAddParameterReader(multipartFile)

		if err := h.tatoebaService.ImportLinks(ctx, iterator); err != nil {
			return xerrors.Errorf("failed to ImportLinks. err: %w", err)
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
