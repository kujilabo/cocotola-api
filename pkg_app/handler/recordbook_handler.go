package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/xerrors"

	"github.com/kujilabo/cocotola-api/pkg_app/application"
	"github.com/kujilabo/cocotola-api/pkg_app/domain"
	"github.com/kujilabo/cocotola-api/pkg_app/handler/converter"
	"github.com/kujilabo/cocotola-api/pkg_app/handler/entity"
	"github.com/kujilabo/cocotola-api/pkg_lib/ginhelper"
	"github.com/kujilabo/cocotola-api/pkg_lib/log"
	user "github.com/kujilabo/cocotola-api/pkg_user/domain"
	"github.com/kujilabo/cocotola-api/pkg_user/handlerhelper"
)

type RecordbookHandler interface {
	FindRecordbook(c *gin.Context)
	SetStudyResult(c *gin.Context)
}

type recordbookHandler struct {
	studyService application.StudyService
}

func NewRecordbookHandler(studyService application.StudyService) RecordbookHandler {
	return &recordbookHandler{
		studyService: studyService,
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
	logger := log.FromContext(ctx)
	logger.Info("FindRecordbook")

	handlerhelper.HandleSecuredFunction(c, func(organizationID user.OrganizationID, operatorID user.AppUserID) error {
		workbookID, err := ginhelper.GetUint(c, "workbookID")
		if err != nil {
			c.Status(http.StatusBadRequest)
			return nil
		}
		studyType := ginhelper.GetString(c, "studyType")

		result, err := h.studyService.FindResults(ctx, organizationID, operatorID, domain.WorkbookID(workbookID), studyType)
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

	handlerhelper.HandleSecuredFunction(c, func(organizationID user.OrganizationID, operatorID user.AppUserID) error {
		workbookID, err := ginhelper.GetUint(c, "workbookID")
		if err != nil {
			c.Status(http.StatusBadRequest)
			return nil
		}
		studyType := ginhelper.GetString(c, "studyType")
		problemID, err := ginhelper.GetUint(c, "problemID")
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

		if err := h.studyService.SetResult(ctx, organizationID, operatorID, domain.WorkbookID(workbookID), studyType, domain.ProblemID(problemID), param.Result, param.Memorized); err != nil {
			return xerrors.Errorf("failed to SetResult. err: %w", err)
		}

		c.Status(http.StatusOK)
		return nil
	}, h.errorHandle)
}

func (h *recordbookHandler) errorHandle(c *gin.Context, err error) bool {
	ctx := c.Request.Context()
	logger := log.FromContext(ctx)
	if xerrors.Is(err, domain.ErrProblemAlreadyExists) {
		c.JSON(http.StatusConflict, gin.H{"message": "Problem already exists"})
		return true
	} else if xerrors.Is(err, domain.ErrWorkbookNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return true
	}
	logger.Errorf("studyHandler error:%v", err)
	return false
}
