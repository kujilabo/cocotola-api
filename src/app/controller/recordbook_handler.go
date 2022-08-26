package controller

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/xerrors"

	"github.com/kujilabo/cocotola-api/src/app/controller/converter"
	"github.com/kujilabo/cocotola-api/src/app/controller/entity"
	"github.com/kujilabo/cocotola-api/src/app/domain"
	"github.com/kujilabo/cocotola-api/src/app/service"
	studentU "github.com/kujilabo/cocotola-api/src/app/usecase/student"
	"github.com/kujilabo/cocotola-api/src/lib/ginhelper"
	"github.com/kujilabo/cocotola-api/src/lib/log"
	controllerhelper "github.com/kujilabo/cocotola-api/src/user/controller/helper"
	userD "github.com/kujilabo/cocotola-api/src/user/domain"
)

type RecordbookHandler interface {
	FindRecordbook(c *gin.Context)

	SetStudyResult(c *gin.Context)

	GetCompletionRate(c *gin.Context)
}

type recordbookHandler struct {
	studentUsecaseStudy studentU.StudentUsecaseStudy
}

func NewRecordbookHandler(studentUsecaseStudy studentU.StudentUsecaseStudy) RecordbookHandler {
	return &recordbookHandler{
		studentUsecaseStudy: studentUsecaseStudy,
	}
}

// FindRecordbook godoc
// @Summary     Find the recordbook
// @Description find results of workbook
// @Tags        study
// @Produce     json
// @Param       workbookID path string true "Workbook ID"
// @Param       studyType  path string true "Study type"
// @Success     200 {object} entity.ProblemWithLevelList
// @Failure     400
// @Router      /v1/study/workbook/{workbookID}/study_type/{studyType} [get]
func (h *recordbookHandler) FindRecordbook(c *gin.Context) {
	ctx := c.Request.Context()

	controllerhelper.HandleSecuredFunction(c, func(organizationID userD.OrganizationID, operatorID userD.AppUserID) error {
		workbookID, err := ginhelper.GetUintFromPath(c, "workbookID")
		if err != nil {
			c.Status(http.StatusBadRequest)
			return nil
		}
		studyType := ginhelper.GetStringFromPath(c, "studyType")

		result, err := h.studentUsecaseStudy.FindResults(ctx, organizationID, operatorID, domain.WorkbookID(workbookID), studyType)
		if err != nil {
			return err
		}

		response, err := converter.ToProblemWithLevelList(ctx, result)
		if err != nil {
			return err
		}

		c.JSON(http.StatusOK, response)
		return nil
	}, h.errorHandle)
}

func (h *recordbookHandler) SetStudyResult(c *gin.Context) {
	ctx := c.Request.Context()
	logger := log.FromContext(ctx)
	logger.Info("SetStudyResult")

	controllerhelper.HandleSecuredFunction(c, func(organizationID userD.OrganizationID, operatorID userD.AppUserID) error {
		workbookID, err := ginhelper.GetUintFromPath(c, "workbookID")
		if err != nil {
			c.Status(http.StatusBadRequest)
			return nil
		}
		studyType := ginhelper.GetStringFromPath(c, "studyType")
		problemID, err := ginhelper.GetUintFromPath(c, "problemID")
		if err != nil {
			return err
		}

		param := entity.StudyResultParameter{}
		if err := c.ShouldBindJSON(&param); err != nil {
			c.Status(http.StatusBadRequest)
			return nil
		}

		// recordbook, err := h.studyService.FindRecordbook(ctx, organizationID, operatorID, domain.WorkbookID(workbookID), studyType)
		// if err != nil {
		// 	return err
		// }

		if err := h.studentUsecaseStudy.SetResult(ctx, organizationID, operatorID, domain.WorkbookID(workbookID), studyType, domain.ProblemID(problemID), param.Result, param.Memorized); err != nil {
			return xerrors.Errorf("failed to SetResult. err: %w", err)
		}

		c.Status(http.StatusOK)
		return nil
	}, h.errorHandle)
}

// GetCompletionRate godoc
// @Summary     Get the completion rate of the workbook
// @Tags        study
// @Produce     json
// @Param       workbookID path string true "Workbook ID"
// @Param       studyType  path string true "Study type"
// @Success     200 {object} entity.ProblemWithLevelList
// @Failure     400
// @Router      /v1/study/workbook/{workbookID}/study_type/{studyType}/completion_rate [get]
func (h *recordbookHandler) GetCompletionRate(c *gin.Context) {
	ctx := c.Request.Context()
	logger := log.FromContext(ctx)
	logger.Info("FindRecordbook")

	controllerhelper.HandleSecuredFunction(c, func(organizationID userD.OrganizationID, operatorID userD.AppUserID) error {
		workbookID, err := ginhelper.GetUintFromPath(c, "workbookID")
		if err != nil {
			c.Status(http.StatusBadRequest)
			return nil
		}

		results, err := h.studentUsecaseStudy.GetCompletionRate(ctx, organizationID, operatorID, domain.WorkbookID(workbookID))
		if err != nil {
			return err
		}

		logger.Infof("FindRecordbook. response: %+v", results)
		c.JSON(http.StatusOK, results)
		return nil
	}, h.errorHandle)
}

func (h *recordbookHandler) errorHandle(c *gin.Context, err error) bool {
	ctx := c.Request.Context()
	logger := log.FromContext(ctx)
	if errors.Is(err, service.ErrProblemAlreadyExists) {
		c.JSON(http.StatusConflict, gin.H{"message": "Problem already exists"})
		return true
	} else if errors.Is(err, service.ErrWorkbookNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return true
	}
	logger.Errorf("studyHandler error:%v", err)
	return false
}
