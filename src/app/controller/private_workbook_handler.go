package controller

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/kujilabo/cocotola-api/src/app/controller/converter"
	"github.com/kujilabo/cocotola-api/src/app/controller/entity"
	"github.com/kujilabo/cocotola-api/src/app/domain"
	"github.com/kujilabo/cocotola-api/src/app/service"
	studentU "github.com/kujilabo/cocotola-api/src/app/usecase/student"
	liberrors "github.com/kujilabo/cocotola-api/src/lib/errors"
	"github.com/kujilabo/cocotola-api/src/lib/ginhelper"
	"github.com/kujilabo/cocotola-api/src/lib/log"
	controllerhelper "github.com/kujilabo/cocotola-api/src/user/controller/helper"
	userD "github.com/kujilabo/cocotola-api/src/user/domain"
)

type PrivateWorkbookHandler interface {
	FindWorkbooks(c *gin.Context)
	FindWorkbookByID(c *gin.Context)
	AddWorkbook(c *gin.Context)
	UpdateWorkbook(c *gin.Context)
	RemoveWorkbook(c *gin.Context)
}

type privateWorkbookHandler struct {
	// repository             gateway.Repository
	studentUsecaseWorkbook studentU.StudentUsecaseWorkbook
}

func NewPrivateWorkbookHandler(studentUsecaseWorkbook studentU.StudentUsecaseWorkbook) PrivateWorkbookHandler {
	return &privateWorkbookHandler{
		studentUsecaseWorkbook: studentUsecaseWorkbook,
	}
}

// FindWorkbooks godoc
// @Summary Find workbooks
// @Produce json
// @Success 200 {object} entity.WorkbookSearchResponse
// @Failure 400
// @Router /v1/private/workbook/search [post]
func (h *privateWorkbookHandler) FindWorkbooks(c *gin.Context) {
	ctx := c.Request.Context()
	logger := log.FromContext(ctx)
	logger.Info("FindWorkbooks")

	id := c.Param("workbookID")
	if id != "search" {
		c.Status(http.StatusNotFound)
		return
	}

	controllerhelper.HandleSecuredFunction(c, func(organizationID userD.OrganizationID, operatorID userD.AppUserID) error {
		result, err := h.studentUsecaseWorkbook.FindWorkbooks(ctx, organizationID, operatorID)
		if err != nil {
			return err
		}

		response, err := converter.ToWorkbookSearchResponse(result)
		if err != nil {
			return err
		}
		c.JSON(http.StatusOK, response)
		return nil
	}, h.errorHandle)
}

func (h *privateWorkbookHandler) FindWorkbookByID(c *gin.Context) {
	ctx := c.Request.Context()
	logger := log.FromContext(ctx)
	logger.Info("FindWorkbookByID")

	controllerhelper.HandleSecuredFunction(c, func(organizationID userD.OrganizationID, operatorID userD.AppUserID) error {
		id := c.Param("workbookID")
		workbookID, err := strconv.Atoi(id)
		if err != nil {
			c.Status(http.StatusBadRequest)
			return nil
		}

		workbook, err := h.studentUsecaseWorkbook.FindWorkbookByID(ctx, organizationID, operatorID, domain.WorkbookID(uint(workbookID)))
		if err != nil {
			return liberrors.Errorf("failed to FindWorkbookByID. err: %w", err)
		}

		workbookResponse, err := converter.ToWorkbookHTTPEntity(workbook)
		if err != nil {
			return liberrors.Errorf("failed to ToWorkbookHTTPEntity. err: %w", err)
		}

		c.JSON(http.StatusOK, workbookResponse)
		return nil
	}, h.errorHandle)
}

// AddWorkbook godoc
// @Summary Create new workbook
// @Produce json
// @Param param body entity.WorkbookAddParameter true "parameter to create new workbook"
// @Success 200 {object} controllerhelper.IDResponse
// @Failure 400
// @Router /v1/private/workbook [post]
func (h *privateWorkbookHandler) AddWorkbook(c *gin.Context) {
	ctx := c.Request.Context()
	logger := log.FromContext(ctx)
	logger.Info("AddWokrbook")

	controllerhelper.HandleSecuredFunction(c, func(organizationID userD.OrganizationID, operatorID userD.AppUserID) error {
		param := entity.WorkbookAddParameter{}
		if err := c.ShouldBindJSON(&param); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			logger.Warnf("failed to BindJSON. err: %v", err)
			return nil
		}

		parameter, err := converter.ToWorkbookAddParameter(&param)
		if err != nil {
			return liberrors.Errorf("failed to ToAdd. err: %w", err)
		}

		workbookID, err := h.studentUsecaseWorkbook.AddWorkbook(ctx, organizationID, operatorID, parameter)
		if err != nil {
			return liberrors.Errorf("failed to addWorkbook. err: %w", err)
		}

		c.JSON(http.StatusOK, controllerhelper.IDResponse{ID: uint(workbookID)})
		return nil
	}, h.errorHandle)
}

// UpdateWorkbook godoc
// @Summary     Update the workbook
// @Description update the workbook
// @Tags        private workbook
// @Accept      json
// @Produce     json
// @Param       workbookID path int true "Workbook ID"
// @Param       param body entity.WorkbookUpdateParameter true "parameter to update the workbook"
// @Success     200 {object} controllerhelper.IDResponse
// @Failure     400
// @Router      /v1/private/workbook/{workbookID} [put]
func (h *privateWorkbookHandler) UpdateWorkbook(c *gin.Context) {
	ctx := c.Request.Context()
	logger := log.FromContext(ctx)
	logger.Info("UpdateWorkbook")

	controllerhelper.HandleSecuredFunction(c, func(organizationID userD.OrganizationID, operatorID userD.AppUserID) error {
		param := entity.WorkbookUpdateParameter{}
		if err := c.BindJSON(&param); err != nil {
			logger.Warnf("failed to BindJSON. err: %v", err)
			return nil
		}
		workbookID, err := ginhelper.GetUintFromPath(c, "workbookID")
		if err != nil {
			c.Status(http.StatusBadRequest)
			return nil
		}

		version, err := ginhelper.GetIntFromQuery(c, "version")
		if err != nil {
			c.Status(http.StatusBadRequest)
			return nil
		}

		parameter, err := converter.ToWorkbookUpdateParameter(&param)
		if err != nil {
			return err
		}

		if err := h.studentUsecaseWorkbook.UpdateWorkbook(ctx, organizationID, operatorID, domain.WorkbookID(workbookID), version, parameter); err != nil {
			logger.WithError(err).Errorf("failed to UpdateWorkbook. err: %v", err)
			return err
		}

		c.JSON(http.StatusOK, controllerhelper.IDResponse{ID: workbookID})
		return nil
	}, h.errorHandle)
}

func (h *privateWorkbookHandler) RemoveWorkbook(c *gin.Context) {
	ctx := c.Request.Context()
	logger := log.FromContext(ctx)
	logger.Info("RemoveWorkbook")

	controllerhelper.HandleSecuredFunction(c, func(organizationID userD.OrganizationID, operatorID userD.AppUserID) error {
		workbookID, err := ginhelper.GetUintFromPath(c, "workbookID")
		if err != nil {
			c.Status(http.StatusBadRequest)
			return nil
		}

		version, err := ginhelper.GetIntFromQuery(c, "version")
		if err != nil {
			c.Status(http.StatusBadRequest)
			return nil
		}

		if err := h.studentUsecaseWorkbook.RemoveWorkbook(ctx, organizationID, operatorID, domain.WorkbookID(workbookID), version); err != nil {
			logger.WithError(err).Errorf("failed to RemoveWorkbook. err: %v", err)
			return err
		}

		c.Status(http.StatusOK)
		return nil
	}, h.errorHandle)
}

func (h *privateWorkbookHandler) errorHandle(c *gin.Context, err error) bool {
	ctx := c.Request.Context()
	logger := log.FromContext(ctx)
	fmt.Println(err)
	if errors.Is(err, service.ErrWorkbookAlreadyExists) {
		logger.Warnf("workbookHandler err: %+v", err)
		c.JSON(http.StatusConflict, gin.H{"message": "Workbook already exists"})
		return true
	} else if errors.Is(err, service.ErrWorkbookNotFound) {
		logger.Warnf("workbookHandler err: %+v", err)
		c.JSON(http.StatusNotFound, gin.H{"message": "Workbook not found"})
		return true
	}
	logger.Errorf("workbookHandler err: %+v", err)
	return false
}
